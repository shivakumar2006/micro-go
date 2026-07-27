package routes

import (
	"api_gateway/internal/config"
	"api_gateway/internal/middleware"
	"api_gateway/internal/proxy"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func Setup(cfg *config.Config, serviceProxy *proxy.ServiceProxy) http.Handler {
	router := chi.NewRouter()

	router.Use(chimiddleware.Logger)
	router.Use(chimiddleware.Recoverer)
	router.Use(chimiddleware.RealIP)
	router.Use(chimiddleware.RequestID)

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// middlewares
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWT.AccessSecret, cfg.JWT.RefreshSecret)
	rateLimiter := middleware.NewRateLimiter(100, 10)

	router.Use(rateLimiter.Middleware)

	// health check
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok", "service": "api_gateway"}`))
	})

	router.Route("/api/v1/auth", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Post("/register", serviceProxy.Auth)
			r.Post("/login", serviceProxy.Auth)
			r.Post("/refresh", serviceProxy.Auth)
		})

		// protected routes
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.Authenticate)

			r.Get("/me", serviceProxy.Auth)
			r.Post("/logout", serviceProxy.Auth)
			r.Post("/logout-all", serviceProxy.Auth)
		})
	})

	router.Route("/api/v1", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Get("/vehicles", serviceProxy.Vehicle)
			r.Get("/vehicles/{id}", serviceProxy.Vehicle)
		})

		// protected routes
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.Authenticate)

			r.Post("/vehicles", serviceProxy.Vehicle)
			r.Put("/vehicles/{id}", serviceProxy.Vehicle)
			r.Delete("/vehicles/{id}", serviceProxy.Vehicle)
		})
	})

	router.Route("/api/v1", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.Authenticate)

			r.Post("/cart", serviceProxy.Cart)
			r.Put("/cart/{id}", serviceProxy.Cart)
			r.Get("/cart", serviceProxy.Cart)
			r.Delete("/cart/{id}", serviceProxy.Cart)
			r.Delete("/cart", serviceProxy.Cart)
			r.Get("/cart/total", serviceProxy.Cart)
			r.Get("/cart/count", serviceProxy.Cart)
		})
	})

	return router
}
