package dbsteps_test

import (
	"encoding/json"

	"github.com/bool64/sqluct"
	"github.com/godogx/dbsteps"
)

func ExampleNewManager() {
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
		Tables: map[string]interface{}{
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
	tableMapper.Decoder.RegisterFunc(func(s string) (interface{}, error) {
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
