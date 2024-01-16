package composer

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Server interface {
	Mux() http.Handler
	Prefix() string
}

type Composer struct {
	servers []Server
	mux     *chi.Mux
}

func NewComposer() *Composer {
	return &Composer{mux: chi.NewRouter()}
}

func (c *Composer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.mux.ServeHTTP(w, r)
}

func (c *Composer) Compose(servers ...Server) error {
	if len(c.servers) > 0 {
		return errors.New("composer is already filled with servers")
	}

	// TODO!
	// error handling needed if s does not have prefix/mux
	// these methods should always return valid results
	for _, s := range servers {
		c.mux.Mount(s.Prefix(), s.Mux())
	}

	return nil
}
