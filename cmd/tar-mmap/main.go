package main

import (
	"log"
	"os"

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
			listCommand(),
			subTarCommand(),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
