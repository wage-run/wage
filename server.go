package wage

import (
	"net/http"
)

type Server struct {
	Root string
}

var _ http.Handler = (*Server)(nil)

func NewServer(root string) *Server {
	s := &Server{
		Root: root,
	}
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
