package app

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"runtime/debug"
	"time"
)

type validator func() error

var ErrInput = errors.New("")

func printVersion() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		fmt.Println("unknown")
	}
	fmt.Println(info.Main.Version)
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

func isSet[T comparable](v *T) bool {
	return *v != *new(T)
}

func mutuallyExclusive[A comparable, B comparable](fs *flag.FlagSet, fA string, vA *A, fB string, vB *B) validator {
	return func() error {
		if *vA != *new(A) && *vB != *new(B) {
			fmt.Fprintf(fs.Output(), "flags mutually exclusive, provide one: -%s, -%s\n", fA, fB)
			fmt.Fprintf(fs.Output(), "Usage of %s:\n", fs.Name())
			fs.PrintDefaults()
			return ErrInput
		}

		return nil
	}
}

func oneRequired[A comparable, B comparable](fs *flag.FlagSet, fA string, vA *A, fB string, vB *B) validator {
	return func() error {
		if *vA == *new(A) && *vB == *new(B) {
			fmt.Fprintf(fs.Output(), "flag required, provide one: -%s, -%s\n", fA, fB)
			fmt.Fprintf(fs.Output(), "Usage of %s:\n", fs.Name())
			fs.PrintDefaults()
			return ErrInput
		}

		return nil
	}
}

func required[T comparable](fs *flag.FlagSet, flag string, v *T) validator {
	return func() error {
		if *v == *new(T) {
			fmt.Fprintf(fs.Output(), "flag required but not provided: -%s\n", flag)
			fmt.Fprintf(fs.Output(), "Usage of %s:\n", fs.Name())
			fs.PrintDefaults()
			return ErrInput
		}

		return nil
	}
}

func inRange(fs *flag.FlagSet, flag string, v *int, min, max int) validator {
	return func() error {
		if *v < min || *v > max {
			fmt.Fprintf(fs.Output(), "value of -%s must be at least %d and at most %d: \n", flag, min, max)
			fmt.Fprintf(fs.Output(), "Usage of %s:\n", fs.Name())
			fs.PrintDefaults()
			return ErrInput
		}

		return nil
	}
}

func flagSet(name string) (*flag.FlagSet, *bytes.Buffer) {
	var buf bytes.Buffer
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(&buf)
	return fs, &buf
}

func parse(fs *flag.FlagSet, buf *bytes.Buffer, args []string, validators ...validator) error {
	err := fs.Parse(args)
	if err != nil {
		return fmt.Errorf("%w%s", ErrInput, buf.String())
	}

	for _, valid := range validators {
		if err := valid(); err != nil {
			return fmt.Errorf("%w%s", ErrInput, buf.String())
		}
	}

	return nil
}
