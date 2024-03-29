package asahi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"

	c "github.com/SpeedSlime/Asahi/config"
	db "github.com/SpeedSlime/Asahi/database"
	"github.com/SpeedSlime/Asahi/middleware"
	"github.com/SpeedSlime/Asahi/router"
	"github.com/go-chi/chi"
	"github.com/rs/zerolog"
)

var log zerolog.Logger

type Server struct {
	status      bool
	idle        chan struct{}
	middlewares []middleware.Middleware
	routers     []router.Router
}

func (s Server) Status() bool {
	return s.status
}

func New() *Server {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006/01/02 15:04:05", NoColor: true}
	output.FormatLevel = func(i interface{}) string {
		switch strings.ToUpper(fmt.Sprintf("%s", i)) {
		case "INFO":
			return fmt.Sprintf("\x1b[36m\"\x1b[0m\x1b[35;1m%s\x1b[0m", strings.ToUpper(fmt.Sprintf("%s", i)))
		case "WARN":
			return fmt.Sprintf("\x1b[36m\"\x1b[0m\x1b[33;1m%s\x1b[0m", strings.ToUpper(fmt.Sprintf("%s", i)))
		case "ERROR":
			return fmt.Sprintf("\x1b[36m\"\x1b[0m\x1b[91;1m%s\x1b[0m", strings.ToUpper(fmt.Sprintf("%s", i)))
		case "FATAL":
			return fmt.Sprintf("\x1b[36m\"\x1b[0m\x1b[31;1m%s\x1b[0m", strings.ToUpper(fmt.Sprintf("%s", i)))
		default:
			return fmt.Sprintf("\x1b[36m\"\x1b[0m\x1b[95;1m%s\x1b[0m", strings.ToUpper(fmt.Sprintf("%s", i)))
		} 
	}
	output.FormatMessage = func(i interface{}) string { return fmt.Sprintf("\x1b[36m%s\x1b[0m\x1b[36m\"\x1b[0m", i) }
	output.FormatFieldName = func(i interface{}) string { return fmt.Sprintf("%s: ", i) }
	output.FormatErrFieldName = func(i interface{}) string { return "Error: " }
	output.FormatErrFieldValue = func(i interface{}) string { 
		v := fmt.Sprint(i)
		if v[0] == '"' {
			v = v[1:]
		}
		if i := len(v)-1; v[i] == '"' {
			v = v[:i]
		}
		return fmt.Sprintf("\"\x1b[31;1m%s\x1b[0m\"", v) 
	}
	output.FormatFieldValue = func(i interface{}) string {
		v := fmt.Sprintf("%s", i)
		if v == "false" {
			return fmt.Sprintf("\"\x1b[31;1m%s\x1b[0m\"", v)
		}
		if v == "true" {
			return fmt.Sprintf("\"\x1b[32;1m%s\x1b[0m\"", v)
		}
		if _, err := strconv.Atoi(v); err == nil {
			return fmt.Sprintf("\"\x1b[34;1m%s\x1b[0m\"", v)
		}
		return fmt.Sprintf("\"\x1b[35;1m%s\x1b[0m\"", v)
	}
	log = zerolog.New(output).With().Timestamp().Logger()

	err := c.New()
	if err != nil {
		log.Fatal().Str("Function", "New").Err(err).Msg("Unexpected error")
	}

	db.Connect()
	if err != nil {
		log.Fatal().Str("Function", "New").Err(err).Msg("Unexpected error")
	}

	return &Server{}
}

func (s *Server) LoadRouter(routers []router.Router) {
	s.routers = append(s.routers, routers...)
}

func (s *Server) LoadMiddleware(middlewares []middleware.Middleware) {
	s.middlewares = append(s.middlewares, middlewares...)
}

func (s *Server) Start() {
	instance := &http.Server{Addr: fmt.Sprintf("%s:%s", c.Address(), c.Port()), Handler: s.handler()}

	log.Info().Str("Port", c.Port()).Str("Address", c.Address()).Bool("Experimental", c.Experimental()).Bool("Status", s.Status()).Msg("Server started")

	s.idle = make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		if err := instance.Shutdown(context.Background()); err != nil {
			log.Warn().Str("Function", "Shutdown").Err(err).Msg("Server shutdown failed!")
		}
		close(s.idle)
	}()

	if err := instance.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal().Str("Function", "ListenAndServe").Err(err).Msg("Unexpected error")
	}

	<-s.idle

	log.Info().Str("Port", c.Port()).Str("Address", c.Address()).Bool("Experimental", c.Experimental()).Bool("Status", s.Status()).Msg("Server stopped")
}

func Handle(err error, function string) {
	if err != nil {
		log.Error().Str("Function", function).Err(err).Msg("Unexpected error")
	}
}

func (s *Server) handler() *chi.Mux {
	m := chi.NewMux()
	for _, middleware := range s.middlewares {
		if middleware.Status() && (middleware.Experimental() == c.Experimental() || !middleware.Experimental()) {
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
				if r.Status() && (r.Experimental() == c.Experimental() || !r.Experimental()) {
					log.Info().Bool("Experimental", r.Experimental()).Bool("Status", r.Status()).Str("Method", r.Method()).Str("Path", r.Path()).Msg("Registering route")
					m.Method(r.Method(), r.Path(), r.Handler())
				}
			}
		}
	}
}

func Parameter(r *http.Request, k string) string {
	return chi.URLParam(r, k)
}

func ParseJSON(r *http.Request, i interface{}) error {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	err := d.Decode(i)
	if err != nil {
		return err
	}
	return nil
}