package main

import (
	"errors"
	"flag"
	"fmt"
	"image/color"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/fogleman/gg"
	generator "github.com/yorologo/go-generator-phrases/go-generator-phrases"
)

const (
	defaultQty      = 1
	defaultWidth    = 1200
	defaultHeight   = 300
	defaultFontPath = "fonts/SedanSC-Regular.ttf"
	defaultOutDir   = "img"
	defaultPrefix   = "img"
)

type config struct {
	qty      int
	width    int
	height   int
	output   string
	prefix   string
	fontPath string
}

func parseConfig(args []string) (config, error) {
	cfg := config{
		qty:      defaultQty,
		width:    defaultWidth,
		height:   defaultHeight,
		output:   defaultOutDir,
		prefix:   defaultPrefix,
		fontPath: defaultFontPath,
	}

	fs := flag.NewFlagSet("go-generator-phrases", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	fs.IntVar(&cfg.qty, "n", cfg.qty, "number of images to generate")
	fs.IntVar(&cfg.width, "width", cfg.width, "image width in pixels")
	fs.IntVar(&cfg.height, "height", cfg.height, "image height in pixels")
	fs.StringVar(&cfg.output, "out", cfg.output, "output directory")
	fs.StringVar(&cfg.prefix, "prefix", cfg.prefix, "output file prefix")
	fs.StringVar(&cfg.fontPath, "font", cfg.fontPath, "font path")

	if err := fs.Parse(args); err != nil {
		return config{}, err
	}

	if fs.NArg() > 1 {
		return config{}, errors.New("expected at most one positional argument for quantity")
	}

	if fs.NArg() == 1 {
		qty, err := strconv.Atoi(fs.Arg(0))
		if err != nil {
			return config{}, fmt.Errorf("invalid quantity %q: %w", fs.Arg(0), err)
		}
		cfg.qty = qty
	}

	if cfg.qty <= 0 {
		return config{}, errors.New("the number of images must be greater than 0")
	}
	if cfg.width <= 0 || cfg.height <= 0 {
		return config{}, errors.New("width and height must be greater than 0")
	}
	if strings.TrimSpace(cfg.output) == "" {
		return config{}, errors.New("output directory cannot be empty")
	}
	if strings.TrimSpace(cfg.prefix) == "" {
		return config{}, errors.New("output prefix cannot be empty")
	}

	return cfg, nil
}

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

func pickFontSize(dc *gg.Context, lines []string, fontPath string, width int, height int) (float64, error) {
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
		if err := dc.LoadFontFace(fontPath, fontSize); err != nil {
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

	return minFontSize, dc.LoadFontFace(fontPath, minFontSize)
}

func createImage(imagePath string, phrase string, cfg config) error {
	lines := strings.Split(phrase, "\n")
	if len(lines) == 0 {
		return errors.New("phrase cannot be empty")
	}

	if err := os.MkdirAll(filepath.Dir(imagePath), 0o755); err != nil {
		return err
	}

	dc := gg.NewContext(cfg.width, cfg.height)
	dc.SetColor(color.White)
	dc.Clear()

	fontSize, err := pickFontSize(dc, lines, cfg.fontPath, cfg.width, cfg.height)
	if err != nil {
		return err
	}

	lineGap := fontSize * 1.35
	totalHeight := float64(len(lines)) * lineGap
	startY := (float64(cfg.height)-totalHeight)/2 + lineGap/2

	dc.SetColor(color.Black)
	for i, line := range lines {
		y := startY + float64(i)*lineGap
		dc.DrawStringAnchored(line, float64(cfg.width)/2, y, 0.5, 0.5)
	}

	return dc.SavePNG(imagePath)
}

func generateImages(cfg config) error {
	g, err := generator.New()
	if err != nil {
		return err
	}

	for i := 1; i <= cfg.qty; i++ {
		imagePath := filepath.Join(cfg.output, fmt.Sprintf("%s%d.png", cfg.prefix, i))
		if err := createImage(imagePath, formatPhrase(g.Generate()), cfg); err != nil {
			return fmt.Errorf("create %s: %w", imagePath, err)
		}
	}

	return nil
}

func main() {
	cfg, err := parseConfig(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := generateImages(cfg); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
