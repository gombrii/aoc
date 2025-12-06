package com_test

import (
	_ "embed"
	"errors"
	"net/http"
	"testing"

	"github.com/gombrii/aoc/internal/com"
)

var (
	//go:embed testdata/answer_success.html
	successful string
	//go:embed testdata/answer_low.html
	low string
)

func TestSubmit(t *testing.T) {
	for name, params := range map[string]struct {
		resp     string
		expected error
	}{
		"success": {
			resp: successful,
		},
		"too low": {
			resp:     low,
			expected: com.ErrAnswerLow,
		},
		//TODO: More cases when I have collected test data
	} {
		t.Run(name, func(t *testing.T) {
			client := &com.Client{
				Client: &http.Client{
					Transport: RT{status: 200, body: params.resp},
				},
			}

			err := com.Submit(client, 0, 0, 0, "")

			if params.expected != nil {
				if err == nil {
					t.Fatalf("Expected error but got none")
				} else if !errors.Is(err, params.expected) {
					t.Fatalf("Got %v\nWant: %v", err, params.expected)
				}
			} else if err != nil {
				t.Fatalf("Got unexpected error: %v", err)
			}
		})
	}
}
