package utils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// It returns the response body as a byte slice or an error if the request fails.
func DoRequest(client *http.Client, ctx context.Context, method, endpoint string, headers map[string]string, params url.Values, body []byte) ([]byte, error) {
	reqURL := buildRequestURL(endpoint, params)
	reqBody := buildRequestBody(body)

	req, err := http.NewRequestWithContext(ctx, method, reqURL, reqBody)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := handleHTTPError(resp); err != nil {
		return nil, err
	}
	return io.ReadAll(resp.Body)
}

// handleHTTPError checks the HTTP response status and returns an error if it's not successful.
func handleHTTPError(resp *http.Response) error {
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("HTTP error: status code %d", resp.StatusCode)
	}
	return nil
}

// buildRequestURL constructs the full request URL with parameters.
func buildRequestURL(endpoint string, params url.Values) string {
	if params == nil {
		return endpoint
	}
	u, _ := url.Parse(endpoint)
	u.RawQuery = u.Query().Encode() + "&" + params.Encode()
	return u.String()
}

// buildRequestBody creates the request body from a byte slice.
func buildRequestBody(body []byte) io.Reader {
	var reqBody io.Reader
	if body != nil {
		reqBody = bytes.NewBuffer(body)
	}
	return reqBody
}
