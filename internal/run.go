package internal

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gombrii/aoc/internal/cache"
	"github.com/gombrii/aoc/internal/gen"
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
		return fmt.Errorf("getting module name: %v", err)
	}

	path, err := getRunnerPath(year, day, part, input, map[string]string{
		"PkgPath":   filepath.Join(mod, year, "solutions", day),
		"PkgName":   day,
		"FuncName":  strings.Replace(part, "p", "P", 1),
		"InputPath": filepath.Join(year, "input", day, inputFile),
	})
	if err != nil {
		return fmt.Errorf("setting up runner: %v", err)
	}

	if err = executeRunner(path); err != nil {
		return fmt.Errorf("executing runner: %v", err)
	}

	return nil
}

func getRunnerPath(year, day, part, input string, data map[string]string) (string, error) {
	cacheKey := fmt.Sprintf("%s-%s-%s-%s", year, day, part, input)

	if cPath, ok := cache.ContainsKey(cacheKey); ok {
		return cPath, nil
	}

	files, err := gen.TempFiles(map[string]string{
		"runner.go": runner,
		"last":      "",
		"lock":      "",
		"dur":       strconv.Itoa(math.MaxInt64),
	}, data)
	if err != nil {
		return "", fmt.Errorf("generating files: %v", err)
	}

	rPath := ""
	for name, path := range files {
		cPath, err := cache.Store(cacheKey, name, path)
		if err != nil {
			return "", fmt.Errorf("caching runner: %v", err)
		}
		if name == "runner.go" {
			rPath = cPath
		}
	}

	return rPath, nil
}

func executeRunner(rPath string) error {
	cmd := exec.Command("go", "run", rPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	var exitErr *exec.ExitError
	if err := cmd.Run(); err != nil && !errors.As(err, &exitErr) {
		return err
	}

	return nil
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
