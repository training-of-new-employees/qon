package pg

import (
	"context"

	"github.com/stretchr/testify/require"

	"github.com/training-of-new-employees/qon/internal/errs"
	"github.com/training-of-new-employees/qon/internal/model"
)

var testCoursesLen = 10

func (suite *storeTestSuite) TestCreateCourse() {
	suite.NotNil(suite.store)

	ca := model.NewTestCreateAdmin()
	uc := model.NewTestUserCreate()
	admin, err := suite.store.UserStorage().CreateAdmin(context.TODO(), uc, ca.Company)
	suite.NoError(err)

	testCases := []struct {
		name    string
		prepare func() model.CourseSet
		err     error
	}{
		{
			name: "success",
			prepare: func() model.CourseSet {
				course := model.NewTestCourseSet()
				course.CreatedBy = admin.ID
				return course
			},
			err: nil,
		},
		{
			name: "empty name",
			prepare: func() model.CourseSet {
				course := model.NewTestCourseSet()
				course.CreatedBy = admin.ID
				course.Name = ""
				return course

			},
			err: errs.ErrCourseNameIsEmpty,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			_, err := suite.store.CourseStorage().CreateCourse(context.TODO(), tc.prepare())
			suite.Equal(tc.err, err)
		})
	}
}

func (suite *storeTestSuite) TestUserCourses() {
	suite.NotNil(suite.store)

	ca1 := model.NewTestCreateAdmin()
	uc1 := model.NewTestUserCreate()
	admin, err := suite.store.UserStorage().CreateAdmin(context.TODO(), uc1, ca1.Company)
	suite.NoError(err)

	p := model.NewTestPositionSet()
	p.CompanyID = admin.CompanyID
	pos, err := suite.store.PositionStorage().CreatePosition(context.TODO(), p)
	suite.NoError(err)

	uc2 := model.NewTestUserCreate()
	uc2.CompanyID = admin.CompanyID
	uc2.PositionID = pos.ID
	user, err := suite.store.UserStorage().CreateUser(context.TODO(), uc2)
	suite.NoError(err)

	courses := make([]model.Course, 0, testCoursesLen)

	for i := 0; i < testCoursesLen; i++ {
		c := model.NewTestCourseSet()
		c.ID = i + 1
		c.CreatedBy = admin.ID
		created, err := suite.store.CourseStorage().CreateCourse(context.TODO(), c)
		suite.NoError(err)
		err = suite.store.PositionStorage().AssignCourse(context.TODO(), pos.ID, created.ID)
		suite.NoError(err)
		courses = append(courses, *created)
	}

	testCases := []struct {
		name    string
		uid     int
		err     error
		courses []model.Course
	}{
		{
			name:    "success",
			uid:     user.ID,
			err:     nil,
			courses: courses,
		},
		{
			name:    "not found",
			uid:     admin.ID,
			err:     nil,
			courses: []model.Course{},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			courses, err := suite.store.CourseStorage().UserCourses(context.TODO(), tc.uid)
			suite.Equal(tc.err, err)
			if err == nil {
				suite.Equal(tc.courses, courses)
			}
		})
	}
}
func (suite *storeTestSuite) TestCompanyCourses() {
	suite.NotNil(suite.store)

	ca1 := model.NewTestCreateAdmin()
	uc1 := model.NewTestUserCreate()
	admin, err := suite.store.UserStorage().CreateAdmin(context.TODO(), uc1, ca1.Company)
	suite.NoError(err)

	ca2 := model.NewTestCreateAdmin()
	uc2 := model.NewTestUserCreate()
	admin2, err := suite.store.UserStorage().CreateAdmin(context.TODO(), uc2, ca2.Company)
	suite.NoError(err)

	courses := make([]model.Course, 0, testCoursesLen)

	for i := 0; i < testCoursesLen; i++ {
		c := model.NewTestCourseSet()
		c.ID = i + 1
		c.CreatedBy = admin.ID
		created, err := suite.store.CourseStorage().CreateCourse(context.TODO(), c)
		suite.NoError(err)
		courses = append(courses, *created)
	}

	testCases := []struct {
		name    string
		cid     int
		err     error
		courses []model.Course
	}{
		{
			name:    "success",
			cid:     admin.CompanyID,
			err:     nil,
			courses: courses,
		},
		{
			name:    "not found",
			cid:     admin2.CompanyID,
			err:     nil,
			courses: []model.Course{},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			courses, err := suite.store.CourseStorage().CompanyCourses(context.TODO(), tc.cid)
			suite.Equal(tc.err, err)
			if err == nil {
				suite.Equal(tc.courses, courses)
			}
		})
	}
}

