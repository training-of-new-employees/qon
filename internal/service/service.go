// Package service - пакет, который содержит бизнес-логику приложения.
package service

// Service - общий сервис.
type Service interface {
	User() ServiceUser
	Position() ServicePosition
	Course() ServiceCourse
	Lesson() ServiceLesson
}
