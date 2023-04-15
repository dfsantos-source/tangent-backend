package http

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	"github.com/dfsantos-source/tangent-backend"
)

// server struct to hold instances needed
// wraps all HTTP functionality
type Server struct {
	server *http.Server
	router *chi.Mux

	MapboxUtil tangent.Util
	YelpUtil   tangent.Util
}

// creates an instance of a server
// using Go http library and go-chi library router
func CreateServer() *Server {
	s := &Server{
		server: &http.Server{},
		router: chi.NewRouter(),
	}

	godotenv.Load("../local.env")
	s.MapboxUtil = *tangent.CreateUtil(os.Getenv("TOKEN_MAPBOX"))
	s.YelpUtil = *tangent.CreateUtil(os.Getenv("TOKEN_YELP"))

	s.registerTangentRoutes(s.router)

	return s
}

func (s *Server) RunServer() error {
	err := http.ListenAndServe(":3000", s.router)
	if err != nil {
		return err
	}
	return nil
}
