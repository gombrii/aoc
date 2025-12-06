package com

import (
	"errors"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func GetExampleInput(client *Client, year, day int) (string, error) {
	resp, err := client.Get(fmt.Sprintf("/%d/day/%d", year, day))
	if err != nil {
		return "", err
	}

	return firstExample(resp)
}

func firstExample(html string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", err
	}

	sel := doc.Find("body main pre code").First()
	if sel.Length() == 0 {
		return "", errors.New("example not found")
	}

	return strings.TrimSpace(sel.Text()), nil
}
