package parser

import (
	"debug/elf"
	"os"

	view "github.com/rxOred/unnatural/view"
)

func LoadElf(v *view.View, pathname string) *elf.File {
	f, err := os.Open(pathname)
	if err != nil {
		v.ShowErrorMsg(err.Error())
	}
	defer f.Close()

	elf, err := elf.NewFile(f)
	if err != nil {
		v.ShowErrorMsg(err.Error())
	}

	var ident [16]uint8
	if _, err := f.ReadAt(ident[0:], 0); err != nil {
		v.ShowErrorMsg("[!] Failed to read elf binary")
	}

	if ident[0] != '\x7f' || ident[1] != 'E' || ident[2] != 'L' || ident[3] != 'F' {
		v.ShowErrorMsg("[!] File has a wrong magic number")
	}

	return elf
}

func ShowInfo()
