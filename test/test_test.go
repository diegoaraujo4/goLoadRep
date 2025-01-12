package test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunLoadTest(t *testing.T) {
	tests := []struct {
		name                 string
		url                  string
		expectedReport       string
		totalRequests        int
		totalSuccessRequests int
		concurrency          int
		expectedErrorStatus  int
	}{
		{
			name:                 "Test with 10 requests and concurrency of 2",
			url:                  "/test",
			totalRequests:        10,
			totalSuccessRequests: 10,
			concurrency:          2,
			expectedReport:       "Total requests made: 10\nRequests with HTTP status 200: 10",
		},
		{
			name:                 "Test with 5 requests and concurrency of 1",
			url:                  "/test",
			totalRequests:        5,
			totalSuccessRequests: 5,
			concurrency:          1,
			expectedReport:       "Total requests made: 5\nRequests with HTTP status 200: 5",
		},
		{
			name:                 "Test with 5 requests and concurrency of 1 expecting failure",
			url:                  "/test",
			totalRequests:        5,
			totalSuccessRequests: 4,
			concurrency:          1,
			expectedReport:       "Total requests made: 5\nRequests with HTTP status 200: 4\nRequests with HTTP status 500: 1",
			expectedErrorStatus:  http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			successCount := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if successCount < tt.totalSuccessRequests {
					w.WriteHeader(http.StatusOK)
					successCount++
				} else {
					w.WriteHeader(tt.expectedErrorStatus)
				}
			}))
			defer server.Close()

			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			err := RunLoadTest(server.URL+tt.url, tt.totalRequests, tt.concurrency)

			w.Close()
			os.Stdout = old
			var buf bytes.Buffer
			io.Copy(&buf, r)

			consoleOutput := buf.String()
			assert.Contains(t, consoleOutput, tt.expectedReport)
			assert.NoError(t, err)
		})
	}
}
