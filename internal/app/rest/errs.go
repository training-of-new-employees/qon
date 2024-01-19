package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/training-of-new-employees/qon/internal/errs"
	"github.com/training-of-new-employees/qon/internal/logger"
)

// errorToCode - преобразование ошибки приложения в http-код.
var errorToCode = map[string]int{
	errs.ErrNoAccess.Error(): http.StatusForbidden,

	errs.ErrNotFound.Error():          http.StatusNotFound,
	errs.ErrUserNotFound.Error():      http.StatusNotFound,
	errs.ErrCompanyNotFound.Error():   http.StatusNotFound,
	errs.ErrPositionNotFound.Error():  http.StatusNotFound,
	errs.ErrCourseNotFound.Error():    http.StatusNotFound,
	errs.ErrPositionReference.Error(): http.StatusNotFound,
	errs.ErrUserReference.Error():     http.StatusNotFound,
	errs.ErrCourseReference.Error():   http.StatusNotFound,
	errs.ErrCompanyReference.Error():  http.StatusNotFound,
	errs.ErrLessonReference.Error():   http.StatusNotFound,

	errs.ErrBadRequest.Error():               http.StatusBadRequest,
	errs.ErrInvalidRequest.Error():           http.StatusBadRequest,
	errs.ErrEmailNotEmpty.Error():            http.StatusBadRequest,
	errs.ErrInvalidEmail.Error():             http.StatusBadRequest,
	errs.ErrPasswordNotEmpty.Error():         http.StatusBadRequest,
	errs.ErrInvalidPassword.Error():          http.StatusBadRequest,
	errs.ErrCompanyNameNotEmpty.Error():      http.StatusBadRequest,
	errs.ErrInvalidCompanyName.Error():       http.StatusBadRequest,
	errs.ErrCompanyIDNotEmpty.Error():        http.StatusBadRequest,
	errs.ErrPositionNameNotEmpty.Error():     http.StatusBadRequest,
	errs.ErrInvalidPositionName.Error():      http.StatusBadRequest,
	errs.ErrPositionIDNotEmpty.Error():       http.StatusBadRequest,
	errs.ErrUserNameNotEmpty.Error():         http.StatusBadRequest,
	errs.ErrInvalidUserName.Error():          http.StatusBadRequest,
	errs.ErrUserSurnameNotEmpty.Error():      http.StatusBadRequest,
	errs.ErrInvalidUserSurname.Error():       http.StatusBadRequest,
	errs.ErrInvalidUserPatronymic.Error():    http.StatusBadRequest,
	errs.ErrCourseNameInvalid.Error():        http.StatusBadRequest,
	errs.ErrCourseNameIsEmpty.Error():        http.StatusBadRequest,
	errs.ErrCourseDescriptionInvalid.Error(): http.StatusBadRequest,

	errs.ErrEmailAlreadyExists.Error(): http.StatusConflict,
	errs.ErrUserActivated.Error():      http.StatusConflict,

	errs.ErrUnauthorized.Error():  http.StatusUnauthorized,
	errs.ErrNotFirstLogin.Error(): http.StatusMethodNotAllowed,
	errs.ErrOnlyAdmin.Error():     http.StatusMethodNotAllowed,

	errs.ErrNotSendEmail.Error(): http.StatusInternalServerError,
	errs.ErrInternal.Error():     http.StatusInternalServerError,
}

// errResponse - тело для ответа.
type errResponse struct {
	Message string `json:"message"`
}

// handleError - преобразование ошибки приложения в http-ответ.
func (r RestServer) handleError(c *gin.Context, err error) {
	logger.Log.Error("handler error", zap.Error(err))

	httpCode, exists := errorToCode[err.Error()]
	if !exists {
		httpCode = http.StatusInternalServerError
	}

	c.AbortWithStatusJSON(httpCode, errResponse{Message: err.Error()})
}
