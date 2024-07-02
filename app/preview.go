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
	tea "github.com/charmbracelet/bubbletea"
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

type PreviewModel struct {
	viewport viewport.Model
	ready    bool
}

func (m PreviewModel) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m PreviewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.ready = true

			m.viewport.YPosition = headerHeight + 1
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}

	case tea.KeyMsg:

		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	newViewport, cmd := m.viewport.Update(msg)
	m.viewport = newViewport
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m PreviewModel) headerView() string {
	title := previewTitleStyle.Render("Test")
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m PreviewModel) footerView() string {
	info := previewInfoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func (m PreviewModel) View() string {
	return "" //fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}

func newPreviewModel() PreviewModel {
	return PreviewModel{
		ready: false,
	}
}

func getImgPreview(fp string) string {
	encodedImg := getEncodedImg(fp)

	// Decode the base64 image
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(encodedImg))
	img, _, err := image.Decode(reader)
	if err != nil {
		log.Fatalf("failed to decode image: %v", err)
	}

	// Resize image to fit terminal
	img = resize.Resize(80, 0, img, resize.Lanczos3) // Increase width to 80

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	// Setup terminal profile
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

	// Decode the image
	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Error decoding image:", err)
	}

	// Encode the image to JPEG (or PNG, GIF depending on your image format)
	var buf bytes.Buffer
	err = jpeg.Encode(&buf, img, nil) // Change to png.Encode for PNG images
	if err != nil {
		fmt.Println("Error encoding image:", err)
	}

	// Convert the bytes buffer to a base64 encoded string
	base64Image := base64.StdEncoding.EncodeToString(buf.Bytes())
	return base64Image
}
