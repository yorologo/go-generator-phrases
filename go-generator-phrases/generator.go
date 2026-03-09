package generator

import (
	"bufio"
	"bytes"
	"embed"
	"errors"
	"math/rand"
	"strings"
	"time"
)

//go:embed dictionaries/*.txt
var dictionariesFS embed.FS

type Generator struct {
	phrases     []string
	auxiliaries []string
	rnd         *rand.Rand
}

func New() (*Generator, error) {
	return newGenerator(rand.New(rand.NewSource(time.Now().UnixNano())))
}

func newGenerator(rnd *rand.Rand) (*Generator, error) {
	phrases, err := loadDictionary("dictionaries/phrases.txt")
	if err != nil {
		return nil, err
	}

	auxiliaries, err := loadDictionary("dictionaries/auxiliaries.txt")
	if err != nil {
		return nil, err
	}

	if len(phrases) == 0 || len(auxiliaries) == 0 {
		return nil, errors.New("dictionaries cannot be empty")
	}

	return &Generator{
		phrases:     phrases,
		auxiliaries: auxiliaries,
		rnd:         rnd,
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

func (g *Generator) Generate() string {
	if g == nil || len(g.phrases) == 0 || len(g.auxiliaries) == 0 {
		return ""
	}

	phrase := g.phrases[g.rnd.Intn(len(g.phrases))]
	parts := strings.Split(phrase, "-")
	if len(parts) == 1 {
		return strings.Join(strings.Fields(phrase), " ")
	}

	var builder strings.Builder
	builder.Grow(len(phrase) + len(parts)*16)
	for i, part := range parts {
		if i > 0 {
			builder.WriteString(g.auxiliaries[g.rnd.Intn(len(g.auxiliaries))])
		}
		builder.WriteString(part)
	}

	return strings.Join(strings.Fields(builder.String()), " ")
}
