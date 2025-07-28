package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// FallbackLookup sends the fallback CURL request and returns a valid response map or error
func FallbackLookup(body []byte) (map[string]interface{}, error) {
	url := "https://cps.combain.com/?key=7b1417e4c9cf9cc1ee9d"

	req, err := http.NewRequestWithContext(context.TODO(), http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("creating fallback request failed: %w", err)
	}

	req.Header.Set("accept", "*/*")
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("origin", "https://portal.combain.com")
	req.Header.Set("referer", "https://portal.combain.com/")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36")

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fallback request failed: %w", err)
	}
	defer res.Body.Close()

	// Step 1: Read body
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading fallback response failed: %w", err)
	}

	// Step 2: Unmarshal the response
	var responseMap map[string]interface{}
	if err := json.Unmarshal(resBody, &responseMap); err != nil {
		log.Println("Error unmarshalling fallback response:", err)
		return nil, err
	}

	return responseMap, nil
}
