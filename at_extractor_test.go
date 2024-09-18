package main

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Mock for ClientInterface
type MockClientInterface struct {
	mock.Mock
}

// Mocking GetTestAlerts method
func (m *MockClientInterface) GetTestAlerts(ctx context.Context, apiKeyFunc func(string) string) (*http.Response, error) {
	args := m.Called(ctx, apiKeyFunc)
	return args.Get(0).(*http.Response), args.Error(1)
}

// Mocking GetTestQueryTerm method
func (m *MockClientInterface) GetTestQueryTerm(ctx context.Context, apiKeyFunc func(string) string) (*http.Response, error) {
	args := m.Called(ctx, apiKeyFunc)
	return args.Get(0).(*http.Response), args.Error(1)
}

// Test for the scenario where the API call fails
func TestGetTestAlerts_APICallFail(t *testing.T) {
	mockClientInterface := new(MockClientInterface)

	apiKeyFunc := func(key string) string {
		return "dummy-api-key"
	}

	// Define the mock behavior: return nil response and an error
	mockClientInterface.On("GetTestAlerts", mock.Anything, mock.Anything).
		Return((*http.Response)(nil), errors.New("API call failed"))

	alerts, err := mockClientInterface.GetTestAlerts(context.TODO(), apiKeyFunc)

	// Assertions
	require.Nil(t, alerts)
	require.Error(t, err)
	require.Equal(t, "API call failed", err.Error())

	// Ensure the mock was called
	mockClientInterface.AssertExpectations(t)
}

func TestGetTestAlerts_Success(t *testing.T) {
	mockClientInterface := new(MockClientInterface)

	apiKeyFunc := func(key string) string {
		return "dummy-api-key"
	}

	// Mock the behavior of GetTestAlerts to return a valid response
	mockResponse := &http.Response{Body: io.NopCloser(bytes.NewReader([]byte(`[   
    	{
        	"id": "1",
        	"contents": [
            	{
                	"text": "test alert",
                	"type": "text",
                	"language": "ên"
            	}
        	],
        	"date": "2024-09-16T11:24:53.318Z",
        	"inputType": "tweet"
    	}  
		]`))),

		StatusCode: http.StatusOK,
	}

	// Set up expectations for the mocked GetTestAlerts method
	mockClientInterface.On("GetTestAlerts", mock.Anything, mock.Anything).Return(mockResponse, nil)

	// Call the method you're testing
	alerts, err := mockClientInterface.GetTestAlerts(context.TODO(), apiKeyFunc)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, alerts)

	// Ensure the mock was called
	mockClientInterface.AssertExpectations(t)
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
