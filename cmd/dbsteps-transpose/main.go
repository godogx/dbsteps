// Package main implements a command-line tool for transposing Gherkin-style tables.
package main

import (
	"bufio"
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

	input, err := readInput()
	if err != nil {
		fatal(err)
	}

	rows, err := parseTable(input)
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

func readInput() (string, error) {
	info, err := os.Stdin.Stat()
	if err != nil {
		return "", err
	}

	if info.Mode()&os.ModeCharDevice == 0 {
		input, err := io.ReadAll(os.Stdin)

		return string(input), err
	}

	fmt.Fprintln(os.Stderr, "Paste table, end with an empty line (or Ctrl-D).")

	var b strings.Builder

	reader := bufio.NewReader(os.Stdin)
	sawData := false

	for {
		line, err := reader.ReadString('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			return "", err
		}

		if err == io.EOF && line == "" {
			break
		}

		trimmed := strings.TrimRight(line, "\r\n")
		if trimmed == "" {
			if sawData {
				break
			}

			if err == io.EOF {
				break
			}

			continue
		}

		sawData = true

		b.WriteString(line)

		if err == io.EOF {
			break
		}
	}

	return b.String(), nil
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

		cells := splitRow(line)

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

func splitRow(line string) []string {
	var (
		cells []string
		b     strings.Builder
	)

	escaped := false

	flush := func() {
		cells = append(cells, strings.TrimSpace(b.String()))
		b.Reset()
	}

	for _, r := range line {
		if escaped {
			switch r {
			case '|', '\\':
				b.WriteRune(r)
			default:
				b.WriteByte('\\')
				b.WriteRune(r)
			}

			escaped = false

			continue
		}

		switch r {
		case '\\':
			escaped = true
		case '|':
			flush()
		default:
			b.WriteRune(r)
		}
	}

	if escaped {
		b.WriteByte('\\')
	}

	flush()

	return cells
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
			escaped := escapeCell(cell)
			if len(escaped) > widths[i] {
				widths[i] = len(escaped)
			}
		}
	}

	var b strings.Builder
	for _, row := range rows {
		b.WriteString("|")

		for i, cell := range row {
			escaped := escapeCell(cell)

			b.WriteString(" ")
			b.WriteString(escaped)

			for pad := widths[i] - len(escaped); pad > 0; pad-- {
				b.WriteByte(' ')
			}

			b.WriteString(" |")
		}

		b.WriteString("\n")
	}

	return b.String()
}

func escapeCell(cell string) string {
	if cell == "" {
		return cell
	}

	var b strings.Builder

	b.Grow(len(cell))

	for _, r := range cell {
		switch r {
		case '\\':
			b.WriteString(`\\`)
		case '|':
			b.WriteString(`\|`)
		default:
			b.WriteRune(r)
		}
	}

	return b.String()
}
