package internalhttp

import (
	"encoding/json"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger/v2"       //nolint:depguard
	"gitlab.wsrubi.ru/go/anti-bruteforce/internal/app"    //nolint:depguard
	"gitlab.wsrubi.ru/go/anti-bruteforce/internal/common" //nolint:depguard
)

type HTTPHandler struct {
	app    app.App
	logger common.LoggerInterface
	mux    *http.ServeMux
}

type JSONErrorResponse struct {
	Code    *string `json:"code,omitempty"`
	Message *string `json:"message,omitempty"`
	Detail  string  `json:"description,omitempty"`
} // @name JSONError .

func NewHandler(app app.App, logger common.LoggerInterface) *HTTPHandler {
	mux := http.NewServeMux()

	handler := &HTTPHandler{app, logger, mux}

	mux.HandleFunc("/ping", handler.pingHandler)
	mux.HandleFunc("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("doc.json"),
	))
	return handler
}

func (h *HTTPHandler) pingHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("pong"))
}

func JSONError(httpcode int, code string, messageError string, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	if httpcode == 503 {
		w.Header().Set("Retry-After", "600")
	}
	w.WriteHeader(httpcode)
	json.NewEncoder(w).Encode(
		JSONErrorResponse{
			Code:    &code,
			Message: &messageError,
			Detail:  err.Error(),
		},
	)
}
