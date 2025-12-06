package com

import (
	"fmt"
	"strings"
)

func GetInput(client *Client, year, day int, session string) (string, error) {
	resp, err := client.get(fmt.Sprintf("/%d/day/%d/input", year, day))
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(resp), nil
}
