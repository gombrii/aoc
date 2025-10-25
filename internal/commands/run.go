package commands

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gombrii/aoc/internal/cache"
	"github.com/gombrii/aoc/internal/exec"
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
		fmt.Printf("Error: could not read file: %v\n", err)
		os.Exit(1)
	}

	return data
}

func write(path string, data string) {
	if err := os.WriteFile(path, []byte(data), 0755); err != nil {
		fmt.Printf("Error: could not write file: %v\n", err)
		os.Exit(1)
	}
}
`

func (c Commands) Run(year, day, part int, input string) error {
	if !files.Exists("go.mod") {
		return errors.New("not in Go module root (no go.mod found)")
	}

	yName := fmt.Sprintf("%d", year)
	dName := fmt.Sprintf("day%d", day)
	pName := fmt.Sprintf("part%d", part)

	if !files.Exists(filepath.Join(yName, "solutions", dName, fmt.Sprintf("%s.go", pName))) {
		return fmt.Errorf("%s does not exist", filepath.Join(yName, dName, pName))
	}
	if !files.Exists(filepath.Join(yName, "input", dName, input)) {
		return fmt.Errorf("input file %s does not exist for %s", input, filepath.Join(yName, dName))
	}

	fmt.Printf("Running %s with %s\n", filepath.Join(yName, dName, pName), input)

	path, err := getRunnerPath(year, day, part, input)
	if err != nil {
		return fmt.Errorf("setting up runner: %v", err)
	}

	if err = exec.BinaryAndPrint(path); err != nil {
		return fmt.Errorf("executing runner: %v", err)
	}

	return nil
}

func getRunnerPath(year, day, part int, input string) (string, error) {
	cacheKey := cache.Key(year, day, part, input)

	if cPath, ok := cache.Contains(cacheKey, files.Runner); ok {
		return cPath, nil
	}

	mod, err := currentModulePath()
	if err != nil {
		return "", fmt.Errorf("getting module name: %v", err)
	}

	yName := fmt.Sprintf("%d", year)
	dName := fmt.Sprintf("day%d", day)
	pName := fmt.Sprintf("part%d", part)

	fPaths, err := files.GenTemp(map[string]string{
		files.Runner: runnerTmpl,
		files.Lock:   strconv.FormatBool(false),
		files.Res:    "",
		files.Dur:    time.Duration(math.MaxInt64).String(),
	}, map[string]string{
		"PkgPath":   filepath.Join(mod, yName, "solutions", dName),
		"PkgName":   dName,
		"FuncName":  strings.Replace(pName, "p", "P", 1),
		"InputPath": filepath.Join(yName, "input", dName, input),
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

func currentModulePath() (string, error) {
	out, err := exec.CommandAndCapture("go", "env", "GOMOD")
	if err != nil {
		return "", err
	}
	gomod := string(bytes.TrimSpace(out))
	if gomod == "" || gomod == "/dev/null" || gomod == "NUL" {
		return "", errors.New("not inside a Go module (no go.mod found)")
	}

	data, err := files.Read(gomod)
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
