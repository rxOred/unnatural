package parser

import (
	"debug/elf"
	"encoding/binary"
	"errors"
	"os"

	"github.com/ghostiam/binstruct"
)

type ElfFile struct {
	File      *os.File // file
	ElfHeader Ehdr     // elf header
	Phdr      []*Phdr  // program header table
	Shdr      []*Shdr  // section header table
	Symtab    []*Sym   // symbol table
	Strtab    []byte   // string table
	shstrtab  []byte   // section header string table.
}

func (e *ElfFile) ParseDynamicSections() error {
	return nil
}

func (e *ElfFile) ParseStringTable() error {
	for i := 0; i < int(e.ElfHeader.EShnum); i++ {
		if e.Shdr[i].ShType == uint32(elf.SHT_STRTAB) {
			if i != int(e.ElfHeader.EShstrndx) {
				e.Strtab = readFile(e.File, int64(e.Shdr[i].ShOffset), uint32(e.Shdr[i].ShSize))
				return nil
			}
		}
	}
	return errors.New("section not found")
}

func (e *ElfFile) ParseSymbolTable() error {
	for i := 0; i < int(e.ElfHeader.EShnum); i++ {
		if e.Shdr[i].ShType == uint32(elf.SHT_SYMTAB) {
			for i := 0; i < (int(e.Shdr[i].ShSize) / int(e.Shdr[i].ShEntsize)); i++ {
				e.File.Seek(int64(e.Shdr[i].ShOffset), os.SEEK_SET)
				decoder := binstruct.NewDecoder(e.File, binary.LittleEndian)
				symtab := new(Sym)
				err := decoder.Decode(symtab)
				if err != nil {
					return err
				}
				e.Symtab = append(e.Symtab, symtab)
			}
		}
	}
	return nil
}

func (e *ElfFile) parseSectionHeaderStringTable() error {
	if e.ElfHeader.EShstrndx <= 0 {
		return errors.New("Failed to find section header string table")
	}
	index := e.ElfHeader.EShstrndx
	e.shstrtab = readFile(e.File, int64(e.Shdr[index].ShOffset), uint32(e.Shdr[index].ShSize))
	return nil
}

func LoadElf(e *ElfFile, pathname string) error {
	f, err := openFile(pathname)
	if err != nil {
		return err
	}
	e.File = f

	// parsing elf header
	decoder := binstruct.NewDecoder(e.File, binary.LittleEndian)
	err = decoder.Decode(&e.ElfHeader)
	if err != nil {
		return err
	}

	for i := 0; i < int(e.ElfHeader.EPhnum); i++ {
		e.File.Seek(int64(e.ElfHeader.EPhoff+uint64(i*int(e.ElfHeader.EPhentsize))), os.SEEK_SET)
		decoder = binstruct.NewDecoder(e.File, binary.LittleEndian)
		ph := new(Phdr)
		err = decoder.Decode(ph)
		if err != nil {
			return err
		}
		e.Phdr = append(e.Phdr, ph)
	}

	for i := 0; i < int(e.ElfHeader.EShnum); i++ {
		e.File.Seek(int64(e.ElfHeader.EShoff+uint64(i*int(e.ElfHeader.EShentsize))), os.SEEK_SET)
		decoder = binstruct.NewDecoder(e.File, binary.LittleEndian)
		sh := new(Shdr)
		err = decoder.Decode(sh)
		if err != nil {
			return err
		}
		e.Shdr = append(e.Shdr, sh)
	}

	err = e.ParseSymbolTable()
	if err != nil {
		e.Symtab = nil
	}

	err = e.ParseStringTable()
	if err != nil {
		e.Strtab = nil
	}

	err = e.parseSectionHeaderStringTable()
	if err != nil {
		e.shstrtab = nil
	}
	return nil
}

func (e *ElfFile) UnloadFile() {
	e.File.Close()
}
