package generator

import (
	"math/rand"
	"strings"
	"testing"
)

func TestNewLoadsEmbeddedDictionaries(t *testing.T) {
	gen, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if len(gen.phrases) == 0 {
		t.Fatal("expected embedded phrases to be loaded")
	}
	if len(gen.auxiliaries) == 0 {
		t.Fatal("expected embedded auxiliaries to be loaded")
	}
}

func TestGenerateReplacesAllPlaceholders(t *testing.T) {
	gen := &Generator{
		phrases:     []string{"- + - = -"},
		auxiliaries: []string{"uno"},
		rnd:         rand.New(rand.NewSource(1)),
	}

	got := gen.Generate()
	if got != "uno + uno = uno" {
		t.Fatalf("Generate() = %q, want %q", got, "uno + uno = uno")
	}
}

func TestGenerateCompactsWhitespace(t *testing.T) {
	gen := &Generator{
		phrases:     []string{"  -   y   -  "},
		auxiliaries: []string{"algo"},
		rnd:         rand.New(rand.NewSource(1)),
	}

	got := gen.Generate()
	if strings.Contains(got, "  ") {
		t.Fatalf("Generate() returned repeated spaces: %q", got)
	}
}
