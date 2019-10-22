package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/cenkalti/backoff"
	"github.com/matsu0228/apiClientExample-GoConference2019Autumn/example/23_testability/usecase"
)

// Colors is response of api
type Colors struct {
	Colors []struct {
		Value string `json:"value"`
	} `json:"colors"`
}

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
		log.Printf("[DEBUG] リクエスト %v", req)
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

func decodeBodyWithRecord(resp *http.Response, out interface{}, f *os.File) error {
	defer resp.Body.Close()

	if f != nil {
		resp.Body = ioutil.NopCloser(io.TeeReader(resp.Body, f))
		defer f.Close()
	}

	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(out)
}

// Get is function of example of http.Get
func (c *NoopClient) Get(url string) (string, error) {
	var colors Colors
	testFile := "testdata/colors.json"
	file, err := os.OpenFile(testFile, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	// resp, err := c.HTTPClient.Do(req)
	resp, err := c.DoWithRetry(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("invalid http status")
	}
	err = decodeBodyWithRecord(resp, &colors, file)
	if err != nil || len(colors.Colors) == 0 {
		return "", err
	}
	return colors.Colors[0].Value, nil
}

func main() {
	url := "https://api.noopschallenge.com/hexbot"

	client := &NoopClient{HTTPClient: &http.Client{}}
	body, err := client.Get(url)
	if err != nil {
		log.Fatalf("[ERROR] %v", err)
	}
	fmt.Printf(" %s \n", body)

	color, err := usecase.DecideColor(client)
	_ = err
	fmt.Printf("usecaseColor: %#v", color)
}
