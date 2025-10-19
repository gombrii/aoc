[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![GitHub tag](https://img.shields.io/github/v/tag/gombrii/aoc)](https://github.com/gombrii/aoc/tags)
![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/gombrii/aoc)


# Advent of Code CLI (`aoc`)

A command-line tool to streamline your **Advent of Code** workflow in **Go**, from scaffolding new solutions to running them with input data, chasing better solutions and execution times.  
It handles project setup, day initialization, and execution of specific parts with minimal friction.
> Last year I wrote something similar to this, albeit much simpler, in my AoC repo. Now I thought it'd be neat to keep the facilitation of AoC challenges as a standalone app, separate from the actual challenges, hence this.
>
> The goal of this application is to impose upon the user as little as possible and let the user run and write his/her own solutions along with any helper code. This app, in its core, only provides a simple way to structure and run the solutions without needing to juggle multiple binaries.

## Features

- **Run** any day and part directly from the CLI
- **Initialize** a new AoC module with shared helpers (`parse`, `exit`, etc.)
- **Scaffold** new solution files for specific days automatically
- **Auto-detects the current Advent of Code year**, defaulting to the year of the most recently started Advent of Code
- **Generate** minimal boilerplate code to get started quickly
- **Lock down** correct answers and use aoc as a test tool to find better puzzle solutions.
- **Chase fast execution times** with execution duration memory

## Getting started

```shell
# Install aoc
go install github.com/gombrii/aoc@latest

# Create and enter a dir for your AoC project
mkdir myaocproject
cd myaocproject

# Initiate go.mod and shared library
aoc init -m myaocproject

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
│   └── parse/input.go
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
aoc [puzzle run] <params>
aoc puzzle {run|status|lock|unlock} -d DAY -p {1|2} [-y YEAR default: Curr AoC year] [-i INPUT default: input.txt]
aoc init   {-d DAY [-y YEAR default: Curr AoC year]] | -m MODULENAME}
aoc check 
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
  check            Run all locked puzzles to verify results
  cache clear      Clear cached runners and metadata
  help             Show this or command-specific help
```

## How it works
- Each day’s solution lives under `YEAR/solutions/dayX/`.
- Each day’s puzzle input lives under `YEAR/input/dayX/`.
- Each part (Part1, Part2) is implemented as a Go function taking a []byte (puzzle input).
- The aoc init command:
    - If provided a mod name, eg. `-m mymodule`, creates a mod file with your system's currently installed Go version as well as a couple of utility packages under `shared/` (can be removed if not needed).
    - If provided a day, eg. `-d 1`, creates the scaffolding for a new day's solutions and input for the given year. If no year, eg. `-y 2023`, is provided the default is the year during which the last Advent of Code started. This means that the default year the majority of time is the previous year. On Dec 1 00:00 UTF-5 when the current years AoC is released the default year flips over to the current year.
- The aoc puzzle run command (default):
    1. Builds a temporary runner or fetches from cache.
    1. Invokes the corresponding function.
    1. Prints the result and execution duration.

How generated day1 part1 looks like:

```go
package day1

func Part1(data []byte) any {
	// in := parse.Lines(data)

	return "NOT IMPLEMENTED!"
}
```

This is where you write your puzzle solution. The puzzle input is provided as raw bytes. To simplify life, the puzzle solution can be returned as is, without needing any type conversion, after which it's printed to the command line. Every initiated day's solution catalogue also gets a `common.go` file, which is simply a convenient place to store code that might be useful for both parts of the challenge.

How running a puzzle looks:
```shell
$ aoc -d 1 -p 1
Running 2024/day1/part1 with input.txt
Res: 2970687
Dur: 304µs
```

`Res` is whatever was returned from the PartX function and `Dur` is the time measured from the moment the PartX function was called to the moment after it returned. The loading of the puzzle input file data happens before time starts recording. Printing (for debug purposes or otherwise) will not interfere with anything, so feel free to do so. Prints will simply appear between "Running year/dayX/partX with X.txt" and the `Res` and `Dur` statements.

Every initiated day's input catalogue gets two empty text files (`input.txt`, `test.txt`) into which you copy paste that day's puzzle input. Run with `-i` to run the puzzle with a specific input file, eg. `aoc -d 1 -p 2 -i test.txt`. The default is `input.txt`. If the challenge presents more than one test input, simply create more test input files.

### Locking
Aoc remembers the results and durations of each solution's last run. Locking a solution does two things. Firstly, it locks down the result for that solution so that it's not overwritten and instead errors if the solution's result deviates from the locked down result. Secondly, it sees to that the duration for the solution only updates if it's shorter than the shortest recorded duration.

The typical usecase for locking a solution is after a correct result has been achieved. Often the first solution is sloppy or naïve. When the solution is locked it gives you the oppertunity to experiment and polish your solution and compare results, while getting clear feedback when something has gone wrong. Effectively you've turned your puzzle solution into a simple unit test and performance test testing itself. 

### Checking
The `check` command will run all locked puzzles at once and check their results. Only puzzles which produce correct results get a golden star (*).

Eg.
```shell
$ aoc check
2024/day1/part1  x      # produced wrong output
2024/day2/part1  *
2024/day3/part1  error  # paniced or exited early
2024/day4/part1  *
2024/day5/part1  *
2024/day6/part1  *
2024/day7/part1  *
2024/day8/part1  *
2024/day9/part1  *
2024/day10/part1 *
```

### Utilities
At least in my mind, Advent of Code solutions are quick and dirty, thus don't need proper code hygiene. To achieve that, two helper packages are included when initiating the module:
- shared/parse — for parsing input data into common formats (Lines, String, Matrix, etc.)
- shared/exit — for exiting quickly in case of error (exit.If(err), exit.PanicIf(err))

## Author's notes
### Feature additions
- I might add a couple of more utilities for parsing input and maybe an iterator.
- I'll probably add one or two printing utilities for matrices. Writing one's own matrix printing loop isn't the most fun thing in the world just to make a debug print.
### Technical debt
- Unit tests to be able to guarantee stability.
- Better app input (cli args). Right now it works, but it's kind of a mess and needs better input validation.