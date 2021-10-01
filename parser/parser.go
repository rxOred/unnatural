package parser

import (
	"errors"

	mmap "github.com/edsrzf/mmap-go"
)

type ElfFile struct {
	pathname  string
	memmap    mmap.MMap
	ElfHeader Ehdr
	Phdr      []Phdr
	Shdr      []Shdr
	symtab    []Sym
	strtab    []byte
}

func (e *ElfFile) GetSectionNames() []string {
	var str []string
	if e.ElfHeader.EShstrndx == 0 {
		str = append(str, "section header string table is empty")
		return str
	}

	shstrtab := &e.memmap[e.elfHeader.EShstrndx]
	for i := 0; i < len(e.Shdr); i++ {
		str = append(str, &shstrtab[e.Shdr[i].ShName])
	}
	return str
}

func (e *ElfFile) GetSymbolNames() []string {
	var str []string
	for i := 0; i < len(e.Sym); i++ {
		// NOTE check which string table
		str = append(str)
	}
}

func (e *ElfFile) GetSectionIndexByName(name string) (int, error) {
	if e.ElfHeader.EShstrndx == 0 {
		return 0, errors.New("section header string table is empty")
	}
	shstrtab := &e.memmap[e.elfHeader.EShstrndx]
	for i := 0; i < len(shstrtab); i++ {
		if shstrtab+i == name {
			return i, nil
		}
	}
	return 0, errors.New("Could not find the section")
}

func LoadElf(e *ElfFile, pathname string) error {
	m, err := openFile(pathname)
	if err != nil {
		return err
	}

	e.memmap = m
	e.pathname = pathname
	e.elfHeader = e.memmap

	if verifyElf(e.ElfHeader.EIdent[:]) == false {
		return errors.New("Not an Elf binary")
	}

	e.Phdr = &e.memmap[e.elfHeader.EPhoff]
	e.Shdr = &e.memmap[e.elfHeader.EShoff]
	symtab_index, err := e.GetSectionIndexByName(".symtab")
	if err == nil && symtab_index != 0 {
		e.symtab = &e.memmap[e.shdr[symtab_index].ShOffset]
	}

	strtab_index, err := e.GetSectionIndexByName(".strtab")
	if err == nil && strtab_index != 0 {
		e.strtab = &e.memmap[e.shdr[strtab_index].ShOffset]
	}

	return nil
}
