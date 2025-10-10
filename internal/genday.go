package internal

import (
	"fmt"
	"os"
	"path"
	"text/template"
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

	solution := map[string]map[string]string{
		day: {
			"part1.go":  part1,
			"part2.go":  part2,
			"common.go": common,
		},
	}

	input := map[string]map[string]string{
		day: {
			"input.txt": "",
			"test.txt":  "",
		},
	}

	data := map[string]any{
		"Day":  day,
		"Year": year,
	}

	if err := create("solutions", year, solution, data); err != nil {
		return fmt.Errorf("generating solution scaffolding: %v", err)
	}
	if err := create("input", year, input, data); err != nil {
		return fmt.Errorf("generating input scaffolding: %v", err)
	}

	fmt.Printf("Initiated %s %s\n", day, year)
	return nil
}

func create(parentDir string, year string, structure map[string]map[string]string, data map[string]any) error {
	for dir, files := range structure {
		err := os.MkdirAll(path.Join(year, parentDir, dir), 0755)
		if err != nil {
			return fmt.Errorf("creating directory %s: %v", dir, err)
		}

		for fileName, templateContent := range files {
			filePath := path.Join(year, parentDir, dir, fileName)

			tmpl, err := template.New(fileName).Parse(templateContent)
			if err != nil {
				return fmt.Errorf("parsing template for file %s: %v", filePath, err)
			}

			file, err := os.Create(filePath)
			if err != nil {
				return fmt.Errorf("creating file %s: %v", filePath, err)
			}
			defer file.Close()

			err = tmpl.Execute(file, data)
			if err != nil {
				return fmt.Errorf("writing to file %s: %v", filePath, err)
			}
		}
	}

	return nil
}
