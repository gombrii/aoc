package internal

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gombrii/aoc/internal/cache"
	"github.com/gombrii/aoc/internal/files"
)

type outcome struct {
	i       int
	success bool
	err     error
}

type printable struct {
	name   string
	result string
}

func Check() error {
	ch := make(chan outcome)
	wg := sync.WaitGroup{}
	puzzles := make([]printable, 0)
	i := 0

	for _, l := range cache.All() {
		data, err := files.ReadAll(map[string]string{
			files.Lock: filepath.Join(l, files.Lock),
		})
		if err != nil {
			return fmt.Errorf("reading cache: %v", err)
		}

		locked, _ := strconv.ParseBool(data[files.Lock])
		if locked {
			wg.Add(1)
			go runnerRoutine(filepath.Join(l, files.Runner), i, ch, &wg)
			printParts := strings.Split(filepath.Base(l), "-")
			printName := strings.Join(printParts[:len(printParts)-1], "/")
			puzzles = append(puzzles, printable{name: printName})
			i++
		}
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	t := time.NewTicker(50 * time.Millisecond)
	defer t.Stop()
	spinner := `-\|/`

	i = 0
	for {
		select {
		case <-t.C:
			i++
		case out, ok := <-ch:
			if !ok {
				print(i, puzzles, spinner)
				return nil
			}
			if out.err != nil {
				puzzles[out.i].result = "\033[38;2;255;0;0merror\033[0m"
			} else if out.success {
				puzzles[out.i].result = "\033[38;2;255;255;103m*\033[0m"
			} else {
				puzzles[out.i].result = "\033[38;2;255;0;0mx\033[0m"
			}
		}

		print(i, puzzles, spinner)
		fmt.Printf("\033[%dA", len(puzzles))
	}
}

func runnerRoutine(path string, i int, ch chan<- outcome, wg *sync.WaitGroup) {
	defer wg.Done()
	out, err := exec.Command("go", "run", path).Output()
	if err != nil {
		ch <- outcome{i, false, err}
	} else if strings.Contains(string(out), "Error") {
		ch <- outcome{i, false, nil}
	} else {
		ch <- outcome{i, true, nil}
	}
}

func print(i int, lines []printable, spinner string) {
	for _, toPrint := range lines {
		if toPrint.result == "" {
			fmt.Printf("%-16s %s\n", toPrint.name, string(spinner[i%len(spinner)]))
		} else {
			fmt.Printf("%-16s %s\n", toPrint.name, toPrint.result)
		}
	}
}
