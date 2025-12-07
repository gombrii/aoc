[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![GitHub tag](https://img.shields.io/github/v/tag/gombrii/aoc)](https://github.com/gombrii/aoc/tags)
![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/gombrii/aoc)


# Advent of Code CLI (`aoc`)

A command-line tool to streamline your **Advent of Code** workflow in **Go**, from scaffolding, running and submitting solutions, to testing and refactoring, chasing better solutions and execution times.  
It handles project setup, day initialization, and execution of specific parts with minimal friction. 

The goal of this app is at its core to simply create a helpful catalogue structure for Advent of Code repos with a method to easily run the solutions without having to juggle multiple binaries. None of the advanced features are needed, but they're there for those who benefit from them.

## Features

- **Run** any puzzle (part 1 and 2) of any day and year directly from the CLI
- **Initialize** a new AoC module with shared helpers (`parse`, `exit`, etc.)
- **Scaffold** new solution files for a given day automatically, with funcs _Part1_ and _Part2_ for you to implement
- **Auto-detects the current Advent of Code year**, defaulting to the year of the most recently started Advent of Code
- **Login** to your AoC account from CLI to pull puzzle inputs straight from server and submit solutions
- **Lock in** correct answers and use aoc as a test tool to polish your puzzle solutions.
- **Chase execution times** with cached runs

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
│   ├── exit/...
│   ├── parse/...
│   └── render/...
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
aoc -d DAY -p {1|2} [-y YEAR] [-i INPUT]
aoc init {-d DAY [-y YEAR] | -m MODULENAME}
aoc submit 
aoc login -s SESSION 
aoc check 
aoc cache clear
aoc help [-v]
aoc version

Run and submit:
  aoc              Run a puzzle solution
  submit           Submit the result of your last run puzzle (requires login)

Project setup:
  init --day       Scaffold solution files for a new day (pull puzzle input from server if logged in)
  init --module    Create a new AoC module structure

Misc:
  login            Enables pulling of puzzle input and submission of solutions to server
  check            Run all locked puzzles to verify results
  cache clear      Delete all data created and kept by aoc
  help             Show this help
  version          Show installed aoc version
```

## How it works
- Each day's solution lives under `YEAR/solutions/dayX/`.
- Each day's puzzle input lives under `YEAR/input/dayX/`.
- Each part (Part1, Part2) is implemented as a Go function taking a []byte (puzzle input).
- The aoc init command:
    - If provided a mod name, eg. `-m mymodule`, creates a mod file with your system's currently installed Go version as well as a couple of utility packages under `shared/` (not mandatory).
    - If provided a day, eg. `-d 1`, creates the scaffolding for a new day's solutions and input for the given year. If no year, eg. `-y 2023`, is provided the default is the year during which the last Advent of Code started. This means that the default year the majority of time is the previous year. On Dec 1 00:00 UTC-5 when the current year's AoC is released the default year flips over to the current year. If logged in as a user, this also pulls puzzle inputs from the server.
- The aoc run command (implicit):
    1. Invokes the corresponding function.
    1. Prints the result and execution duration.
    1. Stores result and execution duration in cache.

How generated day1 part1 looks like:
```go
package day1

func Part1(data []byte) any {
	// in := parse.Lines(data)

	return "NOT IMPLEMENTED!"
}
```

This is where you write your puzzle solution. The puzzle input is provided as raw bytes (that can be parsed using `shared/parse`). To simplify life, the puzzle solution can be returned as is, without needing any type conversion, after which it's printed to the command line. Every initiated day's solution catalogue, apart from `part1.go` and `part2.go`, also gets a `common.go` file, which is simply a convenient place to store code that might be useful for both parts of the puzzle.

How running a puzzle looks:
```shell
$ aoc -d 1 -p 1
Running 2024/day1/part1 with input.txt
Res: 2970687
Dur: 304µs
```

`Res` is whatever was returned from the PartX function and `Dur` is the time measured from the moment the PartX function was called to the moment after it returned. The loading of the puzzle input file data happens before time starts recording. Prints in the puzzle solution (for debug purposes or otherwise) will not interfere with anything, so feel free to use them. Print outputs will simply appear between "Running year/dayX/partX with X.txt" and the `Res` and `Dur` statements.

Every initiated day's input catalogue gets two text files (`input.txt`, `test.txt`). If you are logged in as a user these are pre-filled with the puzzle and example data from the server. Otherwise they are empty for you to paste into. Run a puzzle with `-i` to select a specific input file, eg. `aoc -d 1 -p 2 -i test.txt`. The default is `input.txt`. If the puzzle presents more than one example input, simply create more input files.

### Login
You can login by running `aoc login` and providing your personal AoC session token. It can be found in your web browser's dev tools when logged in to your Advent of Code account. Logging into the aoc CLI lets you interact with the server. When logged in puzzle input and example input are automatically pulled from the server when you initialize a new day.

Being logged in also enables you to submit your puzzle solutions right from the terminal by running `aoc submit`. This submits the previously run puzzle result and outputs the result. Submitting a correct result will also trigger aoc to lock down the result so that future runs of the solution will error if the result differs from the correct one. Moreover, duration will also be continuously updated compared to your fastest execution time for that puzzle since it got locked.

This gives you the opportunity to refactor and polish your solution while getting clear feedback on improved performance and if something alters the result. Effectively your puzzle solution turns into a simple unit- and performance test testing itself. 

### Checking
The `check` command will run all locked puzzles at once and verify their results. Only puzzles which produce correct results get a golden star (*).

```shell
$ aoc check
2024/day1/part1  x      # produced wrong output
2024/day2/part1  *      # produced correct output
2024/day3/part1  error  # panicked, exited early or even didn't compile
2024/day4/part1  *
2024/day5/part1  *
2024/day6/part1  *
2024/day7/part1  *
2024/day8/part1  *
2024/day9/part1  *
2024/day10/part1 *
```

### Utilities
At least in my mind, Advent of Code solutions are quick and dirty, thus don't need proper code hygiene. To achieve that and other things, a few helper packages are included when initiating the module:
- shared/parse — for parsing input data into common formats (Lines, String, Matrix, etc.)
- shared/exit — for exiting quickly in case of error (`exit.If(err)`, `exit.PanicIf(err)`)
- shared/render — for visualizing data, such as printing and animating 2D grids

```
shared/
├── exit/error.go
├── parse/input.go
└── render/
    ├── string.go
    └── print.go
```

### Cache
Aoc uses the OS's default caching location to store data. When aoc runs a puzzle it generates a binary under the hood which is stored in cache for performance reasons. That's why the first couple of runs of a puzzle tend to be slower. The cache also stores results and execution times for each puzzle and keeps track of which puzzles are locked. Configuration data such as your session token is also stored here. Clearing the cache removes every trace of it from your computer and resets aoc's memory. 

## Author's notes
### Feature additions
- Most planned changes are related to code hygiene, among which are:
  - Refactoring the layout of the app's internals to enable better testing and simpler flows
  - Switching over to the stdlib flag package for user input
- _Any feedback_ is welcome, both in terms of bugs or opinions on features or utilities.
