# Cucumber database steps for Go

[![Build Status](https://github.com/godogx/dbsteps/workflows/test-unit/badge.svg)](https://github.com/godogx/dbsteps/actions?query=branch%3Amaster+workflow%3Atest-unit)
[![Coverage Status](https://codecov.io/gh/godogx/dbsteps/branch/master/graph/badge.svg)](https://codecov.io/gh/godogx/dbsteps)
[![GoDevDoc](https://img.shields.io/badge/dev-doc-00ADD8?logo=go)](https://pkg.go.dev/github.com/godogx/dbsteps)
[![Time Tracker](https://wakatime.com/badge/github/godogx/dbsteps.svg)](https://wakatime.com/badge/github/godogx/dbsteps)
![Code lines](https://sloc.xyz/github/godogx/dbsteps/?category=code)
![Comments](https://sloc.xyz/github/godogx/dbsteps/?category=comments)

This module implements database-related step definitions
for [`github.com/cucumber/godog`](https://github.com/cucumber/godog).

## Database Configuration

Databases instances should be configured with `Manager.Instances`.

```
// Initialize database manager with storage and table rows references.
dbm := dbsteps.NewManager()
dbm.Instances["my_db"] = dbsteps.Instance{
    Storage: storage,
    Tables: map[string]interface{}{
        "my_table":         new(repository.MyRow),
        "my_another_table": new(repository.MyAnotherRow),
    },
    // Optionally configure statements to execute after deleting rows from table.
    PostCleanup: map[string][]string{
        "my_table": {"ALTER SEQUENCE my_table_id_seq RESTART"},
    },
}
```

Row types should be structs with `db` field tags, for example:

```go
type MyRow struct {
    ID   int    `db:"id"`
    Name string `db:"name"`
}
```

These structures are used to map data between database and `gherkin` tables.

## Table Mapper Configuration

Table mapper allows customizing decoding string values from godog table cells into Go row structures and back.

```go
tableMapper := dbsteps.NewTableMapper()

// Apply JSON decoding to a particular type.
tableMapper.Decoder.RegisterFunc(func(s string) (interface{}, error) {
    m := repository.Meta{}
    err := json.Unmarshal([]byte(s), &m)
    if err != nil {
        return nil, err
    }
	
    return m, err
}, repository.Meta{})

// Apply string splitting to github.com/lib/pq.StringArray.
tableMapper.Decoder.RegisterFunc(func(s string) (interface{}, error) {
    return pq.StringArray(strings.Split(s, ",")), nil
}, pq.StringArray{})

// Create database manager with custom mapper.
dbm := dbsteps.NewManager()
dbm.TableMapper = tableMapper
```

## Step Definitions

Delete all rows from table.

```gherkin
Given there are no rows in table "my_table" of database "my_db"
```

Populate rows in a database.

```gherkin
And these rows are stored in table "my_table" of database "my_db"
| id | foo   | bar | created_at           | deleted_at           |
| 1  | foo-1 | abc | 2021-01-01T00:00:00Z | NULL                 |
| 2  | foo-1 | def | 2021-01-02T00:00:00Z | 2021-01-03T00:00:00Z |
| 3  | foo-2 | hij | 2021-01-03T00:00:00Z | 2021-01-03T00:00:00Z |
```

```gherkin
And rows from this file are stored in table "my_table" of database "my_db"
 """
 path/to/rows.csv
 """
```

Assert rows existence in a database.

For each row in gherkin table database is queried to find a row with `WHERE` condition that includes provided column
values.

If a column has `NULL` value, it is excluded from `WHERE` condition.

Column can contain variable (any unique string starting with `$` or other prefix configured with `Manager.VarPrefix`).
If variable has not yet been populated, it is excluded from `WHERE` condition and populated with value received from
database. When this variable is used in next steps, it replaces the value of column with value of variable.

Variables can help to assert consistency of dynamic data, for example variable can be populated as ID of one entity and
then checked as foreign key value of another entity. This can be especially helpful in cases of UUIDs.

If column value represents JSON array or object it is excluded from `WHERE` condition, value assertion is done by
comparing Go value mapped from database row field with Go value mapped from gherkin table cell.

```gherkin
Then these rows are available in table "my_table" of database "my_db"
| id   | foo   | bar | created_at           | deleted_at           |
| $id1 | foo-1 | abc | 2021-01-01T00:00:00Z | NULL                 |
| $id2 | foo-1 | def | 2021-01-02T00:00:00Z | 2021-01-03T00:00:00Z |
| $id3 | foo-2 | hij | 2021-01-03T00:00:00Z | 2021-01-03T00:00:00Z |
```

```gherkin
Then rows from this file are available in table "my_table" of database "my_db"
 """
 path/to/rows.csv
 """
```

It is possible to check table contents exhaustively by adding "only" to step statement. Such assertion will also make
sure that total number of rows in database table matches number of rows in gherkin table.

```gherkin
Then only these rows are available in table "my_table" of database "my_db"
| id   | foo   | bar | created_at           | deleted_at           |
| $id1 | foo-1 | abc | 2021-01-01T00:00:00Z | NULL                 |
| $id2 | foo-1 | def | 2021-01-02T00:00:00Z | 2021-01-03T00:00:00Z |
| $id3 | foo-2 | hij | 2021-01-03T00:00:00Z | 2021-01-03T00:00:00Z |
```

```gherkin
Then only rows from this file are available in table "my_table" of database "my_db"
 """
 path/to/rows.csv
 """
```

Assert no rows exist in a database.

```gherkin
And no rows are available in table "my_another_table" of database "my_db"
```

The name of database instance `of database "my_db"` can be omitted in all steps, in such case `"default"` will be used from database instance name.

## Concurrent Usage

Please note, due to centralized nature of database instance, scenarios that work with same tables would conflict.
In order to avoid conflicts, `dbsteps` locks access to a specific scenario until that scenario is finished.
The lock is per table, so if scenarios are operating on different tables, they will not conflict.
It is safe to use concurrent scenarios.
