package internalhttp

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/netip"
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

	mux.HandleFunc("/check/auth", handler.allowAuthHandler)

	mux.HandleFunc("/delete/bucket", handler.deleteBucketHandler)
	mux.HandleFunc("/delete/white-list", handler.deleteWhiteList)
	mux.HandleFunc("/delete/black-list", handler.deleteBlackList)

	mux.HandleFunc("/add/black-list", handler.addBlackList)
	mux.HandleFunc("/add/white-list", handler.addWhiteList)

	mux.HandleFunc("/ping", handler.pingHandler)
	mux.HandleFunc("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("doc.json"),
	))
	return handler
}

func (h *HTTPHandler) parseRequestNet(dtoRequest *dto.RequestNet, w http.ResponseWriter, r *http.Request) bool {
	if err := json.NewDecoder(r.Body).Decode(dtoRequest); err != nil {
		h.logger.Debug("deleteWhiteList-Invalid request body", "err", err)
		JSONError(http.StatusBadRequest, strconv.Itoa(http.StatusBadRequest), "Invalid request body", err, w)
		return false
	}
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(dtoRequest)
	if err != nil {
		h.logger.Debug("deleteWhiteList-Invalid format Event", "err", err)
		JSONError(http.StatusBadRequest, strconv.Itoa(http.StatusBadRequest), "Invalid format Event", err, w)
		return false
	}

	_, err = netip.ParsePrefix(dtoRequest.Net)
	if err != nil {
		JSONError(http.StatusBadRequest, strconv.Itoa(http.StatusBadRequest),
			"Bad Request", common.ErrFormatIp, w)
	}

	return true
}

