package asahi

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/SpeedSlime/Asahi/config"
	"github.com/SpeedSlime/Asahi/middleware"
	"github.com/SpeedSlime/Asahi/router"
	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
)

type Server struct {	
	status       bool
	idle		 chan struct{}
	config		 *config.Config 
	middlewares  []middleware.Middleware
	routers      []router.Router
}

func (s Server) Status() bool {
	return s.status
}

func (s Server) Config() *config.Config {
	return s.config
}

func New() *Server {
	c, err := config.New()
	if err != nil {
		log.Fatal().Str("Function", "New()").Err(err).Msg("Unexpected error")
	}
	return &Server{config: c}
}

func (s *Server) LoadRouter(routers []router.Router) {
	s.routers = append(s.routers, routers...)
}

func (s *Server) LoadMiddleware(middlewares []middleware.Middleware) {
	s.middlewares = append(s.middlewares, middlewares...)
}

func (s *Server) Start() {
	instance := &http.Server{Addr: s.config.Port(), Handler: s.handler()}

	log.Info().Str("Port", s.config.Port()).Str("Address", s.config.Address()).Bool("Experimental", s.config.Experimental()).Bool("Status", s.Status()).Msg("Server started")

	s.idle = make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		if err := instance.Shutdown(context.Background()); err != nil {
			log.Warn().Str("Function", "Shutdown()").Err(err).Msg("Server shutdown failed!")
		}
		close(s.idle)
	}()

	if err := instance.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal().Str("Function", "ListenAndServe()").Err(err).Msg("Unexpected error")
	}

	<-s.idle
	
	log.Info().Str("Port", s.config.Port()).Str("Address", s.config.Address()).Bool("Experimental", s.config.Experimental()).Bool("Status", s.Status()).Msg("Server stopped")
}

func (s *Server) handler() (*chi.Mux) {
	m := chi.NewMux()
	for _, middleware := range s.middlewares {
		if middleware.Status() && (middleware.Experimental() == s.config.Experimental() || !middleware.Experimental()) {
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
				if r.Status() && (r.Experimental() == s.config.Experimental() || !r.Experimental()) {
					log.Info().Bool("Experimental", r.Experimental()).Bool("Status", r.Status()).Str("Method", r.Method()).Str("Path", r.Path()).Msg("Registering route")
					m.Method(r.Method(), r.Path(), r.Handler())
				}
			}
		}
	}
}