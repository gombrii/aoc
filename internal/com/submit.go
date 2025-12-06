package com

import (
	"errors"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var (
	ErrAnswerHigh    = errors.New("answer is too high")
	ErrAnswerLow     = errors.New("answer is too low")
	ErrAlreadySolved = errors.New("puzzle already solved")
)

func Submit(client *Client, year, day, part int, answer string) error {
	resp, err := client.Post(fmt.Sprintf("/%d/day/%d/answer", year, day), fmt.Sprintf("level=%d&answer=%s", part, answer))
	if err != nil {
		return err
	}

	return checkResult(resp)
}

func checkResult(html string) error {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return err
	}

	sel := doc.Find("body main article p").First()
	if sel.Length() == 0 {
		return errors.New("result not found")
	}

	switch {
	case strings.Contains(sel.Text(), "That's the right answer!"):
		return nil
	case strings.Contains(sel.Text(), "That's not the right answer; your answer is too low"):
		return ErrAnswerLow
	case strings.Contains(sel.Text(), "That's not the right answer; your answer is too high"):
		return ErrAnswerHigh
	case strings.Contains(sel.Text(), "Did you already complete it?"):
		return ErrAlreadySolved
	//TODO: more cases
	// - waaaaay to high
	// - waaaaay to low
	// - to recent, send again i xx:xx minutes
	default:
		return fmt.Errorf("unrecognized result: %s", sel.Text())
	}
}

//TODO: Frågan är om detta skall göra från commands. Eller det kanske inte kommer finnas nåt kommando för detta. Det kanske görs som svar Y efter ett resultat
