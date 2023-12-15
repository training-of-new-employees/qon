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
	"github.com/training-of-new-employees/qon/internal/pkg/jwttoken"
)

// CreateAdmin godoc
//
//	@Summary	Создание администратора
//	@Tags		user
//	@Produce	json
//	@Param		object	body		model.CreateAdmin	true	"Create Admin"
//	@Success	201		{array}		model.CreateAdmin
//	@Failure	400		{object}	sErr
//	@Failure	409		{object}	sErr
//	@Failure	500		{object}	sErr
//	@Router		/admin/register [post]
func (r *RestServer) handlerCreateAdminInCache(c *gin.Context) {
	ctx := c.Request.Context()
	createAdmin := model.CreateAdmin{}

	if err := c.ShouldBindJSON(&createAdmin); err != nil {
		c.JSON(http.StatusBadRequest, s().SetError(err))
		return
	}

	if err := createAdmin.Validation(); err != nil {
		c.JSON(http.StatusBadRequest, s().SetError(err))
		return
	}

	if err := createAdmin.ValidatePassword(); err != nil {
		c.JSON(http.StatusBadRequest, s().SetError(err))
		return
	}

	admin, err := r.services.User().WriteAdminToCache(ctx, createAdmin)
	switch {
	case errors.Is(err, model.ErrEmailAlreadyExists):
		c.JSON(http.StatusConflict, s().SetError(err))
		return
	case err != nil:
		c.JSON(http.StatusInternalServerError, s().SetError(err))
		logger.Log.Warn("error: %v", zap.Error(err))
		return
	}

	c.JSON(http.StatusCreated, admin)
}

// CreateUser godoc
//
//	@Summary	Создание пользователя
//	@Tags		user
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
		c.JSON(http.StatusBadRequest, s().SetError(err))
		return
	}

	if err := userReq.Validation(); err != nil {
		c.JSON(http.StatusBadRequest, s().SetError(err))
		return
	}

	user, err := r.services.User().CreateUser(ctx, userReq)
	switch {
	case errors.Is(err, model.ErrEmailAlreadyExists):
		c.JSON(http.StatusConflict, s().SetError(err))
		return
	case err != nil:
		c.JSON(http.StatusInternalServerError, s().SetError(err))
		logger.Log.Warn("error", zap.Error(err))
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
	us := r.getUserSession(c)
	if !us.IsAdmin {
		logger.Log.Warn("Not admin user try to get info about users", zap.Int("id", us.UserID))
		c.JSON(http.StatusForbidden, s().SetError(fmt.Errorf("you can't get users info: %w", errs.ErrOnlyAdmin)))
		return
	}
	users, err := r.services.User().GetUsersByCompany(ctx, us.OrgID)
	switch {
	case errors.Is(err, errs.ErrUserNotFound):
		c.JSON(http.StatusNotFound, s().SetError(err))
		return
	case err != nil:
		c.JSON(http.StatusInternalServerError, s().SetError(err))
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
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		msg := fmt.Errorf("got invalid user id: %s", idParam)
		logger.Log.Warn("error", zap.Error(msg))
		c.JSON(http.StatusBadRequest, s().SetError(msg))
		return
	}
	us := r.getUserSession(c)
	u, err := r.services.User().GetUserByID(ctx, id)
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
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		msg := fmt.Errorf("got invalid user id: %s", idParam)
		logger.Log.Warn("error", zap.Error(msg))
		c.JSON(http.StatusBadRequest, s().SetError(msg))
		return
	}
	edit := &model.UserEdit{}
	if err := c.ShouldBindJSON(&edit); err != nil {
		c.JSON(http.StatusBadRequest, s().SetError(err))
		return
	}
	edit.ID = id
	us := r.getUserSession(c)
	if !access.CanUser(us.IsAdmin, us.OrgID, us.UserID, edit.ID, us.OrgID) {
		logger.Log.Warn("User try edit info without rights", zap.Int("id", us.UserID), zap.Int("edited", edit.ID))
		c.JSON(http.StatusForbidden, s().SetError(fmt.Errorf("you can't edit user")))
		return
	}
	edited, err := r.services.User().EditUser(ctx, edit, us.OrgID)
	switch {
	case errors.Is(err, errs.ErrUserNotFound):
		c.JSON(http.StatusNotFound, s().SetError(err))
		return
	case err != nil:
		c.JSON(http.StatusInternalServerError, s().SetError(err))
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
	if err != nil {
		c.JSON(http.StatusNotFound, s().SetError(err))
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

// SignIn godoc
//
//	@Summary	Вход пользователя
//	@Tags		user
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
//	@Tags		user
//	@Produce	json
//	@Param		object	body		model.Code	true	"User Email Verification"
//	@Success	201		{object}	sEmail
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

	createdAdmin, err := r.services.User().CreateAdmin(ctx, adminFromCache)
	if err != nil {
		c.JSON(http.StatusInternalServerError, s().SetError(err))
		return
	}

	_ = r.services.User().DeleteAdminFromCache(ctx, code.Code)

	c.JSON(http.StatusCreated, s().SetEmail(createdAdmin.Email))

}

// ResetPassword godoc
//
//	@Summary	Сброс пароля пользователя
//	@Tags		user
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
		c.JSON(http.StatusBadRequest, s().SetError(err))
		return
	}

	if err := r.services.User().ResetPassword(ctx, email.Email); err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, s().SetError(err))
			return
		}
		c.JSON(http.StatusInternalServerError, s().SetError(err))
		return
	}

	c.JSON(http.StatusOK, s().SetEmail(email.Email))
}

// AdminEdit godoc
//
//	@Summary	Изменение данных администратора
//	@Tags		user
//	@Produce	json
//	@Param		object	body		model.AdminEdit	true	"Admin Edit"
//	@Success	200		{object}	model.AdminEdit
//	@Failure	400		{object}	sErr
//	@Failure	401		{object}	sErr
//	@Failure	404		{object}	sErr
//	@Failure	500		{object}	sErr
//	@Router		/admin/info [post]
func (r *RestServer) handlerAdminEdit(c *gin.Context) {
	ctx := c.Request.Context()
	edit := &model.AdminEdit{}
	if err := c.ShouldBindJSON(&edit); err != nil {
		c.JSON(http.StatusBadRequest, s().SetError(err))
		return
	}

	token := jwttoken.GetToken(c)
	claims, err := r.tokenVal.ValidateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, s().SetError(err))
		return
	}
	edit.ID = claims.UserID

	edited, err := r.services.User().EditAdmin(ctx, edit)
	switch {
	case errors.Is(err, model.ErrUserNotFound):
		c.JSON(http.StatusNotFound, s().SetError(err))
	case err != nil:
		c.JSON(http.StatusInternalServerError, s().SetError(err))
		return
	}
	c.JSON(http.StatusOK, edited)
}
