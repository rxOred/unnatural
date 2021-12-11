package scanner

import (
	"debug/elf"
	"strconv"

	parser "github.com/rxOred/unnatural/parser"
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

func CheckSegmentInfections(e *parser.ElfFile) Report {
	var r Report
	// text padding infection
	for i := 0; i < len(e.Phdr); i++ {
		if elf.ProgType(e.Phdr[i].PType) == elf.PT_LOAD && elf.ProgFlag(e.Phdr[i].PFlags) == elf.PF_R|elf.PF_X {
			if e.ElfHeader.EEntry > e.Phdr[i].PVaddr {

				// more code here
				// init the structure, assign values, then return
				r = Report{
					R_class: ELF_TEXT_PADDING,
					R_info: []string{
						"[DETECTED] classification : text padding infection",
						"Reasons for above conclusion :",
						"Entry point : 0x" + strconv.FormatUint(uint64(e.ElfHeader.EEntry), 16),
						"Text Segment address : 0x" + strconv.FormatUint(uint64(e.Phdr[i].PVaddr), 16),
						"",
					},
				}
				return r
			}
		}
		// data segment padding infection
		if elf.ProgType(e.Phdr[i].PType) == elf.PT_LOAD && elf.ProgFlag(e.Phdr[i].PFlags) == elf.PF_W|elf.PF_R {
			if e.ElfHeader.EEntry > e.Phdr[i].PVaddr {
				r = Report{
					R_class: ELF_TEXT_PADDING,
					R_info: []string{
						"[DECTED] classification : data segment infection",
						"reasons for above conclusion :",
						"Entry point : 0x" + strconv.FormatUint(uint64(e.ElfHeader.EEntry), 16),
						"Data Segment address : 0x" + strconv.FormatUint(uint64(e.Phdr[i].PVaddr), 16),
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
