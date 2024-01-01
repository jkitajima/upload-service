package httphandler

func (s *fileServer) routes() {
	s.mux.Post("/", s.handleFileCreate())
	s.mux.Get("/{fileID}", s.handleFileFindByID())
	s.mux.Patch("/{fileID}", s.handleFileUpdate())
	s.mux.Delete("/{fileID}", s.handleFileDelete())
}
