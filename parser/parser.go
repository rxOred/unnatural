package parser

import (
	"errors"

	mmap "github.com/edsrzf/mmap-go"
)

type ElfFile struct {
	pathname  string
	memmap    mmap.MMap
	elfHeader Ehdr
	phdr      []Phdr
	shdr      []Shdr
}

func LoadElf(e *ElfFile, pathname string) error {
	m, err := openFile(pathname)
	if err != nil {
		return err
	}

	e.memmap = m
	e.pathname = pathname
	e.elfHeader = e.memmap

	if verifyElf(e.elfHeader.EIdent[:]) == false {
		return errors.New("Not an Elf binary")
	}

	e.phdr = &e.memmap[e.elfHeader.EPhoff]
	e.shdr = &e.memmap[e.elfHeader.EShoff]

	return nil
}
