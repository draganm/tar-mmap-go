package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	tarmmap "github.com/draganm/tar-mmap-go"
	"github.com/urfave/cli/v2"
)

func listCommand() *cli.Command {
	return &cli.Command{
		Name:    "list",
		Aliases: []string{"ls"},
		Usage:   "List contents of a tar file",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name: "human",
				// Aliases: []string{"h"},
				Usage: "Show file sizes in human-readable format",
			},
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "Show detailed information including offsets",
			},
		},
		Action: listAction,
	}
}

func listAction(c *cli.Context) error {
	if c.Args().Len() == 0 {
		return fmt.Errorf("tar file path is required")
	}

	tarFile := c.Args().First()
	mm, err := tarmmap.Open(tarFile)
	if err != nil {
		return fmt.Errorf("failed to open tar file: %w", err)
	}
	defer mm.Close()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	if c.Bool("verbose") {
		fmt.Fprintf(w, "NAME\tSIZE\tTYPE\tMODE\tHEADER OFFSET\tEND OFFSET\n")
		for _, section := range mm.Sections {
			size := formatSize(section.Header.Size, c.Bool("human"))
			fmt.Fprintf(w, "%s\t%s\t%s\t%o\t%d\t%d\n",
				section.Header.Name,
				size,
				formatType(section.Header.Typeflag),
				section.Header.Mode,
				section.HeaderOffset,
				section.EndOfDataOffset)
		}
	} else {
		fmt.Fprintf(w, "NAME\tSIZE\tTYPE\tMODE\n")
		for _, section := range mm.Sections {
			size := formatSize(section.Header.Size, c.Bool("human"))
			fmt.Fprintf(w, "%s\t%s\t%s\t%o\n",
				section.Header.Name,
				size,
				formatType(section.Header.Typeflag),
				section.Header.Mode)
		}
	}

	return nil
}
