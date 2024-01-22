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

	// создаём все должности для самой первой компании
	pos, err := positions(c, api, a[0].CompanyID)
	if err != nil {
		return err
	}

	// назначаем всех пользователей на самую первую должность
	_, err = users(c, api, pos[0].ID, pos[0].CompanyID)
	if err != nil {
		return err
	}

	crs, err := courses(c, api, pos[0].ID, pos[0].CompanyID)
	if err != nil {
		return err
	}
	_, err = lessons(c, api, crs[0].ID)

	return err

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
	err = login(api, admins[0])
	return a, err
}

func login(api *Api, admin model.CreateAdmin) error {
	signIn := model.UserSignIn{
		Email:    admin.Email,
		Password: admin.Password,
	}
	return api.login(signIn)
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
	logger.Log.Info("Positions created")
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

func courses(c *cfg, api *Api, posID, companyID int) ([]model.Course, error) {
	crsFN := path.Join(c.Env, c.Courses)
	logger.Log.Debug("Fetch courses", zap.String("file", crsFN))
	courses, err := fetchEntities[model.CourseSet](crsFN)
	if err != nil {
		return nil, err
	}
	crs, err := api.createCourses(courses)
	if err != nil {
		return nil, err
	}
	logger.Log.Info("Courses created")
	coursesID := make([]int, 0, len(courses))
	for _, c := range crs {
		coursesID = append(coursesID, c.ID)
	}
	err = api.assignCourses(coursesID, posID)
	if err != nil {
		return nil, err
	}
	logger.Log.Info("Courses assigned")

	return crs, nil

}

func lessons(c *cfg, api *Api, courseID int) ([]model.Lesson, error) {
	lesFN := path.Join(c.Env, c.Lessons)
	logger.Log.Debug("Fetch lessons", zap.String("file", lesFN))
	lessons, err := fetchEntities[model.Lesson](lesFN)
	if err != nil {
		return nil, err
	}
	created, err := api.createLessons(lessons, courseID)
	logger.Log.Info("Lessons created")
	return created, err
}
