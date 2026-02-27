package dbsteps_test

import (
	"bytes"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/bool64/sqluct"
	"github.com/cucumber/godog"
	"github.com/godogx/dbsteps"
	"github.com/godogx/vars"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite" // DB driver.
)

func TestNewManager(t *testing.T) {
	sqlDB, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)

	fixture, err := os.ReadFile("_testdata/fixture.sql")
	require.NoError(t, err)

	for _, st := range sqluct.SplitStatements(string(fixture)) {
		_, err := sqlDB.Exec(st)
		require.NoError(t, err, st)
	}

	vs := vars.Steps{}

	m := dbsteps.NewManager()
	m.AddDB(sqlDB)

	buf := bytes.NewBuffer(nil)

	suite := godog.TestSuite{
		Name:                 "DatabaseContext",
		TestSuiteInitializer: nil,
		ScenarioInitializer: func(s *godog.ScenarioContext) {
			m.RegisterSteps(s)
			vs.Register(s)
		},
		Options: &godog.Options{
			Format:    "pretty",
			Output:    buf,
			Paths:     []string{"_testdata/DatabaseE2E.feature"},
			Strict:    true,
			Randomize: time.Now().UTC().UnixNano(),
		},
	}
	status := suite.Run()

	if status != 0 {
		t.Fatal(buf.String())
	}
}
