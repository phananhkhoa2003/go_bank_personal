package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func main() {
	baseURL := "http://localhost:8080"
	
	// Create a unique user
	username := fmt.Sprintf("loadtest_%d", time.Now().Unix())
	userBody := map[string]interface{}{
		"username":  username,
		"password":  "secret123",
		"full_name": "Load Test User",
		"email":     fmt.Sprintf("loadtest_%d@example.com", time.Now().Unix()),
	}
	
	fmt.Printf("Creating test user: %s\n", username)
	
	// Create user
	jsonBody, _ := json.Marshal(userBody)
	resp, err := http.Post(baseURL+"/users", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Printf("Failed to create user: %v\n", err)
		return
	}
	resp.Body.Close()
	
	if resp.StatusCode != 200 {
		fmt.Printf("User creation failed with status: %d\n", resp.StatusCode)
		return
	}
	
	fmt.Println("User created successfully")
	
	// Login to get token
	loginBody := map[string]interface{}{
		"username": username,
		"password": "secret123",
	}
	
	fmt.Println("Logging in to get authentication token...")
	
	jsonBody, _ = json.Marshal(loginBody)
	resp, err = http.Post(baseURL+"/users/login", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Printf("Login failed: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		fmt.Printf("Login failed with status: %d\n", resp.StatusCode)
		return
	}
	
	var loginResp struct {
		AccessToken string `json:"access_token"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		fmt.Printf("Failed to decode login response: %v\n", err)
		return
	}
	
	fmt.Printf("Login successful!\n")
	fmt.Printf("Your auth token: %s\n", loginResp.AccessToken)
	fmt.Printf("\nNow update main.go and replace 'your-jwt-token-here' with this token\n")
	fmt.Printf("Then run: go run main.go\n")
}
