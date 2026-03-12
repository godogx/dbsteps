package dbsteps_test

import (
	"bytes"
	"database/sql"
	"os"
	"strings"
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

func TestDiffClosestRowFeature(t *testing.T) {
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
		Name:                 "DatabaseDiffClosestRow",
		TestSuiteInitializer: nil,
		ScenarioInitializer: func(s *godog.ScenarioContext) {
			m.RegisterSteps(s)
			vs.Register(s)
		},
		Options: &godog.Options{
			Format:    "pretty",
			Output:    buf,
			Paths:     []string{"_testdata/DiffClosestRow.feature"},
			NoColors:  true,
			Strict:    true,
			Randomize: time.Now().UTC().UnixNano(),
		},
	}
	status := suite.Run()

	if status == 0 {
		t.Fatal("expected scenario failure")
	}

	out := buf.String()
	require.Contains(t, out, "Diff vs closest row")
	require.Contains(t, out, "Scenario: Diff With Vars")
	require.Contains(t, out, "Scenario: Diff Transposed")
	require.Contains(t, out, "matched 4/5 columns")
	require.Contains(t, out, "matched 5/6 columns")
	require.Contains(t, out, "id     | 99")
	require.Contains(t, out, "age    | 32")

	if block := diffBlock(out); block != "" {
		require.NotContains(t, block, "$created_at")
	}
}

func diffBlock(out string) string {
	start := strings.Index(out, "Diff vs closest row")
	if start == -1 {
		return ""
	}

	block := out[start:]
	if next := strings.Index(block, "Scenario: "); next > 0 {
		block = block[:next]
	}

	return block
}
