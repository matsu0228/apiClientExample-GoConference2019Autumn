package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// Colors is response of api
type Colors struct {
	Colors []struct {
		Value string `json:"value"`
	} `json:"colors"`
}

func decodeBody(body io.Reader, out interface{}) error {
	decoder := json.NewDecoder(body)
	return decoder.Decode(out)
}

// NoopClient is APIclient
type NoopClient struct {
	URL            *url.URL
	maxRetryNumber uint64
	HTTPClient     *http.Client
	DefaultHeader  http.Header
	authHeader     string
}

// NewNoopClient is constructor
func NewNoopClient(endpointURL, secretKey, userAgent string, maxRetry uint64, httpClient *http.Client) (*NoopClient, error) {
	if len(endpointURL) == 0 {
		return nil, errors.New("invalid url")
	}
	parsedURL, err := url.ParseRequestURI(endpointURL)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse url: %s", endpointURL)
	}
	client := &NoopClient{
		URL:            parsedURL,
		maxRetryNumber: maxRetry,
		HTTPClient:     httpClient,
		DefaultHeader:  make(http.Header),
	}

	client.authHeader = secretKey //TODO

	if userAgent != "" {
		client.DefaultHeader.Set("User-Agent", userAgent)
	}
	client.DefaultHeader.Set("Content-Type", "application/json; charset=utf-8")
	return client, nil
}

// RequestOptions is custom option for each request
type RequestOptions struct {
	Params  map[string]string
	Headers map[string]string
	Body    io.Reader
}

// rawRequest 汎用的なrequest作成
func (c *NoopClient) rawRequest(method, subPath string, ro *RequestOptions) (*http.Request, error) {
	if method == "" {
		return nil, errors.New("missing requestMethod")
	}
	if ro == nil {
		return nil, errors.New("missing RequestOptions")
	}
	//URLに subPath/paramsを追加
	u := *c.URL
	u.Path = path.Join(c.URL.Path, subPath)
	var params = make(url.Values)
	for k, v := range ro.Params {
		params.Add(k, v)
	}
	u.RawQuery = params.Encode()

	request, err := http.NewRequest(method, u.String(), ro.Body)
	if err != nil {
		return nil, err
	}
	// default headers をセット
	for k, v := range c.DefaultHeader {
		request.Header[k] = v
	}
	// request headersの追加分はここでセット
	for k, v := range ro.Headers {
		request.Header.Add(k, v)
	}
	return request, nil
}

// Get is function of example of http.Get
func (c *NoopClient) Get() (Colors, error) {

	var colors Colors
	req, err := c.rawRequest("GET", "hexbot", &RequestOptions{
		Body: nil,
		Headers: map[string]string{
			"Authorization": c.authHeader,
		},
	})
	resp, err := c.HTTPClient.Do(req)

	_ = err
	defer resp.Body.Close()

	err = decodeBody(resp.Body, &colors)
	return colors, err
}

func main() {

	logger, _ := zap.NewDevelopment()
	endpointURL := "https://api.noopschallenge.com"
	secretKey := "secret"

	client, err := NewNoopClient(endpointURL, secretKey, "", 7, &http.Client{})
	_ = err

	colors, err := client.Get()
	_ = err

	logger.Info("get colors", zap.Reflect("colors", colors))
}
