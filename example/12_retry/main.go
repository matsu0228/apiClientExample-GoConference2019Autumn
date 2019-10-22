package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cenkalti/backoff"
)

// NoopClient is noopAPI client
type NoopClient struct {
	HTTPClient *http.Client
}

// IsErrorRetryable :リトライ可能な状態かどうか
func (c *NoopClient) IsErrorRetryable(resp *http.Response) bool {
	if resp.StatusCode >= http.StatusInternalServerError {
		return true
	}
	return false
}

// DoWithRetry :リトライ処理を含めたリクエスト実行
func (c *NoopClient) DoWithRetry(req *http.Request) (*http.Response, error) {
	maxRetryNumber := uint64(7)
	var err error
	var resp *http.Response

	operationWithRetry := func() error {
		resp, err = c.HTTPClient.Do(req)
		if err == nil && c.IsErrorRetryable(resp) {
			err = errors.New("リトライ実施可能")
		}
		return err
	}

	bo := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), maxRetryNumber)
	err = backoff.Retry(operationWithRetry, bo)
	return resp, err
}

// Get is function of example of http.Get
func (c *NoopClient) Get(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	// resp, err := c.HTTPClient.Do(req)
	resp, err := c.DoWithRetry(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("invalid http status")
	}
	return ioutil.ReadAll(resp.Body)
}

func main() {
	url := "https://api.noopschallenge.com/hexbot"

	client := NoopClient{HTTPClient: &http.Client{}}
	body, err := client.Get(url)
	_ = err

	fmt.Printf("\n %s \n", body)
}
