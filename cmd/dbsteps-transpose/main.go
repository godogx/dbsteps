// Package main implements a command-line tool for transposing Gherkin-style tables.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

var errEmptyTable = errors.New("table is empty")

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: dbsteps-transpose < input.txt")
		fmt.Fprintln(os.Stderr, "Reads a Gherkin-style table from stdin and prints its transpose.")
	}
	flag.Parse()

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fatal(err)
	}

	rows, err := parseTable(string(input))
	if err != nil {
		fatal(err)
	}

	out, err := transpose(rows)
	if err != nil {
		fatal(err)
	}

	fmt.Print(formatTable(out))
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "error:", err)
	os.Exit(1)
}

func parseTable(input string) ([][]string, error) {
	lines := strings.Split(input, "\n")
	rows := make([][]string, 0, len(lines))

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if !strings.Contains(line, "|") {
			return nil, fmt.Errorf("invalid row (no pipes): %q", line) //nolint:err113
		}

		parts := strings.Split(line, "|")
		cells := make([]string, 0, len(parts))

		for _, p := range parts {
			cells = append(cells, strings.TrimSpace(p))
		}

		// Drop empty leading/trailing cells from leading/trailing pipes.
		if len(cells) > 0 && cells[0] == "" {
			cells = cells[1:]
		}

		if len(cells) > 0 && cells[len(cells)-1] == "" {
			cells = cells[:len(cells)-1]
		}

		if len(cells) == 0 {
			continue
		}

		rows = append(rows, cells)
	}

	if len(rows) == 0 {
		return nil, errEmptyTable
	}

	cols := len(rows[0])
	for i := 1; i < len(rows); i++ {
		if len(rows[i]) != cols {
			return nil, fmt.Errorf("non-rectangular table: row %d has %d cells, expected %d", i, len(rows[i]), cols) //nolint:err113
		}
	}

	return rows, nil
}

func transpose(in [][]string) ([][]string, error) {
	if len(in) == 0 {
		return nil, errEmptyTable
	}

	cols := len(in[0])
	if cols == 0 {
		return nil, errEmptyTable
	}

	out := make([][]string, cols)

	for c := 0; c < cols; c++ {
		row := make([]string, len(in))

		for r := 0; r < len(in); r++ {
			row[r] = in[r][c]
		}

		out[c] = row
	}

	return out, nil
}

func formatTable(rows [][]string) string {
	widths := make([]int, len(rows[0]))

	for _, row := range rows {
		for i, cell := range row {
			if len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	var b strings.Builder
	for _, row := range rows {
		b.WriteString("|")

		for i, cell := range row {
			b.WriteString(" ")
			b.WriteString(cell)

			for pad := widths[i] - len(cell); pad > 0; pad-- {
				b.WriteByte(' ')
			}

			b.WriteString(" |")
		}

		b.WriteString("\n")
	}

	return b.String()
}
