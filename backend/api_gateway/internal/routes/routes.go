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
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		ExposedHeaders:   []string{"Authorization"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// middlewares
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWT.AccessSecret, cfg.JWT.RefreshSecret)

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

	return router
}
