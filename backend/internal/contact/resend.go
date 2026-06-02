package contact

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Resend struct {
	apiKey string
	client *http.Client
}

func NewResend(apiKey string) *Resend {
	return &Resend{
		apiKey: apiKey,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

type resendRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	ReplyTo string   `json:"reply_to,omitempty"`
	Subject string   `json:"subject"`
	Text    string   `json:"text"`
	HTML    string   `json:"html,omitempty"`
}

type resendResponse struct {
	ID string `json:"id"`
}

type Email struct {
	From    string
	To      string
	ReplyTo string
	Subject string
	Text    string
	HTML    string
}

// Send returns the Resend message ID on success.
func (r *Resend) Send(ctx context.Context, e Email) (string, error) {
	if r.apiKey == "" {
		return "", errors.New("resend: missing api key")
	}
	body, err := json.Marshal(resendRequest{
		From:    e.From,
		To:      []string{e.To},
		ReplyTo: e.ReplyTo,
		Subject: e.Subject,
		Text:    e.Text,
		HTML:    e.HTML,
	})
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://api.resend.com/emails", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+r.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("resend: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		buf, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("resend status %d: %s", resp.StatusCode, string(buf))
	}
	var out resendResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", fmt.Errorf("resend decode: %w", err)
	}
	return out.ID, nil
}