func (suite *storeTestSuite) TestEditCourse() {
	suite.NotNil(suite.store)

	ca1 := model.NewTestCreateAdmin()
	uc1 := model.NewTestUserCreate()
	admin, err := suite.store.UserStorage().CreateAdmin(context.TODO(), uc1, ca1.Company)
	suite.NoError(err)

	c := model.NewTestCourseSet()
	c.ID = 1
	c.CreatedBy = admin.ID
	created, err := suite.store.CourseStorage().CreateCourse(context.TODO(), c)
	suite.NoError(err)
	testCases := []struct {
		name    string
		prepare func() model.CourseSet
		err     error
	}{
		{
			name: "edit full",
			prepare: func() model.CourseSet {
				course := model.NewTestCourseSet()
				course.ID = created.ID
				return course
			},
			err: nil,
		},
		{
			name: "edit name",
			prepare: func() model.CourseSet {
				course := model.NewTestCourseSet()
				course.ID = created.ID
				course.Description = ""
				course.IsArchived = false

				return course
			},
			err: nil,
		},
		{
			name: "archive course",
			prepare: func() model.CourseSet {
				course := model.NewTestCourseSet()
				course.ID = created.ID
				course.Name = ""
				course.Description = ""
				course.IsArchived = true

				return course
			},
			err: nil,
		},
		{
			name: "change nothing name",
			prepare: func() model.CourseSet {
				course := model.NewTestCourseSet()
				course.ID = created.ID
				course.Name = created.Name
				course.Description = ""
				course.IsArchived = true

				return course

			},
			err: nil,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			prepared := tc.prepare()
			c, err := suite.store.CourseStorage().EditCourse(context.TODO(), prepared, admin.CompanyID)
			suite.Equal(tc.err, err)
			if prepared.Name != "" {
				suite.Equal(c.Name, prepared.Name)
			}
			if prepared.Description != "" {
				suite.Equal(c.Description, prepared.Description)
			}
			suite.Equal(c.IsArchived, prepared.IsArchived)
		})
	}
}

func (suite *storeTestSuite) GetUserCoursesStatus() {
	suite.NotNil(suite.store)

	ca1 := model.NewTestCreateAdmin()
	uc1 := model.NewTestUserCreate()
	admin, err := suite.store.UserStorage().CreateAdmin(context.TODO(), uc1, ca1.Company)
	suite.NoError(err)

	p := model.NewTestPositionSet()
	p.CompanyID = admin.CompanyID
	pos, err := suite.store.PositionStorage().CreatePosition(context.TODO(), p)
	suite.NoError(err)

	uc2 := model.NewTestUserCreate()
	uc2.CompanyID = admin.CompanyID
	uc2.PositionID = pos.ID
	user, err := suite.store.UserStorage().CreateUser(context.TODO(), uc2)
	suite.NoError(err)

	for i := 0; i < testCoursesLen; i++ {
		c := model.NewTestCourseSet()
		c.ID = i + 1
		c.CreatedBy = admin.ID
		created, err := suite.store.CourseStorage().CreateCourse(context.TODO(), c)
		suite.NoError(err)
		err = suite.store.PositionStorage().AssignCourse(context.TODO(), pos.ID, created.ID)
		suite.NoError(err)
	}

	lesson, err := suite.store.LessonStorage().CreateLesson(
		context.TODO(),
		model.Lesson{CourseID: 4, Name: "Test Lesson", Content: "Test content", URLPicture: "Test picture"},
		user.ID,
	)
	suite.NoError(err)
	err = suite.store.LessonStorage().UpdateUserLessonStatus(context.TODO(), user.ID, 4, lesson.ID, "done")
	suite.NoError(err)

	testCases := []struct {
		name      string
		uid       int
		err       error
		courseIDs []int
		statuses  map[int]string
	}{
		{
			name:      "success",
			uid:       user.ID,
			err:       nil,
			courseIDs: []int{1, 2, 4},
			statuses:  map[int]string{1: "not-started", 2: "not-started", 4: "done"},
		},
		{
			name:      "success (len of ids is 0)",
			uid:       user.ID,
			err:       nil,
			courseIDs: []int{},
			statuses:  map[int]string{},
		},
		{
			name: "internal error",
			uid:  admin.ID,
			err:  errs.ErrInternal,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			statuses, err := suite.store.CourseStorage().GetUserCoursesStatus(context.TODO(), tc.uid, tc.courseIDs)
			suite.Equal(tc.err, err)
			if err == nil {
				for id := range statuses {
					suite.Equal(statuses[id], tc.statuses[id])
				}
			}
		})
	}
}

