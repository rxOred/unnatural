package disinfect

import (
	"debug/elf"
	"strconv"
)

func DisinfectTextPaddingInfection(f *elf.File) []string {
	sec := f.Section(".init")
	for i := 0; i < len(f.Progs); i++ {
		if f.Progs[i].Type == elf.PT_LOAD && f.Progs[i].Flags == elf.PF_R|elf.PF_X {
			if sec.Addr >= f.Progs[i].Vaddr && sec.Addr < f.Progs[i].Vaddr+f.Progs[i].Memsz {
				f.Entry = sec.Addr
				r := []string{
					"[SUCESSFUL] Disinfected text padding infection",
					"Fixed entry point :", strconv.FormatUint(uint64(f.Entry), 16),
				}
				return r
			}
		}
	}
	return nil
}

func DisinfectDataSegmentInfection(f *elf.File) []string {
	if r := DisinfectTextPaddingInfection(f); r != nil {
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
