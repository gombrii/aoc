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

func GenDay(year, day int) error {
	yName := fmt.Sprintf("%d", year)
	dName := fmt.Sprintf("day%d", day)

	if err := files.Gen(
		map[string]string{
			filepath.Join(yName, "solutions", dName, "part1.go"):  part1Tmpl,
			filepath.Join(yName, "solutions", dName, "part2.go"):  part2Tmpl,
			filepath.Join(yName, "solutions", dName, "common.go"): commonTmpl,
			filepath.Join(yName, "input", dName, "input.txt"):     "",
			filepath.Join(yName, "input", dName, "test.txt"):      "",
		},
		map[string]string{
			"Year":    yName,
			"Day":     fmt.Sprint(day),
			"DayName": dName,
		}); err != nil {
		return fmt.Errorf("generating files: %v", err)
	}

	fmt.Printf("Initiated %s %s\n", dName, yName)

	return nil
}
