package chawk

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	endpoints "github.com/ToyGuy22/chawk/endpoints"
)

var (
	ErrInsufficientPrivileges = errors.New("insufficient privileges")
	ErrClientInitError        = errors.New("clientID, clientSecret, and baseURL are required")
	ErrTokenExpired           = errors.New("token expired")
	ErrTokenMissing           = errors.New("token file missing")
)

const (
	HTTP_TIMEOUT_SECS = 10
)

type BlackboardClient struct {
	mu           sync.RWMutex
	clientID     string
	clientSecret string
	BaseURL      string
	token        *Token
	tokenFile    string
	httpClient   *http.Client

	// Sub Services
	Users        *UserService
	Courses      *CourseService
	Announcement *AnnouncementService
}

// NewClient initializes and returns a new Blackboard API Client.
// It requires a Client ID and Client Secret obtained from the
// Blackboard Developer Portal.
func NewClient(clientID, clientSecret, baseURL, tokenPath string) (*BlackboardClient, error) {
	if clientID == "" || clientSecret == "" || baseURL == "" {
		return nil, ErrClientInitError
	}

	if tokenPath == "" {
		tokenPath = "data/.token.json"
	}

	client := &BlackboardClient{
		clientID:     clientID,
		clientSecret: clientSecret,
		BaseURL:      baseURL,
		tokenFile:    tokenPath,
		httpClient:   &http.Client{Timeout: HTTP_TIMEOUT_SECS * time.Second},
	}

	client.Users = &UserService{client: client}
	client.Courses = &CourseService{client: client}
	client.Announcement = &AnnouncementService{client: client}

	// Attempt to load token from file, ignore error if file missing or expired
	// We will make a new one later
	client.loadToken()

	return client, nil
}

// loadToken loads token data from a JSON file
func (c *BlackboardClient) loadToken() error {
	file, err := os.Open(c.tokenFile)
	if err != nil {
		if os.IsNotExist(err) {
			return ErrTokenMissing
		}
		return fmt.Errorf("failed to open token file: %w", err)
	}
	defer file.Close()

	t := &Token{}

	if err := json.NewDecoder(file).Decode(t); err != nil {
		return fmt.Errorf("failed to decode token file: %w", err)
	}

	if t.IsExpired() {
		return ErrTokenExpired
	}

	c.mu.Lock()
	c.token = t
	c.mu.Unlock()

	return nil
}

// saveToken saves token data to a JSON file
func (c *BlackboardClient) saveToken(t *Token) error {
	// Ensure directory exists
	dir := filepath.Dir(c.tokenFile)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create token directory: %w", err)
	}

	// Marshal and write to file
	data, err := json.Marshal(t)
	if err != nil {
		return fmt.Errorf("failed to marshal token: %w", err)
	}

	if err := os.WriteFile(c.tokenFile, data, 0600); err != nil {
		return fmt.Errorf("failed to write token file: %w", err)
	}

	c.token = t
	return nil
}

// requestNewToken performs OAuth2 client credentials flow to get a new token
func (c *BlackboardClient) requestNewToken(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Do we evern need a new token?
	if c.token != nil && !c.token.IsExpired() {
		return nil
	}

	url := c.BaseURL + endpoints.GetToken()
	form := "grant_type=client_credentials"

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBufferString(form))
	if err != nil {
		return fmt.Errorf("failed to create token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(c.clientID, c.clientSecret)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute token request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, MAX_RESPONSE_SIZE))
		return fmt.Errorf("failed to get token: status %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	// make token
	newToken := &Token{}
	if err := json.NewDecoder(resp.Body).Decode(newToken); err != nil {
		return fmt.Errorf("failed to decode token response: %w", err)
	}

	newToken.Expiry = time.Now().Add(time.Duration(newToken.ExpiresIn) * time.Second)

	if newToken.AccessToken == "" {
		return errors.New("received empty access token from blackboard")
	}

	//c.token = newToken

	// save toke to cache file
	// Failing should be okay. Check this later
	err = c.saveToken(newToken)
	if err != nil {
		fmt.Printf("Warning: failed to save token to file: %v\n", err)
	}

	return nil
}

func (c *BlackboardClient) sendRequest(ctx context.Context, method string, path string, body io.Reader) (*http.Response, error) {
	// If the token is nil OR expired, try to get a new one.
	// requestNewToken handles its own locking, so no need to call it before this.
	if c.token == nil || c.token.IsExpired() {
		if err := c.requestNewToken(ctx); err != nil {
			return nil, fmt.Errorf("auth failure: %w", err)
		}
	}

	// Store token in memeory.
	c.mu.RLock()
	token := c.token
	c.mu.RUnlock()

	url := c.BaseURL + path

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	req.Header.Set("Accept", "application/json")

	if body != nil && (method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch) {
		req.Header.Set("Content-Type", "application/json")
	}

	return c.httpClient.Do(req)
}

func (c *BlackboardClient) Get(ctx context.Context, path string) (*http.Response, error) {
	return c.sendRequest(ctx, "GET", path, nil)
}

func (c *BlackboardClient) Post(ctx context.Context, path string, jsonBody interface{}) (*http.Response, error) {
	body, err := json.Marshal(jsonBody)
	if err != nil {
		return nil, err
	}
	return c.sendRequest(ctx, "POST", path, bytes.NewReader(body))
}

func (c *BlackboardClient) Put(ctx context.Context, path string, jsonBody interface{}) (*http.Response, error) {
	body, err := json.Marshal(jsonBody)
	if err != nil {
		return nil, err
	}
	return c.sendRequest(ctx, "PUT", path, bytes.NewReader(body))
}

func (c *BlackboardClient) Patch(ctx context.Context, path string, jsonBody interface{}) (*http.Response, error) {
	body, err := json.Marshal(jsonBody)
	if err != nil {
		return nil, err
	}
	return c.sendRequest(ctx, http.MethodPatch, path, bytes.NewReader(body))
}

func (c *BlackboardClient) Delete(ctx context.Context, path string) (*http.Response, error) {
	return c.sendRequest(ctx, http.MethodDelete, path, nil)
}

// GetRemainingCalls returns the number of apis calls left of the key. This does consume a call.
func (c *BlackboardClient) GetRemainingCalls(ctx context.Context) (int, error) {
	// Hit a lightweight endpoint to get headers
	resp, err := c.Get(ctx, "/learn/api/public/v1/users?limit=1")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	remaining := resp.Header.Get("X-Rate-Limit-Remaining")
	if remaining == "" {
		return 0, errors.New("X-Rate-Limit-Remaining header not found")
	}

	var remCalls int
	_, err = fmt.Sscanf(remaining, "%d", &remCalls)
	if err != nil {
		return 0, fmt.Errorf("failed to parse remaining calls: %w", err)
	}
	return remCalls, nil
}
