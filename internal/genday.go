package internal

import (
	"fmt"
	"path/filepath"

	"github.com/gombrii/aoc/internal/gen"
)

const part1 = `// Package {{.DayName}} solves puzzle available on https://adventofcode.com/{{.Year}}/day/{{.Day}}
package {{.DayName}}

func Part1(data []byte) any {
	//in := parse.Lines(data)

	return "NOT IMPLEMENTED!"
}`

const part2 = `// Package {{.DayName}} solves puzzle available on https://adventofcode.com/{{.Year}}/day/{{.Day}}
package {{.DayName}}

func Part2(data []byte) any {
	//in := parse.Lines(data)

	return "NOT IMPLEMENTED!"
}`

const common = `// Package {{.DayName}} solves puzzle available on https://adventofcode.com/{{.Year}}/day/{{.Day}}
package {{.DayName}}`

func GenDay(year, day string) error {
	dayName := fmt.Sprintf("day%s", day)

	data := map[string]string{
		"Year":    year,
		"Day":     day,
		"DayName": dayName,
	}

	if err := gen.Files(data, map[string]string{
		filepath.Join(year, "solutions", dayName, "part1.go"):  part1,
		filepath.Join(year, "solutions", dayName, "part2.go"):  part2,
		filepath.Join(year, "solutions", dayName, "common.go"): common,
		filepath.Join(year, "input", dayName, "input.txt"):     "",
		filepath.Join(year, "input", dayName, "test.txt"):      "",
	}); err != nil {
		return fmt.Errorf("generating files: %v", err)
	}

	fmt.Printf("Initiated %s %s\n", dayName, year)

	return nil
}
