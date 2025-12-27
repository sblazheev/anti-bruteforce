package internalhttp

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"                  //nolint:depguard
	httpSwagger "github.com/swaggo/http-swagger/v2"           //nolint:depguard
	"gitlab.wsrubi.ru/go/anti-bruteforce/internal/app"        //nolint:depguard
	"gitlab.wsrubi.ru/go/anti-bruteforce/internal/common"     //nolint:depguard
	"gitlab.wsrubi.ru/go/anti-bruteforce/internal/common/dto" //nolint:depguard
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

type JSONResponse struct {
	Ok bool `json:"ok,omitempty"`
} // @name JSONResponse .

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
// @Param        data body dto.Request true  "Создание события"
// @Success      200  {object} JSONResponse
// @Success      429  {object} JSONResponse
// @Failure		 400  {object} JSONErrorResponse
// @Failure		 503  {object} JSONErrorResponse
// @Router       /auth [post] .
func (h *HTTPHandler) allowAuthHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("X-Content-Type-Options", "nosniff")
	}()
	var dtoRequest *dto.Request
	if err := json.NewDecoder(r.Body).Decode(&dtoRequest); err != nil {
		h.logger.Debug("allowAuthHandler-Invalid request body", "err", err)
		JSONError(http.StatusBadRequest, strconv.Itoa(http.StatusBadRequest), "Invalid request body", err, w)
		return
	}
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(dtoRequest)
	if err != nil {
		h.logger.Debug("allowAuthHandler-Invalid format Event", "err", err)
		JSONError(http.StatusBadRequest, strconv.Itoa(http.StatusBadRequest), "Invalid format Event", err, w)
		return
	}

	allow, err := h.app.CheckWhiteList(dtoRequest.IP)
	if err != nil {
		h.logger.Error("allowAuthHandler", "err", err)
		JSONError(http.StatusServiceUnavailable, strconv.Itoa(http.StatusServiceUnavailable),
			"Service Unavailable", common.ErrServiceUnavailable, w)
		return
	}
	if allow {
		Success(w)
		return
	}

	allow, err = h.app.CheckBlackList(dtoRequest.IP)
	if err != nil {
		h.logger.Error("allowAuthHandler", "err", err)
		JSONError(http.StatusServiceUnavailable, strconv.Itoa(http.StatusServiceUnavailable),
			"Service Unavailable", common.ErrServiceUnavailable, w)
		return
	}
	if allow {
		Failure(w)
		return
	}

	allow, err = h.app.CheckAuthLogin(dtoRequest.Login)
	if err != nil {
		h.logger.Error("allowAuthHandler", "err", err)
		JSONError(http.StatusServiceUnavailable, strconv.Itoa(http.StatusServiceUnavailable),
			"Service Unavailable", common.ErrServiceUnavailable, w)
		return
	}
	if !allow {
		Failure(w)
		return
	}

	allow, err = h.app.CheckAuthPassword(dtoRequest.Password)
	if err != nil {
		h.logger.Error("allowAuthHandler", "err", err)
		JSONError(http.StatusServiceUnavailable, strconv.Itoa(http.StatusServiceUnavailable),
			"Service Unavailable", common.ErrServiceUnavailable, w)
		return
	}
	if !allow {
		Failure(w)
		return
	}

	allow, err = h.app.CheckAuthIP(dtoRequest.Password)
	if err != nil {
		h.logger.Error("allowAuthHandler", "err", err)
		JSONError(http.StatusServiceUnavailable, strconv.Itoa(http.StatusServiceUnavailable),
			"Service Unavailable", common.ErrServiceUnavailable, w)
		return
	}
	if !allow {
		Failure(w)
		return
	}

	Success(w)
}

func JSONError(httpcode int, code string, messageError string, err error, w http.ResponseWriter) {
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

func Failure(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusTooManyRequests)
	w.Write([]byte("{\"ok\":false}"))
}

func Success(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"ok\":true}"))
}
