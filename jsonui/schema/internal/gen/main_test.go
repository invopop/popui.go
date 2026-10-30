package main

import (
	"bytes"
	"os"
	"testing"
)

// TestSchemaUpToDate fails when panel.json doesn't match the registry and
// props, e.g. after adding a component or a prop without regenerating it.
func TestSchemaUpToDate(t *testing.T) {
	t.Chdir("../../../..")
	want, err := render()
	if err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("%s is out of date; run go generate ./jsonui/schema", outPath)
	}
}
