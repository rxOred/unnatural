package analyser

import (
	"debug/elf"
	"strconv"
)

// infection class
type Class int8

const (
	ELF_NONE            Class = 0
	ELF_TEXT_PADDING    Class = 1
	ELF_REVERSE_PADDING Class = 2
)

type Report struct {
	// infection classification
	R_class Class

	// this depends on the infection technique
	R_info []string
}

func CheckSegmentInfections(f *elf.File) Report {
	var r Report

	// text padding infection
	for i := 0; i < len(f.Progs); i++ {
		if f.Progs[i].Type == elf.PT_LOAD && f.Progs[i].Flags == elf.PF_R|elf.PF_X {
			if f.Entry > f.Progs[i].Vaddr {
				// init the structure, assign values, then return
				r = Report{
					R_class: ELF_TEXT_PADDING,
					R_info: []string{
						"[DECTED] classification : text padding infection",
						"Reasons for above conclusion :",
						"Entry point : 0x" + strconv.FormatUint(uint64(f.Entry), 16),
						"Text Segment address : 0x" + strconv.FormatUint(uint64(f.Progs[i].Vaddr), 16),
						"",
					},
				}
				return r
			}
		}

		// data segment padding infection
		if f.Progs[i].Type == elf.PT_LOAD && f.Progs[i].Flags == elf.PF_W|elf.PF_R {
			if f.Entry > f.Progs[i].Vaddr {
				r = Report{
					R_class: ELF_TEXT_PADDING,
					R_info: []string{
						"[DECTED] classification : data segment infection",
						"reasons for above conclusion :",
						"Entry point : 0x" + strconv.FormatUint(uint64(f.Entry), 16),
						"Data Segment address : 0x" + strconv.FormatUint(uint64(f.Progs[i].Vaddr), 16),
						"",
					},
				}
				return r
			}
		}
	}
	r.R_class = ELF_NONE
	r.R_info = []string{"not infected"}
	return r
}

func CheckDataSegmentInfection(f *elf.File) Report {
	var r Report
	for i := 0; i < len(f.Progs); i++ {
	}
	r.R_class = ELF_NONE
	r.R_info = []string{"not infected"}
	return r
}

func CheckReversePaddingInfection(f *elf.File) *Report {
	r := new(Report)
	return r
}
