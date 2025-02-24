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
    
    // Access file headers and contents
    for i, header := range tm.Headers {
        fmt.Printf("File: %s, Size: %d bytes\n", header.Name, header.Size)
        
        // Access file content directly as []byte
        content := tm.Files[i]
        // Do something with content...
    }
}
```

### Command Line Tool

The package includes a simple command-line tool for listing the contents of tar files:

```bash
# Install the command-line tool
go install github.com/draganm/tar-mmap-go/cmd/tar-ls

# List the contents of a tar file
tar-ls archive.tar
```

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