package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gombrii/aoc/internal/cache"
	"github.com/gombrii/aoc/internal/files"
)

func Status(year, day, part, input string) error {
	key := cache.Key(year, fmt.Sprint("day", day), fmt.Sprint("part", part), input)

	if _, exists := cache.ContainsKey(key); !exists {
		fmt.Printf("No record of running %s with %s\n", filepath.Join(year, day, part), fmt.Sprintf("%s.txt", input))
		return nil
	}

	data, err := files.ReadAll(map[string]string{
		files.Lock: cache.MakePath(key, files.Lock),
		files.Res:  cache.MakePath(key, files.Res),
		files.Dur:  cache.MakePath(key, files.Dur),
	})
	if err != nil {
		return err
	}

	locked, _ := strconv.ParseBool(strings.TrimSpace(data[files.Lock]))

	if locked {
		fmt.Printf(`▣ Locked
Lock res: %s
Best dur: %s
`, data[files.Res], data[files.Dur])
	} else {
		fmt.Printf(`□ Unlocked
Last res: %s
Last dur: %s
`, data[files.Res], data[files.Dur])
	}

	return nil
}

func Lock(year, day, part, input string) error {
	res, dur, err := setLock(true, year, day, part, input)
	if err != nil {
		return fmt.Errorf("setting lock to true: %v", err)
	}

	fmt.Printf(`▣ Locked
Lock res: %s
Best dur: %s
`, res, dur)

	return nil
}

func Unock(year, day, part, input string) error {
	_, _, err := setLock(false, year, day, part, input)
	if err != nil {
		return fmt.Errorf("setting lock to false: %v", err)
	}

	fmt.Println(`□ Unlocked`)

	return nil
}

func setLock(lock bool, year, day, part, input string) (res, dur string, err error) {
	cacheKey := cache.Key(year, fmt.Sprint("day", day), fmt.Sprint("part", part), input)

	if _, exists := cache.ContainsKey(cacheKey); !exists {
		fmt.Printf("No record of running %s with %s\n", filepath.Join(year, day, part), fmt.Sprintf("%s.txt", input))
		return "", "", nil
	}

	data, err := files.ReadAll(map[string]string{
		files.Res: cache.MakePath(cacheKey, files.Res),
		files.Dur: cache.MakePath(cacheKey, files.Dur),
	})
	if err != nil {
		return "", "", err
	}

	if err := os.WriteFile(cache.MakePath(cacheKey, files.Lock), []byte(strconv.FormatBool(lock)), 0755); err != nil {
		return "", "", fmt.Errorf("writing to file: %v", err)
	}

	return data[files.Res], data[files.Dur], nil
}
