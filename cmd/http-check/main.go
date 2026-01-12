package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/kuberhealthy/kuberhealthy/v3/pkg/checkclient"
	nodecheck "github.com/kuberhealthy/kuberhealthy/v3/pkg/nodecheck"
	log "github.com/sirupsen/logrus"
)

// APIRequest describes an HTTP request configuration.
type APIRequest struct {
	// URL is the parsed URL for the request.
	URL *url.URL
	// Type is the HTTP method to use.
	Type string
	// Body is the request body.
	Body io.Reader
}

// main wires configuration and executes the HTTP check.
func main() {
	// Enable nodecheck debug output for parity with v2.
	nodecheck.EnableDebugOutput()

	// Parse configuration.
	cfg, err := parseConfig()
	if err != nil {
		reportFailureAndExit(err)
		return
	}

	// Create context for node readiness checks.
	checkTimeLimit := time.Minute * 1
	ctx, _ := context.WithTimeout(context.Background(), checkTimeLimit)

	// Validate URL.
	parsedURL, err := url.Parse(cfg.CheckURL)
	if err != nil {
		log.Errorln("Cannot parse provided URL:", err.Error())
		reportFailureAndExit(err)
		return
	}

	// Wait for Kuberhealthy endpoint readiness.
	err = nodecheck.WaitForKuberhealthy(ctx)
	if err != nil {
		log.Errorln("Error waiting for kuberhealthy endpoint to be contactable by checker pod with error:", err.Error())
	}

	// Calculate passing threshold.
	passingPercentage := float32(cfg.PassingPercent) / 100
	passingScore := passingPercentage * float32(cfg.Count)
	passInt := int(passingScore)
	log.Infoln("Looking for at least", cfg.PassingPercent, "percent of", cfg.Count, "checks to pass")

	// Run the configured checks.
	summary, err := runChecks(cfg, parsedURL)
	if err != nil {
		reportFailureAndExit(err)
		return
	}

	// Log run summary.
	log.Infoln(summary.ChecksRan, "checks ran")
	log.Infoln(summary.ChecksPassed, "checks passed")
	log.Infoln(summary.ChecksFailed, "checks failed")

	// Ensure enough checks passed.
	if summary.ChecksPassed < passInt {
		reportErr := fmt.Errorf("unable to retrieve a valid response (expected status: %d) from %s %s checks failed %d out of %d attempts", cfg.ExpectedStatusCode, cfg.RequestType, parsedURL.Redacted(), summary.ChecksFailed, summary.ChecksRan)
		reportFailureAndExit(reportErr)
		return
	}

	// Report success to Kuberhealthy.
	err = checkclient.ReportSuccess()
	if err != nil {
		log.Fatalln("error when reporting to kuberhealthy:", err.Error())
	}
	log.Infoln("Successfully reported to Kuberhealthy")
}

// checkSummary reports the results of a run.
type checkSummary struct {
	// ChecksRan is the total number of checks.
	ChecksRan int
	// ChecksPassed is the number of successful checks.
	ChecksPassed int
	// ChecksFailed is the number of failed checks.
	ChecksFailed int
}

// runChecks executes the request loop and returns a summary.
func runChecks(cfg *CheckConfig, parsedURL *url.URL) (*checkSummary, error) {
	// Initialize counters.
	log.Infoln("Beginning check.")
	summary := &checkSummary{}

	// Start a ticker if a pause is configured.
	var ticker *time.Ticker
	if cfg.Seconds > 0 {
		ticker = time.NewTicker(time.Duration(cfg.Seconds) * time.Second)
		defer ticker.Stop()
	}

	// Perform the configured number of requests.
	for summary.ChecksRan < cfg.Count {
		response, err := callAPI(APIRequest{
			URL:  parsedURL,
			Type: cfg.RequestType,
			Body: bytes.NewBuffer([]byte(cfg.RequestBody)),
		})
		summary.ChecksRan++

		if err != nil {
			summary.ChecksFailed++
			log.Errorln("Failed to reach URL:", parsedURL.Redacted())
			waitForTicker(ticker)
			continue
		}

		if response.StatusCode != cfg.ExpectedStatusCode {
			log.Errorln("Got a", response.StatusCode, "with a", http.MethodGet, "to", parsedURL.Redacted())
			summary.ChecksFailed++
			waitForTicker(ticker)
			continue
		}

		log.Infoln("Got a", response.StatusCode, "with a", http.MethodGet, "to", parsedURL.Redacted())
		summary.ChecksPassed++

		waitForTicker(ticker)
	}

	return summary, nil
}

// waitForTicker blocks until the ticker fires when configured.
func waitForTicker(ticker *time.Ticker) {
	// Wait for the next tick when configured.
	if ticker == nil {
		return
	}
	if ticker.C == nil {
		return
	}

	<-ticker.C
}

// reportFailureAndExit reports an error to Kuberhealthy and exits the program.
func reportFailureAndExit(err error) {
	// Log the error and report to Kuberhealthy.
	log.Errorln(err)
	reportErr := checkclient.ReportFailure([]string{err.Error()})
	if reportErr != nil {
		log.Fatalln("error when reporting to kuberhealthy:", reportErr.Error())
	}

	os.Exit(0)
}

// callAPI performs an API call on the basis of the request type, body, and URL.
func callAPI(request APIRequest) (*http.Response, error) {
	// Handle GET requests.
	if request.Type == http.MethodGet {
		response, err := http.Get(request.URL.String())
		if err != nil {
			return nil, fmt.Errorf("error occurred while calling %s: %w", request.URL.Redacted(), err)
		}
		return response, nil
	}

	// Handle other request types.
	if request.Type == http.MethodPost || request.Type == http.MethodPut || request.Type == http.MethodDelete || request.Type == http.MethodPatch {
		req, err := http.NewRequest(request.Type, request.URL.String(), request.Body)
		if err != nil {
			return nil, fmt.Errorf("error occurred while calling %s: %w", request.URL.Redacted(), err)
		}
		response, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error occurred while calling %s: %w", request.URL.Redacted(), err)
		}
		return response, nil
	}

	return nil, fmt.Errorf("error occurred while calling %s: wrong request type found", request.URL.Redacted())
}
