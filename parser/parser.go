package parser

import (
	"debug/elf"
	"errors"
	"io"
	"os"
)

var (
	elf_machine string
	elf_class   string
)

func ioReader(file string) (io.Reader, err) {
	r, err := os.Open(file)
	if err != nil {
		return nil, err
	}
}

// i kinda half stole this function from golang source hehe :3
func openElf(file string) (*elf.File, error) {
	f, err := ioReader(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	ident := make([]uint8, 16)
	if _, err = f.ReadAt(ident[0:], 0); err != nil {
		return nil, err
	}
	if ident[0] != '\x7f' || ident[1] != 'E' || ident[2] != 'L' || ident[3] != 'F' {
		return nil, errors.New("Invalid magic number")
	}

	ef := new(elf.File)
	ef.Class = elf.Class(ident[elf.EI_CLASS])
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

}
