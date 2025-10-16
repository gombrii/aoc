package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gombrii/aoc/internal/cache"
)

func Status(year, day, part, input string) error {
	cacheKey := cache.Key(year, fmt.Sprint("day", day), fmt.Sprint("part", part), input)

	if _, exists := cache.ContainsKey(cacheKey); !exists {
		fmt.Printf("No record of running %s with %s\n", filepath.Join(year, day, part), fmt.Sprintf("%s.txt", input))
		return nil
	}

	data, err := readAll(map[string]string{
		"lock": cache.MakePath(cacheKey, "lock"),
		"res":  cache.MakePath(cacheKey, "res"),
		"dur":  cache.MakePath(cacheKey, "dur"),
	})
	if err != nil {
		return err
	}

	locked, _ := strconv.ParseBool(strings.TrimSpace(data["lock"]))

	if locked {
		fmt.Printf(`▣ Locked
Lock res: %s
Best dur: %s
`, data["res"], data["dur"])
	} else {
		fmt.Printf(`□ Unlocked
Last res: %s
Last dur: %s
`, data["res"], data["dur"])
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

	data, err := readAll(map[string]string{
		"res": cache.MakePath(cacheKey, "res"),
		"dur": cache.MakePath(cacheKey, "dur"),
	})
	if err != nil {
		return "", "", err
	}

	if err := os.WriteFile(cache.MakePath(cacheKey, "lock"), []byte(strconv.FormatBool(lock)), 0755); err != nil {
		return "", "", fmt.Errorf("writing to file: %v", err)
	}

	return data["res"], data["dur"], nil
}

func readAll(files map[string]string) (map[string]string, error) {
	data := make(map[string]string, len(files))

	for name, path := range files {
		bytes, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("reading file: %v", err)
		}
		data[name] = string(bytes)
	}

	return data, nil
}
