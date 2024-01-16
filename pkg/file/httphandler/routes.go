package httphandler

import "github.com/go-chi/chi/v5/middleware"

func (s *fileServer) routes() {
	s.mux.Use(middleware.Logger)

	s.mux.Post("/", s.handleFileCreate())
	// s.mux.Get("/{fileID}", s.handleFileFindByID())
	// s.mux.Patch("/{fileID}", s.handleFileUpdate())
	// s.mux.Delete("/{fileID}", s.handleFileDelete())
}
