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
	"time"

	"github.com/gombrii/aoc/internal/cache"
	"github.com/gombrii/aoc/internal/files"
	"golang.org/x/mod/modfile"
)

const runnerTmpl = `package main

import (
	"fmt"
	"time"
	"os"
	"strconv"
	"strings"

	"{{ .PkgPath }}"
)

func main() {
	locked, _ := strconv.ParseBool(strings.TrimSpace(string(read("{{ .LockPath }}"))))
	record, _ := time.ParseDuration(strings.TrimSpace(string(read("{{ .DurPath }}"))))
	lastRes := strings.TrimSpace(string(read("{{ .ResPath }}")))
	data := read("{{ .InputPath }}")
	
	start := time.Now()
	result := {{ .PkgName }}.{{ .FuncName }}(data)	
	duration := time.Since(start)

	correct := fmt.Sprint(result) == lastRes
	switch {
	case locked && !correct:
		fmt.Printf("Error: res: %v, want %v\n", result, lastRes)
		return
	case locked && correct:
		diff := duration - record
		fmt.Println("Res:", result)
		fmt.Printf("Dur: %v (%v, %.0f%%)\n", duration, diff, (float64(diff)/float64(duration))*100.0)
		if diff < 0 {
			write("{{ .DurPath }}", fmt.Sprint(duration))
		}
	case !locked:
		fmt.Println("Res:", result)
		fmt.Println("Dur:", duration)
		write("{{ .ResPath }}", fmt.Sprint(result))
		write("{{ .DurPath }}", fmt.Sprint(duration))
	}
}
	
func read(path string) []byte {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error: could not read file: %v", err)
		os.Exit(1)
	}

	return data
}
	
func write(path string, data string) {
	if err := os.WriteFile(path, []byte(data), 0755); err != nil {
		fmt.Printf("Error: could not write file: %v\n", err)
		os.Exit(1)
	}
}`

func Run(year, day, part, input string) error {
	day = fmt.Sprintf("day%s", day)
	part = fmt.Sprintf("part%s", part)

	if _, err := os.Stat(filepath.Join(year, "solutions", day, fmt.Sprintf("%s.go", part))); err != nil {
		return fmt.Errorf("%v does not exist", filepath.Join(year, day, part))
	}

	fmt.Printf("Running %s with %s\n", filepath.Join(year, day, part), input)

	path, err := getRunnerPath(year, day, part, input)
	if err != nil {
		return fmt.Errorf("setting up runner: %v", err)
	}

	if err = executeRunner(path); err != nil {
		return fmt.Errorf("executing runner: %v", err)
	}

	return nil
}

func getRunnerPath(year, day, part, input string) (string, error) {
	cacheKey := cache.Key(year, day, part, input)

	if cPath, ok := cache.Contains(cacheKey, files.Runner); ok {
		return cPath, nil
	}

	mod, err := currentModulePath()
	if err != nil {
		return "", fmt.Errorf("getting module name: %v", err)
	}

	fPaths, err := files.GenTemp(map[string]string{
		files.Runner: runnerTmpl,
		files.Lock:   strconv.FormatBool(false),
		files.Res:    "",
		files.Dur:    time.Duration(math.MaxInt64).String(),
	}, map[string]string{
		"PkgPath":   filepath.Join(mod, year, "solutions", day),
		"PkgName":   day,
		"FuncName":  strings.Replace(part, "p", "P", 1),
		"InputPath": filepath.Join(year, "input", day, input),
		"LockPath":  cache.MakePath(cacheKey, files.Lock),
		"ResPath":   cache.MakePath(cacheKey, files.Res),
		"DurPath":   cache.MakePath(cacheKey, files.Dur),
	})
	if err != nil {
		return "", fmt.Errorf("generating files: %v", err)
	}

	rPath := ""
	for name, path := range fPaths {
		cPath, err := cache.Store(cacheKey, name, path)
		if err != nil {
			return "", fmt.Errorf("caching files: %v", err)
		}
		if name == files.Runner {
			rPath = cPath
		}
	}

	return rPath, nil
}

func executeRunner(path string) error {
	cmd := exec.Command("go", "run", path)
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
