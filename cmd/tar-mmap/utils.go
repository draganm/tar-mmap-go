package main

import (
	"fmt"
)

// formatSize converts a file size in bytes to a human-readable format if requested
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

// formatType returns a human-readable representation of a tar file type flag
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
