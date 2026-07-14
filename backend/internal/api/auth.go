package api

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/Nikunjsaini07/performx/backend/internal/db"
)

type AuthHandler struct {
	Queries   *db.Queries
	JWTSecret []byte
}

type RegisterRequest struct {
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}

type VerifyOTPRequest struct {
	Email   string `json:"email"`
	OTPCode string `json:"otp_code"`
	Purpose string `json:"purpose"` // 'REGISTER', 'LOGIN'
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type ResetPasswordRequest struct {
	Email       string `json:"email"`
	OTPCode     string `json:"otp_code"`
	NewPassword string `json:"new_password"`
}

type AuthResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	User         UserPayload `json:"user"`
}

type UserPayload struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Bio         string `json:"bio"`
	AvatarURL   string `json:"avatar_url"`
}

// Helper to generate a secure random 6-digit OTP code
func generateOTP() string {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "123456" // safe fallback
	}
	return fmt.Sprintf("%06d", n.Int64())
}

// Helper to generate a cryptographically secure random token (e.g. for Refresh Tokens)
func generateRandomToken() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x", b)
}

// Helper to generate a JWT access token
func (h *AuthHandler) generateAccessToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(15 * time.Minute).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(h.JWTSecret)
}

// POST /auth/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid request body"}`))
		return
	}

	if req.Username == "" || req.DisplayName == "" || req.Email == "" || req.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"All fields are required"}`))
		return
	}

	// Check if user already exists
	_, err := h.Queries.GetUserByEmail(r.Context(), req.Email)
	if err == nil {
		w.WriteHeader(http.StatusConflict)
		_, _ = w.Write([]byte(`{"error":"Conflict","message":"User with this email already exists"}`))
		return
	}
	_, err = h.Queries.GetUserByUsername(r.Context(), req.Username)
	if err == nil {
		w.WriteHeader(http.StatusConflict)
		_, _ = w.Write([]byte(`{"error":"Conflict","message":"User with this username already exists"}`))
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to process password"}`))
		return
	}

	// Create user (unverified by default)
	user, err := h.Queries.CreateUser(r.Context(), db.CreateUserParams{
		Username:      req.Username,
		DisplayName:   req.DisplayName,
		Email:         req.Email,
		Bio:           pgtype.Text{Valid: false},
		AvatarUrl:     pgtype.Text{Valid: false},
		PasswordHash:  pgtype.Text{String: string(hashedPassword), Valid: true},
		EmailVerified: false,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"error":"Internal Error","message":"Failed to create user: %s"}`, err.Error())))
		return
	}

	// Generate and store OTP
	otp := generateOTP()
	expires := time.Now().Add(15 * time.Minute)
	err = h.Queries.CreateOTP(r.Context(), db.CreateOTPParams{
		Email:     req.Email,
		OtpCode:   otp,
		Purpose:   "REGISTER",
		ExpiresAt: pgtype.Timestamptz{Time: expires, Valid: true},
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to generate verification code"}`))
		return
	}

	// Send OTP via Brevo email
	if err := SendOTPEmail(req.Email, req.DisplayName, otp, "REGISTER"); err != nil {
		fmt.Printf("[WARN] Failed to send registration OTP email to %s: %v\n", req.Email, err)
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User registered successfully. Verification code sent to email.",
		"email":   user.Email,
	})
}

// POST /auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid request body"}`))
		return
	}

	user, err := h.Queries.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"Unauthorized","message":"Invalid email or password"}`))
		return
	}

	if !user.PasswordHash.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"Unauthorized","message":"Password login not set up for this user"}`))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash.String), []byte(req.Password))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"Unauthorized","message":"Invalid email or password"}`))
		return
	}

	// Generate and store login verification OTP
	otp := generateOTP()
	expires := time.Now().Add(15 * time.Minute)
	err = h.Queries.CreateOTP(r.Context(), db.CreateOTPParams{
		Email:     user.Email,
		OtpCode:   otp,
		Purpose:   "LOGIN",
		ExpiresAt: pgtype.Timestamptz{Time: expires, Valid: true},
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to generate verification code"}`))
		return
	}

	// Send OTP via Brevo email
	if err := SendOTPEmail(user.Email, user.DisplayName, otp, "LOGIN"); err != nil {
		fmt.Printf("[WARN] Failed to send login OTP email to %s: %v\n", user.Email, err)
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "OTP_REQUIRED",
		"message": "Verification code sent to email",
		"email":   user.Email,
	})
}

