package rest

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/training-of-new-employees/qon/internal/errs"
	"github.com/training-of-new-employees/qon/internal/logger"
	"github.com/training-of-new-employees/qon/internal/model"
	"github.com/training-of-new-employees/qon/internal/pkg/access"
)

// CreateAdmin godoc
//
//	@Summary	Создание администратора
//	@Tags		admin
//	@Produce	json
//	@Param		object	body		model.CreateAdmin	true	"Create Admin"
//	@Success	201		{array}		sEmail
//	@Failure	400		{object}	sErr
//	@Failure	409		{object}	sErr
//	@Failure	500		{object}	sErr
//	@Router		/admin/register [post]
func (r *RestServer) handlerCreateAdminInCache(c *gin.Context) {
	ctx := c.Request.Context()
	createAdmin := model.CreateAdmin{}

	if err := c.ShouldBindJSON(&createAdmin); err != nil {
		r.handleError(c, errs.ErrInvalidRequest)
		return
	}

	if err := createAdmin.Validation(); err != nil {
		r.handleError(c, err)
		return
	}

	admin, err := r.services.User().WriteAdminToCache(ctx, createAdmin)
	if err != nil {
		r.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, s().SetEmail(admin.Email))
}

// CreateUser godoc
//
//	@Summary	Создание пользователя
//	@Tags		admin
//	@Produce	json
//	@Param		object	body		model.UserCreate	true	"User Create"
//	@Success	201		{object}	model.User
//	@Failure	400		{object}	sErr
//	@Failure	409		{object}	sErr
//	@Failure	500		{object}	sErr
//	@Router		/admin/employee [post]
func (r *RestServer) handlerCreateUser(c *gin.Context) {
	ctx := c.Request.Context()
	userReq := model.UserCreate{}

	if err := c.ShouldBindJSON(&userReq); err != nil {
		r.handleError(c, errs.ErrInvalidRequest)
		return
	}

	if err := userReq.Validation(); err != nil {
		r.handleError(c, err)
		return
	}

	user, err := r.services.User().CreateUser(ctx, userReq)
	if err != nil {
		r.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, user)

}

// GetUsers godoc
//
//	@Summary		Получение данных пользователей
//	@Description	Список сотрдуников в компании админа
//	@Tags			user
//	@Produce		json
//	@Success		200	{array}		model.User
//	@Failure		403	{object}	sErr
//	@Failure		404	{object}	sErr
//	@Failure		500	{object}	sErr
//	@Router			/users [get]
func (r *RestServer) handlerGetUsers(c *gin.Context) {
	ctx := c.Request.Context()

	users, err := r.services.User().GetUsersByCompany(ctx, r.getUserSession(c).OrgID)
	if err != nil {
		r.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, users)
}

// GetUser godoc
//
//	@Summary		Получение данных пользователя
//	@Description	Получение по id
//	@Tags			user
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	model.UserInfo
//	@Failure		400	{object}	sErr
//	@Failure		403	{object}	sErr
//	@Failure		404	{object}	sErr
//	@Failure		500	{object}	sErr
//	@Router			/users/{id} [get]
func (r *RestServer) handlerGetUser(c *gin.Context) {
	ctx := c.Request.Context()

	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		r.handleError(c, errs.ErrBadRequest)
		return
	}

	session := r.getUserSession(c)
	user, err := r.services.User().GetUserByID(ctx, userID)
	if err != nil {
		r.handleError(c, err)
		return
	}

	if !access.CanUser(session.IsAdmin, session.OrgID, session.UserID, user.ID, user.CompanyID) {
		r.handleError(c, errs.ErrNoAccess)
		return
	}

	c.JSON(http.StatusOK, user)
}

// EditUser godoc
//
//	@Summary		Изменение данных пользователя
//	@Description	Изменение по id
//	@Tags			user
//	@Produce		json
//	@Param			id		path		int				true	"User ID"
//	@Param			object	body		model.UserEdit	true	"User info"
//	@Success		200		{object}	model.UserEdit
//	@Failure		400		{object}	sErr
//	@Failure		403		{object}	sErr
//	@Failure		404		{object}	sErr
//	@Failure		500		{object}	sErr
//	@Router			/users/{id} [patch]
func (r *RestServer) handlerEditUser(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		r.handleError(c, errs.ErrBadRequest)
		return
	}

	edit := &model.UserEdit{}
	if err := c.ShouldBindJSON(&edit); err != nil {
		r.handleError(c, errs.ErrInvalidRequest)
		return
	}
	edit.ID = id

	session := r.getUserSession(c)
	if !access.CanUser(session.IsAdmin, session.OrgID, session.UserID, edit.ID, session.OrgID) {
		r.handleError(c, errs.ErrNoAccess)
		return
	}

	edited, err := r.services.User().EditUser(ctx, edit, session.OrgID)
	if err != nil {
		r.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, edited)
}

