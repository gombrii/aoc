package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gombrii/aoc/internal/cache"
	"github.com/gombrii/aoc/internal/files"
)

func (c Commands) Status(year, day, part int, input string) error {
	key := cache.Key(year, day, part, input)
	if _, exists := cache.ContainsKey(key); !exists {
		fmt.Printf("No record of running %d/day%d/part%d with %s\n", year, day, part, input)
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

func (c Commands) Lock(year, day, part int, input string) error {
	key := cache.Key(year, day, part, input)
	if _, exists := cache.ContainsKey(key); !exists {
		fmt.Printf("No record of running %d/day%d/part%d with %s\n", year, day, part, input)
		return nil
	}

	data, err := files.ReadAll(map[string]string{
		files.Res: cache.MakePath(key, files.Res),
		files.Dur: cache.MakePath(key, files.Dur),
	})
	if err != nil {
		return err
	}

	err = files.Write(cache.MakePath(key, files.Lock), []byte("true"))
	if err != nil {
		return fmt.Errorf("setting lock to true: %v", err)
	}

	fmt.Printf(`▣ Locked
Lock res: %s
Best dur: %s
`, data[files.Res], data[files.Dur])

	return nil
}

func (c Commands) Unlock(year, day, part int, input string) error {
	key := cache.Key(year, day, part, input)
	if _, exists := cache.ContainsKey(key); !exists {
		fmt.Printf("No record of running %d/day%d/part%d with %s\n", year, day, part, input)
		return nil
	}

	err := files.Write(cache.MakePath(key, files.Lock), []byte("false"))
	if err != nil {
		return fmt.Errorf("setting lock to false: %v", err)
	}

	fmt.Println(`□ Unlocked`)

	return nil
}
