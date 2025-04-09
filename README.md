# tar-mmap-go

A Go library for memory-mapped access to tar files, optimized for fast and efficient reading of tar archives.

## Overview

`tar-mmap-go` provides a way to read tar archives using memory mapping (mmap), which offers several advantages over standard streaming approaches:

- **Fast Access**: Direct memory access to file contents without copying data
- **Efficient Resource Usage**: Reduced memory consumption by sharing memory with the kernel
- **Random Access**: Immediate access to any file in the archive without sequential reading

## Installation
```bash
go get github.com/draganm/tar-mmap-go
```

## Usage

### Basic Example

```go
package main

import (
    "fmt"
    
    tarmmap "github.com/draganm/tar-mmap-go"
)

func main() {
    // Open a tar file using memory mapping
    tm, err := tarmmap.Open("archive.tar")
    if err != nil {
        panic(err)
    }
    defer tm.Close()
    
    // Access file sections (headers and contents)
    for _, section := range tm.Sections {
        fmt.Printf("File: %s, Size: %d bytes\n", section.Header.Name, section.Header.Size)
        
        // Access file content directly as []byte
        content := section.Data
        // Do something with content...
    }
}
```

### Command Line Tool

The package includes a powerful command-line tool for working with tar files:

```bash
# Install the command-line tool
go install github.com/draganm/tar-mmap-go/cmd/tar-mmap

# List the contents of a tar file
tar-mmap list archive.tar

# List with human-readable sizes
tar-mmap list --human archive.tar

# List with detailed information including offsets
tar-mmap list --verbose archive.tar

# Create a new tar from a subset of files (from index 0 to 5)
tar-mmap sub-tar --to 5 --output subset.tar archive.tar

# Create a new tar from a subset of files (from index 3 to 7)
tar-mmap sub-tar --from 3 --to 7 --output subset.tar archive.tar
```

#### Available Commands

- **list (ls)**: List contents of a tar file
  - `--human`: Show file sizes in human-readable format (KB, MB, etc.)
  - `--verbose` (`-v`): Show detailed information including header and data offsets

- **sub-tar**: Create a new tar file from a subset of files in the source tar
  - `--from`: Starting index of files to include (default: 0)
  - `--to`: Ending index of files to include (required)
  - `--output` (`-o`): Output tar file path (required)

The `sub-tar` command uses memory mapping for efficient extraction, copying the raw bytes directly without reparsing the tar structure, resulting in extremely fast operation.

## Features

- Memory-mapped access to tar files for efficient reading
- Direct access to file headers and contents
- Simple, clean API
- Low memory overhead

## Format Support

This library currently only supports tar files with USTAR headers. This is a deliberate design choice made for simplicity and efficiency, rather than a limitation of Go's standard library. The focus of this library is on providing fast, memory-mapped access to standard tar archives while maintaining a clean and straightforward implementation.

## Requirements

- Go 1.20 or later
- Supports platforms where mmap is available

## License

MIT License - see the [LICENSE](LICENSE) file for details.
