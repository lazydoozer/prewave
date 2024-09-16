package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/lazydoozer/prewave/api"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

const prefix string = ""
const indent string = "  "

func init() {
	viper.SetConfigName("application")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println("prewave alert term application started")
	c, ctx, cancelFunc, err := createClient()
	if err != nil {
		fmt.Println(err, "could not create HTTP API client")
		os.Exit(0)
	}
	defer cancelFunc()

	extractor := NewExtractor(c, ctx)

	//1. get query terms from prewave API
	queryTerms, err := extractor.getQueryTerms()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	fmt.Println("prewave query terms successfully retrieved")

	//2. get test alerts from prewave API
	alerts, err := extractor.getTestAlerts()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	fmt.Println("prewave test alerts successfully retrieved")

	//3. scrub query terms for duplicates and lower case all term text to avoid false positives
	uniqueTerms, err := scrubQueryTerms(queryTerms)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	fmt.Println("prewave query terms have been scrubbed")

	//4. compare terms to each alert for matches
	result, err := runMatchAnalysis(alerts, uniqueTerms)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	//5. output results for review
	fileName := viper.GetString("prewave.file-results")
	err = result.saveToFile(fileName)
	if err != nil {
		fmt.Println(err, "could not save prewave analysis results")
		os.Exit(0)
	}
	fmt.Println("prewave query terms to alert analysis is complete, see:", fileName, "for results")

	e := echo.New()

	e.GET("/health", func(c echo.Context) error {
		return c.JSONPretty(http.StatusOK, result, indent)
	})

	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	e.Start(":" + httpPort)
}

func (r result) saveToFile(fileName string) error {
	j, _ := json.MarshalIndent(r, prefix, indent)
	return os.WriteFile(fileName, j, os.ModePerm)
}

func createClient() (*api.ClientWithResponses, context.Context, context.CancelFunc, error) {
	c, err := api.NewClientWithResponses(viper.GetString("prewave.api.path"))
	if err != nil {
		return nil, nil, nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5000*time.Millisecond)
	return c, ctx, cancel, nil
}
