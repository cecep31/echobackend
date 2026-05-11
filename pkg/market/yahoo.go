package market

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	yahooSparkURL        = "https://query1.finance.yahoo.com/v7/finance/spark"
	maxSymbolsPerRequest = 50
)

type QuoteClient interface {
	GetQuotes(ctx context.Context, symbols []string) (map[string]float64, error)
}

type YahooClient struct {
	httpClient *http.Client
	baseURL    string
}

func NewYahooClient(httpClient *http.Client) *YahooClient {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}
	return &YahooClient{
		httpClient: httpClient,
		baseURL:    yahooSparkURL,
	}
}

func (c *YahooClient) GetQuotes(ctx context.Context, symbols []string) (map[string]float64, error) {
	normalized := normalizeSymbols(symbols)
	if len(normalized) == 0 {
		return map[string]float64{}, nil
	}

	quotes := make(map[string]float64, len(normalized))
	for start := 0; start < len(normalized); start += maxSymbolsPerRequest {
		end := start + maxSymbolsPerRequest
		if end > len(normalized) {
			end = len(normalized)
		}

		batchQuotes, err := c.fetchSpark(ctx, normalized[start:end])
		if err != nil {
			return nil, err
		}
		for symbol, price := range batchQuotes {
			quotes[symbol] = price
		}
	}

	return quotes, nil
}

func (c *YahooClient) fetchSpark(ctx context.Context, symbols []string) (map[string]float64, error) {
	query := url.Values{}
	query.Set("symbols", strings.Join(symbols, ","))
	query.Set("range", "1d")
	query.Set("interval", "1d")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"?"+query.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("create yahoo finance request: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request yahoo finance quotes: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read yahoo finance response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("yahoo finance returned status %d", resp.StatusCode)
	}

	var payload sparkResponse
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("decode yahoo finance response: %w", err)
	}
	if payload.Spark.Error != nil {
		return nil, fmt.Errorf("yahoo finance error: %s", payload.Spark.Error.Description)
	}

	quotes := make(map[string]float64, len(payload.Spark.Result))
	for _, result := range payload.Spark.Result {
		if len(result.Response) == 0 {
			continue
		}
		price := result.Response[0].Meta.RegularMarketPrice
		if price <= 0 {
			continue
		}
		quotes[strings.ToUpper(strings.TrimSpace(result.Symbol))] = price
	}

	return quotes, nil
}

func normalizeSymbols(symbols []string) []string {
	seen := make(map[string]struct{}, len(symbols))
	normalized := make([]string, 0, len(symbols))
	for _, symbol := range symbols {
		symbol = strings.ToUpper(strings.TrimSpace(symbol))
		if symbol == "" {
			continue
		}
		if _, ok := seen[symbol]; ok {
			continue
		}
		seen[symbol] = struct{}{}
		normalized = append(normalized, symbol)
	}
	return normalized
}

type sparkResponse struct {
	Spark struct {
		Result []struct {
			Symbol   string `json:"symbol"`
			Response []struct {
				Meta struct {
					RegularMarketPrice float64 `json:"regularMarketPrice"`
				} `json:"meta"`
			} `json:"response"`
		} `json:"result"`
		Error *struct {
			Code        string `json:"code"`
			Description string `json:"description"`
		} `json:"error"`
	} `json:"spark"`
}
