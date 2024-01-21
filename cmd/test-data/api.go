package main

import (
	"errors"
	"net/url"
	"strings"

	"github.com/go-resty/resty/v2"

	"github.com/training-of-new-employees/qon/internal/model"
)

const apiv1 = "/api/v1"

type Api struct {
	Server string
	r      *resty.Client
}

func NewApi(c *cfg) *Api {
	return &Api{
		Server: c.Address,
		r:      resty.New(),
	}
}

func (a *Api) createAdmins(admins []model.CreateAdmin) error {
	reg, err := url.JoinPath(a.Server, apiv1, "admin/register")
	if err != nil {
		return err
	}
	verify, err := url.JoinPath(a.Server, apiv1, "admin/verify")
	for _, adm := range admins {
		r, err := a.r.R().
			SetBody(adm).
			Post(reg)
		if err != nil {
			return err
		}
		s := strings.Split(r.String(), "Код верификации: ")
		if len(s) != 2 {
			return errors.New("Invalid verification")
		}
		code := model.Code{Code: s[1]}
		r, err = a.r.R().
			SetBody(code).
			Post(verify)
		if err != nil {
			return err
		}

	}
	return nil
}
