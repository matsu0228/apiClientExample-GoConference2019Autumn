package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// Get is function of example of http.Get
func Get(url string) {
	resp, err := http.Get(url)
	_ = err
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	_ = err
	fmt.Printf("\n %s \n", body)
}

func main() {
	url := "https://api.noopschallenge.com/hexbot"
	Get(url)
}
