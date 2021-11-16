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

func (e *ElfFile) GetSectionNames() ([]string, error) {
	var str []string
	for i := 0; i < int(e.ElfHeader.EShnum); i++ {
		str = append(str, s)
	}
	return str, nil
}

func (e *ElfFile) getSectionName(index uint32) string {
	return string(e.shstrtab[index])
}

func (e *ElfFile) GetSectionNameByIndex(index uint32) (string, error) {
	if e.ElfHeader.EShstrndx <= 0 {
		return "", errors.New("shstrndx not found")
	}

	if uint32(e.Shdr[e.ElfHeader.EShstrndx].ShSize) < index {
		return "", errors.New("index out of range")
	}

	return e.getSectionName(index), nil
}

func (e *ElfFile) GetSectionIndexByName(name string) (int, error) {
	if e.ElfHeader.EShstrndx <= 0 {
		return -1, errors.New("shstrndx not found")
	}

	for i := 0; i < int(e.ElfHeader.EShnum); i++ {
		str := e.getSectionName(e.Shdr[i].ShName)
		if str == name {
			return i, nil
		}
	}

	return -1, nil
}

func (e *ElfFile) GetSegmentByType(ptype uint32) (*Phdr, error) {
	for i := 0; i < int(e.ElfHeader.EPhnum); i++ {
		if e.Phdr[i].PType == ptype {
			return e.Phdr[i], nil
		}
	}

	return nil, errors.New("Segment not found")
}

func (e *ElfFile) GetNSegmentByType(ptype uint32, n int) (*Phdr, error) {
	k := 0
	for i := 0; i < int(e.ElfHeader.EPhnum); i++ {
		if e.Phdr[i].PType == ptype {
			k++
			if k == n {
				return e.Phdr[i], nil
			}
		}
	}

	return nil, errors.New("Segment not found")
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
