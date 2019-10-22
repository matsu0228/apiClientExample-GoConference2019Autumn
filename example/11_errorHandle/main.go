package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
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

// Get is function of example of http.Get
func (c *NoopClient) Get(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if c.IsErrorRetryable(resp) {
		// TODO: リトライ処理
	}

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
