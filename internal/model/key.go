package model

import (
	"errors"
)

type Key struct {
	Key string `json:"key"`
}

func (k *Key) Validate() error {

	if k.Key == "" {
		return errors.New("error validation Key empty values")
	}

	return nil
}
