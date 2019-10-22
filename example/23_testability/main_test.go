package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	successFunc = func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "success server")
	}
	errorFunc = func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "error server")
	}
)

func TestGetWithRecodedFile(t *testing.T) {

	client := NoopClient{HTTPClient: &http.Client{}}
	testFile := "testdata/colors.json"
	file, err := os.Open(testFile)
	if err != nil {
		t.Errorf("cant open testfile. err:%v", err)
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(file)
	if err != nil {
		t.Errorf("cant read. err:%v", err)
	}

	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, buf.String())
		}))
	defer ts.Close()

	body, err := client.Get(ts.URL)
	if err != nil {
		t.Errorf("invalid request. err:%v", err)
	}
	log.Printf("[DEBUG] response:%s", body)
}

func TestGet(t *testing.T) {
	// 正常なレスポンスが変えるサーバーと、失敗するサーバーを用意
	ts := httptest.NewServer(http.HandlerFunc(successFunc))
	defer ts.Close()
	fs := httptest.NewServer(http.HandlerFunc(errorFunc))
	defer fs.Close()

	client := NoopClient{HTTPClient: &http.Client{}}

	// 成功する場合はエラーがない
	body, err := client.Get(ts.URL)
	if err != nil {
		t.Errorf("invalid request. err:%v", err)
	}
	log.Printf("[DEBUG] response:%s", body)

	// 失敗する場合はエラーがある
	_, err = client.Get(fs.URL)
	if err == nil {
		t.Errorf("invalid request. err:%v", err)
	}
	log.Printf("[DEBUG] error of failserver:%v", err)
}
