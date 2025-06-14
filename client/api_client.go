package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
	"math/rand"

	"github.com/berezovskyivalerii/testtre/model"
)

type APIClient struct {
	http *http.Client
}

func New() *APIClient {
	return &APIClient{
		http: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *APIClient) FetchUsers(ctx context.Context, url string) ([]model.User, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return nil, errors.New(resp.Status)
	}

	body, _ := io.ReadAll(resp.Body)
	var users []model.User
	return users, json.Unmarshal(body, &users)
}

func (c *APIClient) SendUser(ctx context.Context, url string, u model.User) error {
	body, _ := json.Marshal(map[string]string{
		"name":  u.Name,
		"email": u.Email,
	})

	const (
		maxRetries = 3
		baseDelay  = 500 * time.Millisecond
	)

	var lastErr error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		req, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := c.http.Do(req)

		// success
		if err == nil && resp.StatusCode/100 == 2 {
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
			return nil
		}

		// error
		if resp != nil {
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}

		// status 4xx
		if err == nil && resp.StatusCode/100 == 4 {
			return fmt.Errorf("SendUser: hard error %v", resp.Status)
		}
		// network errors
		if ne, ok := err.(net.Error); ok && !ne.Temporary() {
			return err
		}

		lastErr = err

		if ctx.Err() != nil {
			return ctx.Err()
		}

		delay := baseDelay << (attempt - 1)
		jitter := time.Duration(rand.Int63n(int64(delay) / 3))
		time.Sleep(delay + jitter)
	}

	return fmt.Errorf("SendUser: failed after %d attempts, last error: %v", maxRetries, lastErr)
}
