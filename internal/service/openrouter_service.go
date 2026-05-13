package service

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"echobackend/config"
)

type OpenRouterMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenRouterUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type OpenRouterChoice struct {
	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
}

type OpenRouterResponse struct {
	Choices []OpenRouterChoice `json:"choices"`
	Usage   OpenRouterUsage    `json:"usage"`
}

type OpenRouterService interface {
	GenerateResponse(ctx context.Context, messages []OpenRouterMessage, model *string, temperature float64) (*OpenRouterResponse, error)
	GenerateStream(ctx context.Context, messages []OpenRouterMessage, model *string, temperature float64) (<-chan string, <-chan OpenRouterUsage, <-chan error)
}

type openRouterService struct {
	cfg        config.OpenRouterConfig
	httpClient *http.Client
}

func NewOpenRouterService(cfg config.OpenRouterConfig) OpenRouterService {
	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = 90 * time.Second
	}
	return &openRouterService{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (s *openRouterService) GenerateResponse(ctx context.Context, messages []OpenRouterMessage, model *string, temperature float64) (*OpenRouterResponse, error) {
	resp, err := s.callAPI(ctx, messages, model, false, temperature)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var payload OpenRouterResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("failed to decode OpenRouter response: %w", err)
	}
	return &payload, nil
}

func (s *openRouterService) GenerateStream(ctx context.Context, messages []OpenRouterMessage, model *string, temperature float64) (<-chan string, <-chan OpenRouterUsage, <-chan error) {
	chunks := make(chan string)
	usageCh := make(chan OpenRouterUsage, 1)
	errCh := make(chan error, 1)

	go func() {
		defer close(chunks)
		defer close(usageCh)
		defer close(errCh)

		resp, err := s.callAPI(ctx, messages, model, true, temperature)
		if err != nil {
			errCh <- err
			return
		}
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)
		scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

		var usage OpenRouterUsage
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || !strings.HasPrefix(line, "data: ") {
				continue
			}

			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				usageCh <- usage
				return
			}

			var payload struct {
				Choices []struct {
					Delta struct {
						Content string `json:"content"`
					} `json:"delta"`
				} `json:"choices"`
				Usage *OpenRouterUsage `json:"usage"`
			}
			if err := json.Unmarshal([]byte(data), &payload); err != nil {
				continue
			}
			if payload.Usage != nil {
				usage = *payload.Usage
			}
			if len(payload.Choices) == 0 || payload.Choices[0].Delta.Content == "" {
				continue
			}

			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			case chunks <- payload.Choices[0].Delta.Content:
			}
		}

		if err := scanner.Err(); err != nil && err != io.EOF {
			errCh <- fmt.Errorf("failed to read OpenRouter stream: %w", err)
			return
		}
		usageCh <- usage
	}()

	return chunks, usageCh, errCh
}

func (s *openRouterService) callAPI(ctx context.Context, messages []OpenRouterMessage, model *string, stream bool, temperature float64) (*http.Response, error) {
	if s.cfg.APIKey == "" {
		return nil, fmt.Errorf("OPENROUTER_API_KEY is not configured")
	}

	finalModel := strings.TrimSpace(s.cfg.DefaultModel)
	if model != nil && strings.TrimSpace(*model) != "" {
		finalModel = strings.TrimSpace(*model)
	}
	if finalModel == "" {
		return nil, fmt.Errorf("OPENROUTER_DEFAULT_MODEL is not configured")
	}
	if temperature < 0 || temperature > 2 {
		temperature = 0.7
	}

	body, err := json.Marshal(map[string]any{
		"model":       finalModel,
		"messages":    messages,
		"stream":      stream,
		"temperature": temperature,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to encode OpenRouter request: %w", err)
	}

	url := strings.TrimRight(s.cfg.BaseURL, "/") + "/chat/completions"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create OpenRouter request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+s.cfg.APIKey)
	req.Header.Set("Content-Type", "application/json")
	if s.cfg.HTTPReferer != "" {
		req.Header.Set("HTTP-Referer", s.cfg.HTTPReferer)
	}
	if s.cfg.Title != "" {
		req.Header.Set("X-Title", s.cfg.Title)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("OpenRouter request failed: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		defer resp.Body.Close()
		payload, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return nil, fmt.Errorf("OpenRouter API error: %s %s", resp.Status, strings.TrimSpace(string(payload)))
	}

	return resp, nil
}
