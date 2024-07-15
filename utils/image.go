package utils

import (
	"bytes"
	"image/jpeg"
	"image/png"
)

func DataType(data string) string {
	dataBytes := []byte(data)
	reader := bytes.NewReader(dataBytes)

	_, err := png.Decode(reader)
	if err == nil {
		return "png"
	}
	_, err = jpeg.Decode(reader)
	if err == nil {
		return "jpg"
	}

	return "text"
}
