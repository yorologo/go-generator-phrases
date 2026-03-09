package main

import "testing"

func TestFormatPhrase(t *testing.T) {
	got := formatPhrase("hola mundo | adios mundo")
	want := "Hola mundo\nAdios mundo"
	if got != want {
		t.Fatalf("formatPhrase() = %q, want %q", got, want)
	}
}

func TestParseConfigPositionalQuantity(t *testing.T) {
	cfg, err := parseConfig([]string{"3"})
	if err != nil {
		t.Fatalf("parseConfig() returned error: %v", err)
	}

	if cfg.qty != 3 {
		t.Fatalf("cfg.qty = %d, want 3", cfg.qty)
	}
}

func TestParseConfigRejectsInvalidValues(t *testing.T) {
	if _, err := parseConfig([]string{"0"}); err == nil {
		t.Fatal("expected error for zero quantity")
	}
}
