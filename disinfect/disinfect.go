package disinfect

import (
	"debug/elf"
	"strconv"

	parser "github.com/rxOred/unnatural/parser"
)

func DisinfectTextPaddingInfection(e *parser.ElfFile) []string {
	sec, err := e.GetSectionIndexByName(".init")
	if err != nil || sec == -1 {

	}
}

// dis infection
func isinfectTextPaddingInfection(ef *parser.Elf) []string {
	sec := ef.E_file.Section(".init")
	for i := 0; i < len(ef.E_file.Progs); i++ {
		if ef.E_file.Progs[i].Type == elf.PT_LOAD && ef.E_file.Progs[i].Flags == elf.PF_R|elf.PF_X {
			if sec.Addr >= ef.E_file.Progs[i].Vaddr && sec.Addr < ef.E_file.Progs[i].Vaddr+ef.E_file.Progs[i].Memsz {
				ef.E_file.Entry = sec.Addr
				r := []string{
					"[SUCESSFUL] Disinfected text padding infection",
					"Fixed entry point : 0x" + strconv.FormatUint(uint64(ef.E_file.Entry), 16),
				}
				return r
			}
		}
	}
	return nil
}

func DisinfectDataSegmentInfection(ef *parser.Elf) []string {
	if r := DisinfectTextPaddingInfection(ef); r != nil {
		r[0] = "[SUCESSFUL] Disinfected data segment infection"
	}
	return nil
}

/*
func DisinfectReverseTextInfection(f *elf.File) error {

}

func DisinfectFunctionPaddingInfection(f *elf.File) error {

}
*/
