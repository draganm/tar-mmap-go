package tarmmap

import (
	"archive/tar"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/edsrzf/mmap-go"
)

type TarSection struct {
	Header          *tar.Header
	Data            []byte
	HeaderOffset    uint64
	EndOfDataOffset uint64
}

type TarMmap struct {
	Sections []TarSection
	Mmap     mmap.MMap
	f        *os.File
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

	var sections []TarSection

	pos := uint64(0)

	for {
		if pos >= uint64(len(mmap)) {
			break
		}

		headerOffset := pos
		headerBlock := mmap[pos : pos+512]
		hdr, err := tar.NewReader(bytes.NewReader(headerBlock)).Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		dataOffset := pos + 512
		data := mmap[dataOffset : dataOffset+uint64(hdr.Size)]

		blocks := uint64(hdr.Size) / 512
		if uint64(hdr.Size)%512 != 0 {
			blocks++
		}

		// Move to the next header position
		endOfDataOffset := pos + (blocks+1)*512

		section := TarSection{
			Header:          hdr,
			Data:            data,
			HeaderOffset:    headerOffset,
			EndOfDataOffset: endOfDataOffset,
		}

		sections = append(sections, section)

		pos = endOfDataOffset
	}

	return &TarMmap{
		Sections: sections,
		Mmap:     mmap,
		f:        f,
	}, nil
}

func (t *TarMmap) Close() error {
	return errors.Join(
		t.Mmap.Unmap(),
		t.f.Close(),
	)
}