// POST /auth/verify-otp
func (h *AuthHandler) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req VerifyOTPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid request body"}`))
		return
	}

	if req.Purpose != "REGISTER" && req.Purpose != "LOGIN" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid purpose"}`))
		return
	}

	// Retrieve user
	user, err := h.Queries.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"Not Found","message":"User not found"}`))
		return
	}

	// Fetch valid OTP from DB
	otp, err := h.Queries.GetValidOTP(r.Context(), db.GetValidOTPParams{
		Email:   req.Email,
		OtpCode: req.OTPCode,
		Purpose: req.Purpose,
	})
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"Unauthorized","message":"Invalid or expired verification code"}`))
		return
	}

	// Delete verified OTP code
	_ = h.Queries.DeleteOTP(r.Context(), otp.ID)

	// If Register, update verification status
	if req.Purpose == "REGISTER" {
		_ = h.Queries.VerifyUserEmail(r.Context(), req.Email)
		user.EmailVerified = true
	}

	// Generate Access Token (JWT)
	accessToken, err := h.generateAccessToken(user.ID.String())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to generate access token"}`))
		return
	}

	// Generate and save Refresh Token
	refreshToken := generateRandomToken()
	err = h.Queries.CreateRefreshToken(r.Context(), db.CreateRefreshTokenParams{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: pgtype.Timestamptz{Time: time.Now().Add(7 * 24 * time.Hour), Valid: true},
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to save session"}`))
		return
	}

	resp := AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: UserPayload{
			ID:          user.ID.String(),
			Username:    user.Username,
			DisplayName: user.DisplayName,
			Email:       user.Email,
			Bio:         user.Bio.String,
			AvatarURL:   user.AvatarUrl.String,
		},
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

// POST /auth/refresh
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid request body"}`))
		return
	}

	dbToken, err := h.Queries.GetRefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"Unauthorized","message":"Invalid or expired refresh token"}`))
		return
	}

	// Generate new access token
	newAccessToken, err := h.generateAccessToken(dbToken.UserID.String())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to generate access token"}`))
		return
	}

	// Rotate refresh token (highly secure!)
	newRefreshToken := generateRandomToken()

	// Begin Transaction for rotation
	err = h.Queries.DeleteRefreshToken(r.Context(), dbToken.Token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Session rotation failed"}`))
		return
	}

	err = h.Queries.CreateRefreshToken(r.Context(), db.CreateRefreshTokenParams{
		UserID:    dbToken.UserID,
		Token:     newRefreshToken,
		ExpiresAt: pgtype.Timestamptz{Time: time.Now().Add(7 * 24 * time.Hour), Valid: true},
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Session rotation failed"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	})
}

// POST /auth/logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Try to read it from context if it's protected or let them pass in body
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Refresh token is required to log out"}`))
		return
	}

	// Revoke the token session from the database
	_ = h.Queries.DeleteRefreshToken(r.Context(), req.RefreshToken)

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"message":"Successfully logged out"}`))
}

// POST /auth/forgot-password
func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid request body"}`))
		return
	}

	// Verify user exists
	_, err := h.Queries.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		// Return 200 to prevent user enumeration attacks, but do not generate OTP
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message":"If the email exists, a password reset code has been sent."}`))
		return
	}

	// Generate forgot password OTP
	otp := generateOTP()
	expires := time.Now().Add(15 * time.Minute)
	err = h.Queries.CreateOTP(r.Context(), db.CreateOTPParams{
		Email:     req.Email,
		OtpCode:   otp,
		Purpose:   "FORGOT_PASSWORD",
		ExpiresAt: pgtype.Timestamptz{Time: expires, Valid: true},
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to generate reset code"}`))
		return
	}

	// Send OTP via Brevo email
	if err := SendOTPEmail(req.Email, "", otp, "FORGOT_PASSWORD"); err != nil {
		fmt.Printf("[WARN] Failed to send forgot-password OTP email to %s: %v\n", req.Email, err)
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "If the email exists, a password reset code has been sent.",
		"email":   req.Email,
	})
}

