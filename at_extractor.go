package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/lazydoozer/prewave/api"
	"github.com/spf13/viper"
)

type extractor struct {
	client  api.ClientInterface
	context context.Context
}

type term struct {
	Id        int    `json:"id"`
	Target    int    `json:"target,omitempty"`
	Text      string `json:"text,omitempty"`
	Language  string `json:"language,omitempty"`
	KeepOrder bool   `json:"keepOrder,omitempty"`
}

type alert struct {
	Id        string    `json:"id"`
	Contents  []content `json:"contents,omitempty"`
	Date      string    `json:"date,omitempty"`
	InputType string    `json:"inputType,omitempty"`
}

type content struct {
	Text     string `json:"text,omitempty"`
	Type     string `json:"type,omitempty"`
	Language string `json:"language,omitempty"`
}

type prewaveInterceptor struct {
	apikey string
}

func NewPrewaveInterceptor(k string) *prewaveInterceptor {
	return &prewaveInterceptor{
		apikey: k,
	}
}

func (s *prewaveInterceptor) SetAPIKey(ctx context.Context, req *http.Request) error {
	req.URL.RawQuery = fmt.Sprintf(s.apikey)
	return nil
}

func NewExtractor(c api.ClientInterface, ctx context.Context) *extractor {
	return &extractor{client: c, context: ctx}
}

func (e *extractor) getQueryTerms() ([]term, error) {
	var queryTerms []term
	var jsonData []byte

	// Fetch Viper config values only once for efficiency
	apiKey := viper.GetString("prewave.api.key")
	prewaveMode := viper.GetString("prewave.mode")

	if prewaveMode == "test" {
		b, err := getTestData(viper.GetString("prewave.file-test-terms"))
		if err != nil {
			return nil, err
		}
		jsonData = b
	} else {
		resp, err := e.client.GetTestQueryTerm(e.context, NewPrewaveInterceptor(apiKey).SetAPIKey)
		if err != nil {
			return nil, fmt.Errorf("prewave GetTestQueryTerm API call failed: %w", err)
		}
		defer resp.Body.Close()

		b, err := getResponseBody(resp)
		if err != nil {
			return nil, err
		}
		jsonData = b
	}

	if err := json.Unmarshal(jsonData, &queryTerms); err != nil {
		fmt.Println(err, "failed to deserialize list of query terms")
		return nil, err
	}
	return queryTerms, nil
}

func scrubQueryTerms(t []term) ([]term, error) {
	if len(t) == 0 {
		return nil, fmt.Errorf("input query terms is empty")
	}

	// Create a map to store unique term
	uniqueTerms := make(map[string]bool)
	// Pre-process the terms: convert to lowercase and split if needed
	scrubbedTerms := []term{}

	for _, t := range t {
		termTextLower := strings.ToLower(t.Text) // Loop through the terms list, adding elements to the map if they haven't been seen before

		if t.KeepOrder {
			if !uniqueTerms[termTextLower] { // Add unique terms to the map and update proccessed slice
				uniqueTerms[termTextLower] = true

				scrubbedTerms = append(scrubbedTerms, term{
					Id:        t.Id,
					Target:    t.Target,
					Text:      termTextLower,
					Language:  t.Language,
					KeepOrder: t.KeepOrder})
			}
		} else {
			for _, word := range strings.Fields(termTextLower) {
				wordTextLower := strings.ToLower(word)
				if !uniqueTerms[wordTextLower] { // Add unique terms to the map and update results slice
					uniqueTerms[wordTextLower] = true

					scrubbedTerms = append(scrubbedTerms, term{
						Id:        t.Id,
						Target:    t.Target,
						Text:      wordTextLower,
						Language:  t.Language,
						KeepOrder: t.KeepOrder})
				}
			}
		}
	}

	if len(scrubbedTerms) <= 0 {
		return nil, fmt.Errorf("no unique query terms found after scrub")
	}

	return scrubbedTerms, nil
}

func (e *extractor) getTestAlerts() ([]alert, error) {
	var testAlerts []alert
	var jsonData []byte

	apiKey := viper.GetString("prewave.api.key")
	prewaveMode := viper.GetString("prewave.mode")

	if prewaveMode == "test" {
		b, err := getTestData(viper.GetString("prewave.file-test-alerts"))
		if err != nil {
			return nil, err
		}
		jsonData = b
	} else {
		resp, err := e.client.GetTestAlerts(e.context, NewPrewaveInterceptor(apiKey).SetAPIKey)
		if err != nil {
			return nil, fmt.Errorf("prewave API call failed: %w", err)
		}
		defer resp.Body.Close()

		b, err := getResponseBody(resp)
		if err != nil {
			return nil, err
		}
		jsonData = b
	}

	if err := json.Unmarshal(jsonData, &testAlerts); err != nil {
		fmt.Println(err, "failed to deserialize list of test alerts")
		return nil, err
	}
	return testAlerts, nil
}

func getTestData(fn string) ([]byte, error) {
	if len(fn) == 0 {
		return nil, fmt.Errorf("invalid filename provided")
	}

	b, err := os.ReadFile(fn)

	if err != nil {
		return nil, fmt.Errorf("failed to read test data from file: %w", err)
	}
	return b, nil
}

func getResponseBody(r *http.Response) ([]byte, error) {
	if r.Body != http.NoBody {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read prewave API response body: %w", err)
		}
		return b, nil
	}

	return nil, fmt.Errorf("prewave API response body is empty")
}
