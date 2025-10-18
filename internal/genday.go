package internal

import (
	"fmt"
	"path/filepath"

	"github.com/gombrii/aoc/internal/files"
)

const part1Tmpl = `// Package {{.DayName}} solves puzzle available on https://adventofcode.com/{{.Year}}/day/{{.Day}}
package {{.DayName}}

func Part1(data []byte) any {
	//in := parse.Lines(data)

	return "NOT IMPLEMENTED!"
}`

const part2Tmpl = `// Package {{.DayName}} solves puzzle available on https://adventofcode.com/{{.Year}}/day/{{.Day}}
package {{.DayName}}

func Part2(data []byte) any {
	//in := parse.Lines(data)

	return "NOT IMPLEMENTED!"
}`

const commonTmpl = `// Package {{.DayName}} solves puzzle available on https://adventofcode.com/{{.Year}}/day/{{.Day}}
package {{.DayName}}`

func GenDay(year, day string) error {
	dayName := fmt.Sprintf("day%s", day)

	if err := files.Gen(
		map[string]string{
			filepath.Join(year, "solutions", dayName, "part1.go"):  part1Tmpl,
			filepath.Join(year, "solutions", dayName, "part2.go"):  part2Tmpl,
			filepath.Join(year, "solutions", dayName, "common.go"): commonTmpl,
			filepath.Join(year, "input", dayName, "input.txt"):     "",
			filepath.Join(year, "input", dayName, "test.txt"):      "",
		},
		map[string]string{
			"Year":    year,
			"Day":     day,
			"DayName": dayName,
		}); err != nil {
		return fmt.Errorf("generating files: %v", err)
	}

	fmt.Printf("Initiated %s %s\n", dayName, year)

	return nil
}
