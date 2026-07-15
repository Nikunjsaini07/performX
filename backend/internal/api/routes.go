package api

import (
	"encoding/json"
	"net/http"

	"github.com/Nikunjsaini07/performx/backend/internal/db"
	"github.com/Nikunjsaini07/performx/backend/internal/handler"
	"github.com/Nikunjsaini07/performx/backend/internal/middleware"
)

// RegisterRoutes registers all public and protected routes to the main ServeMux router.
func RegisterRoutes(mux *http.ServeMux, queries *db.Queries, jwtSecret []byte, adminPrefix string) {
	// Initialize all sub-handlers
	authHandler := &handler.AuthHandler{
		Queries:   queries,
		JWTSecret: jwtSecret,
	}
	meHandler := &handler.MeHandler{
		Queries: queries,
	}
	usersHandler := &handler.UsersHandler{
		Queries:   queries,
		JWTSecret: jwtSecret,
	}
	playersHandler := &handler.PlayersHandler{
		Queries: queries,
	}
	teamsHandler := &handler.TeamsHandler{
		Queries: queries,
	}
	metadataHandler := &handler.MetadataHandler{
		Queries: queries,
	}
	matchesHandler := &handler.MatchesHandler{
		Queries: queries,
	}
	performancesHandler := &handler.PerformancesHandler{
		Queries: queries,
	}
	ratingsHandler := &handler.RatingsHandler{
		Queries: queries,
	}
	reviewsHandler := &handler.ReviewsHandler{
		Queries: queries,
	}
	commentsHandler := &handler.CommentsHandler{
		Queries: queries,
	}
	likesHandler := &handler.LikesHandler{
		Queries: queries,
	}
	importHandler := &handler.ImportHandler{
		Queries: queries,
	}

	// ----------------------------------------------------
	// 1. Public Auth Endpoints
	// ----------------------------------------------------
	mux.HandleFunc("POST /auth/register", authHandler.Register)
	mux.HandleFunc("POST /auth/login", authHandler.Login)
	mux.HandleFunc("POST /auth/verify-otp", authHandler.VerifyOTP)
	mux.HandleFunc("POST /auth/refresh", authHandler.Refresh)
	mux.HandleFunc("POST /auth/forgot-password", authHandler.ForgotPassword)
	mux.HandleFunc("POST /auth/reset-password", authHandler.ResetPassword)
	mux.HandleFunc("GET /auth/google/callback", authHandler.GoogleCallback)

	// ----------------------------------------------------
	// 2. Public User Endpoints
	// ----------------------------------------------------
	mux.HandleFunc("GET /users", usersHandler.SearchUsers)
	mux.HandleFunc("GET /users/{username}", usersHandler.GetProfile)
	mux.HandleFunc("GET /users/{username}/reviews", usersHandler.GetReviews)
	mux.HandleFunc("GET /users/{username}/ratings", usersHandler.GetRatings)
	mux.HandleFunc("GET /users/{username}/activity", usersHandler.GetActivity)
	mux.HandleFunc("GET /users/{username}/comments", commentsHandler.GetUserComments)

	// ----------------------------------------------------
	// 3. Public Player Endpoints
	// ----------------------------------------------------
	mux.HandleFunc("GET /players", playersHandler.ListPlayers)
	mux.HandleFunc("GET /players/search", playersHandler.SearchPlayers)
	mux.HandleFunc("GET /players/{slug}", playersHandler.GetPlayer)
	mux.HandleFunc("GET /players/{slug}/career", playersHandler.GetPlayerCareer)
	mux.HandleFunc("GET /players/{slug}/current-team", playersHandler.GetPlayerCurrentTeam)
	mux.HandleFunc("GET /players/{slug}/performances", playersHandler.GetPlayerPerformances)
	mux.HandleFunc("GET /players/{slug}/stats", playersHandler.GetPlayerStats)
	mux.HandleFunc("GET /players/{slug}/reviews", playersHandler.GetPlayerReviews)
	mux.HandleFunc("GET /players/{slug}/ratings", playersHandler.GetPlayerRatings)

	// ----------------------------------------------------
	// 4. Public Team Endpoints
	// ----------------------------------------------------
	mux.HandleFunc("GET /teams", teamsHandler.ListTeams)
	mux.HandleFunc("GET /teams/search", teamsHandler.SearchTeams)
	mux.HandleFunc("GET /teams/{slug}", teamsHandler.GetTeam)
	mux.HandleFunc("GET /teams/{slug}/players", teamsHandler.GetTeamPlayers)
	mux.HandleFunc("GET /teams/{slug}/matches", teamsHandler.GetTeamMatches)
	mux.HandleFunc("GET /teams/{slug}/performances", teamsHandler.GetTeamPerformances)
	mux.HandleFunc("GET /teams/{slug}/stats", teamsHandler.GetTeamStats)
	mux.HandleFunc("GET /teams/{slug}/reviews", teamsHandler.GetTeamReviews)
	mux.HandleFunc("GET /teams/{slug}/ratings", teamsHandler.GetTeamRatings)

	// ----------------------------------------------------
	// 5. Public Match Endpoints
	// ----------------------------------------------------
	mux.HandleFunc("GET /matches", matchesHandler.ListMatches)
	mux.HandleFunc("GET /matches/search", matchesHandler.SearchMatches)
	mux.HandleFunc("GET /matches/upcoming", matchesHandler.GetUpcomingMatches)
	mux.HandleFunc("GET /matches/completed", matchesHandler.GetCompletedMatches)
	mux.HandleFunc("GET /matches/{slug}", matchesHandler.GetMatch)
	mux.HandleFunc("GET /matches/{slug}/stats", matchesHandler.GetMatchStats)
	mux.HandleFunc("GET /matches/{slug}/performances", matchesHandler.GetMatchPerformances)
	mux.HandleFunc("GET /matches/{slug}/reviews", matchesHandler.GetMatchReviews)
	mux.HandleFunc("GET /matches/{slug}/reviews/{reviewId}", reviewsHandler.GetMatchReview)
	mux.HandleFunc("GET /matches/{slug}/ratings", matchesHandler.GetMatchRatings)
	mux.HandleFunc("GET /match-reviews/{reviewId}/comments", commentsHandler.GetMatchComments)

	// ----------------------------------------------------
	// 6. Public Performance Endpoints
	// ----------------------------------------------------
	mux.HandleFunc("GET /performances", performancesHandler.ListPerformances)
	mux.HandleFunc("GET /performances/search", performancesHandler.SearchPerformances)
	mux.HandleFunc("GET /performances/top-rated", performancesHandler.GetTopRatedPerformances)
	mux.HandleFunc("GET /performances/recent", performancesHandler.GetRecentPerformances)
	mux.HandleFunc("GET /performances/{id}", performancesHandler.GetPerformance)
	mux.HandleFunc("GET /performances/{id}/stats", performancesHandler.GetPerformanceStats)
	mux.HandleFunc("GET /performances/{id}/reviews", performancesHandler.GetPerformanceReviews)
	mux.HandleFunc("GET /performances/{id}/reviews/{reviewId}", reviewsHandler.GetPerformanceReview)
	mux.HandleFunc("GET /performances/{id}/ratings", performancesHandler.GetPerformanceRatings)
	mux.HandleFunc("GET /performance-reviews/{reviewId}/comments", commentsHandler.GetPerformanceComments)

	// ----------------------------------------------------
	// 7. Public Likes Endpoints
	// ----------------------------------------------------
	mux.HandleFunc("GET /match-reviews/{reviewId}/likes", likesHandler.GetMatchReviewLikes)
	mux.HandleFunc("GET /performance-reviews/{reviewId}/likes", likesHandler.GetPerformanceReviewLikes)
	mux.HandleFunc("GET /match-review-comments/{commentId}/likes", likesHandler.GetMatchCommentLikes)
	mux.HandleFunc("GET /performance-review-comments/{commentId}/likes", likesHandler.GetPerformanceCommentLikes)

	// ----------------------------------------------------
	// 8. Public Trending Endpoints
	// ----------------------------------------------------
	trendingHandler := &handler.TrendingHandler{
		Queries: queries,
	}
	mux.HandleFunc("GET /trending/performances", trendingHandler.GetTrendingPerformances)
	mux.HandleFunc("GET /trending/players", trendingHandler.GetTrendingPlayers)
	mux.HandleFunc("GET /trending/matches", trendingHandler.GetTrendingMatches)
	mux.HandleFunc("GET /trending/reviews", trendingHandler.GetTrendingReviews)

	// ----------------------------------------------------
	// 9. Public Metadata Endpoints
	// ----------------------------------------------------
	// Countries
	mux.HandleFunc("GET /countries", metadataHandler.ListCountries)
	mux.HandleFunc("GET /countries/{code}", metadataHandler.GetCountry)

	// ----------------------------------------------------
	// 10. Setup Middlewares for Protected Endpoints
	// ----------------------------------------------------
	authMiddleware := middleware.RequireAuth(jwtSecret)
	adminMiddleware := func(next http.Handler) http.Handler {
		return authMiddleware(middleware.RequireAdmin(queries)(next))
	}

	// ----------------------------------------------------
	// 11. Protected Auth & Me Endpoints (middleware.RequireAuth)
	// ----------------------------------------------------
	mux.Handle("POST /auth/logout", authMiddleware(http.HandlerFunc(authHandler.Logout)))
	mux.Handle("GET /me", authMiddleware(http.HandlerFunc(meHandler.GetProfile)))
	mux.Handle("PATCH /me", authMiddleware(http.HandlerFunc(meHandler.UpdateProfile)))
	mux.Handle("PATCH /me/username", authMiddleware(http.HandlerFunc(meHandler.UpdateUsername)))
	mux.Handle("PATCH /me/avatar", authMiddleware(http.HandlerFunc(meHandler.UpdateAvatar)))
	mux.Handle("DELETE /me", authMiddleware(http.HandlerFunc(meHandler.DeleteAccount)))

	// ----------------------------------------------------
	// 12. Match & Performance Ratings (Protected)
	// ----------------------------------------------------
	mux.Handle("POST /matches/{slug}/rating", authMiddleware(http.HandlerFunc(ratingsHandler.RateMatch)))
	mux.Handle("PATCH /matches/{slug}/rating", authMiddleware(http.HandlerFunc(ratingsHandler.UpdateMatchRating)))
	mux.Handle("DELETE /matches/{slug}/rating", authMiddleware(http.HandlerFunc(ratingsHandler.DeleteMatchRating)))
	mux.Handle("GET /matches/{slug}/ratings/me", authMiddleware(http.HandlerFunc(ratingsHandler.GetMyMatchRating)))

	mux.Handle("POST /performances/{id}/rating", authMiddleware(http.HandlerFunc(ratingsHandler.RatePerformance)))
	mux.Handle("PATCH /performances/{id}/rating", authMiddleware(http.HandlerFunc(ratingsHandler.UpdatePerformanceRating)))
	mux.Handle("DELETE /performances/{id}/rating", authMiddleware(http.HandlerFunc(ratingsHandler.DeletePerformanceRating)))
	mux.Handle("GET /performances/{id}/ratings/me", authMiddleware(http.HandlerFunc(ratingsHandler.GetMyPerformanceRating)))

	// ----------------------------------------------------
	// 13. Match & Performance Reviews (Protected)
	// ----------------------------------------------------
	mux.Handle("POST /matches/{slug}/reviews", authMiddleware(http.HandlerFunc(reviewsHandler.CreateMatchReview)))
	mux.Handle("PATCH /matches/{slug}/reviews/{reviewId}", authMiddleware(http.HandlerFunc(reviewsHandler.UpdateMatchReview)))
	mux.Handle("DELETE /matches/{slug}/reviews/{reviewId}", authMiddleware(http.HandlerFunc(reviewsHandler.DeleteMatchReview)))

	mux.Handle("POST /performances/{id}/reviews", authMiddleware(http.HandlerFunc(reviewsHandler.CreatePerformanceReview)))
	mux.Handle("PATCH /performances/{id}/reviews/{reviewId}", authMiddleware(http.HandlerFunc(reviewsHandler.UpdatePerformanceReview)))
	mux.Handle("DELETE /performances/{id}/reviews/{reviewId}", authMiddleware(http.HandlerFunc(reviewsHandler.DeletePerformanceReview)))

	// ----------------------------------------------------
	// 14. Match & Performance Comments (Protected)
	// ----------------------------------------------------
	mux.Handle("POST /match-reviews/{reviewId}/comments", authMiddleware(http.HandlerFunc(commentsHandler.CreateMatchComment)))
	mux.Handle("PATCH /match-review-comments/{commentId}", authMiddleware(http.HandlerFunc(commentsHandler.UpdateMatchComment)))
	mux.Handle("DELETE /match-review-comments/{commentId}", authMiddleware(http.HandlerFunc(commentsHandler.DeleteMatchComment)))

	mux.Handle("POST /performance-reviews/{reviewId}/comments", authMiddleware(http.HandlerFunc(commentsHandler.CreatePerformanceComment)))
	mux.Handle("PATCH /performance-review-comments/{commentId}", authMiddleware(http.HandlerFunc(commentsHandler.UpdatePerformanceComment)))
	mux.Handle("DELETE /performance-review-comments/{commentId}", authMiddleware(http.HandlerFunc(commentsHandler.DeletePerformanceComment)))

	// ----------------------------------------------------
	// 15. Match & Performance Likes (Protected)
	// ----------------------------------------------------
	mux.Handle("POST /match-reviews/{reviewId}/like", authMiddleware(http.HandlerFunc(likesHandler.LikeMatchReview)))
	mux.Handle("DELETE /match-reviews/{reviewId}/like", authMiddleware(http.HandlerFunc(likesHandler.UnlikeMatchReview)))

	mux.Handle("POST /performance-reviews/{reviewId}/like", authMiddleware(http.HandlerFunc(likesHandler.LikePerformanceReview)))
	mux.Handle("DELETE /performance-reviews/{reviewId}/like", authMiddleware(http.HandlerFunc(likesHandler.UnlikePerformanceReview)))

	mux.Handle("POST /match-review-comments/{commentId}/like", authMiddleware(http.HandlerFunc(likesHandler.LikeMatchComment)))
	mux.Handle("DELETE /match-review-comments/{commentId}/like", authMiddleware(http.HandlerFunc(likesHandler.UnlikeMatchComment)))

	mux.Handle("POST /performance-review-comments/{commentId}/like", authMiddleware(http.HandlerFunc(likesHandler.LikePerformanceComment)))
	mux.Handle("DELETE /performance-review-comments/{commentId}/like", authMiddleware(http.HandlerFunc(likesHandler.UnlikePerformanceComment)))

	// ----------------------------------------------------
	// 16. Admin-Only CRUD Endpoints (middleware.RequireAdmin)
	// ----------------------------------------------------
	// Players
	mux.Handle("POST /"+adminPrefix+"/players", adminMiddleware(http.HandlerFunc(playersHandler.CreatePlayer)))
	mux.Handle("PATCH /"+adminPrefix+"/players/{slug}", adminMiddleware(http.HandlerFunc(playersHandler.UpdatePlayer)))
	mux.Handle("DELETE /"+adminPrefix+"/players/{slug}", adminMiddleware(http.HandlerFunc(playersHandler.DeletePlayer)))
	// Teams
	mux.Handle("POST /"+adminPrefix+"/teams", adminMiddleware(http.HandlerFunc(teamsHandler.CreateTeam)))
	mux.Handle("PATCH /"+adminPrefix+"/teams/{slug}", adminMiddleware(http.HandlerFunc(teamsHandler.UpdateTeam)))
	mux.Handle("DELETE /"+adminPrefix+"/teams/{slug}", adminMiddleware(http.HandlerFunc(teamsHandler.DeleteTeam)))
	// Countries
	mux.Handle("POST /"+adminPrefix+"/countries", adminMiddleware(http.HandlerFunc(metadataHandler.CreateCountry)))
	mux.Handle("PATCH /"+adminPrefix+"/countries/{code}", adminMiddleware(http.HandlerFunc(metadataHandler.UpdateCountry)))
	mux.Handle("DELETE /"+adminPrefix+"/countries/{code}", adminMiddleware(http.HandlerFunc(metadataHandler.DeleteCountry)))
	// Import
	mux.Handle("POST /"+adminPrefix+"/import-match", http.HandlerFunc(importHandler.ImportMatch))
	// Matches
	mux.Handle("POST /"+adminPrefix+"/matches", adminMiddleware(http.HandlerFunc(matchesHandler.CreateMatch)))
	mux.Handle("PATCH /"+adminPrefix+"/matches/{slug}", http.HandlerFunc(matchesHandler.UpdateMatch))
	mux.Handle("DELETE /"+adminPrefix+"/matches/{slug}", adminMiddleware(http.HandlerFunc(matchesHandler.DeleteMatch)))
	// Performances
	mux.Handle("POST /"+adminPrefix+"/performances", adminMiddleware(http.HandlerFunc(performancesHandler.CreatePerformance)))
	mux.Handle("PATCH /"+adminPrefix+"/performances/{id}", http.HandlerFunc(performancesHandler.UpdatePerformance))
	mux.Handle("DELETE /"+adminPrefix+"/performances/{id}", adminMiddleware(http.HandlerFunc(performancesHandler.DeletePerformance)))

	// Overview stats (public) — powers the homepage stat strip with real totals
	mux.HandleFunc("GET /stats/overview", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		stats, err := queries.GetOverviewStats(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to fetch overview stats"}`))
			return
		}
		_ = json.NewEncoder(w).Encode(stats)
	})

	// Health check endpoint
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"UP","message":"PerformX API Server is healthy"}`))
	})
}
