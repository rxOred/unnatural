package parser

import (
	"debug/elf"
	"io"
	"os"
	"strconv"
)

type Parser interface {
	GetElfHeader() []string
	GetSectionHeaders() []string
	GetSymbols() []string
}

type Elf struct {
	File *elf.File
}

func (e *Elf) GetSectionHeaders() []string {
	var arr []string
	sections := e.File.Sections
	for i := 0; i < len(sections); i++ {
		arr = append(arr, sections[i].Name)
	}
	return arr
}

func (e *Elf) GetSymbols() []string {
	var arr []string

	sym, err := e.File.Symbols()
	if err != nil {
		arr = append(arr, "no symbols found")
		return arr
	}
	for i := 0; i < len(sym); i++ {
		arr = append(arr, sym[i].Name)
	}

	if e.File.Type == elf.ET_DYN {
		dynsym, err := e.File.DynamicSymbols()
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

	arr = append(arr, "class :"+e.File.Class.String())
	arr = append(arr, "data :"+e.File.Data.String())
	arr = append(arr, "version :"+e.File.Version.String())
	arr = append(arr, "os abi :"+e.File.OSABI.String())
	arr = append(arr, "abi version :"+strconv.Itoa(int(e.File.ABIVersion)))
	arr = append(arr, "byteorder :"+e.File.ByteOrder.String())
	arr = append(arr, "type :"+e.File.Type.String())
	arr = append(arr, "machine :"+e.File.Machine.String())
	arr = append(arr, "entry :"+strconv.Itoa(int(e.File.Entry)))

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
	e.File = f
	return e, nil
}
