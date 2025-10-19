package main

import (
	"errors"
	"fmt"
	"iter"
	"strconv"
	"strings"
	"time"
)

const (
	opPuzzleRun    = "puzzle run"
	opPuzzleStatus = "puzzle status"
	opPuzzleLock   = "puzzle lock"
	opPuzzleUnlock = "puzzle unlock"
	opInitModule   = "init -m"
	opInitDay      = "init -d"
	opCacheClear   = "cache clear"
	opCheck        = "check"
)

const usage = `Usage:
  aoc [puzzle run] <params>
  aoc puzzle {run|status|lock|unlock} -d DAY -p {1|2} [-y YEAR default: {{year}}] [-i INPUT default: input.txt]
  aoc init   {-d DAY [-y YEAR default: {{year}}] | -m MODULENAME}
  aoc check 
  aoc cache clear
  aoc help [command]

Puzzle commands:
  run          Execute a puzzle (default when no command is given)
  status       Show last/locked result and last/best duration of of puzzle
  lock         Lock result -> future runs error if they deviate, remembers fastest duration
  unlock       Unlock result -> remember only last run

Project setup:
  init --day       Scaffold solution files for a new day
  init --module    Create a new AoC module structure

Misc:
  check            Run all locked puzzles to verify results
  cache clear      Clear cached runners and metadata
  help             Show this or command-specific help

Concepts:
  • Puzzle = (year, day, part, input).
  • Locked puzzles act like tests: runs must match locked results.
  • Checking makes it easy to verify function of modified shared code.
`

type commands struct {
	Run        func(year, day, part int, input string) error
	Status     func(year, day, part int, input string) error
	Lock       func(year, day, part int, input string) error
	Unlock     func(year, day, part int, input string) error
	GenDay     func(year, day int) error
	GenAoc     func(module string) error
	Check      func() error
	ClearCache func() error
}

type input struct {
	op     string
	year   int
	day    int
	part   int
	input  string
	module string
}

func start(cmd commands, args ...string) error {
	if len(args) == 0 || len(args) == 1 && args[0] == "help" {
		fmt.Println(strings.ReplaceAll(usage, "{{year}}", fmt.Sprint(defaultInput().year)))
		return nil
	}

	in, err := parseInput(args)
	if err != nil {
		return err
	}

	if err := validate(in); err != nil {
		return err
	}

	switch in.op {
	case opPuzzleRun:
		err = cmd.Run(in.year, in.day, in.part, in.input)
	case opPuzzleStatus:
		err = cmd.Status(in.year, in.day, in.part, in.input)
	case opPuzzleLock:
		err = cmd.Lock(in.year, in.day, in.part, in.input)
	case opPuzzleUnlock:
		err = cmd.Unlock(in.year, in.day, in.part, in.input)
	case opInitDay:
		err = cmd.GenDay(in.year, in.day)
	case opInitModule:
		err = cmd.GenAoc(in.module)
	case opCacheClear:
		err = cmd.ClearCache()
	case opCheck:
		err = cmd.Check()
	}

	if err != nil {
		return err
	}

	return nil
}

func parseInput(args []string) (input, error) {
	in := defaultInput()
	op := make([]string, 0)

	i := 0
	for j, arg := range args {
		i = j
		if strings.HasPrefix(arg, "-") {
			break
		}
		op = append(op, arg)
	}

	for param, val := range paramVals(args[i:]) {
		switch param {
		case "-y", "--year":
			numVal, err := strconv.Atoi(val)
			if err != nil {
				return input{}, fmt.Errorf("year (-y) %q must be a number", val)
			}
			in.year = numVal
		case "-d", "--day":
			numVal, err := strconv.Atoi(val)
			if err != nil {
				return input{}, fmt.Errorf("day (-d) %q must be a number", val)
			}
			in.day = numVal
		case "-p", "--part":
			numVal, err := strconv.Atoi(val)
			if err != nil {
				return input{}, fmt.Errorf("part (-p) %q must be a number {1|2}", val)
			}
			in.part = numVal
		case "-i", "--input":
			if val == "" {
				return input{}, errors.New("input (-i) requires a value")
			}
			in.input = val
		case "-m", "--module":
			if val == "" {
				return input{}, errors.New("module (-m) requires a value")
			}
			in.module = val
		default:
			if strings.HasPrefix(param, "-") {
				return input{}, fmt.Errorf("unknown flag %q", param)
			}
			return input{}, fmt.Errorf("stray argument %q", param)
		}
	}

	if len(op) > 0 && op[0] == "init" {
		if in.module != "" {
			op = append(op, "-m")
		} else if in.day != 0 {
			op = append(op, "-d")
		} else {
			return input{}, errors.New("ambiguous call to init, no -d or -m arguments passed")
		}
	}

	if len(op) > 0 {
		in.op = strings.Join(op, " ")
	}

	return in, nil
}

func validate(in input) error {
	switch in.op {
	case opPuzzleRun, opPuzzleStatus, opPuzzleLock, opPuzzleUnlock:
		if in.year == 0 {
			return errors.New("year (-y) is required")
		}
		if in.day == 0 {
			return errors.New("day (-d) is required")
		}
		if in.part > 2 || in.part < 1 {
			return errors.New("part (-p) not set to valid value {1|2}")
		}
	case opInitDay:
		if in.year == 0 {
			return errors.New("year (-y) is required")
		}
		if in.day == 0 {
			return errors.New("day (-d) is required")
		}
	case opInitModule:
		if in.module == "" {
			return errors.New("module name (-m) is required")
		}
	case opCacheClear, opCheck:
		// No args
	default:
		return fmt.Errorf("invalid command %q", in.op)
	}

	return nil
}

func defaultInput() input {
	now := time.Now()
	year := now.Year()

	// Before the start of this year's AoC default to previous year.
	ny, _ := time.LoadLocation("America/New_York")
	if now.Before(time.Date(now.Year(), time.December, 1, 0, 0, 0, 0, ny)) {
		year -= 1
	}

	return input{
		op:    opPuzzleRun,
		year:  year,
		input: "input.txt",
	}
}

func paramVals(args []string) iter.Seq2[string, string] {
	return func(yield func(string, string) bool) {
		for i := 0; i < len(args); i += 2 {
			var second string
			if i+1 < len(args) {
				second = args[i+1]
			}
			if !yield(args[i], second) {
				return
			}
		}
	}
}
