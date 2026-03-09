package generator

import (
	"bufio"
	"bytes"
	"embed"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

//go:embed dictionaries/*.txt
var dictionariesFS embed.FS

// Generator exposes the stable API for phrase generation and reconstruction.
type Generator interface {
	Generate() (string, error)
	GenerateWithIndices() (string, int, []int, error)
	GenerateByID(phraseIndex int, auxiliaryIndices []int) (string, error)
}

type generator struct {
	Phrases     []string
	Auxiliaries []string
	rnd         *rand.Rand
}

// New creates a generator backed by embedded dictionaries.
func New() (Generator, error) {
	phrases, err := loadDictionary("dictionaries/phrases.txt")
	if err != nil {
		return nil, fmt.Errorf("load phrases dictionary: %w", err)
	}

	auxiliaries, err := loadDictionary("dictionaries/auxiliaries.txt")
	if err != nil {
		return nil, fmt.Errorf("load auxiliaries dictionary: %w", err)
	}

	if len(phrases) == 0 {
		return nil, errors.New("phrases dictionary is empty")
	}
	if len(auxiliaries) == 0 {
		return nil, errors.New("auxiliaries dictionary is empty")
	}

	return &generator{
		Phrases:     phrases,
		Auxiliaries: auxiliaries,
		rnd:         rand.New(rand.NewSource(time.Now().UnixNano())),
	}, nil
}

func loadDictionary(path string) ([]string, error) {
	data, err := dictionariesFS.ReadFile(path)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(bytes.NewReader(data))
	entries := make([]string, 0)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		entries = append(entries, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

func (g *generator) Generate() (string, error) {
	phrase, _, _, err := g.GenerateWithIndices()
	return phrase, err
}

func (g *generator) GenerateWithIndices() (string, int, []int, error) {
	if g == nil || len(g.Phrases) == 0 || len(g.Auxiliaries) == 0 {
		return "", 0, nil, errors.New("generator is not initialized")
	}

	phraseIndex := g.rnd.Intn(len(g.Phrases))
	phraseTemplate := g.Phrases[phraseIndex]
	auxiliaryIndices := make([]int, 0, strings.Count(phraseTemplate, "-"))

	var result strings.Builder
	result.Grow(len(phraseTemplate) + len(auxiliaryIndices)*16)

	for _, ch := range phraseTemplate {
		if ch != '-' {
			result.WriteRune(ch)
			continue
		}

		auxIndex := g.rnd.Intn(len(g.Auxiliaries))
		auxiliaryIndices = append(auxiliaryIndices, auxIndex)
		result.WriteString(g.Auxiliaries[auxIndex])
	}

	return compactSpaces(result.String()), phraseIndex, auxiliaryIndices, nil
}

func (g *generator) GenerateByID(phraseIndex int, auxiliaryIndices []int) (string, error) {
	if g == nil || len(g.Phrases) == 0 || len(g.Auxiliaries) == 0 {
		return "", errors.New("generator is not initialized")
	}
	if phraseIndex < 0 || phraseIndex >= len(g.Phrases) {
		return "", errors.New("invalid phrase index")
	}

	phraseTemplate := g.Phrases[phraseIndex]
	requiredAuxiliaries := strings.Count(phraseTemplate, "-")
	if len(auxiliaryIndices) != requiredAuxiliaries {
		return "", fmt.Errorf("expected %d auxiliary indices, got %d", requiredAuxiliaries, len(auxiliaryIndices))
	}

	var result strings.Builder
	auxCounter := 0
	for _, ch := range phraseTemplate {
		if ch != '-' {
			result.WriteRune(ch)
			continue
		}

		auxIndex := auxiliaryIndices[auxCounter]
		if auxIndex < 0 || auxIndex >= len(g.Auxiliaries) {
			return "", errors.New("invalid auxiliary index")
		}

		result.WriteString(g.Auxiliaries[auxIndex])
		auxCounter++
	}

	return compactSpaces(result.String()), nil
}

func compactSpaces(value string) string {
	return strings.Join(strings.Fields(value), " ")
}
