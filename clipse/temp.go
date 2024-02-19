package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"os"
	"strings"
	"time"

	"github.com/atotto/clipboard"
)

func main() {

	fileExt := FileType()
	fileName := fmt.Sprintf("%s.%s", string(time.Now().Second()), strings.ToLower(fileExt))
	cmdString := fmt.Sprintf("wl-paste--type image/png > %s", fileName)
	cmd := os.exec("sh", "-c", cmdString)

	time.Sleep(1000 * time.Second)

}

func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

func btoi(v []byte) uint64 {
	return binary.BigEndian.Uint64(v)
}

func saveImage() {
	imageData, err := clipboard.ReadAll()
	if err != nil {
		fmt.Println("Error reading from clipboard:", err)
		return
	}
	dataBytes := []byte(imageData)

	imageType := "png" //FileType(dataBytes)

	switch imageType {
	case "png":
		newFilename := "newImage.png"
		// Create a new file to save the image
		outFile, err := os.Create(newFilename)
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
		defer outFile.Close()

		reader := bytes.NewReader(dataBytes)
		// Decode the image
		img, _, err := image.Decode(reader)
		if err != nil {
			fmt.Println("Error decoding image:", err)
			return
		}
		// Write the image data to the file
		err = png.Encode(outFile, img)
		if err != nil {
			fmt.Println("Error encoding image:", err)
			return
		}

		fmt.Println("Image saved to clipboard_image.png")

	}
}

func FileType(data []byte) string {
	reader := bytes.NewReader(data)
	_, err := png.Decode(reader)
	if err == nil {
		return "png"
	}
	_, err = jpeg.Decode(reader)
	if err == nil {
		return "jpg"
	}

	return ""

}
