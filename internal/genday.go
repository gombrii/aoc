package internal

import (
	"fmt"
	"path/filepath"

	"github.com/gombrii/aoc/internal/gen"
)

const part1 = `package {{.Day}}

func Part1(data []byte) any {
	//in := input.Lines(data)

	return "NOT IMPLEMENTED!"
}`

const part2 = `package {{.Day}}

func Part2(data []byte) any {
	//in := input.Lines(data)

	return "NOT IMPLEMENTED!"
}`

const common = `package {{.Day}}`

func GenDay(year, day string) error {
	day = fmt.Sprintf("day%s", day)

	data := map[string]string{
		"Day": day,
	}

	if err := gen.Files(data, map[string]string{
		filepath.Join(year, "solutions", day, "part1.go"):  part1,
		filepath.Join(year, "solutions", day, "part2.go"):  part2,
		filepath.Join(year, "solutions", day, "common.go"): common,
		filepath.Join(year, "input", day, "input.txt"):     "",
		filepath.Join(year, "input", day, "test.txt"):      "",
	}); err != nil {
		return fmt.Errorf("generating files: %v", err)
	}

	fmt.Printf("Initiated %s %s\n", day, year)

	return nil
}
