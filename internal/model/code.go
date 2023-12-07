package model

import (
	"errors"
)

type Code struct {
	Code string `json:"code"`
}

func (c *Code) Validate() error {

	if c.Code == "" {
		return errors.New("error validation Code empty values")
	}

	return nil
}
