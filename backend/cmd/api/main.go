package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ecommerce-patrones/internal/catalog"
	"ecommerce-patrones/internal/httpapi"
	"ecommerce-patrones/internal/pricing"
	"ecommerce-patrones/internal/singleton"
)

func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// 1. Patron SINGLETON: Carga de configuracion unica
	cfg := singleton.GetConfig()
	log.Info("Configuracion cargada (Singleton)", "port", cfg.HTTPPort, "db", cfg.DatabaseURL)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 2. Patron SINGLETON: Pool de conexiones a base de datos unico
	pool, err := singleton.GetDB(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Error("Fallo conexion a Postgres", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	// Ejecucion de migraciones
	if err := singleton.Migrate(ctx, pool, log); err != nil {
		log.Error("Fallo al aplicar migraciones", "error", err)
		os.Exit(1)
	}

	// 3. Inicializacion de servicios de los 2 CRUDs
	catalogSvc := catalog.NewService(pool)
	pricingSvc := pricing.NewService(pool, cfg.Currency)

	api := httpapi.New(catalogSvc, pricingSvc)

	srv := &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      api.Handler(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info("Servidor Backend de Patrones escuchando", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("Fallo del servidor HTTP", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	log.Info("Apagando servidor HTTP de forma ordenada...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("Error durante el apagado", "error", err)
	}
}
