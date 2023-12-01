package pg

import (
	"database/sql"
	"log"
	"strings"

	"github.com/integralist/go-findroot/find"
	"github.com/pressly/goose/v3"
)

var migrationPath = getMigrationPath()

// MigrationsUp запускает миграции.
func MigrationsUp(db *sql.DB) error {
	if err := goose.SetDialect("pgx"); err != nil {
		return err
	}

	return goose.Up(db, migrationPath)
}

// MigrationsDown откатывает миграции.
func MigrationsDown(db *sql.DB) error {
	if err := goose.SetDialect("pgx"); err != nil {
		return err
	}

	return goose.Down(db, migrationPath)
}

// getMigrationPath возвращает путь к миграциям.
func getMigrationPath() string {
	// Путь к миграциям по умолчанию
	var defaultMigrationPath = "/migrations"

	// получить получить путь к корню
	rep, err := find.Repo()
	if err != nil {
		log.Printf("cannot get root dir: %s", err.Error())
		log.Println("use default path to migration")
		return defaultMigrationPath
	}

	path := strings.Join([]string{rep.Path, "migrations"}, "/")

	if strings.EqualFold(rep.Path, "./") || strings.EqualFold(rep.Path, "/") {
		path = strings.Join([]string{rep.Path, "migrations"}, "")
	}

	return path
}
