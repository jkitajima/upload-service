package httphandler

import (
	"net/http"

	repo "upload/pkg/file/repo/mongo"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

type fileServer struct {
	mux    *chi.Mux
	prefix string
	db     *repo.FileCollection
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
		prefix: "/files",
		mux:    chi.NewRouter(),
		db:     repo.NewRepo(db),
	}
	s.routes()
	return s
}
