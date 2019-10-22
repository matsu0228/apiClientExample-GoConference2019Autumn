package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
)

func generateAuth(key string) string {
	return fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(key)))
}

func main() {

	url := "https://api.noopschallenge.com/hexbot"
	authKey := "testKey"

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	_ = err
	req.Header.Set("Authorization", generateAuth(authKey))
	resp, err := client.Do(req)
	_ = err

	body, err := ioutil.ReadAll(resp.Body)
	_ = err
	fmt.Printf("\n %s \n", body)
}
