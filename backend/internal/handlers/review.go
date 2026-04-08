package handlers

import (
	"codefix-ai/internal/db"
	"codefix-ai/internal/middleware"
	"codefix-ai/internal/models"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"//router
	"github.com/google/uuid"
)

func CreateReview(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)

	var req models.ReviewRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Code == "" || req.Language == "" {
		http.Error(w, "code and language are required", http.StatusBadRequest)
		return
	}

	aiResult, err := getAIReview(req.Code, req.Language)
	if err != nil {
		http.Error(w, "failed to get AI review: "+err.Error(), http.StatusInternalServerError)
		return
	}

	id := uuid.New().String()
	var review models.Review

	err = db.DB.QueryRow(`
		INSERT INTO reviews (id, user_id, language, original_code, bug_explanation, fixed_code, suggestions)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, user_id, language, original_code, bug_explanation, fixed_code, suggestions, created_at
	`, id, userID, req.Language, req.Code, aiResult.BugExplanation, aiResult.FixedCode, aiResult.Suggestions,
	).Scan(
		&review.ID, &review.UserID, &review.Language,
		&review.OriginalCode, &review.BugExplanation,
		&review.FixedCode, &review.Suggestions, &review.CreatedAt,
	)

	if err != nil {
		http.Error(w, "failed to save review", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(review)
}

func GetReviews(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	language := r.URL.Query().Get("language")

	var rows *sql.Rows // will fix below
	var err error

	if language != "" {
		rows, err = db.DB.Query(`
			SELECT id, user_id, language, original_code, bug_explanation, fixed_code, suggestions, is_helpful, created_at
			FROM reviews WHERE user_id = $1 AND language = $2
			ORDER BY created_at DESC
		`, userID, language)
	} else {
		rows, err = db.DB.Query(`
			SELECT id, user_id, language, original_code, bug_explanation, fixed_code, suggestions, is_helpful, created_at
			FROM reviews WHERE user_id = $1
			ORDER BY created_at DESC
		`, userID)
	}

	if err != nil {
		http.Error(w, "failed to fetch reviews", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	reviews := []models.Review{}
	for rows.Next() {
		var rev models.Review
		err := rows.Scan(
			&rev.ID, &rev.UserID, &rev.Language,
			&rev.OriginalCode, &rev.BugExplanation,
			&rev.FixedCode, &rev.Suggestions,
			&rev.IsHelpful, &rev.CreatedAt,
		)
		if err != nil {
			http.Error(w, "error reading reviews", http.StatusInternalServerError)
			return
		}
		reviews = append(reviews, rev)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reviews)
}

func GetReview(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	reviewID := chi.URLParam(r, "id")

	var review models.Review
	err := db.DB.QueryRow(`
		SELECT id, user_id, language, original_code, bug_explanation, fixed_code, suggestions, is_helpful, created_at
		FROM reviews WHERE id = $1 AND user_id = $2
	`, reviewID, userID).Scan(
		&review.ID, &review.UserID, &review.Language,
		&review.OriginalCode, &review.BugExplanation,
		&review.FixedCode, &review.Suggestions,
		&review.IsHelpful, &review.CreatedAt,
	)

	if err != nil {
		http.Error(w, "review not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(review)
}

func DeleteReview(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	reviewID := chi.URLParam(r, "id")

	result, err := db.DB.Exec(
		"DELETE FROM reviews WHERE id = $1 AND user_id = $2",
		reviewID, userID,
	)
	if err != nil {
		http.Error(w, "failed to delete review", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "review not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "deleted"})
}

func RateReview(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	reviewID := chi.URLParam(r, "id")

	var body struct {
		IsHelpful bool `json:"is_helpful"`
	}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	_, err = db.DB.Exec(
		"UPDATE reviews SET is_helpful = $1 WHERE id = $2 AND user_id = $3",
		body.IsHelpful, reviewID, userID,
	)
	if err != nil {
		http.Error(w, "failed to rate review", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "rated"})
}