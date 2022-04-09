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

// TOOO structure these as a table
// Mock the duration the http requests take

func TestMonitorURL(t *testing.T) {
	mockHTTPClient := NewTestClient(func(req *http.Request) *http.Response {
		assert.Equal(t, "HEAD", req.Method)
		return &http.Response{
			StatusCode: 200,
		}
	})

	err, resp := MonitorURL("https://www.google.com", mockHTTPClient)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}
