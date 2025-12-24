package internalhttp

import (
	"encoding/json"
	"fmt"
	"gitlab.wsrubi.ru/go/anti-bruteforce/internal/common/dto"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
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

	mux.HandleFunc("/auth", handler.allowAuthHandler)

	mux.HandleFunc("/ping", handler.pingHandler)
	mux.HandleFunc("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("doc.json"),
	))
	return handler
}

// @Summary      Пинг-Понг
// @Description  Пинг-Понг
// @Tags         Test
// @Success      200
// @Router       /ping [get] .
func (h *HTTPHandler) pingHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("pong"))
}

// @Summary      Попытка авторизации
// @Description  Попытка авторизации
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        data body dto.Event true  "Создание события"
// @Success      200  {object} dto.Event
// @Failure		 400  {object} JSONErrorResponse
// @Failure		 503  {object} JSONErrorResponse
// @Router       /auth [post] .
func (h *HTTPHandler) allowAuthHandler(w http.ResponseWriter, r *http.Request) {
	var dtoEvent *dto.Event
	if err := json.NewDecoder(r.Body).Decode(&dtoEvent); err != nil {
		h.logger.Debug("createEventHandler-Invalid request body", "err", err)
		JSONError(http.StatusBadRequest, strconv.Itoa(http.StatusBadRequest), "Invalid request body", err, w)
		return
	}
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(dtoEvent)
	if err != nil {
		h.logger.Debug("createEventHandler-Invalid format Event", "err", err)
		JSONError(http.StatusBadRequest, strconv.Itoa(http.StatusBadRequest), "Invalid format Event", err, w)
		return
	}
	allow, err := h.app.CheckAuthLogin(dtoEvent.Login)
	if err != nil {
		h.logger.Error("createEventHandler", "err", err)
		JSONError(http.StatusServiceUnavailable, strconv.Itoa(http.StatusServiceUnavailable),
			"Service Unavailable", common.ErrServiceUnavailable, w)
		return
	}
	if allow {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusTooManyRequests)
	}
	w.Write([]byte(fmt.Sprintf("ok=%v", allow)))
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
