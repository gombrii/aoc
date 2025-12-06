package com

import (
	"fmt"
	"strings"
)

func GetPuzzleInput(client *Client, year, day int) (string, error) {
	resp, err := client.Get(fmt.Sprintf("/%d/day/%d/input", year, day))
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(resp), nil
}
