package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

const mod = `module {{.Mod}}

{{.Version}}`

const exit = `// Package exit does, in the spirit of Advent of Code, provide quick and dirty ways to fail 
package exit

import (
	"fmt"
	"os"
)

// If prints err and exits if err != nil.
func If(err error) {
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
	
// PanicIf panics if err != nil.
func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}`

const input = `// Package input does, in the spirit of Advent of Code, provide some common ways to interpret the 
// challanges' input data.
package input

import (
	"strings"
)

// Lines converts data into text and returns it as a slice where each element represents a line.
func Lines(data []byte) []string {
	return strings.Split(string(data), "\n")
}

// String returns data as a continuous string.
func String(data []byte) string {
	return string(data)
}

// Matrix returns data as a matrix. The delimiter determines what in data separates columns 
// while rows are seperated by linebreaks.
func Matrix(data []byte, delimiter string) [][]string {
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	matrix := make([][]string, len(lines))

	for i, line := range lines {
		matrix[i] = strings.Split(line, delimiter)
	}

	return matrix
}`

// NOTE: Generates project scaffolding, which is really just the shared folder with some nice to haves:
// - input
// - exit
// - render
// - go.mod
// Maybe the init AOC command should also take a single string with the mod name, like go mod init
func GenAoc(module string) error {
	//Generate shared folder with all it contains and go.mod file, I guess with the go version the user has

	out, err := exec.Command("go", "version").Output()
	if err != nil {
		return err
	}
	ver := strings.TrimSpace(string(out))
	ver = regexp.MustCompile(`go\d\.\d{2}\.\d`).FindString(ver)
	ver = strings.Replace(ver, "go", "go ", 1)

	err = createFile("go.mod", mod, map[string]string{
		"Mod":     module,
		"Version": ver,
	})
	if err != nil {
		return fmt.Errorf("initiating go.mod: %v", err)
	}

	err = os.MkdirAll(filepath.Join("shared", "input"), 0755)
	if err != nil {
		return fmt.Errorf("creating directory shared/input: %v", err)
	}

	createFile(filepath.Join("shared", "input", "input.go"), input, nil)

	err = os.MkdirAll(filepath.Join("shared", "exit"), 0755)
	if err != nil {
		return fmt.Errorf("creating directory shared/exit: %v", err)
	}

	createFile(filepath.Join("shared", "exit", "error.go"), exit, nil)

	fmt.Printf("Initiated Advent of Code module %s\n", module)

	return nil
}

func createFile(fileName string, tmpl string, data map[string]string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("creating file %s: %v", fileName, err)
	}
	defer file.Close()

	cmplTmpl, err := template.New(file.Name()).Parse(tmpl)
	if err != nil {
		return fmt.Errorf("parsing template for file go.mod: %v", err)
	}

	err = cmplTmpl.Execute(file, data)
	if err != nil {
		return fmt.Errorf("writing to file %s: %v", fileName, err)
	}

	return nil
}
