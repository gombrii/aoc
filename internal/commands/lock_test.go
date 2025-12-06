package commands_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gombrii/aoc/internal/commands"
)

func TestLock(t *testing.T) {
	_, testCache, wd := prepare(t)
	cmd := commands.Commands{}

	initCache(t, wd, testCache)

	data, _ := os.ReadFile(filepath.Join(testCache, "puzzles", "2024-day1-part1-input", "lock"))
	if strings.TrimSpace(string(data)) != "false" {
		t.Error("lock didn't start as false")
	}

	if err := cmd.Lock(2024, 1, 1, "input.txt"); err != nil {
		t.Errorf("calling Lock: %v", err)
	}

	data, _ = os.ReadFile(filepath.Join(testCache, "puzzles", "2024-day1-part1-input", "lock"))
	if string(data) != "true" {
		t.Error("locking didn't set lock to true")
	}

	if err := cmd.Unlock(2024, 1, 1, "input.txt"); err != nil {
		t.Errorf("calling Unlock: %v", err)
	}

	data, _ = os.ReadFile(filepath.Join(testCache, "puzzles", "2024-day1-part1-input", "lock"))
	if string(data) != "false" {
		t.Error("unlocking didn't set lock to false")
	}
}

func TestStatusNotExists(t *testing.T) {
	_, _, _ = prepare(t)
	if err := (commands.Commands{}).Lock(2024, 1, 1, "input.txt"); err != nil {
		t.Errorf("calling Lock: %v", err)
	}
}

func TestCorruptFiles(t *testing.T) {
	for name, params := range map[string]struct {
		corrupt func(string)
		uut     func(int, int, int, string) error
	}{
		"lock missing res": {
			corrupt: func(testCache string) {
				os.Remove(filepath.Join(testCache, "puzzles", "2024-day1-part1-input", "res"))
			},
			uut: (commands.Commands{}).Lock,
		},
		"status missing res": {
			corrupt: func(testCache string) {
				os.Remove(filepath.Join(testCache, "puzzles", "2024-day1-part1-input", "res"))
			},
			uut: (commands.Commands{}).Status,
		},
	} {
		t.Run(name, func(t *testing.T) {
			_, testCache, wd := prepare(t)

			initCache(t, wd, testCache)

			params.corrupt(testCache)

			if err := params.uut(2024, 1, 1, "input.txt"); err == nil {
				t.Error("Command didn't return error despite corrupt cache")
			}
		})
	}
}
