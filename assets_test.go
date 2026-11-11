package popui

import (
	"strings"
	"sync"
	"testing"
	"testing/fstest"
)

func TestVersionedConcurrent(t *testing.T) {
	content := fstest.MapFS{
		"scripts/a.js": {Data: []byte("alert('a')")},
		"scripts/b.js": {Data: []byte("alert('b')")},
		"styles/c.css": {Data: []byte("body{}")},
	}
	paths := [][]string{
		{"scripts", "a.js"},
		{"scripts", "b.js"},
		{"styles", "c.css"},
		{"missing.js"},
	}

	const workers = 32
	results := make([][]string, workers)
	var wg sync.WaitGroup
	for i := range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range 100 {
				results[i] = append(results[i], Versioned(content, paths[j%len(paths)]...))
			}
		}()
	}
	wg.Wait()

	for i := range workers {
		for j, got := range results[i] {
			want := results[0][j]
			if got != want {
				t.Fatalf("worker %d call %d: got %q, want %q", i, j, got, want)
			}
		}
	}
	if got := results[0][0]; !strings.HasPrefix(got, "scripts/a.js?v=") || len(got) != len("scripts/a.js?v=")+8 {
		t.Errorf("unexpected versioned path %q", got)
	}
	if got := results[0][3]; got != "missing.js" {
		t.Errorf("missing file: got %q, want %q", got, "missing.js")
	}
}
