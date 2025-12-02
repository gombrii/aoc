package app_test

import (
	"runtime"
	"slices"
	"strings"
	"testing"

	"github.com/gombrii/aoc/internal/app"
)

type record map[string][]any

func (r record) save(args ...any) {
	pc, _, _, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)
	r[fn.Name()[strings.LastIndex(fn.Name(), ".")+1:]] = args
}

func (r record) called(fn string) ([]any, bool) {
	res, ok := r[fn]
	return res, ok
}

type commands struct {
	record record
}

func (c *commands) Run(year, day, part int, input string) error {
	c.record.save(year, day, part, input)
	return nil
}
func (c *commands) Status(year, day, part int, input string) error {
	c.record.save(year, day, part, input)
	return nil
}
func (c *commands) Lock(year, day, part int, input string) error {
	c.record.save(year, day, part, input)
	return nil
}
func (c *commands) Unlock(year, day, part int, input string) error {
	c.record.save(year, day, part, input)
	return nil
}
func (c *commands) GenDay(year, day int) error {
	c.record.save(year, day)
	return nil
}
func (c *commands) GenAoc(module string) error {
	c.record.save(module)
	return nil
}
func (c *commands) Check() error {
	c.record.save()
	return nil
}
func (c *commands) ClearCache() error {
	c.record.save()
	return nil
}

func TestSuccessful(t *testing.T) {
	for name, params := range map[string]struct {
		args   string
		called string
		with   []any
	}{
		"Run": {
			args:   "-d 1 -p 1",
			called: "Run",
			with:   []any{2025, 1, 1, "input.txt"},
		},
		"Run other year": {
			args:   "-d 1 -y 2023 -p 1",
			called: "Run",
			with:   []any{2023, 1, 1, "input.txt"},
		},
		"Run other input": {
			args:   "-d 1 -i test.txt -p 1",
			called: "Run",
			with:   []any{2025, 1, 1, "test.txt"},
		},
		"Run other year and input": {
			args:   "-d 1 -i test.txt -p 1 -y 2023",
			called: "Run",
			with:   []any{2023, 1, 1, "test.txt"},
		},
		"Status": {
			args:   "puzzle status -d 1 -p 1",
			called: "Status",
			with:   []any{2025, 1, 1, "input.txt"},
		},
		"Status other year": {
			args:   "puzzle status -d 1 -y 2023 -p 1",
			called: "Status",
			with:   []any{2023, 1, 1, "input.txt"},
		},
		"Status other input": {
			args:   "puzzle status -d 1 -i test.txt -p 1",
			called: "Status",
			with:   []any{2025, 1, 1, "test.txt"},
		},
		"Status other year and input": {
			args:   "puzzle status -d 1 -i test.txt -p 1 -y 2023",
			called: "Status",
			with:   []any{2023, 1, 1, "test.txt"},
		},
		"Lock": {
			args:   "puzzle lock -d 1 -p 1",
			called: "Lock",
			with:   []any{2025, 1, 1, "input.txt"},
		},
		"Lock other year": {
			args:   "puzzle lock -d 1 -y 2023 -p 1",
			called: "Lock",
			with:   []any{2023, 1, 1, "input.txt"},
		},
		"Lock other input": {
			args:   "puzzle lock -d 1 -i test.txt -p 1",
			called: "Lock",
			with:   []any{2025, 1, 1, "test.txt"},
		},
		"Lock other year and input": {
			args:   "puzzle lock -d 1 -i test.txt -p 1 -y 2023",
			called: "Lock",
			with:   []any{2023, 1, 1, "test.txt"},
		},
		"Unlock": {
			args:   "puzzle unlock -d 1 -p 1",
			called: "Unlock",
			with:   []any{2025, 1, 1, "input.txt"},
		},
		"Unlock other year": {
			args:   "puzzle unlock -d 1 -y 2023 -p 1",
			called: "Unlock",
			with:   []any{2023, 1, 1, "input.txt"},
		},
		"Unlock other input": {
			args:   "puzzle unlock -d 1 -i test.txt -p 1",
			called: "Unlock",
			with:   []any{2025, 1, 1, "test.txt"},
		},
		"Unlock other year and input": {
			args:   "puzzle unlock -d 1 -i test.txt -p 1 -y 2023",
			called: "Unlock",
			with:   []any{2023, 1, 1, "test.txt"},
		},
		"GenDay": {
			args:   "init -d 1",
			called: "GenDay",
			with:   []any{2025, 1},
		},
		"GenDay other year": {
			args:   "init -d 1 -y 2023",
			called: "GenDay",
			with:   []any{2023, 1},
		},
		"GenAoc": {
			args:   "init -m mymodule",
			called: "GenAoc",
			with:   []any{"mymodule"},
		},
		"Check": {
			args:   "check",
			called: "Check",
			with:   []any{},
		},
		"ClearCache": {
			args:   "cache clear",
			called: "ClearCache",
			with:   []any{},
		},
	} {
		t.Run(name, func(t *testing.T) {
			cmd := commands{record: record{}}

			err := app.Start(&cmd, strings.Split(params.args, " ")...)
			if err != nil {
				t.Errorf("Start returned error: %v", err)
			}

			args, ok := cmd.record.called(params.called)
			if !ok {
				t.Fatal(params.called, "was not called")
			}
			if !slices.Equal(args, params.with) {
				t.Error("\nGot:", args, "\nWant:", params.with)
			}
		})
	}
}

