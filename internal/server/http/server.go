package internalhttp

import (
	"context"
	"net"
	"net/http"
	"time"

	"gitlab.wsrubi.ru/go/anti-bruteforce/internal/app"              //nolint:depguard
	"gitlab.wsrubi.ru/go/anti-bruteforce/internal/common"           //nolint:depguard
	"gitlab.wsrubi.ru/go/anti-bruteforce/internal/config"           //nolint:depguard
	"gitlab.wsrubi.ru/go/anti-bruteforce/internal/server/http/docs" //nolint:depguard
)

type Server struct {
	Address string
	logger  common.LoggerInterface
	app     app.App
	config  config.HTTPConfig
	server  *http.Server
}

type StatusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *StatusResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func NewStatusResponseWriter(w http.ResponseWriter) *StatusResponseWriter {
	return &StatusResponseWriter{w, http.StatusOK}
}

// @title           Anti Bruteforce API
// @version         0.1
// @description     Anti Bruteforce
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      api.anti-bruteforce.local
// @BasePath  /

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func NewServer(app app.App, config config.HTTPConfig, logger common.LoggerInterface) *Server {
	address := net.JoinHostPort(config.Host, config.Port)

	httpHandler := NewHandler(app, logger)

	docs.SwaggerInfo.Schemes = []string{"http"}

	server := &Server{
		Address: address,
		logger:  logger,
		app:     app,
		config:  config,
		server: &http.Server{
			Addr:           address,
			Handler:        errorJSONMiddleware(loggingMiddleware(httpHandler.mux, logger)),
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
	}

	return server
}

func (s *Server) Start(_ context.Context) error {
	s.server.Addr = s.Address
	err := s.server.ListenAndServe()
	return err
}

func (s *Server) Stop(ctx context.Context) error {
	<-ctx.Done()
	return s.server.Shutdown(ctx)
}
