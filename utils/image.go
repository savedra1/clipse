package utils

import (
	"bytes"
	"image/gif"
	"image/jpeg"
	"image/png"
)

func DataType(data string) string {
	/*
	   Confirms if clipboard data is currently folding a file vs a string
	*/
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
	_, err = gif.Decode(reader)
	if err == nil {
		return "gif"
	}

	return "text"
}
