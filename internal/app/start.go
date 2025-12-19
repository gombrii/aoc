package app

import (
	"flag"
	"fmt"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

const (
	opStatus  = "status"
	opLock    = "lock"
	opUnlock  = "unlock"
	opInit    = "init"
	opCache   = "cache"
	opClear   = "clear"
	opCheck   = "check"
	opLogin   = "login"
	opSubmit  = "submit"
	opVersion = "version"
	opHelp    = "help"
)

const usage = `Usage:
  aoc -d DAY -p {1|2} [-y YEAR def: {{year}}] [-i INPUT def: input.txt]
  aoc init {-d DAY [-y YEAR def: {{year}}] | -m MODULENAME}
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
  version          Show installed aoc version`

const usageAppendix = `
Legacy commands:
  aoc {status|lock|unlock} -d DAY -p {1|2} [-y YEAR current default: {{year}}] [-i INPUT default: input.txt]

  status       Show last/locked result and last/best duration of of puzzle
               (this is also shown after running the puzzle)
  lock         Lock result -> future runs error if result differ, remembers fastest duration
               (sumitting correct result with 'aoc submit' locks the result automatically)
  unlock       Unlock result -> remember only last run`

type Commands interface {
	Run(year, day, part int, input string) error
	Status(year, day, part int, input string) error
	Lock(year, day, part int, input string) error
	Unlock(year, day, part int, input string) error
	GenDay(year, day int) error
	GenAoc(module string) error
	Check() error
	ClearCache() error
	Login(session string) error
	Submit() error
}

func Start(cmd Commands, args ...string) error {
	if len(args) == 0 {
		help()
		return nil
	}

	// Run
	if strings.Contains(args[0], "-") {
		return run(args...)
	}

	switch args[0] {
	case opStatus:
		return status(args[1:]...)
	case opLock:
		return lock(args[1:]...)
	case opUnlock:
		return unlock(args[1:]...)
	case opInit:
		return initialize(args[1:]...)
	case opLogin:
		return login(args[1:]...)
	case opHelp:
		return help(args[1:]...)
	case opVersion:
		return version(args[1:]...)
	case opSubmit:
		return submit(args[1:]...)
	case opCheck:
		return check(args[1:]...)
	case opCache:
		if len(args) < 2 {
			return fmt.Errorf("1unknown command: %s", args[0])
		}
		if args[1] != opClear {
			return fmt.Errorf("2unknown command: %s", args[1])
		}
		return cacheClear(args[2:]...)
	default:
		return fmt.Errorf("3unknown command: %s", args[0])
	}
}

func run(args ...string) error {
	fs := flag.NewFlagSet("run", flag.ExitOnError)

	var year int
	var day int
	var part int
	var input string

	fs.IntVar(&year, "y", defaultYear(), "year of the puzzle to run")
	fs.IntVar(&day, "d", 0, "day of the puzzle")
	fs.IntVar(&part, "p", 0, "which part of the puzzle to run")
	fs.StringVar(&input, "i", "input.txt", "input file to feed the puzzle")

	err := fs.Parse(args)
	if err != nil {
		return err
	}

	validateRequired(fs, "y", year)
	validateRequired(fs, "d", day)
	validateRequired(fs, "p", part)
	validateRequired(fs, "i", input)

	fmt.Println("Run", year, day, part, input) // TODO: Swap for function call
	return nil
}
func initialize(args ...string) error {
	fs := flag.NewFlagSet("init", flag.ExitOnError)

	var year int
	var day int
	var module string

	fs.IntVar(&year, "y", defaultYear(), "year of the puzzle to scaffold. Cannot be used in combination with -m")
	fs.IntVar(&day, "d", 0, "scaffold a puzzle for this day. Cannot be used in combination with -m")
	fs.StringVar(&module, "m", "", "create module with this name. Cannot be used in combination with -d and -y")

	err := fs.Parse(args)
	if err != nil {
		return err
	}

	validateMutEx(fs, "d", "m", day, module)

	fmt.Println("Init", year, day, module) // TODO: Swap for function call
	return nil
}
func login(args ...string) error {
	fs := flag.NewFlagSet("login", flag.ExitOnError)

	var session string

	fs.StringVar(&session, "s", "", "your AoC account session token")

	err := fs.Parse(args)
	if err != nil {
		return err
	}

	validateRequired(fs, "s", session)

	fmt.Println("Login", session) // TODO: Swap for function call
	return nil
}
func cacheClear(args ...string) error {
	fs := flag.NewFlagSet("cache clear", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Println("Usage of cache clear:")
		fmt.Println("Clear aoc cache including all puzzle results, durations and login token")
	}

	err := fs.Parse(args)
	if err != nil {
		return err
	}

	fmt.Println("Cache Clear") // TODO: Swap for function call
	return nil
}
func check(args ...string) error {
	fs := flag.NewFlagSet("check", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Println("Usage of check:")
		fmt.Println("Run and verify correct results from all locked solutions")
	}

	err := fs.Parse(args)
	if err != nil {
		return err
	}

	fmt.Println("Check") // TODO: Swap for function call
	return nil
}
func submit(args ...string) error {
	fs := flag.NewFlagSet("submit", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Println("Usage of submit:")
		fmt.Println("Submit last run result. Requires login.")
	}

	err := fs.Parse(args)
	if err != nil {
		return err
	}

	fmt.Println("Submit") // TODO: Swap for function call
	return nil
}
func help(args ...string) error {
	fs := flag.NewFlagSet("help", flag.ExitOnError)

	var verbose bool

	fs.BoolVar(&verbose, "v", false, "show extended usage")

	err := fs.Parse(args)
	if err != nil {
		return err
	}

	text := usage
	if verbose {
		text += usageAppendix
	}

	fmt.Println(strings.ReplaceAll(text, "{{year}}", fmt.Sprint(defaultYear())))

	return nil
}
func version(args ...string) error {
	fs := flag.NewFlagSet("version", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Println("Usage of version:")
		fmt.Println("Print version information")
	}

	err := fs.Parse(args)
	if err != nil {
		return err
	}

	info, ok := debug.ReadBuildInfo()
	if !ok {
		fmt.Println("unknown")
	}
	fmt.Println(info.Main.Version) // TODO: Swap for function call

	return nil
}

// Legacy
func status(args ...string) error {
	fs := flag.NewFlagSet("status", flag.ExitOnError)

	var year int
	var day int
	var part int
	var input string

	fs.IntVar(&year, "y", defaultYear(), "year of the puzzle")
	fs.IntVar(&day, "d", 0, "day of the puzzle")
	fs.IntVar(&part, "p", 0, "part for which to check")
	// input file implied as input.txt
	//TODO: Make input.txt mandatory and implied instead of default

	err := fs.Parse(args)
	if err != nil {
		return err
	}

	validateRequired(fs, "y", year)
	validateRequired(fs, "d", day)
	validateRequired(fs, "p", part)
	validateRequired(fs, "i", input)

	fmt.Println("Status", year, day, part, input) // TODO: Swap for function call
	return nil
}
func lock(args ...string) error {
	fs := flag.NewFlagSet("status", flag.ExitOnError)

	var year int
	var day int
	var part int
	var input string

	fs.IntVar(&year, "y", defaultYear(), "year of the puzzle")
	fs.IntVar(&day, "d", 0, "day of the puzzle")
	fs.IntVar(&part, "p", 0, "part for which to check")
	// input file implied as input.txt
	//TODO: Make input.txt mandatory and implied instead of default

	err := fs.Parse(args)
	if err != nil {
		return err
	}

	validateRequired(fs, "y", year)
	validateRequired(fs, "d", day)
	validateRequired(fs, "p", part)
	validateRequired(fs, "i", input)

	fmt.Println("Lock", year, day, part, input) // TODO: Swap for function call
	return nil
}
func unlock(args ...string) error {
	fs := flag.NewFlagSet("status", flag.ExitOnError)

	var year int
	var day int
	var part int
	var input string

	fs.IntVar(&year, "y", defaultYear(), "year of the puzzle")
	fs.IntVar(&day, "d", 0, "day of the puzzle")
	fs.IntVar(&part, "p", 0, "part for which to check")
	// input file implied as input.txt
	//TODO: Make input.txt mandatory and implied instead of default

	err := fs.Parse(args)
	if err != nil {
		return err
	}

	validateRequired(fs, "y", year)
	validateRequired(fs, "d", day)
	validateRequired(fs, "p", part)
	validateRequired(fs, "i", input)

	fmt.Println("Unlock", year, day, part, input) // TODO: Swap for function call
	return nil
}

func defaultYear() int {
	now := time.Now()
	year := now.Year()

	// Before the start of this year's AoC default to previous year.
	ny, _ := time.LoadLocation("America/New_York")
	if now.Before(time.Date(year, time.December, 1, 0, 0, 0, 0, ny)) {
		return year - 1
	}

	return year
}

func validateMutEx[A comparable, B comparable](fs *flag.FlagSet, fA, fB string, vA A, vB B) {
	var zeroA A
	var zeroB B
	if (vA != zeroA) && (vB != zeroB) {
		fmt.Fprintf(fs.Output(), "flags mutually exclusive, provide one: -%s, -%s\n", fA, fB)
		fmt.Fprintf(fs.Output(), "Usage of %s:\n", fs.Name())
		fs.PrintDefaults()
		os.Exit(2)
	}
	if (vA == zeroA) && (vB == zeroB) {
		fmt.Fprintf(fs.Output(), "one of flags required, provide one: -%s, -%s\n", fA, fB)
		fmt.Fprintf(fs.Output(), "Usage of %s:\n", fs.Name())
		fs.PrintDefaults()
		os.Exit(2)
	}
}

func validateRequired[T comparable](fs *flag.FlagSet, flag string, v T) {
	var zero T
	if v == zero {
		fmt.Fprintf(fs.Output(), "flag required but not provided: -%s\n", flag)
		fmt.Fprintf(fs.Output(), "Usage of %s:\n", fs.Name())
		fs.PrintDefaults()
		os.Exit(2)
	}
}
