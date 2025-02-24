package tarmmap_test

import (
	"testing"

	tarmmap "github.com/draganm/tar-mmap-go"
)

func TestTwoFilesTar(t *testing.T) {
	tm, err := tarmmap.Open("fixtures/two-files.tar")
	if err != nil {
		t.Fatal(err)
	}

	if len(tm.Headers) != 2 {
		t.Fatalf("Expected 2 headers, got %d", len(tm.Headers))
	}

	if len(tm.Files) != 2 {
		t.Fatalf("Expected 2 files, got %d", len(tm.Files))
	}
}
