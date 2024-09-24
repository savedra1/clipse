package app

import (
	"bytes"
	"fmt"
	"image"
	"image/color/palette"
	"image/draw"
	"os"
	"strings"

	"github.com/BourgeoisBear/rasterm"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/muesli/termenv"
	"github.com/nfnt/resize"

	"github.com/savedra1/clipse/config"
	"github.com/savedra1/clipse/utils"
)

func NewPreview() viewport.Model {
	return viewport.New(20, 40) // default sizing updated on tea.WindowSizeMsg
}

func getImgPreview(fp string, windowWidth int, windowHeight int) string {
	img, err := getDecodedImg(fp)
	if err != nil {
		utils.LogERROR(fmt.Sprintf("failed to decode image file for preview | %s", err))
		return fmt.Sprintf("failed to open image file for preview | %s", err)
	}
	switch config.ClipseConfig.ImageDisplay.Type {
	case "sixel":
		return getSixelString(img, windowWidth, windowHeight)
	case "kitty":
		return getKittyString(img, windowWidth, windowHeight)
	default:
		return getBasicString(img, windowWidth)
	}
}

func getBasicString(img image.Image, windowSize int) string {
	img = resize.Resize(uint(windowSize), 0, img, resize.Lanczos3)

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	p := termenv.ColorProfile()

	var sb strings.Builder

	for y := 0; y < height; y += 2 {
		for x := 0; x < width; x++ {
			upperR, upperG, upperB, _ := img.At(x, y).RGBA()
			lowerR, lowerG, lowerB, _ := img.At(x, y+1).RGBA()

			upperColor := p.Color(
				fmt.Sprintf(
					"#%02x%02x%02x", uint8(upperR>>8), uint8(upperG>>8), uint8(upperB>>8),
				),
			)
			lowerColor := p.Color(
				fmt.Sprintf(
					"#%02x%02x%02x", uint8(lowerR>>8), uint8(lowerG>>8), uint8(lowerB>>8),
				),
			)
			sb.WriteString(termenv.String("▀").Foreground(lowerColor).Background(upperColor).String())
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func getSixelString(img image.Image, windowWidth int, windowHeight int) string {
	img = smartResize(img, windowWidth, windowHeight)
	palettedImg := image.NewPaletted(img.Bounds(), palette.Plan9)
	draw.FloydSteinberg.Draw(palettedImg, img.Bounds(), img, image.Point{})
	var buf bytes.Buffer
	err := rasterm.SixelWriteImage(&buf, palettedImg)
	if err != nil {
		utils.LogERROR(fmt.Sprintf("failed to decode image file to sixel | %s", err))
		return fmt.Sprintf("failed to decode image file to sixel | %s", err)
	}
	return buf.String()
}

func getKittyString(img image.Image, windowWidth int, windowHeight int) string {
	img = smartResize(img, windowWidth, windowHeight)
	var buf bytes.Buffer
	var opts rasterm.KittyImgOpts
	err := rasterm.KittyWriteImage(&buf, img, opts)
	if err != nil {
		utils.LogERROR(fmt.Sprintf("failed to decode image file to kitty | %s", err))
		return fmt.Sprintf("failed to decode image file to kitty | %s", err)
	}
	return buf.String()
}

func smartResize(img image.Image, windowWidth int, windowHeight int) image.Image {
	maxWidth := windowWidth * config.ClipseConfig.ImageDisplay.ScaleX
	maxHeight := windowHeight * config.ClipseConfig.ImageDisplay.ScaleY
	imageWidth := img.Bounds().Dx()
	imageHeight := img.Bounds().Dy()
	if imageWidth/imageHeight > maxWidth/maxHeight {
		return resize.Resize(uint(maxWidth), 0, img, resize.Lanczos3)
	}
	return resize.Resize(0, uint(maxHeight), img, resize.Lanczos3)
}

func getDecodedImg(fp string) (image.Image, error) {
	file, err := os.Open(fp)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}
