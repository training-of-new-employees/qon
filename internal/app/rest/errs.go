package rest

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/training-of-new-employees/qon/internal/errs"
	"github.com/training-of-new-employees/qon/internal/logger"
)

// errorToCode - преобразование ошибки приложения в http-код.
var errorToCode = map[string]int{
	errs.ErrNoAccess.Error():                  http.StatusForbidden,
	errs.ErrEmployeeHasAnotherCompany.Error(): http.StatusForbidden,
	errs.ErrArchiveAdmin.Error():              http.StatusForbidden,

	errs.ErrNotFound.Error():          http.StatusNotFound,
	errs.ErrLessonNotFound.Error():    http.StatusNotFound,
	errs.ErrUserNotFound.Error():      http.StatusNotFound,
	errs.ErrCompanyNotFound.Error():   http.StatusNotFound,
	errs.ErrPositionNotFound.Error():  http.StatusNotFound,
	errs.ErrCourseNotFound.Error():    http.StatusNotFound,
	errs.ErrPositionReference.Error(): http.StatusNotFound,
	errs.ErrUserReference.Error():     http.StatusNotFound,
	errs.ErrCourseReference.Error():   http.StatusNotFound,
	errs.ErrCompanyReference.Error():  http.StatusNotFound,
	errs.ErrLessonReference.Error():   http.StatusNotFound,
	errs.ErrInvalidRoute.Error():      http.StatusNotFound,
	errs.ErrCreaterNotFound.Error():   http.StatusNotFound,
	errs.ErrCompanyNoPosition.Error(): http.StatusNotFound,

	errs.ErrVerifyCodeNotEmpty.Error():       http.StatusBadRequest,
	errs.ErrIncorrectVerifyCode.Error():      http.StatusBadRequest,
	errs.ErrEmailOrPasswordEmpty.Error():     http.StatusBadRequest,
	errs.ErrBadRequest.Error():               http.StatusBadRequest,
	errs.ErrInvalidRequest.Error():           http.StatusBadRequest,
	errs.ErrEmailNotEmpty.Error():            http.StatusBadRequest,
	errs.ErrInviteNotEmpty.Error():           http.StatusBadRequest,
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
	errs.ErrUserPatronymicNotEmpty.Error():   http.StatusBadRequest,
	errs.ErrInvalidCourseName.Error():        http.StatusBadRequest,
	errs.ErrCourseNameIsEmpty.Error():        http.StatusBadRequest,
	errs.ErrInvalidCourseDescription.Error(): http.StatusBadRequest,
	errs.ErrInvalidCourseStatus.Error():      http.StatusBadRequest,
	errs.ErrInvalidLessonStatus.Error():      http.StatusBadRequest,
	errs.ErrInvalidLessonName.Error():        http.StatusBadRequest,
	errs.ErrInvalidTextContent.Error():       http.StatusBadRequest,
	errs.ErrURLPictureLength.Error():         http.StatusBadRequest,
	errs.ErrCreaterNotEmpty.Error():          http.StatusBadRequest,
	errs.ErrCourseIDNotEmpty.Error():         http.StatusBadRequest,
	errs.ErrLessonNameNotEmpty.Error():       http.StatusBadRequest,
	errs.ErrLessonIDNotEmpty.Error():         http.StatusBadRequest,
	errs.ErrTextContentNotEmpty.Error():      http.StatusBadRequest,
	errs.ErrURLPictureNotEmpty.Error():       http.StatusBadRequest,
	errs.ErrUserIDNotEmpty.Error():           http.StatusBadRequest,

	errs.ErrEmailAlreadyExists.Error(): http.StatusConflict,
	errs.ErrUserActivated.Error():      http.StatusConflict,

	errs.ErrUnauthorized.Error():             http.StatusUnauthorized,
	errs.ErrIncorrectEmailOrPassword.Error(): http.StatusUnauthorized,
	errs.ErrInvalidInviteCode.Error():        http.StatusUnauthorized,
	errs.ErrNotFirstLogin.Error():            http.StatusMethodNotAllowed,
	errs.ErrOnlyAdmin.Error():                http.StatusMethodNotAllowed,

	errs.ErrPositionCourseUsed.Error(): http.StatusConflict,
	errs.ErrUserCourseUsed.Error():     http.StatusConflict,
	errs.ErrAssignLessonUsed.Error():   http.StatusConflict,

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
		err = errs.ErrInternal
	}

	// mock-рассылка
	// TODO: конструкция выглядит слишком громоздкой, нужно искать более элегантное решение
	if strings.HasPrefix(err.Error(), "<mock-sender>: ") {
		err = fmt.Errorf(strings.TrimPrefix(err.Error(), "<mock-sender>: "))

		httpCode = http.StatusOK
	}

	c.AbortWithStatusJSON(httpCode, errResponse{Message: err.Error()})
}
