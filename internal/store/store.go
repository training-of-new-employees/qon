// Package store - пакет для работы с хранилищем.
package store

type Storages interface {
	UserStorage() RepositoryUser
	PositionStorage() RepositoryPosition
	CompanyStorage() RepositoryCompany
	CourseStorage() RepositoryCourse
}
