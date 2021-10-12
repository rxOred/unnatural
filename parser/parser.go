package parser

import (
	"debug/elf"
	"encoding/binary"
	"errors"
	"os"
	"strconv"

	"github.com/ghostiam/binstruct"
)

type ElfFile struct {
	pathname  string
	ElfHeader Ehdr
	Phdr      []*Phdr
	Shdr      []*Shdr
	symtab    **Sym
	strtab    []byte
}

func (e *ElfFile) GetElfHeader() []string {
	var str []string
	switch e.ElfHeader.EType {
	case uint16(elf.ET_REL):
		str = append(str, "Type :"+elf.ET_REL.String())
	case uint16(elf.ET_EXEC):
		str = append(str, "Type :"+elf.ET_EXEC.String())
	case uint16(elf.ET_DYN):
		str = append(str, "Type :"+elf.ET_DYN.String())
	case uint16(elf.ET_CORE):
		str = append(str, "Type :"+elf.ET_CORE.String())
	default:
		str = append(str, "Type :"+"none")
	}

	switch e.ElfHeader.EMachine {
	case uint16(elf.EM_386):
		str = append(str, "Machine :"+elf.EM_386.String())
	case uint16(elf.EM_MIPS):
		str = append(str, "Machine :"+elf.EM_MIPS.String())
	case uint16(elf.EM_ARM):
		str = append(str, "Machine :"+elf.EM_ARM.String())
	case uint16(elf.EM_X86_64):
		str = append(str, "Machine :"+elf.EM_X86_64.String())
	default:
		str = append(str, "Machine :"+"none")
	}

	switch e.ElfHeader.EVersion {
	case uint32(elf.EV_CURRENT):
		str = append(str, "Version :"+elf.EV_CURRENT.String())
	default:
		str = append(str, "Version :"+elf.EV_NONE.String())
	}

	str = append(str, "Entry :"+strconv.FormatUint(e.ElfHeader.EEntry, 16))
	str = append(str, "Phoff :"+strconv.FormatUint(e.ElfHeader.EPhoff, 16))
	str = append(str, "Shoff :"+strconv.FormatUint(e.ElfHeader.EShoff, 16))
	str = append(str, "Phnum :"+strconv.FormatUint(uint64(e.ElfHeader.EPhnum), 10))
	str = append(str, "Shnum :"+strconv.FormatUint(uint64(e.ElfHeader.EShnum), 10))
	str = append(str, "Shstrndx :"+strconv.FormatUint(uint64(e.ElfHeader.EShstrndx), 10))
	return str
}

func (e *ElfFile) GetSectionHeaders() []string {
	var str []string
	for i := 0; i< len(e.Shdr); i++ {
		switch e.Shdr[i].ShName {
			str = append(str, )
		}
	}
}

func (e *ElfFile) GetProgHeaders() []string {
	var str []string
	for i := 0; i < len(e.Phdr); i++ {
		switch e.Phdr[i].PType {
		case uint32(elf.PT_LOAD):
			str = append(str, elf.PT_LOAD.String())
		case uint32(elf.PT_DYNAMIC):
			str = append(str, elf.PT_DYNAMIC.String())
		case uint32(elf.PT_INTERP):
			str = append(str, elf.PT_INTERP.String())
		default:
			str = append(str, "none")
		}

		str = append(str, strconv.FormatUint(e.Phdr[i].POffset, 16))
		str = append(str, strconv.FormatUint(e.Phdr[i].PVaddr, 16))
		str = append(str, strconv.FormatUint(e.Phdr[i].PPaddr, 16))
	}

	return str
}

func (e *ElfFile) GetSegmentIndexByType(ptype uint32) (*Phdr, error) {
	for i := 0; i < len(e.Phdr); i++ {
		if e.Phdr[i].PType == ptype {
			return e.Phdr[i], nil
		}
	}
	return nil, errors.New("Segment not found")
}

// whole parser thing should be changed to read from file

func LoadElf(e *ElfFile, pathname string) error {
	e.pathname = pathname
	f, err := openFile(pathname)
	if err != nil {
		return err
	}

	decoder := binstruct.NewDecoder(f, binary.LittleEndian)
	err = decoder.Decode(&e.ElfHeader)
	if err != nil {
		return err
	}

	for i := 0; i < int(e.ElfHeader.EPhnum); i++ {
		f.Seek(int64(e.ElfHeader.EPhoff+uint64(i*int(e.ElfHeader.EPhentsize))), os.SEEK_SET)
		decoder = binstruct.NewDecoder(f, binary.LittleEndian)
		ph := new(Phdr)
		err = decoder.Decode(ph)
		e.Phdr = append(e.Phdr, ph)
	}

	for i := 0; i < int(e.ElfHeader.EShnum); i++ {
		f.Seek(int64(e.ElfHeader.EPhoff+uint64(i*int(e.ElfHeader.EShentsize))), os.SEEK_SET)
		decoder = binstruct.NewDecoder(f, binary.LittleEndian)
		sh := new(Shdr)
		err = decoder.Decode(sh)
		e.Shdr = append(e.Shdr, sh)
	}
	return nil
}
