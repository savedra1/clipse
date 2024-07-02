package app

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/nfnt/resize"
)

var (
	previewTitleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	previewInfoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return previewTitleStyle.BorderStyle(b)
	}()
)

func headerView(v viewport.Model) string {
	title := previewTitleStyle.Render(previewHeader)
	line := strings.Repeat("─", max(0, v.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func footerView(v viewport.Model) string {
	info := previewInfoStyle.Render(fmt.Sprintf("%3.f%%", v.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, v.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func NewPreview() viewport.Model {
	return viewport.New(20, 40)
}

func getImgPreview(fp string) string {
	encodedImg := getEncodedImg(fp)

	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(encodedImg))
	img, _, err := image.Decode(reader)
	if err != nil {
		log.Fatalf("failed to decode image: %v", err)
	}

	// Resize image to fit terminal
	img = resize.Resize(80, 0, img, resize.Lanczos3)

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	p := termenv.ColorProfile()

	var sb strings.Builder

	for y := 0; y < height; y += 2 {
		for x := 0; x < width; x++ {
			upperR, upperG, upperB, _ := img.At(x, y).RGBA()
			lowerR, lowerG, lowerB, _ := img.At(x, y+1).RGBA()

			upperColor := p.Color(fmt.Sprintf("#%02x%02x%02x", uint8(upperR>>8), uint8(upperG>>8), uint8(upperB>>8)))
			lowerColor := p.Color(fmt.Sprintf("#%02x%02x%02x", uint8(lowerR>>8), uint8(lowerG>>8), uint8(lowerB>>8)))

			sb.WriteString(termenv.String("▀").Foreground(lowerColor).Background(upperColor).String())
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func getEncodedImg(fp string) string {
	file, err := os.Open(fp) // Change to your image path
	if err != nil {
		fmt.Println("Error opening file:", err)

	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Error decoding image:", err)
	}

	var buf bytes.Buffer
	err = jpeg.Encode(&buf, img, nil)
	if err != nil {
		fmt.Println("Error encoding image:", err)
	}

	base64Image := base64.StdEncoding.EncodeToString(buf.Bytes())
	return base64Image
}
