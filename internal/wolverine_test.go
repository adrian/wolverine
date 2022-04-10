package wolverine

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

// This RoundTrip function is passed to http.Client and allows us to mock
// the http response.
type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func TestMonitorURL(t *testing.T) {
	testCases := []struct {
		url        string
		statusCode int
	}{
		{"https://httpstat.us/503", 503},
		{"https://httpstat.us/200", 200},
	}
	for _, tc := range testCases {
		mockHTTPClient := NewTestClient(func(req *http.Request) *http.Response {
			assert.Equal(t, "HEAD", req.Method)
			assert.Equal(t, tc.url, req.URL.String())
			return &http.Response{
				StatusCode: tc.statusCode,
			}
		})

		_, resp := MonitorURL(tc.url, mockHTTPClient)
		assert.Equal(t, tc.statusCode, resp.StatusCode)
	}

}
