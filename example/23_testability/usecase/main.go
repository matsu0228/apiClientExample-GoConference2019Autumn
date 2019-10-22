package usecase

import (
	"fmt"
	"image/color"
)

// Noop is interface of api client
type Noop interface {
	Get(url string) (string, error)
}

// DecideColor :色を決定するビジネスロジック
func DecideColor(nClient Noop) (color.RGBA, error) {
	var rgba color.RGBA

	url := "https://api.noopschallenge.com/hexbot"
	color, err := nClient.Get(url)
	if err != nil {
		return rgba, err
	}
	return parseHexColor(color)
}

// parseHexColor is parser for RGBA
func parseHexColor(s string) (color.RGBA, error) {
	var c color.RGBA
	var err error
	c.A = 0xff
	switch len(s) {
	case 7:
		_, err = fmt.Sscanf(s, "#%02x%02x%02x", &c.R, &c.G, &c.B)
	case 4:
		_, err = fmt.Sscanf(s, "#%1x%1x%1x", &c.R, &c.G, &c.B)
		// Double the hex digits:
		c.R *= 17
		c.G *= 17
		c.B *= 17
	default:
		err = fmt.Errorf("invalid length, must be 7 or 4")

	}
	return c, err
}
