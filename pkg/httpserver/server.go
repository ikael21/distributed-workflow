package httpserver

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/gin-contrib/logger"
)

type server struct {
	engine *gin.Engine
	srv    *http.Server
}

type Config struct {
	WriteTimeout time.Duration
	ReadTimeout  time.Duration
	Logger       zerolog.Logger
	Addr         string
	Middlewares  []gin.HandlerFunc
}

func New(cfg Config) *server {
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(logger.SetLogger(
		logger.WithLogger(
			func(_ *gin.Context, _ zerolog.Logger) zerolog.Logger {
				return cfg.Logger
			},
		),
	))

	if len(cfg.Middlewares) > 0 {
		engine.Use(cfg.Middlewares...)
	}

	srv := &http.Server{
		WriteTimeout: cfg.WriteTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		Handler:      engine,
		Addr:         cfg.Addr,
	}

	return &server{
		engine: engine,
		srv:    srv,
	}
}

func (s *server) SetupRoutes(setup func(*gin.Engine)) {
	setup(s.engine)
}

func (s *server) Start() error {
	return s.srv.ListenAndServe()
}

func (s *server) Close(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
