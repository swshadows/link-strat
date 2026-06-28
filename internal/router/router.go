package router

import (
	links "link-strat/internal/service"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.ClientIPFromRemoteAddr)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	linkHandler := links.NewLinkHandler()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, map[string]bool{"ok": true})
	})

	r.Post("/check-links", linkHandler.CheckLinks)

	main := chi.NewRouter()
	main.Mount("/api", r)
	main.Handle("/*", http.FileServer(http.Dir("./static")))

	return main
}
