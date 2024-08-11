# Go-Generator-Phrases

## Description

This project is a web application written in Go that generates and displays images with random phrases. It uses HTML templates for presentation and provides an API for loading new images. **This package generates random phrases based on the _Cards Against Humanity_ card game, known for its **politically incorrect** humor.**

## Project Structure

The project is structured as follows:

```
.
├── README.md
├── fonts
│   ├── Roboto-Regular.ttf
│   └── SedanSC-Regular.ttf
├── generator
│   ├── dictionaries
│   │   ├── auxiliaries.txt
│   │   └── phrases.txt
│   ├── generator.go
│   └── generator_test.go
├── go.mod
├── go.sum
├── image
│   └── image.go
├── main.go
├── main.go.bkp
└── templates
    └── index.html
```

### Files and Directories

- **README.md**: Main documentation of the project.
- **fonts/**: Font files used in the HTML templates.
- **generator/**: Files related to phrase generation.
  - **dictionaries/**: Text files with dictionaries used for generation.
    - **auxiliaries.txt**: Auxiliary words for phrase generation.
    - **phrases.txt**: Predefined phrases for generation.
  - **generator.go**: Code for phrase generation.
  - **generator_test.go**: Tests for the generator.
- **go.mod** and **go.sum**: Go dependency management files.
- **image/**: Source code for handling images.
  - **image.go**: Code for generating and formatting images.
- **main.go**: Main file that sets up and starts the HTTP server.
- **main.go.bkp**: Backup copy of the `main.go` file.
- **templates/**: HTML templates used by the server.
  - **index.html**: HTML template for displaying images.

### Dependencies

The project uses the Go standard library and the `github.com/fogleman/gg` package for image creation, as well as the `github.com/yorologo/GoPhrasesGenerator/generator` package for phrase generation.

## Modules

### `generator/generator.go`

This module handles generating random phrases using two dictionaries.

#### Main Functions

- **`init()`**: Initializes the random number generator with a seed based on the current time.

- **`New()`**: Creates a new instance of the phrase generator. Sets up the dictionaries and counts the lines in each file.

- **`Generate()`**: Generates a random phrase. Selects a random line from the main dictionary and replaces hyphens (`-`) with auxiliary words from the second dictionary.

- **`linesInFile(fileName string)`**: Counts the number of lines in the specified file.

- **`getLine(fileName string, n int)`**: Retrieves the line at position `n` from the specified file.

- **`getPackagePath()`**: Retrieves the path of the current package directory.

### `image/image.go`

This module is responsible for creating and formatting images with generated phrases.

#### Main Functions

- **`textFormater(text string) string`**: Formats the text by capitalizing the first letter of each sentence and joining the sentences with line breaks.

- **`createImage(imageName string, phrase string) error`**: Creates an image with the given text and saves the image with the specified name. Adjusts the font size if the text exceeds the image width.

- **`GenerateImage() (string, error)`**: Generates an image with a random phrase, saves the image, and returns the file path.

- **`GenerateImages(ImagesNumber int) ([]string, error)`**: Generates a specified number of images, each with a random phrase, and returns a list of paths to the generated files.

## Operation

### `main.go`

The `main.go` file is the entry point of the application. It sets up and starts the HTTP server, defines routes for serving the interface, and handling the generation of new images.

#### Routes

- **`/`**: Serves the `index.html` template to display the user interface.

- **`/load-more-images`**: Generates and returns a JSON response with paths to new images.

### Running the Application

1. Make sure you have Go installed.
2. Navigate to the project directory.
3. Run `go run main.go` to start the server.
4. Open a browser and visit `http://localhost:8080`.

