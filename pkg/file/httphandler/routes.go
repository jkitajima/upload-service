package httphandler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
)

func (s *fileServer) routes() {
	s.mux.Use(middleware.Logger)

	// protected routes
	s.mux.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(s.auth))
		r.Use(jwtauth.Authenticator(s.auth))

		r.Post("/", s.handleFileCreate())
		r.Get("/{fileID}", s.handleFileFindByID())
		r.Patch("/{fileID}", s.handleFileUpdate())
		r.Delete("/{fileID}", s.handleFileDelete())
	})
}
