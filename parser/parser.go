package parser

import (
	"debug/elf"
	"fmt"
	"io"
	"os"
	"strconv"
)

type Parser interface {
	Hllo(i string)
	GetElfHeader() []string
	GetSectionHeaders() []string
	GetSymbols() []string
}

type Elf struct {
	f *elf.File
}

func (e *Elf) GetSectionHeaders() []string {
	var arr []string
	sections := e.f.Sections
	for i := 0; i < len(sections); i++ {
		arr = append(arr, sections[i].Name)
	}
	return arr
}

func (e *Elf) GetSymbols() []string {
	var arr []string

	sym, err := e.f.Symbols()
	if err != nil {
		arr = append(arr, "no symbols found")
		return arr
	}
	for i := 0; i < len(sym); i++ {
		arr = append(arr, sym[i].Name)
	}

	if e.f.Type == elf.ET_DYN {
		dynsym, err := e.f.DynamicSymbols()
		if err != nil {
			arr = append(arr, "no dynamic symbols found")
			return arr
		}
		for i := 0; i < len(dynsym); i++ {
			arr = append(arr, dynsym[i].Name)
		}
	}
	return arr
}

func (e *Elf) GetElfHeader() []string {
	var arr []string
	arr = append(arr, e.f.Class.String())
	arr = append(arr, e.f.Data.String())
	arr = append(arr, e.f.Version.String())
	arr = append(arr, e.f.OSABI.String())
	arr = append(arr, strconv.Itoa(int(e.f.ABIVersion)))
	arr = append(arr, e.f.ByteOrder.String())
	arr = append(arr, e.f.Type.String())
	arr = append(arr, e.f.Machine.String())
	arr = append(arr, strconv.Itoa(int(e.f.Entry)))

	return arr
}

func ioReader(file string) (io.ReaderAt, error) {
	r, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func InitElf(file string) (*Elf, error) {
	r, err := ioReader(file)
	if err != nil {
		return nil, err
	}
	f, err := elf.NewFile(r)
	if err != nil {
		return nil, err
	}
	e := new(Elf)
	e.f = f
	return e, nil
}

func (e *Elf) Hllo(i string) {
	fmt.Println("hello")
}
