package main

import (
	"codefix-ai/internal/db"
	"codefix-ai/internal/handlers"
	"codefix-ai/internal/middleware"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	db.Connect()

	r := chi.NewRouter()

r.Use(cors.Handler(cors.Options{
    AllowOriginFunc: func(_ *http.Request, origin string) bool {
        return strings.HasPrefix(origin, "http://localhost:") ||
            strings.HasPrefix(origin, "http://127.0.0.1:") ||
            origin == "https://code-fix-ai-woad.vercel.app"
    },
    AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
    AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
    AllowCredentials: true,
}))

	// public routes
	r.Post("/auth/register", handlers.Register)
	r.Post("/auth/login", handlers.Login)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
    _, err := db.DB.Exec("SELECT 1")
    if err != nil {
        http.Error(w, "DB unavailable", http.StatusInternalServerError)
        return
    }

    w.Write([]byte("ok"))
})

	// protected routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Post("/reviews", handlers.CreateReview)
		r.Get("/reviews", handlers.GetReviews)
		r.Get("/reviews/{id}", handlers.GetReview)
		r.Delete("/reviews/{id}", handlers.DeleteReview)
		r.Patch("/reviews/{id}/rate", handlers.RateReview)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Server running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}