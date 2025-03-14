package tarmmap_test

import (
	"testing"

	tarmmap "github.com/draganm/tar-mmap-go"
	"github.com/stretchr/testify/require"
)

func TestTwoFilesTar(t *testing.T) {
	tm, err := tarmmap.Open("fixtures/two-files.tar")
	require.NoError(t, err)

	// Check that we have 2 sections (previously headers and files)
	require.Len(t, tm.Sections, 2)

	// Test first section
	require.Equal(t, "a", tm.Sections[0].Header.Name)
	require.Equal(t, []byte("a\n"), tm.Sections[0].Data)
	require.Equal(t, uint64(0), tm.Sections[0].HeaderOffset)
	require.Equal(t, uint64(1024), tm.Sections[0].EndOfDataOffset)

	// Test second section
	require.Equal(t, "b", tm.Sections[1].Header.Name)
	require.Equal(t, []byte("b\n"), tm.Sections[1].Data)
	require.Equal(t, uint64(1024), tm.Sections[1].HeaderOffset)

	// Verify offsets are sequential
	require.Equal(t, uint64(2048), tm.Sections[1].EndOfDataOffset)

	err = tm.Close()
	require.NoError(t, err)
}
