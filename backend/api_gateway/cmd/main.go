package main

import (
	"api_gateway/internal/config"
	"api_gateway/internal/proxy"
	"api_gateway/internal/routes"
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.LoadConfig()

	// service proxy
	sp, err := proxy.NewServiceProxy(cfg.Services.Auth.URL, cfg.Services.Cart.URL, cfg.Services.Vehicle.URL, cfg.Services.Orders.URL)
	if err != nil {
		log.Printf("Failed to setup proxies : %v", err)
	}

	r := routes.Setup(cfg, sp)

	server := &http.Server{
		Addr:         cfg.Server.Addr,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
		IdleTimeout:  100 * time.Second,
	}

	quit := make(chan os.Signal, 1)

	slog.Info("Server started successfully", slog.String("address : ", cfg.Server.Addr))

	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		slog.Info("HTTP server listening.....")

		log.Printf("🌐 api gateway running on port : %s", cfg.Server.Addr)
		log.Printf("🔑 auth service : %s", cfg.Services.Auth.URL)
		log.Printf("🚗 vehicle service : %s", cfg.Services.Vehicle.URL)
		log.Printf("🛒 cart service : %s", cfg.Services.Cart.URL)
		log.Printf("📦 order service : %s", cfg.Services.Orders.URL)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("gateway error : %v", err)
		}
	}()

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	slog.Info("server shutting down gracefully...")

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("failed to shutdonw gateway : %v", err)
	}

	slog.Info("api gateway exited successfully")
}
