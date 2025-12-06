package com_test

import (
	_ "embed"
	"net/http"
	"strings"
	"testing"

	"github.com/gombrii/aoc/internal/com"
)

var (
	//go:embed testdata/puzzle2024d10.html
	puzzle2024d10 string
	//go:embed testdata/expectedtest2024d10.txt
	expected2024d10 string
	//go:embed testdata/puzzle2025d1.html
	puzzle2025d1 string
	//go:embed testdata/expectedtest2025d1.txt
	expected2025d1 string
	//go:embed testdata/puzzle2025d2.html
	puzzle2025d2 string
	//go:embed testdata/expectedtest2025d2.txt
	expected2025d2 string
	//go:embed testdata/puzzle2025d3.html
	puzzle2025d3 string
	//go:embed testdata/expectedtest2025d3.txt
	expected2025d3 string
	//go:embed testdata/puzzle2025d4.html
	puzzle2025d4 string
	//go:embed testdata/expectedtest2025d4.txt
	expected2025d4 string
)

// TODO: Find all? The rest of them?
func TestGetExampleInput(t *testing.T) {
	for name, params := range map[string]struct {
		resp     string
		expected string
	}{
		"2024 d10": {
			resp:     puzzle2024d10,
			expected: expected2024d10,
		},
		"2025 d1": {
			resp:     puzzle2025d1,
			expected: expected2025d1,
		},
		"2025 d2": {
			resp:     puzzle2025d2,
			expected: expected2025d2,
		},
		"2025 d3": {
			resp:     puzzle2025d3,
			expected: expected2025d3,
		},
		"2025 d4": {
			resp:     puzzle2025d4,
			expected: expected2025d4,
		},
	} {
		t.Run(name, func(t *testing.T) {
			client := &com.Client{
				Client: &http.Client{
					Transport: RT{status: 200, body: params.resp},
				},
			}

			res, err := com.GetExampleInput(client, 0, 0)
			if err != nil {
				t.Fatal(err)
			}

			if res != strings.TrimSpace(params.expected) {
				t.Fatalf("Got %s\nWant: %s", res, params.expected)
			}
		})
	}
}
