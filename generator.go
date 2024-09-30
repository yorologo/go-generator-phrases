package generator

import (
    "bufio"
    "errors"
    "math/rand"
    "os"
    "path"
    "runtime"
    "strings"
    "time"
)

func init() {
    rand.Seed(time.Now().UnixNano())
}

// Generator interface
type Generator interface {
    // Generate a random phrase
    Generate() (string, error)
    // Generate a random phrase and return indices
    GenerateWithIndices() (string, int, []int, error)
    // Generate a phrase by given indices
    GenerateByID(phraseIndex int, auxiliaryIndices []int) (string, error)
}

// generator implements Generator interface
type generator struct {
    Phrases     []string
    Auxiliaries []string
}

// New creates a new instance of generator
func New() (Generator, error) {
    g := new(generator)

    // Set the dictionaries paths
    pkgPath := getPackagePath()
    phrasesPath := pkgPath + "/dictionaries/phrases.txt"
    auxiliariesPath := pkgPath + "/dictionaries/auxiliaries.txt"

    // Read the dictionaries
    var err error
    g.Phrases, err = readLines(phrasesPath)
    if err != nil {
        return nil, err
    }

    g.Auxiliaries, err = readLines(auxiliariesPath)
    if err != nil {
        return nil, err
    }

    if len(g.Phrases) == 0 {
        return nil, errors.New("no phrases found in phrases.txt")
    }
    if len(g.Auxiliaries) == 0 {
        return nil, errors.New("no auxiliaries found in auxiliaries.txt")
    }

    return g, nil
}

// Generate a random phrase
func (g *generator) Generate() (string, error) {
    phrase, _, _, err := g.GenerateWithIndices()
    return phrase, err
}

// GenerateWithIndices generates a random phrase and returns the indices used
func (g *generator) GenerateWithIndices() (string, int, []int, error) {
    // Select a random phrase index
    phraseIndex := rand.Intn(len(g.Phrases))
    phraseTemplate := g.Phrases[phraseIndex]
    var result strings.Builder
    var auxiliaryIndices []int

    for _, ch := range phraseTemplate {
        if ch == '-' {
            // Select a random auxiliary index
            auxIndex := rand.Intn(len(g.Auxiliaries))
            auxiliaryIndices = append(auxiliaryIndices, auxIndex)
            auxiliary := g.Auxiliaries[auxIndex]
            result.WriteString(auxiliary)
        } else {
            result.WriteRune(ch)
        }
    }

    return result.String(), phraseIndex, auxiliaryIndices, nil
}

// GenerateByID reconstructs a phrase using the given indices
func (g *generator) GenerateByID(phraseIndex int, auxiliaryIndices []int) (string, error) {
    if phraseIndex < 0 || phraseIndex >= len(g.Phrases) {
        return "", errors.New("invalid phrase index")
    }

    phraseTemplate := g.Phrases[phraseIndex]
    var result strings.Builder
    auxCounter := 0

    for _, ch := range phraseTemplate {
        if ch == '-' {
            if auxCounter >= len(auxiliaryIndices) {
                return "", errors.New("not enough auxiliary indices provided")
            }
            auxIndex := auxiliaryIndices[auxCounter]
            if auxIndex < 0 || auxIndex >= len(g.Auxiliaries) {
                return "", errors.New("invalid auxiliary index")
            }
            auxiliary := g.Auxiliaries[auxIndex]
            result.WriteString(auxiliary)
            auxCounter++
        } else {
            result.WriteRune(ch)
        }
    }

    return result.String(), nil
}

// readLines reads a file and returns a slice of lines
func readLines(fileName string) ([]string, error) {
    f, err := os.Open(fileName)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    var lines []string
    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    if err := scanner.Err(); err != nil {
        return nil, err
    }
    return lines, nil
}

// getPackagePath gets the current path of this package
func getPackagePath() string {
    _, filename, _, ok := runtime.Caller(0)
    if !ok {
        return ""
    }
    return path.Dir(filename)
}
