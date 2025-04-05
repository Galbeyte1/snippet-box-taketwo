package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"

	"github.com/Galbeyte1/snippet-box-taketwo/internal/config"
	"github.com/Galbeyte1/snippet-box-taketwo/internal/models"
	"github.com/Galbeyte1/snippet-box-taketwo/internal/transport"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	var cfg config.Config

	flag.StringVar(&cfg.Addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.StaticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.StringVar(&cfg.DSN, "dsn", "web:YES@/snippetbox?parseTime=true", "web:YES@/snippetbox?parseTime=true")
	flag.Parse()

	app := &config.Application{
		Logger: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		})),
	}

	db, err := openDB(cfg.DSN)
	if err != nil {
		app.Logger.Error(err.Error())
		os.Exit(1)
	}
	app.Snippets = &models.SnippetModel{DB: db}

	defer db.Close()

	app.Logger.Info("starting server", slog.String("addr", cfg.Addr))

	err = http.ListenAndServe(cfg.Addr, transport.Routes(app, cfg))
	app.Logger.Error(err.Error())
	os.Exit(1)

}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
