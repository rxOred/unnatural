package parser

import (
	"debug/elf"
	"io"
	"os"
)

var (
	elf_machine string
	elf_class   string
)

func getSectionHeaders(f *elf.File) []string {
	var arr []string
	sections := f.Sections
	for i := 0; i < len(sections); i++ {
		arr = append(arr, sections[i].Name)
	}
	return arr
}

func getSymbols(f *elf.File) ([]string, error) {
	sym, err := f.Symbols()
	if err != nil {
		return nil, err
	}
	var arr []string
	for i := 0; i < len(sym); i++ {
		arr = append(arr, sym[i].Name)
	}

	if f.Type == elf.ET_DYN {
		dynsym, err := f.DynamicSymbols()
		if err != nil {
			return arr, err
		}
		for i := 0; i < len(dynsym); i++ {
			arr = append(arr, dynsym[i].Name)
		}
	}
	return arr, nil
}

func ioReader(file string) (io.ReaderAt, error) {
	r, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func InitElf(file string) error {
	r, err := ioReader(file)
	if err != nil {
		return err
	}
	f, err := elf.NewFile(r)
	if err != nil {

	}
}
