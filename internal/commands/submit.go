package commands

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gombrii/aoc/internal/cache"
	"github.com/gombrii/aoc/internal/com"
	"github.com/gombrii/aoc/internal/files"
)

func (c Commands) Submit() error {
	session, ok := LoggedIn()
	if !ok {
		return errors.New("no logged in user")
	}

	res, puzzleKey, err := fetchLastResult()
	if err != nil {
		return err
	}

	fmt.Printf("Submit answer %q for %d/day%d/part%d? [y/N]:", res, puzzleKey.Year, puzzleKey.Day, puzzleKey.Part)
	var s string
	_, err = fmt.Scanln(&s)
	if err != nil {
		return fmt.Errorf("reading user submit input: %v", err)
	}

	if strings.ToLower(s) != "y" && strings.ToLower(s) != "yes" {
		return nil
	}

	if err := com.Submit(com.NewClient(session), puzzleKey.Year, puzzleKey.Day, puzzleKey.Part, res); err != nil {
		switch {
		case errors.Is(err, com.ErrAnswerHigh):
			fmt.Println("Incorrect! Answer is too high.")
		case errors.Is(err, com.ErrAnswerHigh):
			fmt.Println("Incorrect! Answer is too low.")
		case errors.Is(err, com.ErrAlreadySolved):
			fmt.Println("This puzzle has already been solved. Go ahead and continue your quest. :)")
		default:
		}
	}

	fmt.Println("Correct answer! \033[38;2;255;255;103m*\033[0m")
	fmt.Println("This answer is now locked in. Deviating results from future runs of this puzzle will error.")
	fmt.Println("To check correctness of all locked puzzle results, run 'aoc check'.")

	return nil
}

func fetchLastResult() (string, cache.PuzzleKey, error) {
	path, ok := cache.Contains(cache.ConfigKey{Domain: User}, files.LastRun)
	if !ok {
		return "", cache.PuzzleKey{}, errors.New("no puzzle has yet been run")
	}

	data, err := files.Read(path)
	if err != nil {
		return "", cache.PuzzleKey{}, fmt.Errorf("checking last run: %v", err)
	}

	key, err := cache.ParsePuzzleKey(string(data))
	if err != nil {
		return "", cache.PuzzleKey{}, fmt.Errorf("parsing cache key: %v", err)
	}
	if !strings.HasSuffix(key.ID(), "-input") {
		return "", cache.PuzzleKey{}, errors.New("last run was not with input file input.txt")
	}

	path, ok = cache.Contains(key, files.Res)
	if !ok {
		return "", cache.PuzzleKey{}, errors.New("no result for run")
	}

	data, err = files.Read(path)
	if err != nil {
		return "", cache.PuzzleKey{}, fmt.Errorf("checking result of last run: %v", err)
	}

	return string(data), key, nil
}