func (suite *storeTestSuite) TestGetUserCourse() {
	suite.NotNil(suite.store)

	ca1 := model.NewTestCreateAdmin()
	uc1 := model.NewTestUserCreate()
	admin, err := suite.store.UserStorage().CreateAdmin(context.TODO(), uc1, ca1.Company)
	suite.NoError(err)

	p := model.NewTestPositionSet()
	p.CompanyID = admin.CompanyID
	pos, err := suite.store.PositionStorage().CreatePosition(context.TODO(), p)
	suite.NoError(err)

	uc2 := model.NewTestUserCreate()
	uc2.CompanyID = admin.CompanyID
	uc2.PositionID = pos.ID
	user, err := suite.store.UserStorage().CreateUser(context.TODO(), uc2)
	suite.NoError(err)

	for i := 0; i < testCoursesLen; i++ {
		c := model.NewTestCourseSet()
		c.ID = i + 1
		c.CreatedBy = admin.ID
		created, err := suite.store.CourseStorage().CreateCourse(context.TODO(), c)
		suite.NoError(err)
		err = suite.store.PositionStorage().AssignCourse(context.TODO(), pos.ID, created.ID)
		suite.NoError(err)
	}

	testCases := []struct {
		name     string
		uid      int
		err      error
		courseID int
	}{
		{
			name:     "success",
			uid:      user.ID,
			err:      nil,
			courseID: 1,
		},
		{
			name:     "not found",
			uid:      admin.ID,
			err:      errs.ErrNotFound,
			courseID: testCoursesLen + 2,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			course, err := suite.store.CourseStorage().GetUserCourse(context.TODO(), tc.courseID, tc.uid)
			suite.Equal(tc.err, err)
			if err == nil {
				suite.Equal(tc.courseID, course.ID)
			}
		})
	}
}

func (suite *storeTestSuite) TestGetCompanyCourse() {
	suite.NotNil(suite.store)

	ca1 := model.NewTestCreateAdmin()
	uc1 := model.NewTestUserCreate()
	admin, err := suite.store.UserStorage().CreateAdmin(context.TODO(), uc1, ca1.Company)
	suite.NoError(err)

	for i := 0; i < testCoursesLen; i++ {
		c := model.NewTestCourseSet()
		c.ID = i + 1
		c.CreatedBy = admin.ID
		_, err := suite.store.CourseStorage().CreateCourse(context.TODO(), c)
		suite.NoError(err)
	}

	testCases := []struct {
		name     string
		cid      int
		err      error
		courseID int
	}{
		{
			name:     "success",
			cid:      admin.CompanyID,
			err:      nil,
			courseID: 1,
		},
		{
			name:     "not found",
			cid:      admin.CompanyID,
			err:      errs.ErrNotFound,
			courseID: testCoursesLen + 2,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			course, err := suite.store.CourseStorage().CompanyCourse(context.TODO(), tc.courseID, tc.cid)
			require.Equal(suite.T(), tc.err, err)
			if tc.err == nil {
				require.Equal(suite.T(), tc.courseID, course.ID)
			}
		})
	}
}
