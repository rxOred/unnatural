package analyser

import (
	"debug/elf"
	"strconv"
)

// infection class
type Class int

var (
	ELF_TEXT_PADDING    Class = 1
	ELF_REVERSE_PADDING Class = 2
)

type Report struct {
	// infection classification
	R_class Class

	// this depends on the infection technique
	R_info []string
}

func CheckTextPaddingInfection(f *elf.File) *Report {
	for i := 0; i < len(f.Progs); i++ {
		if f.Progs[i].Type == elf.ProgType(elf.PF_X|elf.PF_R) {
			// text segment found
			if f.Entry > f.Progs[i].Vaddr {
				// init the structure, assign values, then return
				r := &Report{
					R_class: ELF_TEXT_PADDING,
					R_info: []string{
						"[DECTED] classification : text padding infection",
						"reasons for above conclution :",
						"Entry point :" + strconv.Itoa(int(f.Entry)),
						"Text Segment address :" + strconv.Itoa(int(f.Progs[i].Vaddr)),
					},
				}
				return r
			}
		}
	}
	return nil
}

func CheckReversePaddingInfection(f *elf.File) *Report {
	r := new(Report)
	return r
}
