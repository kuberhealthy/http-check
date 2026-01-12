package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	// defaultCount is used when COUNT is unset.
	defaultCount = 0
	// defaultSeconds is used when SECONDS is unset.
	defaultSeconds = 0
	// defaultPassingPercent is used when PASSING_PERCENT is unset.
	defaultPassingPercent = 100
	// defaultRequestType is used when REQUEST_TYPE is unset.
	defaultRequestType = "GET"
	// defaultRequestBody is used when REQUEST_BODY is unset.
	defaultRequestBody = "{}"
	// defaultExpectedStatusCode is used when EXPECTED_STATUS_CODE is unset.
	defaultExpectedStatusCode = 200
)

// CheckConfig stores configuration for the HTTP check.
type CheckConfig struct {
	// CheckURL is the URL to query.
	CheckURL string
	// Count is the number of requests to perform.
	Count int
	// Seconds is the pause between requests.
	Seconds int
	// PassingPercent is the percent of successful responses required.
	PassingPercent int
	// RequestType is the HTTP method to use.
	RequestType string
	// RequestBody is the body payload for non-GET requests.
	RequestBody string
	// ExpectedStatusCode is the HTTP status code to expect.
	ExpectedStatusCode int
}

// parseConfig loads environment variables into a CheckConfig.
func parseConfig() (*CheckConfig, error) {
	// Start with defaults.
	cfg := &CheckConfig{}
	cfg.Count = defaultCount
	cfg.Seconds = defaultSeconds
	cfg.PassingPercent = defaultPassingPercent
	cfg.RequestType = defaultRequestType
	cfg.RequestBody = defaultRequestBody
	cfg.ExpectedStatusCode = defaultExpectedStatusCode

	// Read the check URL.
	checkURL := os.Getenv("CHECK_URL")
	if len(checkURL) == 0 {
		return nil, fmt.Errorf("empty CHECK_URL specified. Please update your CHECK_URL environment variable")
	}
	if !strings.HasPrefix(checkURL, "http") {
		return nil, fmt.Errorf("given URL does not declare a supported protocol. (http | https)")
	}
	cfg.CheckURL = checkURL

	// Parse COUNT.
	count := os.Getenv("COUNT")
	if len(count) != 0 {
		countValue, err := strconv.Atoi(count)
		if err != nil {
			return nil, fmt.Errorf("error converting COUNT to int: %w", err)
		}
		cfg.Count = countValue
	}

	// Parse SECONDS.
	seconds := os.Getenv("SECONDS")
	if len(seconds) != 0 {
		secondsValue, err := strconv.Atoi(seconds)
		if err != nil {
			return nil, fmt.Errorf("error converting SECONDS to int: %w", err)
		}
		cfg.Seconds = secondsValue
	}

	// Parse PASSING_PERCENT.
	passing := os.Getenv("PASSING_PERCENT")
	if len(passing) != 0 {
		passingValue, err := strconv.Atoi(passing)
		if err != nil {
			return nil, fmt.Errorf("error converting PASSING_PERCENT to int: %w", err)
		}
		cfg.PassingPercent = passingValue
	}
	if cfg.PassingPercent == 0 {
		cfg.PassingPercent = defaultPassingPercent
	}

	// Parse REQUEST_TYPE.
	requestType := os.Getenv("REQUEST_TYPE")
	if len(requestType) != 0 {
		cfg.RequestType = requestType
	}

	// Parse REQUEST_BODY.
	requestBody := os.Getenv("REQUEST_BODY")
	if len(requestBody) != 0 {
		cfg.RequestBody = requestBody
	}

	// Parse EXPECTED_STATUS_CODE.
	expectedStatusCode := os.Getenv("EXPECTED_STATUS_CODE")
	if len(expectedStatusCode) != 0 {
		statusValue, err := strconv.Atoi(expectedStatusCode)
		if err != nil {
			return nil, fmt.Errorf("error converting EXPECTED_STATUS_CODE to int: %w", err)
		}
		cfg.ExpectedStatusCode = statusValue
	}
	if cfg.ExpectedStatusCode == 0 {
		cfg.ExpectedStatusCode = defaultExpectedStatusCode
	}

	return cfg, nil
}
