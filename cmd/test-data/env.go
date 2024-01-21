package main

import (
	"encoding/json"
	"os"
	"path"

	"go.uber.org/zap"

	"github.com/training-of-new-employees/qon/internal/logger"
	"github.com/training-of-new-employees/qon/internal/model"
)

func upTestEnv(c *cfg) error {
	api := NewApi(c)

	admFN := path.Join(c.Env, c.Admins)
	logger.Log.Debug("Fetch admins", zap.String("file", admFN))
	admins, err := fetchAdmins(admFN)
	if err != nil {
		return err
	}
	err = api.createAdmins(admins)
	if err != nil {
		return err
	}
	logger.Log.Info("Admins created")
	usersFN := path.Join(c.Env, c.Users)
	logger.Log.Debug("Fetch users", zap.String("file", usersFN))
	posFN := path.Join(c.Env, c.Positions)
	logger.Log.Debug("Fetch positions", zap.String("file", posFN))
	crsFN := path.Join(c.Env, c.Courses)
	logger.Log.Debug("Fetch courses", zap.String("file", crsFN))
	lesFN := path.Join(c.Env, c.Lessons)
	logger.Log.Debug("Fetch lessons", zap.String("file", lesFN))
	return nil

}

func fetchAdmins(n string) ([]model.CreateAdmin, error) {
	f, err := os.Open(n)
	if err != nil {
		return nil, err
	}
	d := json.NewDecoder(f)
	admins := []model.CreateAdmin{}
	err = d.Decode(&admins)
	return admins, err
}
