package usecase

import (
	"log"
	"testing"
)

type clientMock struct{}

func (c *clientMock) Get(url string) (string, error) {
	return "#6D5C1C", nil
}

func TestDecideColor(t *testing.T) {

	client := &clientMock{}
	color, err := DecideColor(client)
	if err != nil {
		t.Errorf("invalid. err:%v", err)
	}

	log.Printf("[TEST] color: %#v", color)
}
