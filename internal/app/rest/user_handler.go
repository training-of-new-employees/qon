package rest

import (
	//"errors"
	"github.com/gin-gonic/gin"
	//"github.com/training-of-new-employees/qon/internal/model"
	//"net/http"
)

func (r *RestServer) handlerRegisterAdmin(c *gin.Context) {
	//ctx := c.Request.Context()
	//createAdmin := model.CreateAdmin{}
	//
	//if err := c.ShouldBindJSON(&createAdmin); err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	return
	//}
	//
	//if err := createAdmin.Validation(); err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	return
	//}
	//
	//if err := createAdmin.ValidatePassword(); err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	return
	//}
	//
	//_, err := u.registerAdmin(ctx, createAdmin)
	//switch {
	//case errors.Is(err, model.ErrEmailAlreadyExists):
	//	c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	//	return
	//case err != nil:
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//	u.logger.Warn("error creating admin: %v", err)
	//
	//	return
}
