package singleton

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"ecommerce-patrones/migrations"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	dbPoolInstance *pgxpool.Pool
	dbInitErr      error
	dbOnce         sync.Once
)

// GetDB implementa el PATRON SINGLETON para el Pool de conexiones a PostgreSQL.
// Garantiza que en toda la aplicacion solo exista UNA unica instancia del pool,
// evitando la sobrecarga de crear conexiones repetidas por cada peticion HTTP.
func GetDB(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	dbOnce.Do(func() {
		cfg, err := pgxpool.ParseConfig(databaseURL)
		if err != nil {
			dbInitErr = fmt.Errorf("error al parsear conexion db: %w", err)
			return
		}

		cfg.MaxConns = 25
		cfg.MinConns = 2
		cfg.MaxConnLifetime = 1 * time.Hour

		pool, err := pgxpool.NewWithConfig(ctx, cfg)
		if err != nil {
			dbInitErr = fmt.Errorf("error al conectar con base de datos: %w", err)
			return
		}

		// Ping de comprobacion
		pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
		if err := pool.Ping(pingCtx); err != nil {
			dbInitErr = fmt.Errorf("fallo ping a postgres: %w", err)
			return
		}

		dbPoolInstance = pool
	})

	return dbPoolInstance, dbInitErr
}

// Migrate ejecuta el script SQL embebido de migraciones.
func Migrate(ctx context.Context, pool *pgxpool.Pool, log *slog.Logger) error {
	data, err := migrations.FS.ReadFile("001_init.sql")
	if err != nil {
		return fmt.Errorf("no se pudo leer 001_init.sql: %w", err)
	}

	if _, err := pool.Exec(ctx, string(data)); err != nil {
		return fmt.Errorf("error ejecutando migracion 001_init.sql: %w", err)
	}

	log.Info("Migracion de base de datos ejecutada con exito")
	return nil
}
