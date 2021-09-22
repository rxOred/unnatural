package disinfect

import (
	"debug/elf"
	"errors"
)

func DisinfectTextPaddingInfection(f *elf.File) error {
	sec := f.Section(".init")
	for i := 0; i < len(f.Progs); i++ {
		if f.Progs[i].Type == elf.PT_LOAD && f.Progs[i].Flags == elf.PF_R|elf.PF_X {
			if sec.Addr >= f.Progs[i].Vaddr && sec.Addr < f.Progs[i].Vaddr+f.Progs[i].Memsz {
				f.Entry = sec.Addr
				return nil
			}
		}
	}
	return errors.New("Disinfection Failed")
}

func DisinfectDataSegmentInfection(f *elf.File) error {
	if err := DisinfectTextPaddingInfection(f); err != nil {
		return err
	}
	return nil
}

func DisinfectReverseTextInfection(f *elf.File) error {

}

func DisinfectFunctionPaddingInfection(f *elf.File) error {

}
