package generator

import (
	"math/rand"
	"strings"
	"testing"
)

func TestGenerateReplacesAllPlaceholders(t *testing.T) {
	g := &Generator{
		phrases:     []string{"- + - = -"},
		auxiliaries: []string{"uno"},
		rnd:         rand.New(rand.NewSource(1)),
	}

	got := g.Generate()
	want := "uno + uno = uno"
	if got != want {
		t.Fatalf("Generate() = %q, want %q", got, want)
	}
}

func TestNewLoadsEmbeddedDictionaries(t *testing.T) {
	g, err := newGenerator(rand.New(rand.NewSource(1)))
	if err != nil {
		t.Fatalf("newGenerator() returned error: %v", err)
	}

	if len(g.phrases) == 0 {
		t.Fatal("expected phrases to be loaded")
	}
	if len(g.auxiliaries) == 0 {
		t.Fatal("expected auxiliaries to be loaded")
	}
}

func TestGenerateReturnsCompactWhitespace(t *testing.T) {
	g := &Generator{
		phrases:     []string{"  -   y   -  "},
		auxiliaries: []string{"algo"},
		rnd:         rand.New(rand.NewSource(1)),
	}

	got := g.Generate()
	if strings.Contains(got, "  ") {
		t.Fatalf("Generate() returned repeated spaces: %q", got)
	}
}