func TestError(t *testing.T) {
	for name, params := range map[string]struct {
		args string
	}{
		"Run missing part": {
			args: "-p 1",
		},
		"Run missing day": {
			args: "-d 1",
		},
		"Run with module": {
			args: "-d 1 -p 1 -m mymodule",
		},
		"Status missing part": {
			args: "puzzle status -d 1",
		},
		"Status missing day": {
			args: "puzzle status -p 1",
		},
		"Status with module": {
			args: "puzzle status -d 1 -p 1 -m mymodule",
		},
		"Lock missing part": {
			args: "puzzle lock -d 1",
		},
		"Lock missing day": {
			args: "puzzle lock -p 1",
		},
		"Lock with module": {
			args: "puzzle lock -d 1 -p 1 -m mymodule",
		},
		"Unlock missing part": {
			args: "puzzle unlock -d 1",
		},
		"Unlock missing day": {
			args: "puzzle unlock -p 1",
		},
		"Unlock with module": {
			args: "puzzle unlock -d 1 -p 1 -m mymodule",
		},
		"GenDay missing day": {
			args: "init",
		},
		"GenDay with module": {
			args: "init -d 1 -m mymodule",
		},
		"GenDay with part": {
			args: "init -d 1 -p 1",
		},
		"GenDay with input": {
			args: "init -d 1 -i test.txt",
		},
		"GenAoc missing module": {
			args: "init",
		},
		"GenAoc with year": {
			args: "init -m mymodule -y 2023",
		},
		"GenAoc with day": {
			args: "init -m mymodule -d 1",
		},
		"GenAoc with part": {
			args: "init -m mymodule -p 1",
		},
		"GenAoc with input": {
			args: "init -m mymodule -i test.txt",
		},
		"year wrong format": {
			args: "-y senap",
		},
		"day wrong format": {
			args: "-d senap",
		},
		"part wrong format": {
			args: "-p senap",
		},
		"input with digit": {
			args: "-i 2",
		},
		"input missing file ending": {
			args: "-i test",
		},
		"module param without init": {
			args: "-m mymodule",
		},
		"unknown command": {
			args: "start",
		},
		"unknown second command": {
			args: "cache start",
		},
		"missing arg": {
			args: "-d -p 1",
		},
		"stray arg": {
			args: "-d 1 2 -p 1",
		},
	} {
		t.Run(name, func(t *testing.T) {
			cmd := commands{record: record{}}

			err := app.Start(&cmd, strings.Split(params.args, " ")...)
			if err == nil {
				t.Error("Start didn't return an error")
			}
		})
	}
}
