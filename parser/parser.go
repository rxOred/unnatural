package parser

import (
	"debug/elf"
	"fmt"
	"io"
	"os"
	"strconv"
)

var (
	elf_machine string
	elf_class   string
)

type Parser interface {
	Hllo(i string)
	InitElf(file string) (*elf.File, error)
	GetElfHeader(f *elf.File) []string
	GetSectionHeaders(f *elf.File) []string
	GetSymbols(f *elf.File) ([]string, error)
}

func GetSectionHeaders(f *elf.File) []string {
	var arr []string
	sections := f.Sections
	for i := 0; i < len(sections); i++ {
		arr = append(arr, sections[i].Name)
	}
	return arr
}

func GetSymbols(f *elf.File) ([]string, error) {
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

func GetElfHeader(f *elf.File) []string {
	var arr []string
	arr = append(arr, f.Class.String())
	arr = append(arr, f.Data.String())
	arr = append(arr, f.Version.String())
	arr = append(arr, f.OSABI.String())
	arr = append(arr, strconv.Itoa(int(f.ABIVersion)))
	arr = append(arr, f.ByteOrder.String())
	arr = append(arr, f.Type.String())
	arr = append(arr, f.Machine.String())
	arr = append(arr, strconv.Itoa(int(f.Entry)))

	return arr
}

func ioReader(file string) (io.ReaderAt, error) {
	r, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func InitElf(file string) (*elf.File, error) {
	r, err := ioReader(file)
	if err != nil {
		return nil, err
	}
	f, err := elf.NewFile(r)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func Hllo(i string) {
	fmt.Println("hello")
}
