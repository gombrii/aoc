package main

import (
	"errors"
	"fmt"
	"iter"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gombrii/aoc/internal"
)

const (
	opPuzzleRun    = "puzzle run"
	opPuzzleStatus = "puzzle status"
	opPuzzleLock   = "puzzle lock"
	opPuzzleUnlock = "puzzle unlock"
	opInitModule   = "init -m"
	opInitDay      = "init -d"
	opCacheClear   = "cache clear"
)

const usage = `Usage:
  aoc [puzzle run] <params>
  aoc puzzle {run|status|lock|unlock} -d DAY -p {1|2} [-y YEAR default: {{year}}] [-i INPUT default: input.txt]
  aoc init   {-d DAY [-y YEAR] | -m MODULENAME}
  aoc cache clear
  aoc help [command]

Puzzle commands:
  run          Execute a puzzle (default when no subcommand is given)
  status       Show result and duration of last run of puzzle
  lock         Lock result; future runs error if they deviate; keep fastest duration
  unlock       Unlock result; remember only last run

Project setup:
  init --day       Scaffold solution files for a day
  init --module    Create a new AoC module structure

Misc:
  cache clear      Clear cached runners and metadata
  help             Show this or command-specific help

Concepts:
  • Puzzle = (year, day, part, input).
  • Locked puzzles act like tests: runs must match locked results.
  • Input files must include their extension (e.g. test.txt).
`

type input struct {
	op     string
	year   string
	day    string
	part   string
	input  string
	module string
}

func main() {
	if len(os.Args) == 1 || len(os.Args) == 2 && os.Args[1] == "help"{
		fmt.Println(strings.ReplaceAll(usage, "{{year}}", defaultInput().year))
		return
	}

	in, err := parseInput(os.Args[1:])
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(2)
	}

	if err := validate(in); err != nil {
		fmt.Println("Error:", err)
		os.Exit(2)
	}

	switch in.op {
	case opPuzzleRun:
		err = internal.Run(in.year, in.day, in.part, in.input)
	case opPuzzleStatus:
		err = internal.Status(in.year, in.day, in.part, in.input)
	case opPuzzleLock:
		err = internal.Lock(in.year, in.day, in.part, in.input)
	case opPuzzleUnlock:
		err = internal.Unlock(in.year, in.day, in.part, in.input)
	case opInitDay:
		err = internal.GenDay(in.year, in.day)
	case opInitModule:
		err = internal.GenAoc(in.module)
	case opCacheClear:
		err = internal.ClearCache()
	}

	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(2)
	}
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
			in.year = val
		case "-d", "--day":
			in.day = val
		case "-p", "--part":
			in.part = val
		case "-i", "--input":
			in.input = val
		case "-m", "--module":
			in.module = val
		}
	}

	if len(op) > 0 && op[0] == "init" {
		if in.module != "" {
			op = append(op, "-m")
		} else if in.day != "" {
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
		if _, err := strconv.Atoi(in.year); err != nil {
			return fmt.Errorf("year %q must be a number", in.year)
		}
		if _, err := strconv.Atoi(in.day); err != nil {
			return fmt.Errorf("day %q must be a number", in.day)
		}
		if i, err := strconv.Atoi(in.part); err != nil || i > 2 || i < 1 {
			return fmt.Errorf("part %q must be either 1 or 2", in.part)
		}
	case opInitDay:
		if _, err := strconv.Atoi(in.year); err != nil {
			return fmt.Errorf("year %q must be a number", in.year)
		}
		if _, err := strconv.Atoi(in.day); err != nil {
			return fmt.Errorf("day %q must be a number", in.day)
		}
	case opInitModule:
		if in.module == "" {
			return errors.New("no module name (-m) provided")
		}
	case opCacheClear:
	default:
		return fmt.Errorf("invalid command %q", in.op)
	}

	return nil
}

func defaultInput() input {
	now := time.Now()
	year := now.Year()

	// Before the start of this year's AoC default to previous year.
	if now.Before(time.Date(now.Year(), time.December, 1, 0, 0, 0, 0, time.FixedZone("UTC-5", int(-5*time.Hour)))) {
		year -= 1
	}

	return input{
		op:    opPuzzleRun,
		year:  strconv.Itoa(year),
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
