# Go-Generator-Phrases

This package generates random phrases based on the game Cards Against Humanity, known for its politically incorrect humor. 

**Note**: This project currently generates phrases in Spanish only.

## Supported Go Versions & Installation

Go-Generator-Phrases requires version 1.14 or higher of Go. You can download Go from golang.org.

To download and install Go-Generator-Phrases, use the following command:

```bash
go get -u github.com/yorologo/go-generator-phrases
```

## Example Usage

```go
package main

import (
    "fmt"
    "github.com/yorologo/go-generator-phrases/generator"
)

func main() {
    // Create a new instance of the generator
    gen, err := generator.New()
    if err != nil {
        fmt.Println("Error creating the generator:", err)
        return
    }

    // Generate a random phrase
    phrase, err := gen.Generate()
    if err != nil {
        fmt.Println("Error generating the phrase:", err)
        return
    }
    fmt.Println("Generated phrase:", phrase)

    // Generate a phrase and get the used indices
    phraseWithIndices, phraseIndex, auxiliaryIndices, err := gen.GenerateWithIndices()
    if err != nil {
        fmt.Println("Error generating the phrase with indices:", err)
        return
    }
    fmt.Println("Generated phrase with indices:", phraseWithIndices)
    fmt.Println("Phrase index:", phraseIndex)
    fmt.Println("Auxiliary indices:", auxiliaryIndices)

    // Reconstruct the same phrase using the indices
    reconstructedPhrase, err := gen.GenerateByID(phraseIndex, auxiliaryIndices)
    if err != nil {
        fmt.Println("Error reconstructing the phrase:", err)
        return
    }
    fmt.Println("Reconstructed phrase:", reconstructedPhrase)
}
```

## Features

- **Random Phrase Generation**: Generates phrases by combining base phrases and auxiliary words.
- **Phrase Identification**: Each generated phrase can be identified and reconstructed using unique indices.
- **Reproducibility**: Allows exact reconstruction of the same phrase based on the provided indices.

## Function Documentation

### `generator.New() (Generator, error)`

Creates a new instance of the generator. Loads phrases and auxiliaries from the dictionary files.

- **Returns**:
  - `Generator`: The instance of the generator.
  - `error`: Error if the loading of the dictionaries fails.

### `Generator.Generate() (string, error)`

Generates a random phrase.

- **Returns**:
  - `string`: The generated phrase.
  - `error`: Error if the generation fails.

### `Generator.GenerateWithIndices() (string, int, []int, error)`

Generates a random phrase and returns the indices used to construct it.

- **Returns**:
  - `string`: The generated phrase.
  - `int`: Index of the base phrase in the dictionary.
  - `[]int`: List of auxiliary indices used.
  - `error`: Error if the generation fails.

### `Generator.GenerateByID(phraseIndex int, auxiliaryIndices []int) (string, error)`

Reconstructs a phrase using the provided indices.

- **Parameters**:
  - `phraseIndex int`: Index of the base phrase.
  - `auxiliaryIndices []int`: List of auxiliary indices.
- **Returns**:
  - `string`: The reconstructed phrase.
  - `error`: Error if the reconstruction fails.

## Dictionary Files

The phrase and auxiliary dictionaries must be located in the `dictionaries` directory within the package. The files should be named `phrases.txt` and `auxiliaries.txt`, respectively.

- `phrases.txt`: Contains the base phrases. Use the `-` character to indicate where an auxiliary word will be inserted.
- `auxiliaries.txt`: Contains the auxiliary words that will replace the dashes in the base phrases.

## Unit Tests

To run the unit tests, use the following command:

```bash
go test -v
```

The tests cover the following aspects:

- Generation of random phrases.
- Generation of phrases with indices.
- Reconstruction of phrases from indices.
- Error handling when providing invalid indices.

## Contributions

Contributions are welcome. If you find any issues or have suggestions for improvements, feel free to open an issue or a pull request.

## License

This project is licensed under the MIT License. Please see the LICENSE file for more information.
