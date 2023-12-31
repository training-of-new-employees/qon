package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SetPassword(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		user := CreateAdmin{Password: "password"}
		err := user.SetPassword()
		assert.NoError(t, err)
		assert.NotEmpty(t, user.Password)
	})
}

func Test_CheckPassword(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		userResp := User{}
		userReq := CreateAdmin{Password: "password"}
		err := userReq.SetPassword()
		assert.NoError(t, err)
		assert.NotEmpty(t, userReq.Password)
		userResp.Password = userReq.Password

		err = userResp.CheckPassword("password")
		assert.NoError(t, err)
	})
}
