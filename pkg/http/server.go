package http

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/gin-contrib/logger"
	"github.com/getkin/kin-openapi/openapi3"
	oapiMiddleware "github.com/oapi-codegen/gin-middleware"
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
	Swagger      *openapi3.T
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
	if cfg.Swagger != nil {
		engine.Use(oapiMiddleware.OapiRequestValidator(cfg.Swagger))
	}

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
