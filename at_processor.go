package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type result struct {
	AlertsProcessed int          `json:"alerts_processed,omitempty"`
	TermsScanned    int          `json:"terms_scanned_per_alert,omitempty"`
	Created         string       `json:"created_date,omitempty"`
	Results         []alertMatch `json:"results,omitempty"`
}

type alertMatch struct {
	AlertId          string           `json:"alert_id,omitempty"`
	AlertTermMatches []alertTermMatch `json:"alert_term_matches,omitempty"`
}

type alertTermMatch struct {
	AlertText   string      `json:"text,omitempty"`
	TermMatches []termMatch `json:"term_matches,omitempty"`
}

type termMatch struct {
	TermId     string `json:"id,omitempty"`
	TermText   string `json:"text,omitempty"`
	Occurances int    `json:"occurances,omitempty"`
}

func runMatchAnalysis(a []alert, t []term) (result, error) {
	if len(a) == 0 || len(t) == 0 {
		return result{}, fmt.Errorf("invalid input, cannot perform match analysis")
	}

	//create a blank results list
	result := result{
		AlertsProcessed: len(a),
		TermsScanned:    len(t),
		Created:         time.Now().Format(time.RFC3339),
		Results:         []alertMatch{}}

	//for each alert, run scan for unique substring of term
	for _, alert := range a {
		alertMatch := processAlert(alert, t)

		if len(alertMatch.AlertTermMatches) > 0 {
			result.AddItem(alertMatch)
		}
	}

	if len(result.Results) == 0 {
		return result, fmt.Errorf("match analysis yielded no results")
	}

	return result, nil
}

func processAlert(a alert, t []term) alertMatch {
	alertMatch := alertMatch{
		AlertId: a.Id}

	for _, content := range a.Contents {
		alertTermMatch := alertTermMatch{
			AlertText: content.Text}

		contentLower := strings.ToLower(content.Text)
		for _, term := range t { //cycle through each term per alert content

			if containsWholeWord(contentLower, term.Text) { //case is being ignored
				termMatch := termMatch{
					TermId:     strconv.Itoa(term.Id),
					TermText:   term.Text,
					Occurances: strings.Count(contentLower, term.Text), // Count occurrences of the term in content
				}
				alertTermMatch.AddItem(termMatch)
			}
		}

		//ensure only add alert match if instances found
		if len(alertTermMatch.TermMatches) > 0 {
			alertMatch.AddItem(alertTermMatch)
		}
	}

	return alertMatch
}

func (r *result) AddItem(am alertMatch) []alertMatch {
	r.Results = append(r.Results, am)
	return r.Results
}

func (r *alertMatch) AddItem(atm alertTermMatch) []alertTermMatch {
	r.AlertTermMatches = append(r.AlertTermMatches, atm)
	return r.AlertTermMatches
}

func (r *alertTermMatch) AddItem(tm termMatch) []termMatch {
	r.TermMatches = append(r.TermMatches, tm)
	return r.TermMatches
}

func containsWholeWord(c string, w string) bool {
	// regular expression to match the word with word boundaries (\b)
	// \b matches at the start or end of a word
	pattern := fmt.Sprintf(`\b%s\b`, regexp.QuoteMeta(w))

	// Compile the regular expression
	re := regexp.MustCompile(pattern)

	// Check if the word exists in the text as a whole word
	return re.MatchString(c)
}
