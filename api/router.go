package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/shrinkUrl/db"
	"gopkg.in/go-playground/validator.v9"
)

var validate *validator.Validate

// API structure object
type API struct {
	*db.Client
	hostname string
}

// New returns an Api
func New(c *db.Client, hostname string) *API {
	validate = validator.New()
	api := &API{c, hostname}
	return api
}

// Start initialize router and start server
func (a *API) Start(port int) error {
	log.Printf("Starting server on port: %d", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), a.initRouter())
}

// Init returns a router for the api
func (a *API) initRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// Redirect
	r.Get("/{slug}", a.redirectShortURL)

	// API
	r.Route("/api", func(r chi.Router) {
		r.Post("/url", a.createNewShortURL)
		r.Get("/url/{slug}", a.loadShortURL)
	})

	r.NotFound(a.notFound)
	r.MethodNotAllowed(a.methodNotAllowed)

	return r
}

func (a *API) rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("URL Shortener API"))
}

func (a *API) methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not allowed"))
}

func (a *API) notFound(w http.ResponseWriter, r *http.Request) {
	render.Render(w, r, errNotFound(errors.New("Not Found")))
}
