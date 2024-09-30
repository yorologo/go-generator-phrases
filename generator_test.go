package generator

import (
    "testing"
)

func TestGenerate(t *testing.T) {
    gen, err := New()
    if err != nil {
        t.Fatalf("Failed to create generator: %v", err)
    }
    phrase, err := gen.Generate()
    if err != nil {
        t.Errorf("Failed to generate phrase: %v", err)
    }
    if phrase == "" {
        t.Error("Expected a non-empty phrase")
    }
    t.Log("Generated phrase:", phrase)
}

func TestGenerateWithIndices(t *testing.T) {
    gen, err := New()
    if err != nil {
        t.Fatalf("Failed to create generator: %v", err)
    }
    phrase, phraseIndex, auxiliaryIndices, err := gen.GenerateWithIndices()
    if err != nil {
        t.Errorf("Failed to generate phrase with indices: %v", err)
    }
    if phrase == "" {
        t.Error("Expected a non-empty phrase")
    }
    if phraseIndex < 0 || phraseIndex >= len(gen.(*generator).Phrases) {
        t.Errorf("Phrase index out of bounds: %d", phraseIndex)
    }
    if len(auxiliaryIndices) == 0 {
        t.Error("Expected at least one auxiliary index")
    }
    t.Logf("Generated phrase: %s (Phrase index: %d, Auxiliary indices: %v)", phrase, phraseIndex, auxiliaryIndices)
}

func TestGenerateByID(t *testing.T) {
    gen, err := New()
    if err != nil {
        t.Fatalf("Failed to create generator: %v", err)
    }
    phrase, phraseIndex, auxiliaryIndices, err := gen.GenerateWithIndices()
    if err != nil {
        t.Fatalf("Failed to generate phrase with indices: %v", err)
    }
    reconstructedPhrase, err := gen.GenerateByID(phraseIndex, auxiliaryIndices)
    if err != nil {
        t.Errorf("Failed to generate phrase by ID: %v", err)
    }
    if phrase != reconstructedPhrase {
        t.Errorf("Expected reconstructed phrase to match original. Original: %s, Reconstructed: %s", phrase, reconstructedPhrase)
    }
    t.Log("Reconstructed phrase matches original")
}

func TestGenerateByID_InvalidPhraseIndex(t *testing.T) {
    gen, err := New()
    if err != nil {
        t.Fatalf("Failed to create generator: %v", err)
    }
    _, err = gen.GenerateByID(-1, []int{0})
    if err == nil {
        t.Error("Expected error due to invalid phrase index, but got none")
    } else {
        t.Logf("Received expected error: %v", err)
    }
}

func TestGenerateByID_InvalidAuxiliaryIndex(t *testing.T) {
    gen, err := New()
    if err != nil {
        t.Fatalf("Failed to create generator: %v", err)
    }
    _, phraseIndex, auxiliaryIndices, err := gen.GenerateWithIndices()
    if err != nil {
        t.Fatalf("Failed to generate phrase with indices: %v", err)
    }
    auxiliaryIndices[0] = -1
    _, err = gen.GenerateByID(phraseIndex, auxiliaryIndices)
    if err == nil {
        t.Error("Expected error due to invalid auxiliary index, but got none")
    } else {
        t.Logf("Received expected error: %v", err)
    }
}

func TestGenerateByID_InsufficientAuxiliaryIndices(t *testing.T) {
    gen, err := New()
    if err != nil {
        t.Fatalf("Failed to create generator: %v", err)
    }
    _, phraseIndex, _, err := gen.GenerateWithIndices()
    if err != nil {
        t.Fatalf("Failed to generate phrase with indices: %v", err)
    }
    _, err = gen.GenerateByID(phraseIndex, []int{})
    if err == nil {
        t.Error("Expected error due to insufficient auxiliary indices, but got none")
    } else {
        t.Logf("Received expected error: %v", err)
    }
}

func TestConsistency(t *testing.T) {
    gen, err := New()
    if err != nil {
        t.Fatalf("Failed to create generator: %v", err)
    }
    phrase1, err := gen.Generate()
    if err != nil {
        t.Errorf("Failed to generate first phrase: %v", err)
    }
    phrase2, err := gen.Generate()
    if err != nil {
        t.Errorf("Failed to generate second phrase: %v", err)
    }

    if phrase1 == phrase2 {
        t.Error("Expected different phrases on consecutive generations, but got the same")
    } else {
        t.Logf("Generated phrases are different as expected: \"%s\" vs \"%s\"", phrase1, phrase2)
    }
}
