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
	ELF_DATA_PADDING    Class = 3
)

type Report struct {
	// infection classification
	R_class Class
	// description + possibility of disinfection (T/F)
	R_info []string
}

func CheckSegmentInfections(e *parser.ElfFile) Report {
	var r Report
	var (
		data_seg_index = 0
		text_seg_index = 0
	)

	for i := 0; i < len(e.Phdr); i++ {
		// text padding infection
		if elf.ProgType(e.Phdr[i].PType) == elf.PT_LOAD && elf.ProgFlag(e.Phdr[i].PFlags) == elf.PF_R|elf.PF_X {
			text_seg_index = i
		}
		// data segment padding infection
		if elf.ProgType(e.Phdr[i].PType) == elf.PT_LOAD && elf.ProgFlag(e.Phdr[i].PFlags) == elf.PF_W|elf.PF_R {
			data_seg_index = i
		}
	}

	if e.ElfHeader.EEntry > e.Phdr[text_seg_index].PVaddr && e.ElfHeader.EEntry < e.Phdr[data_seg_index].PVaddr {
		for i := 0; i < len(e.Shdr); i++ {
			sec_name, err := e.GetSectionNameByIndex(uint32(i))
			if err != nil {
				continue
			}
			if sec_name == ".fini" && e.Shdr[i].ShFlags == uint64(elf.SHF_ALLOC|elf.SHF_EXECINSTR) {
				if e.ElfHeader.EEntry > e.Shdr[i].ShAddr {
					// init the structure, assign values, then return
					r = Report{
						R_class: ELF_TEXT_PADDING,
						R_info: []string{
							"[DETECTED] classification : text padding infection",
							"Reasons for the above conclusion :",
							"Entry point : 0x" + strconv.FormatUint(uint64(e.ElfHeader.EEntry), 16),
							"Text Segment address : 0x" + strconv.FormatUint(uint64(e.Phdr[text_seg_index].PVaddr), 16),
							"",
						},
					}
					return r
				}
			}
		}
	}

	if e.ElfHeader.EEntry > e.Phdr[data_seg_index].PVaddr {
		r = Report{
			R_class: ELF_DATA_PADDING,
			R_info: []string{
				"[DECTED] classification : data segment infection",
				"reasons for above conclusion :",
				"Entry point : 0x" + strconv.FormatUint(uint64(e.ElfHeader.EEntry), 16),
				"Data Segment address : 0x" + strconv.FormatUint(uint64(e.Phdr[data_seg_index].PVaddr), 16),
				"",
			},
		}
		return r
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
