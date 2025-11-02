package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"

	"llm-inference/config"
	"llm-inference/models"
)

// OllamaService handles communication with Ollama API
type OllamaService struct {
	config     *config.Config
	client     *http.Client
	authCookie *http.Cookie
}

// NewOllamaService creates a new OllamaService instance
func NewOllamaService(cfg *config.Config) *OllamaService {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Don't follow redirects automatically
			return http.ErrUseLastResponse
		},
	}

	return &OllamaService{
		config: cfg,
		client: client,
	}
}

// authenticate performs initial authentication to get the session cookie
func (s *OllamaService) authenticate() error {
	if s.config.OllamaToken == "" {
		return fmt.Errorf("OLLAMA_TOKEN is required for authentication")
	}

	// Make initial request with token to get auth cookie
	req, err := http.NewRequest("GET", s.config.OllamaURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create auth request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.config.OllamaToken)

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to authenticate: %v", err)
	}
	defer resp.Body.Close()

	// Check if we got a redirect with Set-Cookie
	if resp.StatusCode == http.StatusFound {
		// Extract the auth cookie from the response
		for _, cookie := range resp.Cookies() {
			if cookie.Name == "C.27500673_auth_token" {
				s.authCookie = cookie
				fmt.Printf("Authentication successful, got cookie: %s\n", cookie.Value[:20]+"...")
				return nil
			}
		}
	}

	// If no redirect or cookie found, check if we're already authenticated
	if resp.StatusCode == http.StatusOK {
		fmt.Println("Already authenticated or no authentication required")
		return nil
	}

	return fmt.Errorf("authentication failed with status: %d", resp.StatusCode)
}

// CallAPI makes a request to the Ollama API
// CallAPI makes a request to the Ollama API
func (s *OllamaService) CallAPI(prompt string) (models.OllamaResponse, error) {
	// Ensure we're authenticated
	if s.authCookie == nil && s.config.OllamaToken != "" {
		if err := s.authenticate(); err != nil {
			return models.OllamaResponse{}, fmt.Errorf("authentication failed: %v", err)
		}
	}

	reqBody := models.OllamaRequest{
		Model:  s.config.ModelName,
		Prompt: prompt,
		Stream: false,
		Token:  s.config.OllamaToken,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return models.OllamaResponse{}, fmt.Errorf("failed to marshal request: %v", err)
	}

	req, err := http.NewRequest("POST", s.config.APIUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return models.OllamaResponse{}, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Add authentication cookie if we have one
	if s.authCookie != nil {
		req.AddCookie(s.authCookie)
	}

	// Also add Bearer token as fallback
	if s.config.OllamaToken != "" {
		req.Header.Set("Authorization", "Bearer "+s.config.OllamaToken)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return models.OllamaResponse{}, fmt.Errorf("failed to call Ollama API: %v", err)
	}
	defer resp.Body.Close()

	// Handle authentication retry if we get unauthorized
	if resp.StatusCode == http.StatusFound || resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("Re-authenticating due to invalid session...")
		if err := s.authenticate(); err != nil {
			return models.OllamaResponse{}, fmt.Errorf("re-authentication failed: %v", err)
		}

		// Retry the request with new authentication
		req, err = http.NewRequest("POST", s.config.APIUrl, bytes.NewBuffer(jsonData))
		if err != nil {
			return models.OllamaResponse{}, fmt.Errorf("failed to create retry request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")
		if s.authCookie != nil {
			req.AddCookie(s.authCookie)
		}
		if s.config.OllamaToken != "" {
			req.Header.Set("Authorization", "Bearer "+s.config.OllamaToken)
		}

		resp, err = s.client.Do(req)
		if err != nil {
			return models.OllamaResponse{}, fmt.Errorf("failed to retry Ollama API call: %v", err)
		}
		defer resp.Body.Close()
	}

	if resp.StatusCode != http.StatusOK {
		return models.OllamaResponse{}, fmt.Errorf("ollama API returned status %d", resp.StatusCode)
	}

	var ollamaResp models.OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return models.OllamaResponse{}, fmt.Errorf("failed to decode response: %v", err)
	}
	ollamaResp.ParseJSON()

	return ollamaResp, nil
}
