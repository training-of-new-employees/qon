package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"

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
	regURL, _ := url.JoinPath(a.Server, apiv1, "admin/register")
	verifyURL, _ := url.JoinPath(a.Server, apiv1, "admin/verify")
	created := make([]model.User, 0, len(admins))

	for _, adm := range admins {
		r, err := a.r.R().
			SetBody(adm).
			Post(regURL)
		if err != nil {
			return nil, err
		}
		msg := errResp{}
		err = json.Unmarshal(r.Body(), &msg)
		_ = err
		s := strings.Split(msg.Message, "Код верификации: ")
		if len(s) != 2 {
			logger.Log.Sugar().Debugf("can't create admin: %s\n\t%v", msg.Message, adm)
			continue
		}
		code := model.Code{Code: strings.TrimSpace(s[1])}
		_, err = a.r.R().
			SetBody(code).
			Post(verifyURL)
		if err != nil {
			return nil, err
		}
	}
	err := a.login(
		model.UserSignIn{
			Email:    admins[0].Email,
			Password: admins[0].Password,
		},
	)
	if err != nil {
		return nil, err
	}
	u, err := a.getInfo()
	if err != nil {
		return nil, err
	}
	created = append(created, u)

	return created, nil
}
func (a *Api) createUsers(users []model.UserCreate) ([]model.User, error) {
	employeeURL, _ := url.JoinPath(a.Server, apiv1, "admin/employee")
	created := make([]model.User, 0, len(users))

	for _, u := range users {
		r, err := a.r.R().
			SetBody(u).
			Post(employeeURL)
		if err != nil {
			return nil, err
		}
		msg := errResp{}
		err = json.Unmarshal(r.Body(), &msg)
		_ = err
		s := strings.Split(msg.Message, "Пригласительная cсылка: ")
		if len(s) != 2 {
			logger.Log.Sugar().Debugf("can't create user: %s\n\t%v", msg.Message, u)
			continue
		}
		link := s[1]
		r, err = a.r.R().
			Get(link)
		if err != nil {
			return nil, err
		}
		user := model.User{}
		err = json.Unmarshal(r.Body(), &user)
		_ = err
		created = append(created, user)
	}
	return created, nil
}

func (a *Api) createPositions(positions []model.PositionSet) ([]model.Position, error) {
	positionsURL, _ := url.JoinPath(a.Server, apiv1, "positions")
	created := make([]model.Position, 0, len(positions))

	for _, p := range positions {
		r, err := a.r.R().
			SetBody(p).
			Post(positionsURL)
		if err != nil {
			return nil, err
		}
		msg := errResp{}
		err = json.Unmarshal(r.Body(), &msg)
		_ = err
		if msg.Message != "" {
			logger.Log.Sugar().Debugf("can't create position: %s\n\t%v", msg.Message, p)
			continue
		}
		if err != nil {
			return nil, err
		}
		pos := model.Position{}
		err = json.Unmarshal(r.Body(), &pos)
		_ = err
		created = append(created, pos)
	}
	return created, nil
}

func (a *Api) createCourses(courses []model.CourseSet) ([]model.Course, error) {
	coursesURL, _ := url.JoinPath(a.Server, apiv1, "admin/courses")
	created := make([]model.Course, 0, len(courses))
	for _, c := range courses {
		r, err := a.r.R().
			SetBody(c).
			Post(coursesURL)
		if err != nil {
			return nil, err
		}
		logger.Log.Sugar().Debugf("Course created:%s", r.String())
		msg := errResp{}
		err = json.Unmarshal(r.Body(), &msg)
		_ = err
		if msg.Message != "" {
			logger.Log.Sugar().Debugf("can't create course: %s\n\t%v", msg.Message, c)
			continue
		}
		course := model.Course{}
		err = json.Unmarshal(r.Body(), &course)
		_ = err

		created = append(created, course)
	}

	return created, nil
}

func (a *Api) assignCourses(coursesID []int, posID int) error {
	assignURL, _ := url.JoinPath(a.Server, apiv1, fmt.Sprintf("positions/%d/courses", posID))
	logger.Log.Sugar().Debugf("Courses ID: %v", coursesID)
	courses := model.PositionAssignCourses{CoursesID: coursesID}
	r, err := a.r.R().
		SetBody(courses).
		Patch(assignURL)
	if err != nil {
		return err
	}
	msg := errResp{}
	err = json.Unmarshal(r.Body(), &msg)
	_ = err
	if msg.Message != "" {
		logger.Log.Sugar().Debugf("can't assign course: %s", msg.Message)
	}
	return nil
}

func (a *Api) createLessons(lessons []model.Lesson, courseID int) ([]model.Lesson, error) {
	lessonsURL, _ := url.JoinPath(a.Server, apiv1, "lessons")
	created := make([]model.Lesson, 0, len(lessons))
	for _, l := range lessons {
		l.CourseID = courseID
		r, err := a.r.R().
			SetBody(l).
			Post(lessonsURL)
		if err != nil {
			return nil, err
		}
		msg := errResp{}
		err = json.Unmarshal(r.Body(), &msg)
		if err == nil && msg.Message != "" {
			logger.Log.Sugar().Debugf("can't create lesson: %s\n\t%v", msg.Message, l)
		}
		lesson := model.Lesson{}
		err = json.Unmarshal(r.Body(), &lesson)
		_ = err

		created = append(created, lesson)
	}

	return created, nil
}

func (a *Api) login(signIn model.UserSignIn) error {
	loginURL, _ := url.JoinPath(a.Server, apiv1, "login")
	resp, err := a.r.R().
		SetBody(signIn).
		Post(loginURL)
	logger.Log.Info("Login", zap.String("msg", resp.String()))
	a.r.SetHeader("Authorization", resp.Header().Get("Authorization"))

	return err

}

func (a *Api) getInfo() (model.User, error) {
	infoURL, _ := url.JoinPath(a.Server, apiv1, "users/info")
	resp, err := a.r.R().
		Get(infoURL)
	logger.Log.Info("Get user info", zap.String("msg", resp.String()))
	_ = err
	var u model.User
	err = json.Unmarshal(resp.Body(), &u)

	return u, err

}
