package tarmmap

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/edsrzf/mmap-go"
)

type TarMmap struct {
	Headers []*tar.Header
	Files   [][]byte
	mmap    mmap.MMap
	f       *os.File
}

func Open(fileName string) (*TarMmap, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	// defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	size := fi.Size()
	if size == 0 {
		return nil, fmt.Errorf("tar file is empty")
	}

	mmap, err := mmap.Map(f, mmap.RDONLY, 0)
	if err != nil {
		return nil, err
	}
	// defer mmap.Unmap()

	// var (
	// 	buf     = bytes.NewBuffer(mmap)
	headers := []*tar.Header{}
	var files = [][]byte{}
	// )

	pos := int64(0)

	for {

		headerBlock := mmap[pos : pos+512]
		hdr, err := tar.NewReader(bytes.NewReader(headerBlock)).Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		headers = append(headers, hdr)

		data := mmap[pos+512 : pos+512+hdr.Size]
		// data := make([]byte, hdr.Size)
		// if _, err := io.ReadFull(buf, data); err != nil {
		// 	return nil, err
		// }

		files = append(files, data)

		blocks := hdr.Size / 512
		if hdr.Size%512 != 0 {
			blocks++
		}

		pos += (blocks + 1) * 512

		if pos >= int64(len(mmap)) {
			break
		}
	}

	return &TarMmap{
		Headers: headers,
		Files:   files,
		mmap:    mmap,
		f:       f,
	}, nil
}
