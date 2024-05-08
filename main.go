package main

import (
	"fmt"
	"image/color"
	"os"
	"strconv"
	"strings"

	"github.com/fogleman/gg"
	"github.com/yorologo/GoPhrasesGenerator/go-generator-phrases"
)

func textFormater(text string) string {
	// Split the text into sentences
	sentences := strings.FieldsFunc(text, func(r rune) bool {
		return r == '|'
	})

	// Capitalize the first letter of each sentence
	for i, sentence := range sentences {
		sentences[i] = strings.TrimSpace(strings.Title(sentence))
	}

	// Join the sentences with a line break
	return strings.Join(sentences, "\n")
}

func createImage(imageName string, phrase string) error {
	// Width heigth
	imgWidth := 1200
	imgHeight := 300

	// Create a new canvas with a white background
	dc := gg.NewContext(imgWidth, imgHeight)
	dc.SetColor(color.White)
	dc.Clear()

	// Define the text and its initial style
	fontSize := 40.0
	fontPath := "fonts/SedanSC-Regular.ttf" // Path of the font you want to use

	// Separate the text into individual lines
	lines := strings.Split(phrase, "\n")

	// Calculate the total height of the text
	totalHeight := float64(len(lines)) * fontSize * 1.5 // Line spacing

	// Calculate the position y to vertically center the text
	y := (float64(imgHeight)-totalHeight)/2 + 25

	// Draw each line on the canvas
	for _, line := range lines {
		// Try to draw the text with the current font size
		err := dc.LoadFontFace(fontPath, fontSize)
		if err != nil {
			return err
		}
		width, _ := dc.MeasureString(line)
		if width > float64(imgWidth) {
			// If the text extends beyond the width of the canvas, reduce the font size
			for width > float64(imgWidth) {
				fontSize -= 1.0
				err := dc.LoadFontFace(fontPath, fontSize)
				if err != nil {
					return err
				}
				width, _ = dc.MeasureString(line)
			}
		}

		// Draw the text on the canvas
		dc.SetColor(color.Black)                                      // Text color
		dc.DrawStringAnchored(line, float64(imgWidth)/2, y, 0.5, 0.5) // x, y coordinates of the center of the text

		//Increment position for next line
		y += fontSize * 1.5 // Adjust the value based on the desired spacing between lines
	}

	// Save the image to a file with the given name
	err := dc.SavePNG(imageName)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	qty := 1               // Default value
	if len(os.Args) == 2 { // Check if an argument was provided
		// Convert the argument to an integer
		var err error
		qty, err = strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Println("The argument must be a valid integer.")
			return
		}
		// Validate that qty is greater than 0
		if qty <= 0 {
			fmt.Println("The number of images must be greater than 0.")
			return
		}
	}

	// Setup the generator
	g := generator.New()

	// Generate the specified number of images
	for i := 1; i <= qty; i++ {
		// Call the function to create the image
		err := createImage(fmt.Sprintf("img/img%d.png", i), textFormater(g.Generate()))
		if err != nil {
			panic(err)
		}
	}
}
