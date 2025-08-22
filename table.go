package dbsteps

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/godogx/vars"
	"github.com/swaggest/form/v5"
)

const null = "NULL"

// TableMapper maps data from Go value to string and back.
type TableMapper struct {
	Decoder *form.Decoder
	Encoder *form.Encoder
}

func isNil(v interface{}) bool {
	if v == nil {
		return true
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr && rv.IsZero() {
		return true
	}

	return false
}

// Encode converts Go value to string.
func (m *TableMapper) Encode(v interface{}) (string, error) {
	if m.Encoder == nil {
		m.Encoder = form.NewEncoder()
	}

	if isNil(v) {
		return null, nil
	}

	vv, err := m.Encoder.Encode(v)
	if err != nil {
		return "", fmt.Errorf("failed to stringify variable value of type %T: %w", v, err)
	}

	return vv[""][0], nil
}

// SliceFromTable creates a slice from gherkin table, item type is used as slice element type.
func (m *TableMapper) SliceFromTable(data [][]string, item interface{}) (interface{}, error) {
	itemType := reflect.TypeOf(item)
	if itemType == nil {
		return nil, errNilItemStruct
	}

	if itemType.Kind() == reflect.Ptr {
		itemType = itemType.Elem()
	}

	result := reflect.MakeSlice(reflect.SliceOf(itemType), len(data)-1, len(data)-1)

	err := m.IterateTable(IterateConfig{
		Data: data, Item: item,
		ReceiveRow: func(index int, row interface{}, colNames []string, rawValues []string) error {
			result.Index(index).Set(reflect.Indirect(reflect.ValueOf(row)))

			return nil
		},
	})
	if err != nil {
		return nil, err
	}

	return result.Interface(), nil
}

// IterateConfig controls behavior of TableMapper.IterateTable.
type IterateConfig struct {
	Data       [][]string
	SkipDecode func(column, value string) bool
	Item       interface{}
	Replaces   map[string]string
	ReceiveRow func(index int, row interface{}, colNames []string, rawValues []string) error
}

var (
	errNilItemStruct = errors.New("nil item struct received")
	errRowRequired   = errors.New("header and at least one row required in table")
)

func itemType(v interface{}) (reflect.Type, error) {
	itemType := reflect.TypeOf(v)
	if itemType == nil {
		return nil, errNilItemStruct
	}

	if itemType.Kind() == reflect.Ptr {
		itemType = itemType.Elem()
	}

	return itemType, nil
}

// IterateTable walks gherkin table calling row receiver with mapped row.
// If receiver returns error iteration stops and error is propagated.
func (m *TableMapper) IterateTable(c IterateConfig) error {
	if m.Decoder == nil && c.Item != nil {
		m.Decoder = form.NewDecoder()
	}

	if len(c.Data) < 2 {
		return errRowRequired
	}

	var (
		it       reflect.Type
		err      error
		itemBuf  reflect.Value
		colNames = c.Data[0]
		values   = make(map[string][]string, len(colNames))
		rowMap   = make(map[string]interface{})
	)

	if c.Item != nil {
		it, err = itemType(c.Item)
		if err != nil {
			return err
		}
	}

	for rowIndex, row := range c.Data[1:] {
		if c.Item != nil {
			itemBuf = reflect.New(it)
		}

		raw := make([]string, 0, len(colNames))

		for i, cell := range row {
			raw = append(raw, cell)

			if c.SkipDecode != nil && c.SkipDecode(colNames[i], cell) {
				continue
			}

			cell = strings.TrimSuffix(cell, "::string")

			if v, found := c.Replaces[cell]; found {
				cell = v
			}

			if cell != null {
				rowMap[colNames[i]] = vars.Infer(cell)
				values[colNames[i]] = []string{cell}
			} else {
				rowMap[colNames[i]] = nil

				delete(values, colNames[i])
			}
		}

		if err = c.ReceiveRow(rowIndex, m.decode(itemBuf, values, rowMap), colNames, raw); err != nil {
			return err
		}
	}

	return nil
}

func (m *TableMapper) decode(itemBuf reflect.Value, values map[string][]string, rowMap map[string]interface{}) interface{} {
	if !itemBuf.IsValid() {
		return rowMap
	}

	val := itemBuf.Interface()

	err := m.Decoder.Decode(val, values)
	if err != nil {
		return err
	}

	return itemBuf.Interface()
}
