package main

import (
	"encoding/json"
	"io"
	"net/http"

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

// Get is function of example of http.Get
func Get(url string) {

	var colors Colors
	resp, err := http.Get(url)
	_ = err
	defer resp.Body.Close()

	err = decodeBody(resp.Body, &colors)
	_ = err

	logger, _ := zap.NewDevelopment()
	logger.Info("get colors", zap.Reflect("colors", colors))
}

func main() {
	url := "https://api.noopschallenge.com/hexbot"
	Get(url)
}
