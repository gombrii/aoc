package internal

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"golang.org/x/mod/modfile"
)

var runner = `package main

import (
	"fmt"
	"time"

	"{{ .PkgPath }}"
)

func main() {
	start := time.Now()
	result := {{ .PkgName }}.{{ .FuncName }}("{{ .FuncArg }}")
	fmt.Println("Res:", result)
	fmt.Println("Dur:", time.Since(start))
}`

type data struct {
	PkgPath  string
	PkgName  string
	FuncName string
	FuncArg  string
}

func Run(year, day, part, input string) error {
	day = fmt.Sprintf("day%s", day)
	part = fmt.Sprintf("part%s", part)

	inputPath := fmt.Sprintf("%s/input/%s", year, day)
	inputFile := fmt.Sprintf("%s.txt", input)

	fmt.Printf("Running %s/%s/%s with %s\n", year, day, part, inputFile)

	//TODO: New attempt
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("starting runner: %v", err)
	}
	mod, err := currentModulePath()
	if err != nil {
		return fmt.Errorf("starting runner: %v", err)
	}

	dir, path, file, err := createRunner(data{
		PkgPath:  path.Join(mod, year, "solutions", day),
		PkgName:  day,
		FuncName: strings.Replace(part, "p", "P", -1),
		FuncArg:  filepath.Join(wd, inputPath, inputFile),
	})
	if err != nil {
		return fmt.Errorf("starting runner: %v", err)
	}
	defer os.RemoveAll(dir)
	defer file.Close()

	if err = executeRunner(path); err != nil {
		return fmt.Errorf("executing runner: %v", err)
	}

	return nil
}

func createRunner(data data) (dirPath, filePath string, file *os.File, err error) {
	dir, err := os.MkdirTemp("", "aoc-runner-*")
	if err != nil {
		return "", "", nil, fmt.Errorf("setting up runner: %v", err)
	}

	path := filepath.Join(dir, "runner.go")
	file, err = os.Create(path)
	if err != nil {
		return "", "", nil, fmt.Errorf("setting up runner: %v", err)
	}

	if err = template.Must(template.New("runner").Parse(runner)).Execute(file, data); err != nil {
		return "", "", nil, fmt.Errorf("setting up runner: %v", err)
	}

	return dir, path, file, nil
}

func executeRunner(path string) error {
	cmd := exec.Command("go", "run", path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

func currentModulePath() (string, error) {
	out, err := exec.Command("go", "env", "GOMOD").Output()
	if err != nil {
		return "", err
	}
	gomod := string(bytes.TrimSpace(out))
	if gomod == "" || gomod == "/dev/null" || gomod == "NUL" {
		return "", errors.New("not inside a Go module (no go.mod found)")
	}

	data, err := os.ReadFile(gomod)
	if err != nil {
		return "", err
	}
	f, err := modfile.Parse(gomod, data, nil)
	if err != nil {
		return "", err
	}
	if f.Module == nil || f.Module.Mod.Path == "" {
		return "", errors.New("module path not found in go.mod")
	}
	return f.Module.Mod.Path, nil
}
