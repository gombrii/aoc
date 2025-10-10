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
	opRun     = "run"
	opInit    = "init"
	opInitAoc = "initAoc"
	opInitDay = "initDay"
)

const usage = `Usage:
  aoc --day DAY --part PART [--year YEAR default: {{year}}] [--input INPUT default: input]
  aoc init <module>
  aoc init --day DAY [--year YEAR default: {{year}}]

Commands:
  run (default)      Run a solution for a given day and part
  init               Initialize either a new AoC module or a new day

Examples:
  aoc -d 1 -p 1                    # Run part 1 of day 1 for the default year ({{year}})
  aoc -y 2023 -d 5 -p 2 -i test    # Run part 2 of day 5 year 2023 using input file "test.txt"
  aoc init mymodule                # Initialize a new AoC module named mymodule (creates project structure)
  aoc init --day 3                 # Initialize solution files for day 3 for the default year ({{year}})`

type input struct {
	op     string
	year   string
	day    string
	part   string
	input  string
	module string
}

func main() {
	if len(os.Args) == 1 {
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
	}

	if err != nil {
		fmt.Println("Error:", err)
	}
}

func parseInput(args []string) (input, error) {
	in := defaultInput()
	i := 1

	if args[i] == "init" {
		i++
		in.op = opInit
	}

	for param, val := range argIter(args[i:]) {
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

func argIter(args []string) iter.Seq2[string, string] {
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
