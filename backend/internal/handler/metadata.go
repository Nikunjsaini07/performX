package handler

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5"

	"github.com/Nikunjsaini07/performx/backend/internal/db"
)

type MetadataHandler struct {
	Queries *db.Queries
}

// ==========================================
// COUNTRIES HANDLERS
// ==========================================

type CreateCountryRequest struct {
	Name string `json:"name"`
	Iso2 string `json:"iso2"`
	Iso3 string `json:"iso3"`
}

type UpdateCountryRequest struct {
	Name *string `json:"name"`
	Iso2 *string `json:"iso2"`
	Iso3 *string `json:"iso3"`
}

func (h *MetadataHandler) ListCountries(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	countries, err := h.Queries.ListCountries(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to list countries"}`))
		return
	}
	if countries == nil {
		countries = []db.ListCountriesRow{}
	}
	_ = json.NewEncoder(w).Encode(countries)
}

func (h *MetadataHandler) GetCountry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	code := r.PathValue("code")
	country, err := h.Queries.GetCountryByCode(r.Context(), code)
	if err != nil {
		if err == pgx.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error":"Not Found","message":"Country not found"}`))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to fetch country"}`))
		return
	}
	_ = json.NewEncoder(w).Encode(country)
}

func (h *MetadataHandler) CreateCountry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req CreateCountryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" || req.Iso2 == "" || req.Iso3 == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"name, iso2, and iso3 are required"}`))
		return
	}
	country, err := h.Queries.CreateCountry(r.Context(), db.CreateCountryParams{
		Name: req.Name,
		Iso2: req.Iso2,
		Iso3: req.Iso3,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to create country"}`))
		return
	}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(country)
}

func (h *MetadataHandler) UpdateCountry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	code := r.PathValue("code")
	country, err := h.Queries.GetCountryByCode(r.Context(), code)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"Not Found","message":"Country to update not found"}`))
		return
	}
	var req UpdateCountryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"Bad Request","message":"Invalid request body"}`))
		return
	}
	nameVal := country.Name
	if req.Name != nil {
		nameVal = *req.Name
	}
	iso2Val := country.Iso2
	if req.Iso2 != nil {
		iso2Val = *req.Iso2
	}
	iso3Val := country.Iso3
	if req.Iso3 != nil {
		iso3Val = *req.Iso3
	}
	updated, err := h.Queries.UpdateCountry(r.Context(), db.UpdateCountryParams{
		ID:   country.ID,
		Name: nameVal,
		Iso2: iso2Val,
		Iso3: iso3Val,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to update country"}`))
		return
	}
	_ = json.NewEncoder(w).Encode(updated)
}

func (h *MetadataHandler) DeleteCountry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	code := r.PathValue("code")
	country, err := h.Queries.GetCountryByCode(r.Context(), code)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"Not Found","message":"Country to delete not found"}`))
		return
	}
	err = h.Queries.DeleteCountry(r.Context(), country.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal Error","message":"Failed to delete country"}`))
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"message":"Country successfully deleted"}`))
}