// SetPassword godoc
//
//	@Summary	Активация пользователя и установка ему пароля
//	@Tags		user
//	@Produce	json
//	@Param		object	body		model.UserSignIn	true	"User Set Password"
//	@Success	200		{object}	sToken
//	@Failure	400		{object}	sErr
//	@Failure	401		{object}	sErr
//	@Failure	404		{object}	sErr
//	@Failure	500		{object}	sErr
//	@Router		/users/set-password [post]
func (r *RestServer) handlerSetPassword(c *gin.Context) {
	ctx := c.Request.Context()
	userReq := model.UserSignIn{}

	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(http.StatusBadRequest, s().SetError(err))
		return
	}

	if err := userReq.Validation(); err != nil {
		c.JSON(http.StatusBadRequest, s().SetError(err))
		return
	}

	user, err := r.services.User().GetUserByEmail(ctx, userReq.Email)
	switch {
	case errors.Is(err, errs.ErrUserNotFound):
		c.JSON(http.StatusNotFound, s().SetError(err))
		return
	case err != nil:
		c.JSON(http.StatusInternalServerError, s().SetError(err))
		return
	}

	if user.IsActive {
		c.JSON(http.StatusUnauthorized, s().SetError(errs.ErrNotFirstLogin))
		return
	}

	if err := r.services.User().UpdatePasswordAndActivateUser(ctx, userReq.Email, userReq.Password); err != nil {
		c.JSON(http.StatusInternalServerError, s().SetError(err))
		return
	}

	tokens, err := r.services.User().GenerateTokenPair(ctx, user.ID, user.IsAdmin, user.CompanyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, s().SetError(err))
		return
	}

	c.Header("Authorization", "Bearer "+tokens.AccessToken)

	c.JSON(http.StatusOK, s().SetToken(tokens.AccessToken))
}

// ArchiveUser godoc
//
//	@Summary	Архивирование пользователя по id
//	@Tags		user
//	@Produce	json
//	@Param		id	path	int	true	"User ID"
//	@Success	200
//	@Failure	400	{object}	sErr
//	@Failure	403	{object}	sErr
//	@Failure	404	{object}	sErr
//	@Failure	500	{object}	sErr
//	@Router		/users/archive/{id} [patch]
func (r *RestServer) handlerArchiveUser(c *gin.Context) {
	ctx := c.Request.Context()
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		msg := fmt.Errorf("got invalid user id: %s", idParam)
		logger.Log.Warn("error", zap.Error(msg))
		c.JSON(http.StatusBadRequest, s().SetError(msg))
		return
	}
	us := r.getUserSession(c)
	err = r.services.User().ArchiveUser(ctx, id, us.OrgID)
	switch {
	case errors.Is(err, errs.ErrUserNotFound):
		c.JSON(http.StatusNotFound, s().SetError(err))
		return
	case err != nil:
		c.JSON(http.StatusInternalServerError, s().SetError(err))
		return
	}

	c.Status(http.StatusOK)
}

// SignIn godoc
//
//	@Summary	Вход пользователя
//	@Produce	json
//	@Param		object	body		model.UserSignIn	true	"User SignIn"
//	@Success	200		{object}	sToken
//	@Failure	400		{object}	sErr
//	@Failure	401		{object}	sErr
//	@Failure	500		{object}	sErr
//	@Router		/login [post]
func (r *RestServer) handlerSignIn(c *gin.Context) {
	ctx := c.Request.Context()
	userReq := model.UserSignIn{}

	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(http.StatusBadRequest, s().SetError(err))
		return
	}

	if err := userReq.Validation(); err != nil {
		c.JSON(http.StatusBadRequest, s().SetError(err))
	}

	user, err := r.services.User().GetUserByEmail(ctx, userReq.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, s().SetError(err))
		return
	}

	if err = user.CheckPassword(userReq.Password); err != nil {
		c.JSON(http.StatusUnauthorized, s().SetError(err))
		return
	}

	tokens, err := r.services.User().GenerateTokenPair(ctx, user.ID, user.IsAdmin, user.CompanyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, s().SetError(err))
		return
	}

	c.Header("Authorization", "Bearer "+tokens.AccessToken)

	c.JSON(http.StatusOK, s().SetToken(tokens.AccessToken))
}

