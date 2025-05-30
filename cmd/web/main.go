package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/Galbeyte1/snippet-box-taketwo/internal/config"
	"github.com/Galbeyte1/snippet-box-taketwo/internal/database"
	"github.com/Galbeyte1/snippet-box-taketwo/internal/models"
	"github.com/Galbeyte1/snippet-box-taketwo/internal/templates"
	"github.com/Galbeyte1/snippet-box-taketwo/internal/transport"
	"github.com/boj/redistore"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
)

func main() {

	var cfg config.Config

	cfg = config.LoadConfigFromEnv()

	flag.StringVar(&cfg.APIAddr, "addr", cfg.APIAddr, "HTTP network address")
	flag.StringVar(&cfg.StaticDir, "static-dir", cfg.StaticDir, "Path to static assets")
	flag.Parse()

	if cfg.Database.DSN() == "" {
		fmt.Fprintln(os.Stderr, "error: DSN must not be empty. Check your environment variables or flags.")
		os.Exit(1)
	}

	if cfg.SessionKey == "" {
		fmt.Fprintln(os.Stderr, "error: SESSION_SECRET must not be empty. Check your enviornment variables.")
		os.Exit(1)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		// AddSource: true,
	}))

	app := &transport.Application{
		Logger:      logger,
		FormDecoder: form.NewDecoder(),
	}

	db, err := database.OpenDB(cfg.Database.DSN())
	if err != nil {
		app.Logger.Error("database failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer db.Close()

	runner := database.NewMigrationsRunner(db)
	if err := runner.Run(); err != nil {
		app.Logger.Error("migration failed", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// if err := database.VerifyDBMigrations(db); err != nil {
	// 	app.Logger.Error("database migration verification failed", slog.String("error", err.Error()))
	// 	os.Exit(1)
	// }

	app.Snippets = &models.SnippetModel{DB: db}

	pool := database.OpenRedis(*cfg.Redis)

	store, err := redistore.NewRediStoreWithPool(pool, []byte(cfg.SessionKey))
	if err != nil {
		app.Logger.Error("cannot create Redis session store", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer store.Close()

	store.Options = &sessions.Options{
		Path:     cfg.SessionOpts.Path,
		MaxAge:   cfg.SessionOpts.MaxAge,
		HttpOnly: cfg.SessionOpts.HttpOnly,
		Secure:   cfg.SessionOpts.Secure,
		SameSite: cfg.SessionOpts.SameSite,
	}

	app.SessionStore = store

	templateCache, err := templates.NewTemplateCache()
	if err != nil {
		app.Logger.Error(err.Error())
		os.Exit(1)
	}
	app.TemplateCache = templateCache

	srv := &http.Server{
		Addr:         cfg.APIAddr,
		Handler:      app.Routes(cfg),
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		TLSConfig:    cfg.TLSOpts.Config,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	app.Logger.Info("starting server", slog.String("addr", srv.Addr))

	go func() {
		err := http.ListenAndServe(":80", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			target := "https://" + r.Host + r.URL.RequestURI()
			http.Redirect(w, r, target, http.StatusMovedPermanently)
		}))
		if err != nil {
			app.Logger.Error("HTTP redirect server failed: " + err.Error())
		}
	}()
	err = srv.ListenAndServeTLS("/tls/cert.pem", "/tls/key.pem")
	app.Logger.Error(err.Error())
	os.Exit(1)

}
