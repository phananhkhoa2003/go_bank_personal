package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

type LoadTestConfig struct {
	BaseURL         string
	NumRequests     int
	ConcurrentUsers int
	AuthToken       string
}

type TestResult struct {
	TotalRequests  int
	SuccessCount   int
	ErrorCount     int
	AverageLatency time.Duration
	MaxLatency     time.Duration
	MinLatency     time.Duration
	RequestsPerSec float64
}

func main() {
	config := LoadTestConfig{
		BaseURL:         "http://localhost:8081", // Changed to mock server port
		NumRequests:     1000,
		ConcurrentUsers: 10,
		AuthToken:       "your-jwt-token-here", // Replace with actual token
	}

	// Create results file with timestamp
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("../loadtest_results_%s.txt", timestamp)

	// Redirect output to both console and file
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	// Print to both console and file
	output := func(format string, args ...interface{}) {
		message := fmt.Sprintf(format, args...)
		fmt.Print(message)
		file.WriteString(message)
	}

	output("Load Test Report - %s\n", time.Now().Format("2006-01-02 15:04:05"))
	output("================================\n")
	output("URL: %s\n", config.BaseURL)
	output("Total Requests: %d\n", config.NumRequests)
	output("Concurrent Users: %d\n", config.ConcurrentUsers)
	output("================================\n\n")

	// Test different endpoints
	endpoints := []struct {
		name   string
		method string
		path   string
		body   interface{}
	}{
		{"Get Accounts", "GET", "/accounts", nil},
		{"Create Account", "POST", "/accounts", map[string]string{"currency": "USD"}},
		{"Create Transfer", "POST", "/transfers", map[string]interface{}{
			"from_account_id": 1,
			"to_account_id":   2,
			"amount":          10,
			"currency":        "USD",
		}},
	}

	for _, endpoint := range endpoints {
		output("Testing %s (%s %s)...\n", endpoint.name, endpoint.method, endpoint.path)
		result := runLoadTest(config, endpoint.method, endpoint.path, endpoint.body)
		printResults(endpoint.name, result, output)
		output("\n")
	}
}

func runLoadTest(config LoadTestConfig, method, path string, body interface{}) TestResult {
	var wg sync.WaitGroup
	results := make(chan time.Duration, config.NumRequests)
	errors := make(chan error, config.NumRequests)

	startTime := time.Now()
	requestsPerUser := config.NumRequests / config.ConcurrentUsers

	// Start concurrent users
	for i := 0; i < config.ConcurrentUsers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < requestsPerUser; j++ {
				latency, err := makeRequest(config.BaseURL+path, method, body, config.AuthToken)
				if err != nil {
					errors <- err
				} else {
					results <- latency
				}
			}
		}()
	}

	wg.Wait()
	close(results)
	close(errors)

	totalTime := time.Since(startTime)

	// Collect results
	var latencies []time.Duration
	successCount := 0
	errorCount := 0

	for latency := range results {
		latencies = append(latencies, latency)
		successCount++
	}

	for range errors {
		errorCount++
	}

	// Calculate statistics
	var totalLatency time.Duration
	maxLatency := time.Duration(0)
	minLatency := time.Hour // Start with a large value

	for _, latency := range latencies {
		totalLatency += latency
		if latency > maxLatency {
			maxLatency = latency
		}
		if latency < minLatency {
			minLatency = latency
		}
	}

	averageLatency := time.Duration(0)
	if len(latencies) > 0 {
		averageLatency = totalLatency / time.Duration(len(latencies))
	}

	requestsPerSec := float64(config.NumRequests) / totalTime.Seconds()

	return TestResult{
		TotalRequests:  config.NumRequests,
		SuccessCount:   successCount,
		ErrorCount:     errorCount,
		AverageLatency: averageLatency,
		MaxLatency:     maxLatency,
		MinLatency:     minLatency,
		RequestsPerSec: requestsPerSec,
	}
}

func makeRequest(url, method string, body interface{}, authToken string) (time.Duration, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return 0, err
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	start := time.Now()
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	if authToken != "" {
		req.Header.Set("Authorization", "Bearer "+authToken)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	latency := time.Since(start)

	if err != nil {
		return latency, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return latency, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	return latency, nil
}

func printResults(testName string, result TestResult, output func(string, ...interface{})) {
	output("Results for %s:\n", testName)
	output("  Total Requests: %d\n", result.TotalRequests)
	output("  Successful: %d\n", result.SuccessCount)
	output("  Failed: %d\n", result.ErrorCount)
	output("  Success Rate: %.2f%%\n", float64(result.SuccessCount)/float64(result.TotalRequests)*100)
	output("  Requests/sec: %.2f\n", result.RequestsPerSec)
	output("  Average Latency: %v\n", result.AverageLatency)
	output("  Min Latency: %v\n", result.MinLatency)
	output("  Max Latency: %v\n", result.MaxLatency)
}
