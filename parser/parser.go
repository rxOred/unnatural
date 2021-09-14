package parser

import (
	"debug/elf"
	"encoding/binary"
	"errors"
	"io"
	"os"
)

var (
	elf_machine string
	elf_class   string
)

func ioReader(file string) (io.ReaderAt, error) {
	r, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// i kinda half stole this function from golang source hehe :3
func openElf(file string) (*elf.File, error) {

	r, err := ioReader(file)
	if err != nil {
		return nil, err
	}

	section_reader := io.NewSectionReader(r, 0, 1<<63-1)

	ident := make([]uint8, 16)
	if _, err = r.ReadAt(ident[0:], 0); err != nil {
		return nil, err
	}
	if ident[0] != '\x7f' || ident[1] != 'E' || ident[2] != 'L' || ident[3] != 'F' {
		return nil, errors.New("Invalid magic number")
	}

	ef := new(elf.File)
	ef.Class = elf.Class(ident[elf.EI_CLASS])

	// we need elf class to select whether to use header32 or header64
	switch ef.Class.String() {
	case "ELFCLASS64":
		elf_class = "64b"
	case "ELFCLASS32":
		elf_class = "32b"
	}
	ef.OSABI = elf.OSABI(ident[elf.EI_OSABI])
	switch ef.OSABI.String() {
	case "":
		// ok
	default:
		return nil, errors.New("Invalid OSABI")
	}

	// parse elf header
	var (
		phoff, shoff         int64
		phentsize, shentsize int
		phnum, shnum         int
		shstrndx             int // this is the thing that is most important for us to parse section names
	)
	switch elf_class {
	case "32b":
		hdr := new(elf.Header32)
		section_reader.Seek(0, io.SeekStart)
		if err := binary.Read(section_reader, ef.ByteOrder, hdr); err != nil {
			return nil, errors.New(err.Error())
		}

		ef.Type = elf.Type(hdr.Type)
		ef.Machine = elf.Machine(hdr.Machine)
		ef.Entry = uint64(hdr.Entry)

	case "64b":
	}
}
