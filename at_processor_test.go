package main

import (
	"testing"
)

func TestProcessAlert_BasicMatch(t *testing.T) {
	alert := alert{
		Id: "1",
		Contents: []content{
			{Text: "This is a test message"},
		},
	}

	terms := []term{
		{Id: 1, Text: "test"},
	}
	result := processAlert(alert, terms)

	if len(result.AlertTermMatches) != 1 {
		t.Errorf("Expected 1 match, got %d", len(result.AlertTermMatches))
	}
}

func TestProcessAlert_BasicMatchMultipleOccurances(t *testing.T) {
	alert := alert{
		Id: "1",
		Contents: []content{
			{Text: "This is a test message test"},
		},
	}

	terms := []term{
		{Id: 1, Text: "test"},
	}
	result := processAlert(alert, terms)

	if result.AlertTermMatches[0].TermMatches[0].Occurances != 2 {
		t.Errorf("Expected 2 occurances, got %d", result.AlertTermMatches[0].TermMatches[0].Occurances)
	}
}

func TestProcessAlert_NoMatch(t *testing.T) {
	alert := alert{
		Id: "1",
		Contents: []content{
			{Text: "text with no match"},
		},
	}

	terms := []term{
		{Id: 1, Text: "test"},
	}
	result := processAlert(alert, terms)

	if len(result.AlertTermMatches) != 0 {
		t.Errorf("Expected 0 matches, got %d", len(result.AlertTermMatches))
	}
}

func TestProcessAlert_MultipleContents(t *testing.T) {
	alert := alert{
		Id: "1",
		Contents: []content{
			{Text: "First content with test"},
			{Text: "Second content with test again"},
		},
	}

	terms := []term{
		{Id: 1, Text: "test"},
	}
	result := processAlert(alert, terms)

	if len(result.AlertTermMatches) != 2 {
		t.Errorf("Expected 2 term matches, got %d", len(result.AlertTermMatches))
	}
}

func TestProcessAlert_EmptyAlertContent(t *testing.T) {
	alert := alert{
		Id:       "1",
		Contents: []content{},
	}

	terms := []term{
		{Id: 1, Text: "test"},
	}
	result := processAlert(alert, terms)

	if len(result.AlertTermMatches) != 0 {
		t.Errorf("Expected 0 matches, got %d", len(result.AlertTermMatches))
	}
}

func TestProcessAlert_EmptyTerms(t *testing.T) {
	alert := alert{
		Id: "1",
		Contents: []content{
			{Text: "This is a test message"},
		},
	}

	terms := []term{}
	result := processAlert(alert, terms)

	if len(result.AlertTermMatches) != 0 {
		t.Errorf("Expected 0 matches, got %d", len(result.AlertTermMatches))
	}
}
