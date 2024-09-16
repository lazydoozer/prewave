package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type result struct {
	AlertsProcessed int              `json:"alerts_processed,omitempty"`
	TermsScanned    int              `json:"terms_scanned_per_alert,omitempty"`
	Created         string           `json:"created_date,omitempty"`
	Matches         []alertTermMatch `json:"matches,omitempty"`
}

type alertTermMatch struct {
	AlertId   string      `json:"alert_id,omitempty"`
	AlertText string      `json:"alert_text,omitempty"`
	TermMatch []termMatch `json:"terms,omitempty"`
}

type termMatch struct {
	TermId     string `json:"term_id,omitempty"`
	TermText   string `json:"term_text,omitempty"`
	Occurances int    `json:"occurances,omitempty"`
}

func runMatchAnalysis(alerts []alert, terms []term) (result, error) {
	if len(alerts) || len(terms) == 0 {
		return nil, fmt.Errorf("invalid input, cannot perform match analysis")
	}

	//create a blank results list
	result := result{
		AlertsProcessed: len(alerts),
		TermsScanned:    len(terms),
		Created:         time.Now().Format(time.RFC3339),
		Matches:         []alertTermMatch{}}

	//for each alert, run scan for unique substring of term
	for _, alert := range alerts {
		alertMatch := processAlert(alert, terms)

		if len(alertMatch.TermMatch) > 0 {
			result.AddItem(alertMatch)
		}
	}

	if len(result.Matches) == 0 {
		return result, fmt.Errorf("match analysis yielded no results")
	}

	return result, nil
}

func processAlert(alert alert, terms []term) alertTermMatch {
	alertTermMatch := alertTermMatch{
		AlertId:   alert.Id,
		AlertText: alert.Contents[0].Text}

	for _, content := range alert.Contents {
		contentLower := strings.ToLower(content.Text)
		for _, term := range terms { //cycle through each term per alert content
			if containsWholeWord(contentLower, term.Text) { //case is being ignored
				termMatch := termMatch{
					TermId:     strconv.Itoa(term.Id),
					TermText:   term.Text,
					Occurances: strings.Count(contentLower, term.Text), // Count occurrences of the term in content
				}
				alertTermMatch.AddItem(termMatch)
			}
		}
	}

	return alertTermMatch
}

func (r *result) AddItem(alertMatch alertTermMatch) []alertTermMatch {
	r.Matches = append(r.Matches, alertMatch)
	return r.Matches
}

func (r *alertTermMatch) AddItem(tm termMatch) []termMatch {
	r.TermMatch = append(r.TermMatch, tm)
	return r.TermMatch
}

func containsWholeWord(text string, word string) bool {
	// regular expression to match the word with word boundaries (\b)
	// \b matches at the start or end of a word
	pattern := fmt.Sprintf(`\b%s\b`, regexp.QuoteMeta(word))

	// Compile the regular expression
	re := regexp.MustCompile(pattern)

	// Check if the word exists in the text as a whole word
	return re.MatchString(text)
}
