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

func TestGeneratePassword(t *testing.T) {
	testCases := []struct {
		name    string
		isValid bool
	}{
		{
			name:    "valid password",
			isValid: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			password := GeneratePassword()
			err := validatePassword(password)(password)

			if testCase.isValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
