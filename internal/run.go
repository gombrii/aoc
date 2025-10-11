package internal

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"golang.org/x/mod/modfile"
)

const runner = `package main

import (
	"fmt"
	"time"
	"os"

	"{{ .PkgPath }}"
)

func main() {
	data := read("{{ .InputPath }}")
	start := time.Now()
	result := {{ .PkgName }}.{{ .FuncName }}(data)
	dur := time.Since(start)
	fmt.Println("Res:", result)
	fmt.Println("Dur:", dur)
}
	
func read(path string) []byte {
	_, err := os.Stat(path)
	if err != nil {
		fmt.Printf("Error: invalid file path: %v", err)
		os.Exit(1)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error: could not read file: %v", err)
		os.Exit(1)
	}

	return data
}`

func Run(year, day, part, input string) error {
	day = fmt.Sprintf("day%s", day)
	part = fmt.Sprintf("part%s", part)
	inputFile := fmt.Sprintf("%s.txt", input)

	if _, err := os.Stat(filepath.Join(year, "solutions", day, fmt.Sprintf("%s.go", part))); err != nil {
		return fmt.Errorf("%v does not exist", filepath.Join(year, day, part))
	}

	fmt.Printf("Running %s with %s\n", filepath.Join(year, day, part), inputFile)

	mod, err := currentModulePath()
	if err != nil {
		return fmt.Errorf("starting runner: %v", err)
	}

	file, err := getRunner(year, day, part, input, map[string]string{
		"PkgPath":   filepath.Join(mod, year, "solutions", day),
		"PkgName":   day,
		"FuncName":  strings.Replace(part, "p", "P", 1),
		"InputPath": filepath.Join(year, "input", day, inputFile),
	})
	if err != nil {
		return fmt.Errorf("setting up runner: %v", err)
	}

	if err = executeRunner(file); err != nil {
		return fmt.Errorf("executing runner: %v", err)
	}

	return nil
}

func getRunner(year, day, part, input string, data map[string]string) (string, error) {
	osCache, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("finding cache location: %v", err)
	}

	file := filepath.Join(osCache, "gombrii-aoc", fmt.Sprintf("%s-%s-%s-%s", year, day, part, input), "runner.go")

	if _, err := os.Stat(file); err == nil{
		return file, nil
	}

	return file, createRunner(file, data)
}

func createRunner(fileName string, data map[string]string) error {
	dir := filepath.Dir(fileName)
	
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		os.RemoveAll(dir)
		return fmt.Errorf("creating temp dir: %v", err)
	}

	file, err := os.Create(fileName)
	if err != nil {
		os.RemoveAll(dir)
		return fmt.Errorf("creating runner file: %v", err)
	}
	defer file.Close()

	if err = template.Must(template.New("runner").Parse(runner)).Execute(file, data); err != nil {
		os.RemoveAll(dir)
		return fmt.Errorf("compiling runner: %v", err)
	}

	return nil
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
