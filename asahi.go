package asahi

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/SpeedSlime/Asahi/middleware"
	"github.com/SpeedSlime/Asahi/router"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
)

type Server struct {
	port         string
	address      string
	status       bool
	experimental bool
	idle		 chan struct{}
	middlewares  []middleware.Middleware
	routers      []router.Router
}

func (s Server) Port() string {
	return s.port
}

func (s Server) Address() string {
	return s.address
}

func (s Server) Experimental() bool {
	return s.experimental
}

func (s Server) Status() bool {
	return s.status
}

func New(port string, address string, experimental bool) *Server {
	return &Server{
		port:         port,
		experimental: experimental,
		address:      address,
	}
}

func (s *Server) LoadRouter(routers []router.Router) {
	s.routers = append(s.routers, routers...)
}

func (s *Server) LoadMiddleware(middlewares []middleware.Middleware) {
	s.middlewares = append(s.middlewares, middlewares...)
}

func (s *Server) Start() {
	instance := &http.Server{Addr: s.Port(), Handler: s.handler()}
	log.Info().Str("Port", s.Port()).Str("Address", s.Address()).Bool("Experimental", s.Experimental()).Bool("Status", s.Status()).Msg("Server started")

	s.idle = make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		if err := instance.Shutdown(context.Background()); err != nil {
			log.Fatal().Msgf("Shutdown: server shutdown failed: %w", err)
		}
		close(s.idle)
	}()

	if err := instance.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal().Msgf("ListenAndServe: unexpected error: %w", err)
	}

	<-s.idle
	
	log.Info().Str("Port", s.Port()).Str("Address", s.Address()).Bool("Experimental", s.Experimental()).Bool("Status", s.Status()).Msg("Server stopped")
}

func (s *Server) handler() (*chi.Mux) {
	m := chi.NewMux()
	for _, middleware := range s.middlewares {
		if middleware.Status() && (middleware.Experimental() == s.Experimental() || !middleware.Experimental()) {
			log.Info().Bool("Experimental", middleware.Experimental()).Bool("Status", middleware.Status()).Msg("Registering middlewear")
			m.Use(middleware.Method())
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
