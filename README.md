[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![GitHub tag](https://img.shields.io/github/v/tag/gombrii/aoc)](https://github.com/gombrii/aoc/tags)
![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/gombrii/aoc)


# Advent of Code CLI (`aoc`)

A simple command-line tool to streamline your **Advent of Code** workflow in **Go**, from scaffolding new solutions to running them with input data.  
It handles project setup, day initialization, and execution of specific parts with minimal friction.
> That was the professional pitch. I wrote the same kind of functionality last year, but as part of my actual Advent of Code repository. I just thought it'd be neat to keep the facilitation of AoC challenges as a standalone app, separate from the actual challenges, hence this.
>
> The goal of this application is to impose upon the user as little as possible and let the user run and write his/her own solutions along with any helper code. The app only provides a simple way to structure and run the solutions without needing to juggle multiple binaries.

## Features

- **Run** any day and part directly from the CLI
- **Initialize** a new AoC module with shared helpers (`input`, `exit`, etc.)
- **Scaffold** new solution files for specific days automatically
- **Auto-detects the current Advent of Code year**, defaulting to the year of the most recently started Advent of Code
- Generates minimal boilerplate code to get started quickly

## Getting started

```shell
# Install aoc
go install github.com/gombrii/aoc@latest

# Create and enter a dir for your AoC project
mkdir myaocproject
cd myaocproject

# Initiate go.mod and shared library
aoc init myaocproject

# Initiate a day
aoc init -d 1

# Now you're ready to write and run your day 1 solution
```

This is how your project structure should (year may differ) look at this point.

```
myaocproject/
├── go.mod
├── shared/
│   ├── exit/error.go
│   └── input/input.go
├── 2025/
│   ├── input/
│   │   └── day1/
│   │       ├── input.txt
│   │       └── test.txt
│   └── solutions/
│       └── day1/
│           ├── part1.go
│           ├── part2.go
│           └── common.go
```
## Usage
```
aoc --day DAY --part PART [--year YEAR] [--input INPUT]
aoc init --day DAY [--year YEAR]
aoc init <module>

Commands:
  run (default)      Run a solution for a given day and part
  init               Initialize either a new AoC module or a new day
```

```shell
# Examples
aoc -d 1 -p 1
# Run part 1 of day 1 for the current year's AoC with puzzle input from input.txt

aoc -y 2023 -d 5 -p 2 -i test
# Run part 2 of day 5, year 2023, using input file test.txt

aoc init mymodule
# Initialize a new Advent of Code module (creates go.mod and shared utilities)

aoc init --day 3
# Create scaffolding for day 3 (solutions + input directories)
```

## How it works
- Each day’s solution lives under `YEAR/solutions/dayX/`.
- Each day’s puzzle input lives under `YEAR/input/dayX/`.
- Each part (Part1, Part2) is implemented as a Go function taking a []byte (puzzle input).
- The aoc init command:
    - If provided a mod name creates a mod file with your system's currently installed Go version as well as a couple of utility packages under `shared/` (can be removed if not needed).
    - If provided a day, eg. `-d 1`, creates the scaffolding for a new day's solutions and input.
- The aoc run command (default):
    1. Builds a temporary runner.
    1. Invokes the corresponding function.
    1. Prints the result and execution duration.

How generated day1 part1 looks like:

```go
package day1

func Part1(data []byte) any {
	// in := input.Lines(data)

	return "NOT IMPLEMENTED!"
}
```

This is where you write your puzzle solution. The puzzle input is provided as raw bytes. To simplify life, the puzzle solution is returned as is, without needing any type conversion, and then printed to the command line. Every initiated day's solution catalogue also gets a `common.go` file, which is simply a convenient place to store code that might be useful for both parts of the challenge.

How running a puzzle looks:
```shell
$ aoc -d 1 -p 1
Running 2024/day1/part1 with input.txt
Res: 2970687
Dur: 304µs
```

"Res" is whatever was returned from the PartX function and "Dur" is the time measured from the moment the PartX function was called to the moment after it returned. The loading of the puzzle input file data happens before time starts recording. No printing (for debug purposes or otherwise) will interfere with anything, so feel free to do so. Prints will simply appear between "Running year/dayX/partX with X.txt" and the "Res" and "Dur" statements.

Every initiated day's input catalogue gets two empty text files (`input.txt`, `test.txt`) into which you copy paste that day's puzzle input. Run with `-i` to run the puzzle with a specific input file, eg. `aoc -d 1 -p 2 -i test` to run with `test.txt`. The default is `input`, which runs the puzzle with `input.txt`. If the challenge presents more than one test input, simply create more test input files.

### Utilities
At least in my mind, Advent of Code solutions are quick and dirty, thus don't need proper code hygiene. To achieve that, two helper packages are included when initiating the module:
- shared/input — for parsing input data into common formats (Lines, String, Matrix, etc.)
- shared/exit — for exiting quickly in case of error (exit.If(err), exit.PanicIf(err))

## Roadmap
- I will add unit tests to be able to guarantee stability
- I might add a couple of more utilities for parsing input and maybe an iterator.
- One downside of this whole project is that the nature of generating temporary runners prevents the OS from applying common optimizations, such as caches. Currently, an AoC solution runs about 60-100% slower using this tool as compared to running it using the same binary again and again. Hopefully I can find some solutions in terms of warmup and preemptive memory allocations to mitigate this.