package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestSchema(t *testing.T) {
	files, err := filepath.Glob("testdata/*.json")
	if err != nil {
		t.Fatalf("Failed to glob ./testdata: %v", err)
	}
	for _, fn := range files {
		b, err := os.ReadFile(fn)
		if err != nil {
			t.Errorf("Failed to open(%q): %v", fn, err)
			continue
		}
		var ch CharacterResponse
		if err := json.Unmarshal(b, &ch); err != nil {
			t.Errorf("Failed to parse %q: %v", fn, err)
		}
	}
}
