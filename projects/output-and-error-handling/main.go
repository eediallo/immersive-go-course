package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	url := "http://localhost:8080"
	resp, err := getWithRetry(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	// Print the weather response body if successful
	var bodyBytes []byte
	bodyBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	fmt.Println("Weather:")
	fmt.Println(string(bodyBytes))
}

// getWithRetry performs a GET request and handles 429 responses with appropriate retry logic.
func getWithRetry(url string) (*http.Response, error) {
	var resp *http.Response
	var err error
	for {
		resp, err = http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("Connection terminated by server. Can't get you the weather.")
		}
		if resp.StatusCode != 429 {
			return resp, nil
		}
		// Handle 429 Retry-After logic
		sleepDuration, sleepReason := parseRetryAfter(resp.Header.Get("Retry-After"))
		if sleepDuration > 5*time.Second {
			resp.Body.Close()
			return nil, fmt.Errorf("Server asked us to wait more than 5 seconds. Can't get you the weather.")
		}
		if sleepDuration > 1*time.Second {
			fmt.Printf("Notice: Retrying after %.0f seconds. Things may be a bit slow.\n", sleepDuration.Seconds())
		}
		resp.Body.Close()
		time.Sleep(sleepDuration)
		// If we used a default sleep, only retry once
		if sleepReason == "default" {
			resp, err = http.Get(url)
			if err != nil {
				return nil, fmt.Errorf("Connection terminated by server. Can't get you the weather.")
			}
			if resp.StatusCode == 429 {
				resp.Body.Close()
				return nil, fmt.Errorf("Still received 429 after default retry. Can't get you the weather.")
			}
			return resp, nil
		}
	}
}

// parseRetryAfter parses the Retry-After header and returns the sleep duration and reason.
func parseRetryAfter(retryAfter string) (time.Duration, string) {
	if retryAfter == "a while" {
		// Can't determine how long to sleep for
		// Decision: sleep for 1 second and retry once, then give up if still 429
		fmt.Println("429 received with invalid Retry-After header. Will sleep for 1 second and retry once.")
		return 1 * time.Second, "default"
	}
	// Try to parse as seconds
	var sleepSeconds int
	_, err := fmt.Sscanf(retryAfter, "%d", &sleepSeconds)
	if err == nil {
		return time.Duration(sleepSeconds) * time.Second, "seconds"
	}
	// Try to parse as HTTP-date
	retryTime, err := http.ParseTime(retryAfter)
	if err != nil {
		fmt.Printf("Could not parse Retry-After header: %s. Will sleep for 1 second and retry once.\n", err)
		return 1 * time.Second, "default"
	}
	now := time.Now()
	wait := retryTime.Sub(now)
	if wait < 0 {
		wait = 0
	}
	return wait, "http-date"
}
