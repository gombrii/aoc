package com_test

import (
	_ "embed"
	"errors"
	"net/http"
	"testing"

	"github.com/gombrii/aoc/internal/com"
)

var (
	//go:embed testdata/settings.html
	settings string
)

func TestPingSuccess(t *testing.T) {
	client := &com.Client{
		Client: &http.Client{
			Transport: RT{status: 200, body: settings},
		},
	}

	res, err := com.Ping(client)
	if err != nil {
		t.Fatal(err)
	}

	if res != "Simon Gombrii" {
		t.Fatalf("Got %s\nWant: Simon Gombrii", res)
	}
}

func TestPingFail(t *testing.T) {
	client := &com.Client{
		Client: &http.Client{
			Transport: RT{status: 303},
		},
	}

	_, err := com.Ping(client)
	if err == nil {
		t.Fatal("Expected an error but got none")
	} else if !errors.Is(err, com.ErrUnauthorized) {
		t.Fatalf("Got %v\nWant: %v", err, com.ErrUnauthorized)
	}
}
