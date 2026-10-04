package main

import (
	"context"
	"gmailer/internal/database"
	"gmailer/internal/handlers"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Print("Ошибка подрузги .env")
	}

	ctx := context.Background()
	appPort := os.Getenv("APP_PORT")
	databaseUrl := os.Getenv("DB_URL")

	db, err := database.Connect(ctx, databaseUrl)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе: %v", err)
	}

	defer db.Close()

	authStore := database.NewAuthStore(db)
	authHandler := handlers.NewAuthHandler(authStore)

	contactStore := database.NewContactStore(db)
	contactHandler := handlers.NewContactHandler(contactStore) 

	r := chi.NewRouter()
	r.Use(middleware.Logger, middleware.Recoverer)

	r.Mount("/api/auth", handlers.AuthRouter(authHandler))
	r.Mount("/api/contact", handlers.ContactRouter(contactHandler))

	log.Printf("Сервер запущен на порту %s", appPort)
	if err := http.ListenAndServe(":"+appPort, r); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
