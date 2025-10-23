package main

import (
	"fmt"
	"time"
	"os"
	"strconv"
	"strings"

	"github.com/gombrii/Advent-of-code/2024/solutions/day1"
)

func main() {
	locked, _ := strconv.ParseBool(strings.TrimSpace(string(read("/Users/simongombrii/code/aoctest/aoc-cache/2024-day1-part1-input/lock"))))
	record, _ := time.ParseDuration(strings.TrimSpace(string(read("/Users/simongombrii/code/aoctest/aoc-cache/2024-day1-part1-input/dur"))))
	lastRes := strings.TrimSpace(string(read("/Users/simongombrii/code/aoctest/aoc-cache/2024-day1-part1-input/res")))
	data := read("2024/input/day1/input.txt")
	
	start := time.Now()
	result := day1.Part1(data)	
	duration := time.Since(start)

	correct := fmt.Sprint(result) == lastRes
	switch {
	case locked && !correct:
		fmt.Printf("Error: res: %v, want %v\n", result, lastRes)
		return
	case locked && correct:
		diff := duration - record
		fmt.Println("Res:", result)
		fmt.Printf("Dur: %v (%v, %.0f%%)\n", duration, diff, (float64(diff)/float64(duration))*100.0)
		if diff < 0 {
			write("/Users/simongombrii/code/aoctest/aoc-cache/2024-day1-part1-input/dur", fmt.Sprint(duration))
		}
	case !locked:
		fmt.Println("Res:", result)
		fmt.Println("Dur:", duration)
		write("/Users/simongombrii/code/aoctest/aoc-cache/2024-day1-part1-input/res", fmt.Sprint(result))
		write("/Users/simongombrii/code/aoctest/aoc-cache/2024-day1-part1-input/dur", fmt.Sprint(duration))
	}
}
	
func read(path string) []byte {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error: could not read file: %v", err)
		os.Exit(1)
	}

	return data
}
	
func write(path string, data string) {
	if err := os.WriteFile(path, []byte(data), 0755); err != nil {
		fmt.Printf("Error: could not write file: %v\n", err)
		os.Exit(1)
	}
}