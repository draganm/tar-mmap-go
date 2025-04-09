package main

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	tarmmap "github.com/draganm/tar-mmap-go"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "tar-mmap",
		Usage: "Fast tar file operations using memory mapping",
		Description: `A command-line tool for efficient tar file operations using memory mapping.
   Provides fast access to tar contents without extracting the entire archive.`,
		Version: "1.0.0",
		Commands: []*cli.Command{
			{
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
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
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

func formatSize(size int64, human bool) string {
	if !human {
		return fmt.Sprintf("%d", size)
	}

	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

func formatType(typeflag byte) string {
	switch typeflag {
	case '0', '\x00':
		return "file"
	case '1':
		return "link"
	case '2':
		return "symlink"
	case '3':
		return "char"
	case '4':
		return "block"
	case '5':
		return "dir"
	case '6':
		return "fifo"
	default:
		return fmt.Sprintf("unknown(%c)", typeflag)
	}
}
