// Made By Sharkeyes | website https://sharkeyes.dev | Protect your web forms
package sharkeyes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	apiKey  string
	apiURL  = "https://api.sharkeyes.dev/api/v2/server-verify"
	timeout = 5 * time.Second
)

func Configure(key string, url string, t time.Duration) {
	apiKey = key
	if url != "" {
		apiURL = url
	}
	if t > 0 {
		timeout = t
	}
}

func Verify(token string, r *http.Request) (bool, string) {
	if apiKey == "" {
		panic("RuntimeError: First, call sharkeyes.Configure(apiKey, ...)")
	}

	token = strings.TrimSpace(token)

	if token == "" {
		return false, "Enable JavaScript."
	}

	bodyMap := map[string]string{
		"verification_token": token,
	}
	jsonBody, err := json.Marshal(bodyMap)
	if err != nil {
		log.Printf("[SharkEyes] %v", err)
		return false, "Technical error during verification."
	}

	client := &http.Client{Timeout: timeout}
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Printf("[SharkEyes] %v", err)
		return false, "Technical error during verification."
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[SharkEyes] %v", err)
		return false, "Technical error during verification."
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Sprintf("Verification error (Status: %d)", resp.StatusCode)
	}

	var data struct {
		IsBot  *bool  `json:"is_bot"`
		Reason string `json:"reason"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Printf("[SharkEyes] %v", err)
		return false, "Technical error during verification."
	}

	if data.IsBot != nil && *data.IsBot == false {
		return true, ""
	}

	if data.Reason != "" {
		return false, data.Reason
	}

	return false, "Verification failed."
}
