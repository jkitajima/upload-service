package httphandler

import (
	"net/http"

	repo "upload/pkg/file/repo/mongo"
	"upload/shared/blob"
	"upload/shared/zombiekiller"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"
)

type fileServer struct {
	mux       *chi.Mux
	prefix    string
	db        *repo.FileCollection
	blobstg   blob.Storager
	thrash    chan<- zombiekiller.KillOperation
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

func NewServer(db *mongo.Collection, blobstg blob.Storager, thrashChan chan<- zombiekiller.KillOperation) *fileServer {
	s := &fileServer{
		prefix:    "/files",
		mux:       chi.NewRouter(),
		db:        repo.NewRepo(db),
		blobstg:   blobstg,
		thrash:    thrashChan,
		validator: validator.New(validator.WithRequiredStructEnabled()),
	}
	s.routes()
	return s
}
