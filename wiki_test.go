package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sebdah/goldie/v2"
)

func TestWikiPage(t *testing.T) {
	g := goldie.New(t)
	files, err := filepath.Glob("testdata/*.golden")
	if err != nil {
		t.Fatalf("Failed to glob ./testdata: %v", err)
	}
	for _, fn := range files {
		b, err := os.ReadFile(strings.TrimSuffix(fn, ".golden") + ".json")
		if err != nil {
			t.Errorf("Failed to open(%q): %v", fn, err)
			continue
		}
		var ch CharacterResponse
		if err := json.Unmarshal(b, &ch); err != nil {
			t.Errorf("Failed to parse %q: %v", fn, err)
		}
		text := characterToWikiPage(ch)
		g.Assert(t, strings.TrimSuffix(filepath.Base(fn), ".golden"), []byte(text))
	}
}
