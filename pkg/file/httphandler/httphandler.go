package httphandler

import (
	"net/http"

	repo "upload/pkg/file/repo/mongo"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"
)

type fileServer struct {
	mux       *chi.Mux
	prefix    string
	db        *repo.FileCollection
	validator *validator.Validate
}

func (s *fileServer) Prefix() string {
	return s.prefix
}

func (s *fileServer) Mux() http.Handler {
	return s.mux
}

func (s *fileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func NewServer(db *mongo.Collection) *fileServer {
	s := &fileServer{
		prefix:    "/files",
		mux:       chi.NewRouter(),
		db:        repo.NewRepo(db),
		validator: validator.New(validator.WithRequiredStructEnabled()),
	}
	s.routes()
	return s
}
