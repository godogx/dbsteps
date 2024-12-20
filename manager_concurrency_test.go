package dbsteps //nolint:testpackage

import (
	"context"
	"database/sql/driver"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bool64/sqluct"
	"github.com/cucumber/godog"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestNewManager_concurrent(t *testing.T) {
	dbm := NewManager()

	db1, mock1, err := sqlmock.New()
	assert.NoError(t, err)
	db2, mock2, err := sqlmock.New()
	assert.NoError(t, err)
	db3, mock3, err := sqlmock.New()
	assert.NoError(t, err)

	mock1.ExpectExec(`DELETE FROM t1`).
		WillReturnResult(driver.ResultNoRows)
	mock2.ExpectExec(`DELETE FROM t2`).
		WillReturnResult(driver.ResultNoRows)
	mock3.ExpectExec(`DELETE FROM t3`).
		WillReturnResult(driver.ResultNoRows)

	dbm.Instances = map[string]Instance{
		"db1": {
			Storage: sqluct.NewStorage(sqlx.NewDb(db1, "sqlmock")),
			Tables:  map[string]interface{}{"t1": nil},
		},
		"db2": {
			Storage: sqluct.NewStorage(sqlx.NewDb(db2, "sqlmock")),
			Tables:  map[string]interface{}{"t2": nil},
		},
		"db3": {
			Storage: sqluct.NewStorage(sqlx.NewDb(db3, "sqlmock")),
			Tables:  map[string]interface{}{"t3": nil},
		},
	}

	suite := godog.TestSuite{
		ScenarioInitializer: func(s *godog.ScenarioContext) {
			dbm.RegisterSteps(s)
			s.Step(`^I should not be blocked for "([^"]*)"$`, func(ctx context.Context, key string) error {
				if dbm.lock.IsLocked(ctx, key) {
					return fmt.Errorf("%s is locked", key)
				}

				return nil
			})
			s.Step("^I sleep$", func() {
				time.Sleep(time.Millisecond * time.Duration(100+rand.Int63n(100)))
			})
		},
		Options: &godog.Options{
			Format:      "pretty",
			Strict:      true,
			Paths:       []string{"_testdata/DatabaseConcurrent.feature"},
			Concurrency: 10,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("test failed")
	}
}

func TestNewManager_concurrent_blocked(t *testing.T) {
	dbm := NewManager()

	db1, mock1, err := sqlmock.New()
	assert.NoError(t, err)
	db3, mock3, err := sqlmock.New()
	assert.NoError(t, err)

	mock1.ExpectExec(`DELETE FROM t1`).
		WillReturnResult(driver.ResultNoRows)
	mock1.ExpectExec(`DELETE FROM t1`).
		WillReturnResult(driver.ResultNoRows)
	mock3.ExpectExec(`DELETE FROM t3`).
		WillReturnResult(driver.ResultNoRows)

	dbm.Instances = map[string]Instance{
		"db1": {
			Storage: sqluct.NewStorage(sqlx.NewDb(db1, "sqlmock")),
			Tables:  map[string]interface{}{"t1": nil},
		},
		"db3": {
			Storage: sqluct.NewStorage(sqlx.NewDb(db3, "sqlmock")),
			Tables:  map[string]interface{}{"t3": nil},
		},
	}

	suite := godog.TestSuite{
		ScenarioInitializer: func(s *godog.ScenarioContext) {
			dbm.RegisterSteps(s)
			s.Step(`^I should not be blocked for "([^"]*)"$`, func(ctx context.Context, key string) error {
				if dbm.lock.IsLocked(ctx, key) {
					return fmt.Errorf("%s is locked", key)
				}

				return nil
			})
			s.Step(`^I should be blocked for "([^"]*)"$`, func(ctx context.Context, key string) error {
				if !dbm.lock.IsLocked(ctx, key) {
					return fmt.Errorf("%s is not locked", key)
				}

				return nil
			})
			s.Step("^I sleep$", func() {
				time.Sleep(time.Millisecond * time.Duration(100+rand.Int63n(100)))
			})
		},
		Options: &godog.Options{
			Format:      "pretty",
			Strict:      true,
			Paths:       []string{"_testdata/DatabaseConcurrentBlocked.feature"},
			Concurrency: 10,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("test failed")
	}
}
