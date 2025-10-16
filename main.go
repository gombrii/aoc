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
	opRun        = "run"
	opInit       = "init"
	opInitAoc    = "initAoc"
	opInitDay    = "initDay"
	opClearCache = "clearCache"
	opStatus     = "status"
	opLock       = "lock"
	opUnlock     = "unlock"
)

const usage = `Usage:
  aoc [run] --day DAY --part {1|2} [--year YEAR default: {{year}}] [--input INPUT default: input]
  aoc init --day DAY [--year YEAR default: {{year}}]
  aoc init <module>
  aoc cache clear
  aoc status --day DAY --part {1|2} [--year YEAR default: {{year}}] [--input INPUT default: input]
  aoc lock   --day DAY --part {1|2} [--year YEAR default: {{year}}] [--input INPUT default: input]
  aoc unlock --day DAY --part {1|2} [--year YEAR default: {{year}}] [--input INPUT default: input]
  aoc help

Commands:
  run (default)         Run solution for a given day and part
  init                  Initialize an AoC module or a new day
  status                Show cached status for a specific puzzle
  lock                  Lock the result for a specific puzzle
  unlock                Unlock the result for a specific puzzle
  cache clear           Clear puzzle solutions cache

Examples:
  aoc -d 1 -p 1                    # Run part 1 of day 1 with "input.txt" for the default year ({{year}})
  aoc -y 2023 -d 5 -p 2 -i test    # Run part 2 of day 5 year 2023 using input file "test.txt"
  aoc init mymodule                # Initialize a new AoC module named mymodule (creates project structure)
  aoc init --day 3                 # Initialize solution files for day 3 for the default year ({{year}})
  aoc status -d 1 -p 2             # Show status for solution {{year}} day 1 part 2 with input.txt
  aoc lock -d 2 -p 1               # Lock the result for {{year}} day 2 part 1 with input.txt
  aoc unlock -d 2 -p 1             # Unlock the result for {{year}} day 1 part 2 with input.txt
  aoc cache clear                  # Clear cache (runners, results, execution times, locks)
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
	if len(os.Args) == 1 || len(os.Args) == 2 && (os.Args[1] == "h" || os.Args[1] == "help" || os.Args[1] == "-h" || os.Args[1] == "--help") {
		fmt.Println(strings.ReplaceAll(usage, "{{year}}", defaultInput().year))
		return
	}

	in, err := parseInput(os.Args)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if err := validate(in); err != nil {
		fmt.Println("Error:", err)
		return
	}

	switch in.op {
	case opRun:
		err = internal.Run(in.year, in.day, in.part, in.input)
	case opInitDay:
		err = internal.GenDay(in.year, in.day)
	case opInitAoc:
		err = internal.GenAoc(in.module)
	case opClearCache:
		err = internal.ClearCache()
	case opStatus:
		err = internal.Status(in.year, in.day, in.part, in.input)
	case opLock:
		err = internal.Lock(in.year, in.day, in.part, in.input)
	case opUnlock:
		err = internal.Unock(in.year, in.day, in.part, in.input)
	}

	if err != nil {
		fmt.Println("Error:", err)
	}
}

func parseInput(args []string) (input, error) {
	in := defaultInput()
	i := 1

	switch args[i] {
	case opInit, opStatus, opLock, opUnlock:
		in.op = args[i]
		i++
	case opRun:
		i++
	case "cache":
		if len(args) > 2 && args[i+1] == "clear" {
			in.op = opClearCache
			return in, nil
		}
	}

	for param, val := range paramVals(args[i:]) {
		switch param {
		case "-y", "--year":
			in.year = val
		case "-d", "--day":
			in.day = val
			if in.op == opInit {
				in.op = opInitDay
			}
		case "-p", "--part":
			in.part = val
		case "-i", "--input":
			in.input = val
		default:
			if in.op != opInit {
				if strings.HasPrefix(param, "-") {
					return input{}, fmt.Errorf("unknown parameter %q", param)
				}
				return input{}, fmt.Errorf("loose argument %q", param)
			}
			in.module = param
			in.op = opInitAoc
		}
	}

	return in, nil
}

func validate(input input) error {
	switch input.op {
	case opInitAoc:
		if input.module == "" {
			return errors.New("no module name provided")
		}
	case opRun:
		if input.part == "" {
			return errors.New("no part (-p) provided")
		}
		if i, err := strconv.Atoi(input.part); err != nil || i > 2 || i < 1 {
			return fmt.Errorf("part %q must be either 1 or 2", input.part)
		}
		fallthrough
	case opInitDay:
		if _, err := strconv.Atoi(input.year); err != nil {
			return fmt.Errorf("year %q must be a number", input.year)
		}
		if input.day == "" {
			return errors.New("no day (-d) provided")
		}
		if _, err := strconv.Atoi(input.day); err != nil {
			return fmt.Errorf("day %q must be a number", input.day)
		}
	case opInit:
		return errors.New("missing init argument")
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
		op:    opRun,
		year:  strconv.Itoa(year),
		input: "input",
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
