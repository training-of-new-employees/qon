package pg

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/suite"
)

// Тестовый адрес БД.
var testDatabaseDSN = getTestDsn()

type storeTestSuite struct {
	suite.Suite
	store *Store
}

// SetupSuite - запуск до начала выполнения набора тестов.
func (suite *storeTestSuite) SetupSuite() {
	// ожидание готовности контейнера с БД
	time.Sleep(2 * time.Second)

	// тестирование будет пропущено если БД недоступно
	if err := suite.isDBAvailable(testDatabaseDSN); err != nil {
		suite.T().Skipf("skip db tests: database is not available: %v", err)
		return
	}
}

// TearDownSuite - запуск после выполнения всех тестов.
func (suite *storeTestSuite) TearDownSuite() {
}

// SetupTest - выполнение перед каждым тест-кейсом.
func (suite *storeTestSuite) SetupTest() {
	// очищаем БД
	err := suite.clearDB(testDatabaseDSN)
	suite.NoError(err)

	suite.store, err = NewStore(testDatabaseDSN)

	suite.NoError(err)
}

// TearDownTest - запуск после каждого тест-кейса.
func (suite *storeTestSuite) TearDownTest() {
	suite.store.Close()
}

// getTestDsn - получение адреса тестовой БД.
func getTestDsn() string {
	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		dsn = "postgres://qon:qon@localhost:15439/qon_test"

		log.Printf("env TEST_DB_DSN is empty, used default value: %s", dsn)
	}
	return dsn
}

// isDBAvailable - проверка соединения с БД.
func (suite *storeTestSuite) isDBAvailable(dsn string) error {
	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return err
	}

	defer db.Close()

	return db.Ping()
}

// clearDB - откат миграций и завершение соединения с БД.
func (suite *storeTestSuite) clearDB(dsn string) error {
	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return err
	}

	defer db.Close()

	if err = goose.SetDialect("pgx"); err != nil {
		return err
	}
	ver, err := goose.GetDBVersion(db.DB)
	if err != nil {
		return err
	}
	if ver == 0 {
		return nil
	}

	return MigrationsDown(db)
}

// TestStoreTestSuite - точка входа для тестирования store.
func TestStoreTestSuite(t *testing.T) {
	suite.Run(t, new(storeTestSuite))
}
