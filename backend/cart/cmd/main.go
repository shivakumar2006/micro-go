package cmd

import (
	"cart/internal/config"
	"cart/internal/db"
	"cart/internal/handler"
	"cart/internal/repository"
	"cart/internal/service"
	"log"
	"log/slog"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {
	//load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config file : %v", err.Error())
	}

	// setup db
	database, err := db.NewCartDatabase(*cfg)
	if err != nil {
		log.Fatalf("failed to setup db : %v", err.Error())
	}
	defer database.DB.Close()

	slog.Info("Database connected successfully", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))

	// layers
	repo := repository.NewCartRepository(database.DB)
	service := service.NewCartService(repo)
	handler := handler.NewCartHandler(service)

	// add routes
	router := chi.NewRouter()

	router.Use(chimiddleware.Logger)
	router.Use(chimiddleware.Recoverer)
	router.Use(chimiddleware.RequestID)
	router.Use(chimiddleware.RealIP)
	router.Use(chimiddleware.Timeout(10 * time.Second))

	router.Post("/api/v1/cart", handler.AddToCart)
	router.Put("/api/v1/cart/{id}", handler.UpdateCartQuantity)
	router.Get("/api/v1/cart/{id}", handler.GetUserCart)
	router.Delete("/api/v1/cart/{id}", handler.DeleteCartItem)
	router.Delete("/api/v1/cart/{id}", handler.DeleteCart)
	router.Get("/api/v1/cart/total/{id}", handler.GetCartTotal)
	router.Get("/api/v1/cart/count/{id}", handler.CountItems)

	// server
}
