package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/training-of-new-employees/qon/internal/errs"
	"github.com/training-of-new-employees/qon/internal/logger"
	"go.uber.org/zap"
)

// errorToCode - мапа соответствия ошибки http-коду.
var errorToCode = map[string]int{

	errs.ErrNotFound.Error():     http.StatusNotFound,
	errs.ErrUserNotFound.Error(): http.StatusNotFound,

	errs.ErrBadRequest.Error():           http.StatusBadRequest,
	errs.ErrInvalidRequest.Error():       http.StatusBadRequest,
	errs.ErrEmailNotEmpty.Error():        http.StatusBadRequest,
	errs.ErrInvalidEmail.Error():         http.StatusBadRequest,
	errs.ErrPasswordNotEmpty.Error():     http.StatusBadRequest,
	errs.ErrInvalidPassword.Error():      http.StatusBadRequest,
	errs.ErrCompanyNameNotEmpty.Error():  http.StatusBadRequest,
	errs.ErrIncorrectCompanyName.Error(): http.StatusBadRequest,

	errs.ErrEmailAlreadyExists.Error(): http.StatusConflict,

	errs.ErrUnauthorized.Error():  http.StatusUnauthorized,
	errs.ErrNotFirstLogin.Error(): http.StatusMethodNotAllowed,
	errs.ErrOnlyAdmin.Error():     http.StatusMethodNotAllowed,

	errs.ErrNotSendEmail.Error(): http.StatusInternalServerError,
	errs.ErrInternal.Error():     http.StatusInternalServerError,
}

// errResponse - тело для ответа.
type errResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// handleError - преобразование ошибки приложения в http-ответ.
func (r RestServer) handleError(c *gin.Context, err error) {
	logger.Log.Error("handler error", zap.Error(err))

	httpCode, exists := errorToCode[err.Error()]
	if !exists {
		httpCode = http.StatusInternalServerError
	}

	c.JSON(httpCode, errResponse{Code: httpCode, Message: err.Error()})
}
