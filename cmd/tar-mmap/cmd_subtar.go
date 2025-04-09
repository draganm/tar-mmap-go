package main

import (
	"fmt"
	"os"

	tarmmap "github.com/draganm/tar-mmap-go"
	"github.com/urfave/cli/v2"
)

func subTarCommand() *cli.Command {
	return &cli.Command{
		Name:  "sub-tar",
		Usage: "Create a new tar file from a subset of files in the source tar",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  "from",
				Usage: "Starting index of files to include (default: 0)",
				Value: 0,
			},
			&cli.IntFlag{
				Name:     "to",
				Usage:    "Ending index of files to include (required)",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "output",
				Aliases:  []string{"o"},
				Usage:    "Output tar file path (required)",
				Required: true,
			},
		},
		Action: subTarAction,
	}
}

func subTarAction(c *cli.Context) error {
	if c.Args().Len() == 0 {
		return fmt.Errorf("tar file path is required")
	}

	sourcePath := c.Args().First()
	outputPath := c.String("output")
	fromIndex := c.Int("from")
	toIndex := c.Int("to")

	// Open source tar file
	mm, err := tarmmap.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to open source tar file: %w", err)
	}
	defer mm.Close()

	// Validate indices
	if fromIndex < 0 {
		return fmt.Errorf("from index must be >= 0")
	}

	if toIndex < fromIndex {
		return fmt.Errorf("to index must be >= from index")
	}

	if toIndex >= len(mm.Sections) {
		return fmt.Errorf("to index out of range, max index is %d", len(mm.Sections)-1)
	}

	// Create output tar file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	// Get first and last section to extract
	firstSection := mm.Sections[fromIndex]
	lastSection := mm.Sections[toIndex]

	// Calculate total range to copy
	startOffset := firstSection.HeaderOffset
	endOffset := lastSection.EndOfDataOffset

	// Get the slice of raw memory-mapped data for the entire range
	rawData := mm.Mmap[startOffset:endOffset]

	// Write the raw data to output file in one operation
	bytesWritten, err := outputFile.Write(rawData)
	if err != nil {
		return fmt.Errorf("failed to write tar data: %w", err)
	}

	// Calculate number of files copied
	numFiles := toIndex - fromIndex + 1

	// Add trailing blocks of zeros to properly terminate the tar file
	// (standard tar files end with at least two zero blocks)
	trailer := make([]byte, 1024) // 2 blocks of 512 zeros
	_, err = outputFile.Write(trailer)
	if err != nil {
		return fmt.Errorf("failed to write tar trailer: %w", err)
	}

	fmt.Printf("Wrote %d bytes from %d files to %s\n", bytesWritten, numFiles, outputPath)
	fmt.Printf("Files included: \n")

	// Display list of files included
	for i := fromIndex; i <= toIndex; i++ {
		section := mm.Sections[i]
		fmt.Printf("  %s (%d bytes)\n", section.Header.Name, section.Header.Size)
	}

	return nil
}
