package com

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var ErrUnauthorized = errors.New("unauthorized")

type Client struct {
	*http.Client
	Address string
	Session string
}

func NewClient(session string) *Client {
	return &Client{
		Client:  &http.Client{},
		Address: "https://adventofcode.com",
		Session: session,
	}
}

func (c *Client) Post(path, body string) (string, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", c.Address, path), strings.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return send(c, req, c.Session)
}

func (c *Client) get(path string) (string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", c.Address, path), nil)
	if err != nil {
		return "", fmt.Errorf("creating request: %v", err)
	}

	return send(c, req, c.Session)
}

func send(client *Client, req *http.Request, session string) (string, error) {
	req.Header.Set("Cookie", fmt.Sprintf("session=%s", session))

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("sending request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", ErrUnauthorized
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading body: %v", err)
	}

	return string(data), nil
}