// EmailVerification godoc
//
//	@Summary	Верификация email'a пользователя
//	@Tags		admin
//	@Produce	json
//	@Param		object	body		model.Code	true	"User Email Verification"
//	@Success	201		{object}	model.User
//	@Failure	400		{object}	sErr
//	@Failure	401		{object}	sErr
//	@Failure	500		{object}	sErr
//	@Router		/admin/verify [post]
func (r *RestServer) handlerAdminEmailVerification(c *gin.Context) {
	ctx := c.Request.Context()
	code := model.Code{}

	if err := c.ShouldBindJSON(&code); err != nil {
		c.JSON(http.StatusBadRequest, s().SetError(err))
		return
	}

	if err := code.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, s().SetError(err))
		return
	}

	adminFromCache, err := r.services.User().GetAdminFromCache(ctx, code.Code)
	if err != nil {
		c.JSON(http.StatusUnauthorized, s().SetError(err))
		return
	}

	logger.Log.Info("admin from cache: %v", zap.String("email", adminFromCache.Email))

	createdAdmin, err := r.services.User().CreateAdmin(ctx, *adminFromCache)
	if err != nil {
		c.JSON(http.StatusInternalServerError, s().SetError(err))
		return
	}

	_ = r.services.User().DeleteAdminFromCache(ctx, code.Code)

	c.JSON(http.StatusCreated, createdAdmin)

}

// ResetPassword godoc
//
//	@Summary	Сброс пароля пользователя
//	@Produce	json
//	@Param		object	body		model.EmailReset	true	"User Reset Password"
//	@Success	200		{object}	sEmail
//	@Failure	400		{object}	sErr
//	@Failure	404		{object}	sErr
//	@Failure	500		{object}	sErr
//	@Router		/password [post]
func (r *RestServer) handlerResetPassword(c *gin.Context) {
	ctx := c.Request.Context()
	email := model.EmailReset{}

	if err := c.ShouldBindJSON(&email); err != nil {
		r.handleError(c, errs.ErrInvalidRequest)
		return
	}

	if err := r.services.User().ResetPassword(ctx, email.Email); err != nil {
		r.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, s().SetEmail(email.Email))
}

// AdminEdit godoc
//
//	@Summary	Изменение данных администратора
//	@Tags		admin
//	@Produce	json
//	@Param		object	body		model.AdminEdit	true	"Admin Edit"
//	@Success	200		{object}	model.AdminEdit
//	@Failure	400		{object}	sErr
//	@Failure	401		{object}	sErr
//	@Failure	404		{object}	sErr
//	@Failure	500		{object}	sErr
//	@Router		/admin/info [post]
func (r *RestServer) handlerEditAdmin(c *gin.Context) {
	ctx := c.Request.Context()

	edit := model.AdminEdit{}
	if err := c.ShouldBindJSON(&edit); err != nil {
		r.handleError(c, errs.ErrInvalidRequest)
		return
	}
	edit.ID = r.getUserSession(c).UserID

	edited, err := r.services.User().EditAdmin(ctx, edit)
	if err != nil {
		r.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, edited)
}

// @Summary		Регенерация пригласительной ссылки
// @Description	Изменение по email сотрудника
// @Tags			admin
// @Produce		json
// @Param			object	body		model.InvitationLinkRequest	true	"User email"
// @Success		200		{object}	model.InvitationLinkResponse
// @Failure		400		{object}	sErr
// @Failure		401		{object}	sErr
// @Failure		403		{object}	sErr
// @Failure		404		{object}	sErr
// @Failure		409		{object}	sErr
// @Failure		500		{object}	sErr
// @Router			/invitation-link [patch]
func (r *RestServer) handlerRegenerationInvitationLink(c *gin.Context) {
	ctx := c.Request.Context()

	session := r.getUserSession(c)

	invitationLinkRequest := model.InvitationLinkRequest{}

	if err := c.ShouldBindJSON(&invitationLinkRequest); err != nil {
		r.handleError(c, errs.ErrInvalidRequest)
		return
	}

	if err := invitationLinkRequest.Validate(); err != nil {
		r.handleError(c, errs.ErrInvalidRequest)
		return
	}

	response, err := r.services.User().RegenerationInvitationLinkUser(ctx, invitationLinkRequest.Email, session.OrgID)
	if err != nil {
		r.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetUser godoc
//
//	@Summary		Получение данные авторизованного пользователя
//	@Description	Получение по сесии авторизованного пользователя
//	@Tags			user
//	@Produce		json
//	@Success		200	{object}	model.UserInfo
//	@Failure		400	{object}	sErr
//	@Failure		401	{object}	sErr
//	@Failure		403	{object}	sErr
//	@Failure		404	{object}	sErr
//	@Failure		500	{object}	sErr
//	@Router			/users/info [get]
func (r *RestServer) handlerUserInfo(c *gin.Context) {
	ctx := c.Request.Context()
	us := r.getUserSession(c)
	if us.UserID == 0 {
		c.Status(http.StatusNotFound)
		return
	}

	u, err := r.services.User().GetUserByID(ctx, us.UserID)

	switch {
	case errors.Is(err, errs.ErrUserNotFound):

		c.JSON(http.StatusNotFound, s().SetError(err))
		return
	case err != nil:
		c.JSON(http.StatusInternalServerError, s().SetError(err))
		return
	}

	if !access.CanUser(us.IsAdmin, us.OrgID, us.UserID, u.ID, u.CompanyID) {
		c.Status(http.StatusForbidden)
		return
	}

	c.JSON(http.StatusOK, u)
}
