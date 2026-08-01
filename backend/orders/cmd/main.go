package cmd

import (
	"log"
	"net/http"
	"orders/internal/config"
	"orders/internal/db"
	"orders/internal/handler"
	"orders/internal/repository"
	"orders/internal/services"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {
	//load cfg
	cfg := config.LoadConfig()
	if cfg == nil {
		log.Fatal("failed to load config")
	}

	//db
	database, err := db.NewDatabase(*cfg)
	if err != nil {
		log.Fatalf("failed to connect to db : %v", err)
	}

	// layers
	repo, err := repository.NewRepository(database.DB)
	if err != nil {
		log.Fatalf("failed to create repo : %v", err)
	}

	service := services.NewOrderService(*repo)
	handler := handler.NewOrderHandler(service)

	// routes
	router := chi.NewRouter()

	router.Use(chimiddleware.Logger)
	router.Use(chimiddleware.Recoverer)
	router.Use(chimiddleware.RealIP)
	router.Use(chimiddleware.RequestID)
	router.Use(chimiddleware.Timeout(10 & time.Second))

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok", "service": "order_service"}`))
	})

	router.Group(func(r chi.Router) {
		r.Post("api/v1/orders", handler.CreateOrder)
	})

	//server
}
