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

// this needs to be exported
type Elf struct {
	E_file *elf.File
}

func (e *Elf) GetSectionHeaders() []string {
	var arr []string
	sections := e.E_file.Sections
	for i := 0; i < len(sections); i++ {
		arr = append(arr, sections[i].Name)
	}
	return arr
}

func (e *Elf) GetSymbols() []string {
	var arr []string

	sym, err := e.E_file.Symbols()
	if err != nil {
		arr = append(arr, "no symbols found")
		return arr
	}
	for i := 0; i < len(sym); i++ {
		arr = append(arr, sym[i].Name)
	}

	if e.E_file.Type == elf.ET_DYN {
		dynsym, err := e.E_file.DynamicSymbols()
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

	arr = append(arr, "class :"+e.E_file.Class.String())
	arr = append(arr, "data :"+e.E_file.Data.String())
	arr = append(arr, "version :"+e.E_file.Version.String())
	arr = append(arr, "os abi :"+e.E_file.OSABI.String())
	arr = append(arr, "abi version :"+strconv.Itoa(int(e.E_file.ABIVersion)))
	arr = append(arr, "byteorder :"+e.E_file.ByteOrder.String())
	arr = append(arr, "type :"+e.E_file.Type.String())
	arr = append(arr, "machine :"+e.E_file.Machine.String())
	arr = append(arr, "entry :"+strconv.Itoa(int(e.E_file.Entry)))

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
	e.E_file = f
	return e, nil
}
