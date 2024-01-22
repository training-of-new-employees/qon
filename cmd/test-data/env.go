package main

import (
	"encoding/json"
	"errors"
	"os"
	"path"

	"go.uber.org/zap"

	"github.com/training-of-new-employees/qon/internal/logger"
	"github.com/training-of-new-employees/qon/internal/model"
)

func upTestEnv(c *cfg) error {
	api := NewApi(c)
	a, err := admins(c, api)
	if err != nil {
		return err
	}

	err = login(api, a[0])
	if err != nil {
		return err
	}
	// создаём все должности для самой первой компании
	pos, err := positions(c, api, a[0].CompanyID)
	if err != nil {
		return err
	}

	// назначаем всех пользователей на самую первую должность
	u, err := users(c, api, pos[0].ID, pos[0].CompanyID)
	if err != nil {
		return err
	}
	_ = u

	crsFN := path.Join(c.Env, c.Courses)
	logger.Log.Debug("Fetch courses", zap.String("file", crsFN))
	lesFN := path.Join(c.Env, c.Lessons)
	logger.Log.Debug("Fetch lessons", zap.String("file", lesFN))
	return nil

}

func fetchEntities[T any](n string) ([]T, error) {
	f, err := os.Open(n)
	if err != nil {
		return nil, err
	}
	d := json.NewDecoder(f)
	entities := []T{}
	err = d.Decode(&entities)
	return entities, err
}

func admins(c *cfg, api *Api) ([]model.User, error) {
	admFN := path.Join(c.Env, c.Admins)
	logger.Log.Debug("Fetch admins", zap.String("file", admFN))
	admins, err := fetchEntities[model.CreateAdmin](admFN)
	if err != nil {
		return nil, err
	}
	if len(admins) == 0 {
		return nil, errors.New("admins not found")
	}
	a, err := api.createAdmins(admins)
	if err != nil {
		return nil, err
	}
	logger.Log.Info("Admins created")
	return a, nil
}

func login(api *Api, admin model.User) error {
	signIn := model.UserSignIn{
		Email:    admin.Email,
		Password: admin.Password,
	}
	return api.Login(signIn)
}

func positions(c *cfg, api *Api, companyID int) ([]model.Position, error) {
	posFN := path.Join(c.Env, c.Positions)
	logger.Log.Debug("Fetch positions", zap.String("file", posFN))
	positions, err := fetchEntities[model.PositionSet](posFN)
	if err != nil {
		return nil, err
	}
	sumPos := make([]model.Position, 0, len(positions))
	for i := range positions {
		positions[i].CompanyID = companyID
	}
	p, err := api.createPositions(positions)
	if err != nil {
		return nil, err
	}
	sumPos = append(sumPos, p...)
	return sumPos, nil
}

func users(c *cfg, api *Api, posID, companyID int) ([]model.User, error) {
	usersFN := path.Join(c.Env, c.Users)
	logger.Log.Debug("Fetch users", zap.String("file", usersFN))
	users, err := fetchEntities[model.UserCreate](usersFN)
	if err != nil {
		return nil, err
	}
	for i := range users {
		users[i].CompanyID = companyID
		users[i].PositionID = posID
	}
	u, err := api.createUsers(users)
	if err != nil {
		return nil, err
	}
	logger.Log.Info("Users created")
	return u, nil
}
