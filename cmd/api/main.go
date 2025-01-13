package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn         string
		maxConns    int
		minConns    int
		maxIdleTime time.Duration
	}
}

type application struct {
	config config
	logger *slog.Logger
	// validate *validator.Validate
}

func main() {

	var cfg config
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	err := godotenv.Load()

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	flag.IntVar(&cfg.port, "port", 4000, "api server")
	flag.StringVar(&cfg.env, "env", "dev", "Environment (dev, staging, prod)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("DSN"), "PostgreSQL DSN")

	//connection pool settings
	flag.IntVar(&cfg.db.maxConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.minConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.DurationVar(&cfg.db.maxIdleTime, "db-max-idle-time", 15*time.Minute, "PostgreSQL max connection idle time")

	flag.Parse()
	// Initialize validator
	// validate := validator.New()

	db, err := openDB(cfg)

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	logger.Info("database connection pool established")

	app := application{
		config: cfg,
		logger: logger,
		// validate: validate,
	}

	srv := http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("starting server", "addr", srv.Addr, "env", cfg.env)

	err = srv.ListenAndServe()

	logger.Error(err.Error())
	os.Exit(1)

}

func openDB(cfg config) (*pgxpool.Pool, error) {
	// Create a context with a 5-second timeout deadline.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Parse the connection string into a pgxpool.Config
	poolConfig, err := pgxpool.ParseConfig(cfg.db.dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to parse connection string: %w", err)
	}

	// Optional: Customize pool settings
	poolConfig.MaxConns = int32(cfg.db.maxConns)    // Maximum number of connections in the pool
	poolConfig.MinConns = int32(cfg.db.minConns)    // Minimum number of connections in the pool
	poolConfig.MaxConnLifetime = 30 * time.Minute   // Maximum lifetime of a connection
	poolConfig.MaxConnIdleTime = cfg.db.maxIdleTime // Maximum idle time for a connection

	// Create a new connection pool
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Use Ping to verify the connection pool is working
	err = pool.Ping(ctx)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	// Return the connection pool
	return pool, nil
}
