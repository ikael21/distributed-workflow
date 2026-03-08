package http

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/gin-contrib/logger"
)

type Server struct {
	engine *gin.Engine
	srv    *http.Server
}

type Config struct {
	WriteTimeout time.Duration
	ReadTimeout  time.Duration
	Logger       zerolog.Logger
	Addr         string
}

func NewServer(cfg Config) *Server {
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(logger.SetLogger(
		logger.WithLogger(
			func(_ *gin.Context, _ zerolog.Logger) zerolog.Logger {
				return cfg.Logger
			},
		),
	))

	srv := &http.Server{
		WriteTimeout: cfg.WriteTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		Handler:      engine,
		Addr:         cfg.Addr,
	}

	return &Server{
		engine: engine,
		srv:    srv,
	}
}

func (s *Server) Router() *gin.Engine {
	return s.engine
}

func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Close(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