// POST /auth/reset-password
func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid request body"}`))
		return
	}

	if req.NewPassword == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"New password is required"}`))
		return
	}

	// Verify valid reset OTP
	otp, err := h.Queries.GetValidOTP(r.Context(), db.GetValidOTPParams{
		Email:   req.Email,
		OtpCode: req.OTPCode,
		Purpose: "FORGOT_PASSWORD",
	})
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"Unauthorized","message":"Invalid or expired verification code"}`))
		return
	}

	// Delete verified OTP code
	_ = h.Queries.DeleteOTP(r.Context(), otp.ID)

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to process password"}`))
		return
	}

	// Update user password and invalidate refresh token sessions for security
	err = h.Queries.UpdateUserPassword(r.Context(), db.UpdateUserPasswordParams{
		Email:        req.Email,
		PasswordHash: pgtype.Text{String: string(hashedPassword), Valid: true},
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to update password"}`))
		return
	}

	// Optional: Revoke all existing refresh sessions for security
	user, err := h.Queries.GetUserByEmail(r.Context(), req.Email)
	if err == nil {
		_ = h.Queries.DeleteUserRefreshTokens(r.Context(), user.ID)
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"message":"Password has been reset successfully. Please log in with your new password."}`))
}

// Helper to get Google OAuth config
func getGoogleOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("API_URL") + "/auth/google/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

// GoogleUserInfo represents the response from Google's userinfo endpoint
type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}

// GET /auth/google/callback
func (h *AuthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	// Get authorization code from query params
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Redirect(w, r, os.Getenv("FRONTEND_URL")+"?auth_error=missing_code", http.StatusTemporaryRedirect)
		return
	}

	// Exchange code for token
	config := getGoogleOAuthConfig()
	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		fmt.Printf("Error exchanging code for token: %v\n", err)
		http.Redirect(w, r, os.Getenv("FRONTEND_URL")+"?auth_error=token_exchange_failed", http.StatusTemporaryRedirect)
		return
	}

	// Get user info from Google
	client := config.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		fmt.Printf("Error getting user info: %v\n", err)
		http.Redirect(w, r, os.Getenv("FRONTEND_URL")+"?auth_error=user_info_failed", http.StatusTemporaryRedirect)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		http.Redirect(w, r, os.Getenv("FRONTEND_URL")+"?auth_error=read_failed", http.StatusTemporaryRedirect)
		return
	}

	var googleUser GoogleUserInfo
	if err := json.Unmarshal(body, &googleUser); err != nil {
		fmt.Printf("Error unmarshaling user info: %v\n", err)
		http.Redirect(w, r, os.Getenv("FRONTEND_URL")+"?auth_error=parse_failed", http.StatusTemporaryRedirect)
		return
	}

	// Check if user exists by email
	existingUser, err := h.Queries.GetUserByEmail(r.Context(), googleUser.Email)
	if err != nil {
		// User doesn't exist, create new user
		// Generate a unique username from email (part before @)
		atIndex := 0
		for i, ch := range googleUser.Email {
			if ch == '@' {
				atIndex = i
				break
			}
		}
		username := googleUser.Email[:atIndex]
		
		// Check if username exists and append random suffix if needed
		_, usernameErr := h.Queries.GetUserByUsername(r.Context(), username)
		if usernameErr == nil {
			// Username exists, append random number
			randomNum, _ := rand.Int(rand.Reader, big.NewInt(9999))
			username = fmt.Sprintf("%s%d", username, randomNum.Int64())
		}

		// Create user with Google OAuth (no password, email verified by default)
		createdUser, err := h.Queries.CreateUser(r.Context(), db.CreateUserParams{
			Username:      username,
			DisplayName:   googleUser.Name,
			Email:         googleUser.Email,
			Bio:           pgtype.Text{Valid: false},
			AvatarUrl:     pgtype.Text{String: googleUser.Picture, Valid: true},
			PasswordHash:  pgtype.Text{Valid: false}, // No password for OAuth users
			EmailVerified: true,                      // Google emails are already verified
		})

		if err != nil {
			fmt.Printf("Error creating user: %v\n", err)
			http.Redirect(w, r, os.Getenv("FRONTEND_URL")+"?auth_error=user_creation_failed", http.StatusTemporaryRedirect)
			return
		}

		// Convert CreateUserRow to GetUserByEmailRow
		existingUser = db.GetUserByEmailRow{
			ID:            createdUser.ID,
			Username:      createdUser.Username,
			DisplayName:   createdUser.DisplayName,
			Email:         createdUser.Email,
			Bio:           createdUser.Bio,
			AvatarUrl:     createdUser.AvatarUrl,
			PasswordHash:  pgtype.Text{Valid: false}, // Not returned by CreateUser
			EmailVerified: createdUser.EmailVerified,
			CreatedAt:     createdUser.CreatedAt,
			UpdatedAt:     createdUser.UpdatedAt,
		}
	} else {
		// User exists - update avatar if needed
		if googleUser.Picture != "" && (!existingUser.AvatarUrl.Valid || existingUser.AvatarUrl.String == "") {
			updatedUser, err := h.Queries.UpdateAvatar(r.Context(), db.UpdateAvatarParams{
				ID:        existingUser.ID,
				AvatarUrl: pgtype.Text{String: googleUser.Picture, Valid: true},
			})
			if err == nil {
				existingUser.AvatarUrl = updatedUser.AvatarUrl
			}
		}

		// Mark email as verified if not already
		if !existingUser.EmailVerified {
			_ = h.Queries.VerifyUserEmail(r.Context(), existingUser.Email)
			existingUser.EmailVerified = true
		}
	}

	// Generate Access Token (JWT)
	accessToken, err := h.generateAccessToken(existingUser.ID.String())
	if err != nil {
		fmt.Printf("Error generating access token: %v\n", err)
		http.Redirect(w, r, os.Getenv("FRONTEND_URL")+"?auth_error=token_generation_failed", http.StatusTemporaryRedirect)
		return
	}

	// Generate and save Refresh Token
	refreshToken := generateRandomToken()
	err = h.Queries.CreateRefreshToken(r.Context(), db.CreateRefreshTokenParams{
		UserID:    existingUser.ID,
		Token:     refreshToken,
		ExpiresAt: pgtype.Timestamptz{Time: time.Now().Add(7 * 24 * time.Hour), Valid: true},
	})
	if err != nil {
		fmt.Printf("Error saving refresh token: %v\n", err)
		http.Redirect(w, r, os.Getenv("FRONTEND_URL")+"?auth_error=session_failed", http.StatusTemporaryRedirect)
		return
	}

	// Redirect to frontend with tokens
	frontendURL := os.Getenv("FRONTEND_URL")
	redirectURL := fmt.Sprintf("%s/auth/callback?access_token=%s&refresh_token=%s&user_id=%s&username=%s&display_name=%s&email=%s",
		frontendURL,
		accessToken,
		refreshToken,
		existingUser.ID.String(),
		existingUser.Username,
		existingUser.DisplayName,
		existingUser.Email,
	)

	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}
