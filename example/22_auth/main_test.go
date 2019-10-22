package main

import "testing"

func TestGenerateAuth(t *testing.T) {
	key := "aabbcc:ddeecc"
	want := "Basic YWFiYmNjOmRkZWVjYw=="
	got := generateAuth(key)
	if got != want {
		t.Errorf("generated invalid key. want:%v, got:%v", want, got)
	}
}
