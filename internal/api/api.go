package api

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"testovoe/internal/config"
	"testovoe/internal/handlers"
	"testovoe/internal/storage"
	"testovoe/internal/utils"
	"testovoe/internal/validators"
    "time"

	"github.com/go-playground/validator/v10"
	_ "github.com/lib/pq"
)

type Api struct {
	cfg *config.Config
}

func New(cfg *config.Config) *Api {
	return &Api{
		cfg: cfg,
	}
}

func (a *Api) Run() error {
	mux := http.NewServeMux()

	db, err := sql.Open("postgres",
		fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			a.cfg.Database.Host,
			a.cfg.Database.Port,
			a.cfg.Database.User,
			a.cfg.Database.Password,
			a.cfg.Database.Database,
		),
	)
	if err != nil {
		return err
	}
	store := storage.NewPostgresStorage(db)

	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterValidation("mm_yyyy", validators.MonthYearValidator)

	subHandler := handlers.NewSubscriptionHandler(store, validate)
	mux.HandleFunc("POST /subscriptions", subHandler.Create)
	mux.HandleFunc("GET /subscriptions/{id}", subHandler.Get)
	mux.HandleFunc("PUT /subscriptions/{id}", subHandler.Update)
	mux.HandleFunc("DELETE /subscriptions/{id}", subHandler.Delete)
	mux.HandleFunc("GET /subscriptions", subHandler.List)

	// Wrap the mux with gzip compression to reduce payload sizes
	handler := utils.GzipMiddleware(mux)

	srv := &http.Server{
		Addr:              net.JoinHostPort(a.cfg.Api.Host, fmt.Sprint(a.cfg.Api.Port)),
		Handler:           handler,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	slog.Info("start server", "host", a.cfg.Api.Host, "port", a.cfg.Api.Port)
	return srv.ListenAndServe()
}
