package app

import (
	"fmt"
	"strings"
)

const (
	opRun     = "run"
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
  aoc -d DAY -p {1|2} [-y YEAR def: {{year}}] [{-i INPUT def: input.txt | -t}]
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
  init -d DAY      Scaffold solution files for a new day (pull puzzle input from server if logged in)
  init -m MODULE   Create a new AoC module structure

Misc:
  login            Enables pulling of puzzle input and submission of solutions to server
  check            Run all locked puzzles to verify results
  cache clear      Delete all data created and kept by aoc
  help             Show this help
  version          Show installed aoc version`

const usageAppendix = `

Legacy commands:
  aoc <command> -d DAY -p {1|2} [-y YEAR current default: {{year}}]

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
		return run(cmd, args...)
	}

	switch args[0] {
	case opStatus:
		return status(cmd, args[1:]...)
	case opLock:
		return lock(cmd, args[1:]...)
	case opUnlock:
		return unlock(cmd, args[1:]...)
	case opInit:
		return initialize(cmd, args[1:]...)
	case opLogin:
		return login(cmd, args[1:]...)
	case opHelp:
		return help(args[1:]...)
	case opVersion:
		return version(args[1:]...)
	case opSubmit:
		return submit(cmd, args[1:]...)
	case opCheck:
		return check(cmd, args[1:]...)
	case opCache:
		if len(args) < 2 {
			return fmt.Errorf("unknown command: %s", args[0])
		}
		if args[1] != opClear {
			return fmt.Errorf("unknown command: %s", args[1])
		}
		return cacheClear(cmd, args[2:]...)
	default:
		return fmt.Errorf("unknown command: %s", args[0])
	}
}

func run(cmd Commands, args ...string) error {
	fs, buf := flagSet(opRun)

	year := fs.Int("y", defaultYear(), "year of the puzzle to run")
	day := fs.Int("d", 0, "day of the puzzle")
	part := fs.Int("p", 0, "which part of the puzzle to run")
	input := fs.String("i", "", `input file to feed the puzzle (default "input.txt"). Mutually exclusive with -t`)
	test := fs.Bool("t", false, `shorthand for "-i test.txt". Mutually exclusive with -i`)

	err := parse(fs, buf, args,
		required(fs, "y", year),
		required(fs, "d", day),
		required(fs, "p", part),
		mutuallyExclusive(fs, "i", input, "t", test),
		inRange(fs, "p", part, 1, 2),
	)
	if err != nil {
		return err
	}

	if *test {
		i := "test.txt"
		input = &i
	} else if *input == "" {
		i := "input.txt"
		input = &i
	}

	return cmd.Run(*year, *day, *part, *input)
}
func initialize(cmd Commands, args ...string) error {
	fs, buf := flagSet(opInit)

	year := fs.Int("y", defaultYear(), "year of the puzzle to scaffold. Mutually exclusive with -m")
	day := fs.Int("d", 0, "scaffold a puzzle for this day. Mutually exclusive with -m")
	module := fs.String("m", "", "create module with this name. Mutually exclusive with -d and -y")

	err := parse(fs, buf, args,
		oneRequired(fs, "d", day, "m", module),
		mutuallyExclusive(fs, "d", day, "m", module),
	)
	if err != nil {
		return err
	}

	if *module != "" {
		return cmd.GenAoc(*module)
	} else {
		return cmd.GenDay(*year, *day)
	}
}
func login(cmd Commands, args ...string) error {
	fs, buf := flagSet(opLogin)

	session := fs.String("s", "", "your AoC account session token")

	err := parse(fs, buf, args, required(fs, "s", session))
	if err != nil {
		return err
	}

	return cmd.Login(*session)
}
func cacheClear(cmd Commands, args ...string) error {
	fs, buf := flagSet(opCache + " " + opClear)
	fs.Usage = func() {
		fmt.Println("Usage of cache clear:")
		fmt.Println("Clear aoc cache including all puzzle results, durations and login token")
	}

	err := parse(fs, buf, args)
	if err != nil {
		return err
	}

	return cmd.ClearCache()
}
func check(cmd Commands, args ...string) error {
	fs, buf := flagSet(opCheck)
	fs.Usage = func() {
		fmt.Println("Usage of check:")
		fmt.Println("Run and verify correct results from all locked solutions")
	}

	err := parse(fs, buf, args)
	if err != nil {
		return err
	}

	return cmd.Check()
}
func submit(cmd Commands, args ...string) error {
	fs, buf := flagSet(opSubmit)
	fs.Usage = func() {
		fmt.Println("Usage of submit:")
		fmt.Println("Submit last run result. Requires login.")
	}

	err := parse(fs, buf, args)
	if err != nil {
		return err
	}

	return cmd.Submit()
}
func help(args ...string) error {
	fs, buf := flagSet(opHelp)

	verbose := fs.Bool("v", false, "show extended usage")

	err := parse(fs, buf, args)
	if err != nil {
		return err
	}

	text := usage
	if *verbose {
		text += usageAppendix
	}

	fmt.Println(strings.ReplaceAll(text, "{{year}}", fmt.Sprint(defaultYear())))
	return nil
}
func version(args ...string) error {
	fs, buf := flagSet(opVersion)
	fs.Usage = func() {
		fmt.Println("Usage of version:")
		fmt.Println("Print installed aoc version")
	}

	err := parse(fs, buf, args)
	if err != nil {
		return err
	}

	printVersion()
	return nil
}

func status(cmd Commands, args ...string) error {
	fs, buf := flagSet(opStatus)

	year := fs.Int("y", defaultYear(), "year of the puzzle")
	day := fs.Int("d", 0, "day of the puzzle")
	part := fs.Int("p", 0, "part for which to check")

	err := parse(fs, buf, args,
		required(fs, "y", year),
		required(fs, "d", day),
		required(fs, "p", part),
		inRange(fs, "p", part, 1, 2),
	)
	if err != nil {
		return err
	}

	return cmd.Status(*year, *day, *part, "input.txt")
}
func lock(cmd Commands, args ...string) error {
	fs, buf := flagSet(opLock)

	year := fs.Int("y", defaultYear(), "year of the puzzle")
	day := fs.Int("d", 0, "day of the puzzle")
	part := fs.Int("p", 0, "part for which to check")

	err := parse(fs, buf, args,
		required(fs, "y", year),
		required(fs, "d", day),
		required(fs, "p", part),
		inRange(fs, "p", part, 1, 2),
	)
	if err != nil {
		return err
	}

	return cmd.Lock(*year, *day, *part, "input.txt")
}
func unlock(cmd Commands, args ...string) error {
	fs, buf := flagSet(opUnlock)

	year := fs.Int("y", defaultYear(), "year of the puzzle")
	day := fs.Int("d", 0, "day of the puzzle")
	part := fs.Int("p", 0, "part for which to check")

	err := parse(fs, buf, args,
		required(fs, "y", year),
		required(fs, "d", day),
		required(fs, "p", part),
		inRange(fs, "p", part, 1, 2),
	)
	if err != nil {
		return err
	}

	return cmd.Unlock(*year, *day, *part, "input.txt")
}
