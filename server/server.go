package server

import (
	"github.com/SpeedSlime/Covalence/server/middleware"
	"github.com/SpeedSlime/Covalence/server/router"
	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
)

type Server struct {
	version      string
	port         string
	name         string
	address      string
	connection   string
	status       string
	experimental bool
	middlewares  []middleware.Middleware
	routers      []router.Router
}

func (s Server) Version() string {
	return s.version
}

func (s Server) Port() string {
	return s.port
}

func (s Server) Name() string {
	return s.name
}

func (s Server) Connection() string {
	return s.connection
}

func (s Server) Address() string {
	return s.address
}

func (s Server) Experimental() bool {
	return s.experimental
}

func (s Server) Status() string {
	return s.status
}

func Create(version string, port string, name string, address string, connection string, experimental bool) *Server {
	return &Server{
		version:      version,
		port:         port,
		experimental: experimental,
		name:         name,
		address:      address,
		connection:   connection,
	}
}

func (s *Server) LoadRouters(routers ...router.Router) {
	s.routers = append(s.routers, routers...)
}

func (s *Server) LoadMiddleware(middlewares ...middleware.Middlware) {
	s.routers = append(s.middlewares, middlewares...)
}

func (s *Server) Handler() *chi.Mux {
	m := chi.NewMux()
	for _, w := range s.middlewares {
		if w.Status() && (w.Experimental() == s.Experimental() || !w.Experimental()) {
			log.Info().Bool("Experimental", w.Experimental()).Bool("Status", w.Status()).Str("Method", w.Method()).Msg("Registering middlewear")
			m.Use(w.Method())
		}
	}
	s.initMux(m)
	return m
}

func (s *Server) initMux(m *chi.Mux) {
	for _, router := range s.routers {
		if router.Status() {
			for _, r := range router.Routes() {
				if r.Status() && (r.Experimental() == s.Experimental() || !r.Experimental()) {
					log.Info().Bool("Experimental", r.Experimental()).Bool("Status", r.Status()).Str("Method", r.Method()).Str("Path", r.Path()).Msg("Registering route")
					m.Method(r.Method(), r.Path(), r.Handler())
				}
			}
		}
	}
}
