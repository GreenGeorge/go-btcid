package btcid

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	stubAPIKey     = "71239123s"
	stubSecret     = "hush"
	stubBaseURL    = "https://www.test.com"
	stubHTTPClient = http.Client{Timeout: time.Second * 5}
)

func TestNew(t *testing.T) {
	tests := map[string]struct {
		inputAPIKey        string
		inputSecret        string
		inputDomain        string
		inputHTTPClient    *http.Client
		expectedHTTPClient *http.Client
	}{
		"Valid params test": {
			inputAPIKey:        stubAPIKey,
			inputSecret:        stubSecret,
			inputDomain:        stubBaseURL,
			inputHTTPClient:    &stubHTTPClient,
			expectedHTTPClient: &stubHTTPClient,
		},
		"Nil http client test": {
			inputAPIKey:        stubAPIKey,
			inputSecret:        stubSecret,
			inputDomain:        stubBaseURL,
			inputHTTPClient:    nil,
			expectedHTTPClient: http.DefaultClient,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := New(test.inputAPIKey, test.inputSecret, test.inputHTTPClient)
			client.Domain = test.inputDomain
			assert.Equal(t, test.inputAPIKey, client.APIKey)
			assert.Equal(t, test.inputSecret, client.Secret)
			assert.Equal(t, test.expectedHTTPClient, client.HTTPClient)
		})
	}
}

func TestGetTicker(t *testing.T) {
	ticker := struct {
		Ticker Ticker `json:"ticker"`
	}{
		Ticker: Ticker{
			High: "2500",
			Low:  "2500",
			Last: "2500",
			Buy:  "2500",
			Sell: "2500",
		},
	}
	tickerJSON, _ := json.Marshal(ticker)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, string(tickerJSON))
	})
	ts := httptest.NewTLSServer(handler)
	defer ts.Close()

	client := New(stubAPIKey, stubSecret, ts.Client())
	client.Domain = ts.URL
	tests := map[string]struct {
	}{
		"Valid request": {},
	}
	for name := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := client.GetTicker()
			assert.Nil(t, err)
		})
	}

}
