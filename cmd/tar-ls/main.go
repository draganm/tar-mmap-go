package main

import (
	"fmt"
	"log"
	"os"

	tarmmap "github.com/draganm/tar-mmap-go"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Action: func(c *cli.Context) error {
			tarFile := c.Args().First()
			if tarFile == "" {
				return fmt.Errorf("tar file is required")
			}
			mm, err := tarmmap.Open(tarFile)
			if err != nil {
				return fmt.Errorf("failed to open tar file: %w", err)
			}
			for _, f := range mm.Headers {
				fmt.Println(f.Name)
			}

			return nil
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
