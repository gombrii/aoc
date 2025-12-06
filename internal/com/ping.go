package com

import (
	"errors"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func Ping(client *Client) (string, error) {
	resp, err := client.Get("/2025/settings")
	if err != nil {
		return "", err
	}

	return findUsername(resp)
}

func findUsername(html string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", err
	}

	sel := doc.Find("body header div div.user").First()
	if sel.Length() == 0 {
		return "", errors.New("name not found")
	}

	name := sel.Clone().Children().Remove().End().Text()

	return strings.TrimSpace(strings.TrimSpace(name)), nil
}
