package dbsteps_test

import (
	"database/sql"
	"encoding/json"

	"github.com/bool64/sqluct"
	"github.com/cucumber/godog"
	"github.com/godogx/dbsteps"
	"github.com/godogx/vars"
)

func ExampleNewManager() {
	var db *sql.DB

	vs := &vars.Steps{}

	// Initialize database manager.
	dbm := dbsteps.NewManager()
	dbm.VS = vs // Setup shared vars steps.
	dbm.AddDB(db)

	suite := godog.TestSuite{
		Name:                 "DatabaseContext",
		TestSuiteInitializer: nil,
		ScenarioInitializer: func(s *godog.ScenarioContext) {
			dbm.RegisterSteps(s)
			vs.Register(s)
		},
		Options: &godog.Options{
			Format: "pretty",
			Paths:  []string{"features"},
			Strict: true,
		},
	}
	status := suite.Run()

	println(status)
}

func ExampleNewManager_advanced() {
	// Set up a storage adapter for your database connection.
	var storage *sqluct.Storage

	// Define database row structures in your repository packages.
	type MyRow struct {
		Foo string `db:"foo"`
	}

	type MyAnotherRow struct {
		Bar int     `db:"bar"`
		Baz float64 `db:"baz"`
	}

	// ...........

	// Initialize database manager with storage and table rows references.
	dbm := dbsteps.NewManager()
	dbm.Instances["my_db"] = dbsteps.Instance{
		Storage: storage,
		Tables: map[string]any{
			"my_table":         new(MyRow),
			"my_another_table": new(MyAnotherRow),
		},
		// Optionally configure statements to execute after deleting rows from table.
		PostCleanup: map[string][]string{
			"my_table": {"ALTER SEQUENCE my_table_id_seq RESTART"},
		},
	}
}

func ExampleNewTableMapper() {
	type jsonData struct {
		Foo string `json:"foo"`
	}

	tableMapper := dbsteps.NewTableMapper()

	// Apply JSON decoding to a particular type.
	tableMapper.Decoder.RegisterFunc(func(s string) (any, error) {
		data := jsonData{}

		err := json.Unmarshal([]byte(s), &data)
		if err != nil {
			return nil, err
		}

		return data, err
	}, jsonData{})

	// Create database manager with custom mapper.
	dbm := dbsteps.NewManager()
	dbm.TableMapper = tableMapper
}
