package generator

import (
	"strings"
	"testing"
)

func TestGenerateReturnsPhrase(t *testing.T) {
	gen, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	phrase, err := gen.Generate()
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if strings.TrimSpace(phrase) == "" {
		t.Fatal("Generate() returned an empty phrase")
	}
}

func TestGenerateWithIndicesCanBeRebuilt(t *testing.T) {
	gen, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	phrase, phraseIndex, auxiliaryIndices, err := gen.GenerateWithIndices()
	if err != nil {
		t.Fatalf("GenerateWithIndices() error = %v", err)
	}

	rebuilt, err := gen.GenerateByID(phraseIndex, auxiliaryIndices)
	if err != nil {
		t.Fatalf("GenerateByID() error = %v", err)
	}

	if phrase != rebuilt {
		t.Fatalf("rebuilt phrase mismatch: got %q want %q", rebuilt, phrase)
	}
}

func TestGenerateByIDRejectsInvalidPhraseIndex(t *testing.T) {
	gen, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if _, err := gen.GenerateByID(-1, nil); err == nil {
		t.Fatal("expected error for invalid phrase index")
	}
}

func TestGenerateByIDRejectsInvalidAuxiliaryIndex(t *testing.T) {
	gen, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	_, phraseIndex, auxiliaryIndices, err := gen.GenerateWithIndices()
	if err != nil {
		t.Fatalf("GenerateWithIndices() error = %v", err)
	}
	auxiliaryIndices[0] = -1

	if _, err := gen.GenerateByID(phraseIndex, auxiliaryIndices); err == nil {
		t.Fatal("expected error for invalid auxiliary index")
	}
}

func TestGenerateByIDRejectsWrongAuxiliaryCount(t *testing.T) {
	gen, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	_, phraseIndex, auxiliaryIndices, err := gen.GenerateWithIndices()
	if err != nil {
		t.Fatalf("GenerateWithIndices() error = %v", err)
	}

	if len(auxiliaryIndices) == 0 {
		t.Fatal("expected at least one auxiliary index")
	}

	if _, err := gen.GenerateByID(phraseIndex, auxiliaryIndices[:len(auxiliaryIndices)-1]); err == nil {
		t.Fatal("expected error for wrong auxiliary count")
	}
}