// @Summary      Удалить сеть из белого списка
// @Description  Удалить сеть из белого списка
// @Tags         Delete
// @Accept       json
// @Produce      json
// @Param        data body dto.RequestNet true  "Удалить сеть из белого списка"
// @Success      200
// @Failure		 400  {object} JSONErrorResponse
// @Failure		 503  {object} JSONErrorResponse
// @Router       /delete/white-list [post] .
func (h *HTTPHandler) deleteWhiteList(w http.ResponseWriter, r *http.Request) {
	var dtoRequest dto.RequestNet
	if !h.parseRequestNet(&dtoRequest, w, r) {
		return
	}
	err := h.app.DeleteWhiteList(dtoRequest.Net)
	if err != nil {
		h.logger.Error("DeleteWhiteList", "err", err)
		JSONError(http.StatusServiceUnavailable, strconv.Itoa(http.StatusServiceUnavailable),
			"Service Unavailable", common.ErrServiceUnavailable, w)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// @Summary      Удалить сеть из черного списка
// @Description  Удалить сеть из черного списка
// @Tags         Delete
// @Accept       json
// @Produce      json
// @Param        data body dto.RequestNet true  "Удалить сеть из черного списка"
// @Success      200
// @Failure		 400  {object} JSONErrorResponse
// @Failure		 503  {object} JSONErrorResponse
// @Router       /delete/black-list [post] .
func (h *HTTPHandler) deleteBlackList(w http.ResponseWriter, r *http.Request) {
	var dtoRequest dto.RequestNet
	if !h.parseRequestNet(&dtoRequest, w, r) {
		return
	}
	err := h.app.DeleteBlackList(dtoRequest.Net)
	if err != nil {
		h.logger.Error("DeleteBlackList", "err", err)
		JSONError(http.StatusServiceUnavailable, strconv.Itoa(http.StatusServiceUnavailable),
			"Service Unavailable", common.ErrServiceUnavailable, w)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// @Summary      Добавить сеть в черный список
// @Description  Добавить сеть в черный список
// @Tags         Add
// @Accept       json
// @Produce      json
// @Param        data body dto.RequestNet true  "Добавить сеть в черный список"
// @Success      200
// @Failure		 400  {object} JSONErrorResponse
// @Failure		 503  {object} JSONErrorResponse
// @Router       /add/black-list [post] .
func (h *HTTPHandler) addBlackList(w http.ResponseWriter, r *http.Request) {
	var dtoRequest dto.RequestNet
	if !h.parseRequestNet(&dtoRequest, w, r) {
		return
	}
	_, err := h.app.AddBlackList(dtoRequest.Net)
	if err != nil {
		if errors.Is(err, common.ErrIPSubnetOverlapped) {
			JSONError(http.StatusBadRequest, strconv.Itoa(http.StatusBadRequest),
				"Bad Request", common.ErrIPSubnetOverlapped, w)
			return
		}
		h.logger.Error("AddBlackList", "err", err)
		JSONError(http.StatusServiceUnavailable, strconv.Itoa(http.StatusServiceUnavailable),
			"Service Unavailable", common.ErrServiceUnavailable, w)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// @Summary      Добавить сеть в белый список
// @Description  Добавить сеть в белый список
// @Tags         Add
// @Accept       json
// @Produce      json
// @Param        data body dto.RequestNet true  "Добавить сеть в белый список"
// @Success      200
// @Failure		 400  {object} JSONErrorResponse
// @Failure		 503  {object} JSONErrorResponse
// @Router       /add/white-list [post] .
func (h *HTTPHandler) addWhiteList(w http.ResponseWriter, r *http.Request) {
	var dtoRequest dto.RequestNet
	if !h.parseRequestNet(&dtoRequest, w, r) {
		return
	}
	_, err := h.app.AddWhiteList(dtoRequest.Net)
	if err != nil {
		if errors.Is(err, common.ErrIPSubnetOverlapped) {
			JSONError(http.StatusBadRequest, strconv.Itoa(http.StatusBadRequest),
				"Bad Request", common.ErrIPSubnetOverlapped, w)
			return
		}
		h.logger.Error("addWhiteList", "err", err)
		JSONError(http.StatusServiceUnavailable, strconv.Itoa(http.StatusServiceUnavailable),
			"Service Unavailable", common.ErrServiceUnavailable, w)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// @Summary      Удалить бакет
// @Description  Удалить бакет
// @Tags         Delete
// @Accept       json
// @Produce      json
// @Param        data body dto.RequestDeleteBucket true  "Удалить бакет"
// @Success      200
// @Failure		 400  {object} JSONErrorResponse
// @Failure		 503  {object} JSONErrorResponse
// @Router       /delete/bucket [post] .
func (h *HTTPHandler) deleteBucketHandler(w http.ResponseWriter, r *http.Request) {
	var dtoRequest *dto.RequestDeleteBucket
	if err := json.NewDecoder(r.Body).Decode(&dtoRequest); err != nil {
		h.logger.Debug("deleteBucketHandler-Invalid request body", "err", err)
		JSONError(http.StatusBadRequest, strconv.Itoa(http.StatusBadRequest), "Invalid request body", err, w)
		return
	}
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(dtoRequest)
	if err != nil {
		h.logger.Debug("deleteBucketHandler-Invalid format Event", "err", err)
		JSONError(http.StatusBadRequest, strconv.Itoa(http.StatusBadRequest), "Invalid format Event", err, w)
		return
	}
	if len(dtoRequest.IP) > 0 {
		_ = h.app.DeleteIPBucket(dtoRequest.IP)
	}
	if len(dtoRequest.Login) > 0 {
		_ = h.app.DeleteIPBucket(dtoRequest.Login)
	}
	w.WriteHeader(http.StatusOK)
}

// @Summary      Попытка авторизации
// @Description  Попытка авторизации
// @Tags         Check
// @Accept       json
// @Produce      json
// @Param        data body dto.RequestCheck true  "Создание события"
// @Success      200  {object} JSONResponse
// @Success      429  {object} JSONResponse
// @Failure		 400  {object} JSONErrorResponse
// @Failure		 503  {object} JSONErrorResponse
// @Router       /check/auth [post] .
func (h *HTTPHandler) allowAuthHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("X-Content-Type-Options", "nosniff")
	}()
	var dtoRequest *dto.RequestCheck
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
		h.logger.Error("CheckWhiteList", "err", err)
		JSONError(http.StatusBadRequest, strconv.Itoa(http.StatusBadRequest),
			"Bad Request", common.ErrFormatIp, w)
		return
	}
	if allow {
		Success(w, "white_list")
		return
	}

	allow, err = h.app.CheckBlackList(dtoRequest.IP)
	if err != nil {
		h.logger.Error("CheckBlackList", "err", err)
		JSONError(http.StatusBadRequest, strconv.Itoa(http.StatusBadRequest),
			"Bad Request", common.ErrFormatIp, w)
		return
	}
	if allow {
		Failure(w, "black_list")
		return
	}

	allow, err = h.app.CheckAuthLogin(dtoRequest.Login)
	if err != nil {
		h.logger.Error("CheckAuthLogin", "err", err)
		JSONError(http.StatusServiceUnavailable, strconv.Itoa(http.StatusServiceUnavailable),
			"Service Unavailable", common.ErrServiceUnavailable, w)
		return
	}
	if !allow {
		Failure(w, "login")
		return
	}

	allow, err = h.app.CheckAuthPassword(dtoRequest.Password)
	if err != nil {
		h.logger.Error("CheckAuthPassword", "err", err)
		JSONError(http.StatusServiceUnavailable, strconv.Itoa(http.StatusServiceUnavailable),
			"Service Unavailable", common.ErrServiceUnavailable, w)
		return
	}
	if !allow {
		Failure(w, "password")
		return
	}

	allow, err = h.app.CheckAuthIP(dtoRequest.Password)
	if err != nil {
		h.logger.Error("CheckAuthIP", "err", err)
		JSONError(http.StatusServiceUnavailable, strconv.Itoa(http.StatusServiceUnavailable),
			"Service Unavailable", common.ErrServiceUnavailable, w)
		return
	}
	if !allow {
		Failure(w, "ip")
		return
	}

	Success(w, "all")
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

func Failure(w http.ResponseWriter, rule string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Rule", rule)
	w.WriteHeader(http.StatusTooManyRequests)
	w.Write([]byte("{\"ok\":false}"))
}

func Success(w http.ResponseWriter, rule string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Rule", rule)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"ok\":true}"))
}
