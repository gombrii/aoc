package com_test

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/gombrii/aoc/internal/com"
)

type RT struct {
	body   string
	status int
}

func (m RT) RoundTrip(r *http.Request) (*http.Response, error) {
	return &http.Response{
		Status:     "OK",
		StatusCode: m.status,
		Body:       io.NopCloser(strings.NewReader(m.body)),
		Request:    r,
	}, nil
}

func TestSuccess(t *testing.T) {
	client := &com.Client{
		Client: &http.Client{
			Transport: RT{status: 200, body: settings},
		},
	}

	_, err := com.Ping(client)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFail(t *testing.T) {
	for name, params := range map[string]struct {
		respStatus  int
		expectedErr error
	}{
		"unauthorized": {
			respStatus:  400,
			expectedErr: com.ErrUnauthorized,
		},
	} {
		t.Run(name, func(t *testing.T) {
			client := &com.Client{
				Client: &http.Client{
					Transport: RT{status: params.respStatus},
				},
			}

			_, err := com.Ping(client)
			if err == nil {
				t.Fatal("Expected error but got none")
			} else if !errors.Is(err, params.expectedErr) {
				t.Fatalf("Got %v\nWant: %v", err, params.expectedErr)
			}
		})
	}
}
