package image

import (
	"errors"
	"fmt"
	"image/color"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/fogleman/gg"
	"github.com/yorologo/go-generator-phrases/generator"
)

const (
	defaultWidth    = 1200
	defaultHeight   = 300
	defaultFontPath = "fonts/SedanSC-Regular.ttf"
	defaultOutput   = "img"
)

func formatPhrase(text string) string {
	parts := strings.FieldsFunc(text, func(r rune) bool {
		return r == '|'
	})

	lines := make([]string, 0, len(parts))
	for _, part := range parts {
		line := strings.TrimSpace(part)
		if line == "" {
			continue
		}
		lines = append(lines, capitalizeFirst(line))
	}

	if len(lines) == 0 {
		return strings.TrimSpace(text)
	}

	return strings.Join(lines, "\n")
}

func capitalizeFirst(text string) string {
	first, size := utf8.DecodeRuneInString(text)
	if first == utf8.RuneError && size == 0 {
		return text
	}

	return string(unicode.ToUpper(first)) + text[size:]
}

func pickFontSize(dc *gg.Context, lines []string, width int, height int) (float64, error) {
	const (
		maxFontSize = 56.0
		minFontSize = 18.0
		step        = 1.0
		paddingX    = 64.0
		paddingY    = 48.0
		lineGap     = 1.35
	)

	availableWidth := float64(width) - paddingX*2
	availableHeight := float64(height) - paddingY*2

	for fontSize := maxFontSize; fontSize >= minFontSize; fontSize -= step {
		if err := dc.LoadFontFace(defaultFontPath, fontSize); err != nil {
			return 0, err
		}

		maxWidth := 0.0
		for _, line := range lines {
			lineWidth, _ := dc.MeasureString(line)
			maxWidth = math.Max(maxWidth, lineWidth)
		}

		totalHeight := float64(len(lines)) * fontSize * lineGap
		if maxWidth <= availableWidth && totalHeight <= availableHeight {
			return fontSize, nil
		}
	}

	return minFontSize, dc.LoadFontFace(defaultFontPath, minFontSize)
}

func createImage(imagePath string, phrase string) error {
	lines := strings.Split(phrase, "\n")
	if len(lines) == 0 {
		return errors.New("phrase cannot be empty")
	}

	if err := os.MkdirAll(filepath.Dir(imagePath), 0o755); err != nil {
		return err
	}

	dc := gg.NewContext(defaultWidth, defaultHeight)
	dc.SetColor(color.White)
	dc.Clear()

	fontSize, err := pickFontSize(dc, lines, defaultWidth, defaultHeight)
	if err != nil {
		return err
	}

	lineGap := fontSize * 1.35
	totalHeight := float64(len(lines)) * lineGap
	startY := (float64(defaultHeight)-totalHeight)/2 + lineGap/2

	dc.SetColor(color.Black)
	for i, line := range lines {
		y := startY + float64(i)*lineGap
		dc.DrawStringAnchored(line, float64(defaultWidth)/2, y, 0.5, 0.5)
	}

	return dc.SavePNG(imagePath)
}

func GenerateImage() (string, error) {
	images, err := GenerateImages(1)
	if err != nil {
		return "", err
	}
	return images[0], nil
}

func GenerateImages(imagesNumber int) ([]string, error) {
	if imagesNumber <= 0 {
		return nil, errors.New("imagesNumber must be greater than 0")
	}

	gen, err := generator.New()
	if err != nil {
		return nil, err
	}

	imagePaths := make([]string, 0, imagesNumber)
	baseTime := time.Now().UnixNano()
	for i := 0; i < imagesNumber; i++ {
		imageName := filepath.Join(defaultOutput, fmt.Sprintf("img_%d.png", baseTime+int64(i)))
		if err := createImage(imageName, formatPhrase(gen.Generate())); err != nil {
			return nil, err
		}
		imagePaths = append(imagePaths, imageName)
	}

	return imagePaths, nil
}
