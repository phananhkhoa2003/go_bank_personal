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
	// Create timestamped filename
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("../loadtest_results_%s.txt", timestamp)

	// Create file for saving results
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Warning: Could not create output file: %v\n", err)
		file = nil
	}
	defer func() {
		if file != nil {
			file.Close()
		}
	}()

	// Output function that writes to both console and file
	output := func(format string, args ...interface{}) {
		text := fmt.Sprintf(format, args...)
		fmt.Print(text)
		if file != nil {
			file.WriteString(text)
		}
	}

	config := LoadTestConfig{
		BaseURL:         "http://localhost:8081", // Mock server port
		NumRequests:     1000,                    // Updated to 1000 requests per endpoint
		ConcurrentUsers: 10,                      // Increased concurrent users for better load simulation
		AuthToken:       "",                      // Will be obtained from mock server
	}

	output("Load Test Report - %s\n", time.Now().Format("2006-01-02 15:04:05"))
	output("================================\n")
	output("URL: %s (Mock Database Server)\n", config.BaseURL)
	output("Total Requests: %d\n", config.NumRequests)
	output("Concurrent Users: %d\n", config.ConcurrentUsers)
	output("Results saved to: %s\n", filename)
	output("================================\n\n")

	// Step 1: Get authentication token from mock server
	output("Step 1: Obtaining authentication token from mock server...\n")
	authToken, err := getAuthToken(config.BaseURL, output)
	if err != nil {
		output("Failed to get auth token: %v\n", err)
		output("Running public endpoint tests only...\n\n")
	} else {
		config.AuthToken = authToken
		output("Successfully obtained auth token!\n\n")
	}

	// Step 2: Test public endpoints
	output("Step 2: Testing public endpoints...\n")
	publicTests := []struct {
		name   string
		method string
		path   string
		body   interface{}
	}{
		{"User Creation", "POST", "/users", map[string]interface{}{
			"username":  fmt.Sprintf("user_%d", time.Now().Unix()),
			"password":  "secret123",
			"full_name": "Load Test User",
			"email":     fmt.Sprintf("user_%d@test.com", time.Now().Unix()),
		}},
		{"User Login", "POST", "/users/login", map[string]interface{}{
			"username": "testuser",
			"password": "secret123",
		}},
	}

	for _, test := range publicTests {
		output("Testing %s (%s %s)...\n", test.name, test.method, test.path)
		result := runLoadTest(config, test.method, test.path, test.body, false)
		printResults(test.name, result, output)
		output("\n")
	}

	// Step 3: Test authenticated endpoints (if we have token)
	if config.AuthToken != "" {
		output("Step 3: Testing authenticated endpoints...\n")
		authTests := []struct {
			name   string
			method string
			path   string
			body   interface{}
		}{
			{"List Accounts", "GET", "/accounts", nil},
			{"Create Account", "POST", "/accounts", map[string]interface{}{
				"owner":    "testuser",
				"currency": "USD",
			}},
			{"Create Transfer", "POST", "/transfers", map[string]interface{}{
				"from_account_id": 1,
				"to_account_id":   2,
				"amount":          10,
				"currency":        "USD",
			}},
		}

		for _, test := range authTests {
			output("Testing %s (%s %s)...\n", test.name, test.method, test.path)
			result := runLoadTest(config, test.method, test.path, test.body, true)
			printResults(test.name, result, output)
			output("\n")
		}
	}

	output("================================\n")
	output("Load test completed! Results saved to: %s\n", filename)
	output("Note: This test used MOCK DATABASE for pure performance measurement\n")
	output("Real database performance will be different due to I/O constraints\n")
}

func getAuthToken(baseURL string, output func(string, ...interface{})) (string, error) {
	// First create a user
	userBody := map[string]interface{}{
		"username":  "testuser",
		"password":  "secret123",
		"full_name": "Test User",
		"email":     "test@example.com",
	}

	output("Creating test user...\n")
	jsonBody, _ := json.Marshal(userBody)
	resp, err := http.Post(baseURL+"/users", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create user: %v", err)
	}
	resp.Body.Close()

	// Login to get token (even if user creation failed, login with existing user should work)
	output("Logging in to get token...\n")
	loginBody := map[string]interface{}{
		"username": "testuser",
		"password": "secret123",
	}

	jsonBody, _ = json.Marshal(loginBody)
	resp, err = http.Post(baseURL+"/users/login", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("login request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("login failed with status: %d", resp.StatusCode)
	}

	var loginResp struct {
		AccessToken string `json:"access_token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		return "", fmt.Errorf("failed to decode login response: %v", err)
	}

	return loginResp.AccessToken, nil
}

func runLoadTest(config LoadTestConfig, method, path string, body interface{}, useAuth bool) TestResult {
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
				authToken := ""
				if useAuth {
					authToken = config.AuthToken
				}
				latency, err := makeRequest(config.BaseURL+path, method, body, authToken)
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
