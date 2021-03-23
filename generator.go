package generator

import (
	"bufio"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Generator struct
type Generator interface {
	// Generate a random phrase
	Generate() string
}

// generator implements Generator interface
type generator struct {
	Dictionary1 string
	Dictionary2 string

	Dictionary1Lenght int
	Dictionary2Lenght int
}

//New creates a new instance of generator
func New() Generator {
	g := new(generator)

	// set the dictionaries directions
	dir := filepath.Dir(".")
	g.Dictionary1 = dir + "/dictionaries/phrases.txt"
	g.Dictionary2 = dir + "/dictionaries/auxiliaries.txt"

	// get the dictionaries lines
	g.Dictionary1Lenght = linesInFile(g.Dictionary1)
	g.Dictionary2Lenght = linesInFile(g.Dictionary2)

	return g
}

// Generate a random phrase
func (g *generator) Generate() string {
	// get a random phrase
	phrase := getLine(
		g.Dictionary1,
		(rand.Int()%g.Dictionary1Lenght)+1)

	var result string
	var buffer string

	// get all runes in the phrase and parse to string
	for _, r := range phrase {
		buffer = string(r)

		// replace the horizontal line (-) in the phrase to the auxiliary
		if buffer == "-" {
			buffer = ""
			// get a random auxiliary
			auxiliary := getLine(
				g.Dictionary2,
				(rand.Int()%g.Dictionary2Lenght)+1)
			// get all runes in the phrase and parse to string
			for _, r := range auxiliary {
				buffer += string(r)
			}
		}
		result += buffer
	}

	return result
}

// linesInFile counts the lines in the file
func linesInFile(fileName string) int {
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	f.Seek(0, os.SEEK_SET)
	var count int = 0

	for scanner.Scan() {
		scanner.Text()
		count++
	}

	return count
}

// getLine get the line in the n position
func getLine(fileName string, n int) string {
	var result string

	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	f.Seek(0, os.SEEK_SET)

	for i := 0; i < n && scanner.Scan(); i++ {
		result = scanner.Text()
	}

	return result
}
