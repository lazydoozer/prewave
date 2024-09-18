package main

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/lazydoozer/prewave/api"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Mock for ClientInterface
type MockClient struct {
	mock.Mock
}

// Mocking GetTestAlerts method
func (m *MockClient) GetTestAlerts(ctx context.Context, reqEditors ...api.RequestEditorFn) (*http.Response, error) {
	args := m.Called(ctx, reqEditors)
	return args.Get(0).(*http.Response), args.Error(1)
}

// Mocking GetTestQueryTerm method
func (m *MockClient) GetTestQueryTerm(ctx context.Context, reqEditors ...api.RequestEditorFn) (*http.Response, error) {
	args := m.Called(ctx, reqEditors)
	return args.Get(0).(*http.Response), args.Error(1)
}

// Test for the scenario where the API call succeeds
func TestGetTestAlerts_API_Success(t *testing.T) {
	// Create a new mock client
	mockClient := new(MockClient)

	// Create the extractor with mock client
	extractor := &extractor{
		client:  mockClient,
		context: context.TODO(),
	}

	viper.Set("prewave.mode", "production")
	expectedResponse := &http.Response{
		Body: io.NopCloser(bytes.NewReader([]byte(`[
    		{
				"id": "1",
				"contents": [
					{
						"text": "Another disaster for Ferrari.",
						"type": "text",
						"language": "ên"
					}
				],
				"date": "2024-09-18T12:51:08.758Z",
				"inputType": "tweet"
    		}
		]`))),
		StatusCode: http.StatusOK,
	}

	// mock the client response
	mockClient.On("GetTestAlerts", mock.Anything, mock.Anything).Return(expectedResponse, nil)

	// Call the method with the injected mock
	alerts, err := extractor.getTestAlerts()

	require.NoError(t, err)
	require.NotNil(t, alerts)
	require.Len(t, alerts, 1)
	require.Equal(t, "1", alerts[0].Id)
	// Ensure the mock was called
	mockClient.AssertExpectations(t)
}

// Test for the scenario where the API call fails
func TestGetTestAlerts_API_Fail(t *testing.T) {
	// Create a new mock client
	mockClient := new(MockClient)

	// Create the extractor with mock client
	extractor := &extractor{
		client:  mockClient,
		context: context.TODO(),
	}

	viper.Set("prewave.mode", "production")

	mockResp := &http.Response{
		Body:       io.NopCloser(bytes.NewReader([]byte(`invalid json`))),
		StatusCode: http.StatusOK,
	}

	// Mock the client to return an error
	mockClient.On("GetTestAlerts", mock.Anything, mock.Anything).Return(mockResp, errors.New("Prewave API failure"))

	_, err := extractor.getTestAlerts()
	require.Error(t, err)
	require.Contains(t, err.Error(), "Prewave API failure")
	// Ensure the mock was called
	mockClient.AssertExpectations(t)
}

func TestGetTestData_EmptyFileName(t *testing.T) {
	_, err := getTestData("")
	if err == nil {
		t.Error("expected an error for empty filename, got none")
	}
}

func TestGetTestData_FileNotExist(t *testing.T) {
	_, err := getTestData("nofile.json")
	if err == nil {
		t.Error("expected an error for non existent file, got none")
	} else if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("expected file does not exist error, got: %v", err)
	}
}
