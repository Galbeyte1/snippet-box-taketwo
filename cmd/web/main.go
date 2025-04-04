package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"

	"github.com/Galbeyte1/snippet-box-taketwo/internal/config"
	"github.com/Galbeyte1/snippet-box-taketwo/internal/transport"
)

func main() {

	var cfg config.Config

	flag.StringVar(&cfg.Addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.StaticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.Parse()

	app := &config.Application{
		Logger: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		})),
	}

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(cfg.StaticDir))

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", transport.Home(app))
	mux.HandleFunc("GET /snippet/view", transport.SnippetView(app))
	mux.HandleFunc("GET /snippet/create", transport.SnippetCreate(app))
	mux.HandleFunc("POST /snippet/create", transport.SnippetCreatePost(app))

	app.Logger.Info("starting server", slog.String("addr", cfg.Addr))

	err := http.ListenAndServe(cfg.Addr, mux)
	app.Logger.Error(err.Error())
	os.Exit(1)

}
