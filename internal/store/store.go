// Package store - пакет для работы с хранилищем.
package store

// Storages - интерфейс слоя хранилище, состоящий из репозиториев.
type Storages interface {
	// UserStorage - интерфейс репозитория пользователя
	UserStorage() RepositoryUser
	// PositionStorage - интерфейс репозитория должности
	PositionStorage() RepositoryPosition
	// CompanyStorage - интерфейс репозитория компании/организации
	CompanyStorage() RepositoryCompany
	// CourseStorage - интерфейс репозитория курсов
	CourseStorage() RepositoryCourse
	// LessonStorage - интерфейс репозитория уроков
	LessonStorage() RepositoryLesson
}
