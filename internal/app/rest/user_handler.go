package rest

import (
	"errors"

	"github.com/training-of-new-employees/qon/internal/logger"
	"go.uber.org/zap"

	//"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/training-of-new-employees/qon/internal/model"
	//"github.com/training-of-new-employees/qon/internal/model"
	//"net/http"
)

func (r *RestServer) handlerCreateAdminInCache(c *gin.Context) {
	ctx := c.Request.Context()
	createAdmin := model.CreateAdmin{}

	if err := c.ShouldBindJSON(&createAdmin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := createAdmin.Validation(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := createAdmin.ValidatePassword(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	admin, err := r.services.User().WriteAdminToCache(ctx, createAdmin)
	switch {
	case errors.Is(err, model.ErrEmailAlreadyExists):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		logger.Log.Warn("error: %v", zap.Error(err))
		return
	}

	c.JSON(http.StatusCreated, gin.H{"admin": admin.Email})
}

func (r *RestServer) handlerCreateUser(c *gin.Context) {
	ctx := c.Request.Context()
	userReq := model.UserCreate{}

	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := r.services.User().CreateUser(ctx, userReq)
	switch {
	case errors.Is(err, model.ErrEmailAlreadyExists):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		logger.Log.Warn("error", zap.Error(err))
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": user})

}

func (r *RestServer) handlerSignIn(c *gin.Context) {
	ctx := c.Request.Context()
	userReq := model.UserSignIn{}

	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := userReq.Validation(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	user, err := r.services.User().GetUserByEmail(ctx, userReq.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if err = user.CheckPassword(userReq.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	tokens, err := r.services.User().GenerateTokenPair(ctx, user.ID, user.IsAdmin, user.CompanyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Authorization", "Bearer "+tokens.AccessToken)

	c.JSON(http.StatusOK, gin.H{"access token": tokens.AccessToken})
}

func (r *RestServer) handlerAdminEmailVerification(c *gin.Context) {
	ctx := c.Request.Context()
	key := model.Key{}

	if err := c.ShouldBindJSON(&key); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := key.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminFromCache, err := r.services.User().GetAdminFromCache(ctx, key.Key)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	logger.Log.Info("admin from cache: %v", zap.String("email", adminFromCache.Email))

	createdAdmin, err := r.services.User().CreateAdmin(ctx, adminFromCache)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_ = r.services.User().DeleteAdminFromCache(ctx, key.Key)

	c.JSON(http.StatusCreated, gin.H{"admin created": createdAdmin.Email})

}
func (r *RestServer) handlerResetPassword(c *gin.Context) {
	ctx := c.Request.Context()
	email := model.EmailReset{}
	if err := c.ShouldBindJSON(&email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := r.services.User().ResetPassword(ctx, email.Email); err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
