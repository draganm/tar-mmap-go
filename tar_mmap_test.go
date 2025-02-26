package tarmmap_test

import (
	"testing"

	tarmmap "github.com/draganm/tar-mmap-go"
	"github.com/stretchr/testify/require"
)

func TestTwoFilesTar(t *testing.T) {
	tm, err := tarmmap.Open("fixtures/two-files.tar")
	require.NoError(t, err)

	require.Len(t, tm.Headers, 2)
	require.Len(t, tm.Files, 2)

	err = tm.Close()
	require.NoError(t, err)
}
