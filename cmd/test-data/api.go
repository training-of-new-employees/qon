package main

import (
	"encoding/json"
	"net/url"
	"strings"

	"github.com/go-resty/resty/v2"

	"github.com/training-of-new-employees/qon/internal/logger"
	"github.com/training-of-new-employees/qon/internal/model"
)

const apiv1 = "/api/v1"

type errResp struct {
	Message string `json:"message"`
}

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

func (a *Api) createAdmins(admins []model.CreateAdmin) ([]model.User, error) {
	reg, _ := url.JoinPath(a.Server, apiv1, "admin/register")
	verify, _ := url.JoinPath(a.Server, apiv1, "admin/verify")
	created := make([]model.User, 0, len(admins))

	for _, adm := range admins {
		r, err := a.r.R().
			SetBody(adm).
			Post(reg)
		if err != nil {
			return nil, err
		}
		msg := errResp{}
		json.Unmarshal(r.Body(), &msg)
		s := strings.Split(msg.Message, "Код верификации: ")
		if len(s) != 2 {
			logger.Log.Sugar().Debugf("can't create admin: %s", msg.Message)
			continue
		}
		code := model.Code{Code: strings.TrimSpace(s[1])}
		r, err = a.r.R().
			SetBody(code).
			Post(verify)
		if err != nil {
			return nil, err
		}
		user := model.User{}
		json.Unmarshal(r.Body(), &user)
		created = append(created, user)
	}
	return created, nil
}
func (a *Api) createUsers(users []model.UserCreate) ([]model.User, error) {
	empl, _ := url.JoinPath(a.Server, apiv1, "admin/employee")
	created := make([]model.User, 0, len(users))

	for _, u := range users {
		r, err := a.r.R().
			SetBody(u).
			Post(empl)
		if err != nil {
			return nil, err
		}
		msg := errResp{}
		json.Unmarshal(r.Body(), &msg)
		s := strings.Split(msg.Message, "Пригласительная cсылка: ")
		if len(s) != 2 {
			logger.Log.Sugar().Debugf("can't create user: %s", msg.Message)
			continue
		}
		link := s[1]
		r, err = a.r.R().
			Get(link)
		if err != nil {
			return nil, err
		}
		user := model.User{}
		json.Unmarshal(r.Body(), &user)
		created = append(created, user)
	}
	return created, nil
}

func (a *Api) createPositions(positions []model.PositionSet) ([]model.Position, error) {
	pos, _ := url.JoinPath(a.Server, apiv1, "positions")
	created := make([]model.Position, 0, len(positions))

	for _, p := range positions {
		r, err := a.r.R().
			SetBody(p).
			Post(pos)
		if err != nil {
			return nil, err
		}
		msg := errResp{}
		json.Unmarshal(r.Body(), &msg)
		if msg.Message != "" {
			logger.Log.Sugar().Debugf("can't create position: %s", msg.Message)
			continue
		}
		if err != nil {
			return nil, err
		}
		pos := model.Position{}
		json.Unmarshal(r.Body(), &pos)
		created = append(created, pos)

	}
	return created, nil
}

func (a *Api) Login(signIn model.UserSignIn) error {
	login, _ := url.JoinPath(a.Server, apiv1, "login")
	resp, err := a.r.R().
		SetBody(signIn).
		Post(login)
	a.r.SetHeader("Authorization", resp.Header().Get("Authorization"))

	return err

}
