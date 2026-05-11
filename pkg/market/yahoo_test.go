package market

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestYahooClientGetQuotes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("symbols") != "AAPL,BBCA.JK" {
			t.Fatalf("unexpected symbols query: %q", r.URL.Query().Get("symbols"))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"spark": {
				"result": [
					{
						"symbol": "AAPL",
						"response": [{"meta": {"regularMarketPrice": 292.68}}]
					},
					{
						"symbol": "BBCA.JK",
						"response": [{"meta": {"regularMarketPrice": 6150}}]
					}
				],
				"error": null
			}
		}`))
	}))
	defer server.Close()

	client := NewYahooClient(server.Client())
	client.baseURL = server.URL

	quotes, err := client.GetQuotes(context.Background(), []string{"aapl", "BBCA.JK", "aapl"})
	if err != nil {
		t.Fatalf("GetQuotes returned error: %v", err)
	}

	if got := quotes["AAPL"]; got != 292.68 {
		t.Fatalf("AAPL price = %v, want 292.68", got)
	}
	if got := quotes["BBCA.JK"]; got != 6150 {
		t.Fatalf("BBCA.JK price = %v, want 6150", got)
	}
}
