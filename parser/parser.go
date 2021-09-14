package parser

import (
	"debug/elf"
	"encoding/binary"
	"errors"
	"io"
	"os"
)

type Header struct {
	Type                 elf.Type
	Machine              elf.Machine
	Entry                uint64
	Version              elf.Version
	Phoff, Shoff         int64
	Phentsize, Shentsize int
	Phnum, Shnum         int
	Shstrndx             int // need to parse section names
}

var (
	header      *Header
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
func ParseHeader(file string) (*Header, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := io.ReaderAt(f)
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
	default:
		return nil, errors.New("Invalid class")
	}
	ef.OSABI = elf.OSABI(ident[elf.EI_OSABI])
	switch ef.OSABI.String() {
	case "":
		// ok
	default:
		return nil, errors.New("Invalid OSABI")
	}
	ef.Version = elf.Version(ident[elf.EI_VERSION])
	if ef.Version != elf.EV_CURRENT {
		return nil, errors.New("Invalid version")
	}

	header = new(Header)

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
		if header.Version = elf.Version(hdr.Version); header.Version != ef.Version {
			return nil, errors.New("Mismatching versions")
		}

		// assigning values to header
		header.Type = ef.Type
		header.Machine = ef.Machine
		header.Entry = ef.Entry
		header.Phoff = int64(hdr.Phoff)
		header.Shoff = int64(hdr.Shoff)
		header.Phentsize = int(hdr.Phentsize)
		header.Shentsize = int(hdr.Shentsize)
		header.Phnum = int(hdr.Phnum)
		header.Shnum = int(hdr.Shnum)
		header.Shstrndx = int(hdr.Shstrndx)
	case "64b":
		hdr := new(elf.Header64)
		section_reader.Seek(0, io.SeekStart)
		if err := binary.Read(section_reader, ef.ByteOrder, hdr); err != nil {
			return nil, errors.New(err.Error())
		}

		ef.Type = elf.Type(hdr.Type)
		ef.Machine = elf.Machine(hdr.Machine)
		ef.Entry = uint64(hdr.Entry)
		if header.Version = elf.Version(hdr.Version); header.Version != ef.Version {
			return nil, errors.New("Mismatching versions")
		}

		// assigning values to header
		header.Type = ef.Type
		header.Machine = ef.Machine
		header.Entry = ef.Entry
		header.Phoff = int64(hdr.Phoff)
		header.Shoff = int64(hdr.Shoff)
		header.Phentsize = int(hdr.Phentsize)
		header.Shentsize = int(hdr.Shentsize)
		header.Phnum = int(hdr.Phnum)
		header.Shnum = int(hdr.Shnum)
		header.Shstrndx = int(hdr.Shstrndx)
	}
	return header, nil
}
