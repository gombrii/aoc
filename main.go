package main

import (
	"errors"
	"fmt"
	"iter"
	"os"
	"strconv"
	"time"

	"github.com/gombrii/aoc/internal"
)

const (
	opRun     = "run"
	opInitAoc = "initAoc"
	opInitDay = "initDay"
)

// TODO: Update usage
const usage = `Usage: aoc [--year YEAR] --day DAY --part PART [--input INPUT]

Options:
  --year YEAR, -y YEAR [default: %s]
  --day DAY, -d DAY
  --part PART, -p PART
  --input INPUT, -i INPUT [default: input]`

type input struct {
	op    string
	year  string
	day   string
	part  string
	input string
}

func main() {
	if len(os.Args) == 1 {
		fmt.Printf(usage, defaultInput().year)
		return
	}

	in, err := parseInput(os.Args)
	if err != nil {
		fmt.Println("Error:", err)
	}

	if err := validate(in); err != nil {
		fmt.Println("Error:", err)
	}

	switch in.op {
	case opRun:
		err = internal.Run(in.year, in.day, in.part, in.input)
	case opInitDay:
		err = internal.GenDay(in.year, in.day)
	case opInitAoc:
		err = internal.GenAoc()
	}

	if err != nil {
		fmt.Println("Error:", err)
	}
}

func parseInput(args []string) (input, error) {
	in := defaultInput()
	i := 1

	if args[i] == "init" {
		if len(args) == 2 {
			in.op = opInitAoc
			return in, nil
		}
		i++
		in.op = opInitDay
	}

	//TODO: Better feedback
	if (len(args)-i)%2 != 0 {
		return input{}, errors.New("uneven number of arguments")
	}

	for param, val := range argIter(args, i) {
		switch param {
		case "-y", "--year":
			in.year = val
		case "-d", "--day":
			in.day = val
		case "-p", "--part":
			in.part = val
		case "-i", "--input":
			in.input = val
		default:
			return input{}, fmt.Errorf("unknown parameter (%s)", param)
		}
	}

	return in, nil
}

func validate(input input) error {
	if _, err := strconv.Atoi(input.year); err != nil {
		return fmt.Errorf("year (%s) must be a number", input.year)
	}
	if _, err := strconv.Atoi(input.day); err != nil {
		return fmt.Errorf("day (%s) must be a number", input.day)
	}
	if i, err := strconv.Atoi(input.part); err != nil || i > 2 || i < 1 {
		return fmt.Errorf("part (%s) must be either 1 or 2", input.part)
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

func argIter(args []string, n int) iter.Seq2[string, string] {
	return func(yield func(string, string) bool) {
		for i := n; i < len(args); i += 2 {
			//TODO: Maybe return empty second value if uneven length
			if !yield(args[i], args[i+1]) {
				return
			}
		}
	}
}
