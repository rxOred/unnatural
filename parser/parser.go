package parser

import (
	"debug/elf"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/ghostiam/binstruct"
)

type ElfFile struct {
	pathname  string  // pathname
	ElfHeader Ehdr    // elf header
	Phdr      []*Phdr // program header table
	Shdr      []*Shdr // section header table
	Symtab    []*Sym  // symbol table
	Shstrtab  []byte  // string table -- should not be accesss outside
}

// return name of a section from the index
/*
func (e *ElfFile) GetSectionName(index int) (string, error) {
	if e.ElfHeader.EShstrndx == 0 {
		return "fail", errors.New("shstrndx not found")
    }
    return &e.strtab[index], nil
}
*/

func (e *ElfFile) ParseSectionHeaderStringTable(f *os.File) error {
	if e.ElfHeader.EShstrndx <= 0 {
		return errors.New("Failed to find section header string table")
	}

	fmt.Println("offset is here ma nigga !!! ", e.Shdr[e.ElfHeader.EShstrndx].ShOffset)
	//var b [e.Shdr[shstrndx].ShSize]byte
	f.Seek(int64(e.Shdr[e.ElfHeader.EShstrndx].ShOffset), os.SEEK_SET)
	e.Shstrtab = make([]byte, e.Shdr[e.ElfHeader.EShstrndx].ShSize)
	f.Read(e.Shstrtab)

	return nil
}

// return elf header in a string array
func (e *ElfFile) GetElfHeader() []string {
	var str []string
	switch e.ElfHeader.EType {
	case uint16(elf.ET_REL):
		str = append(str, "Type :"+elf.ET_REL.String())
	case uint16(elf.ET_EXEC):
		str = append(str, "Type :"+elf.ET_EXEC.String())
	case uint16(elf.ET_DYN):
		str = append(str, "Type :"+elf.ET_DYN.String())
	case uint16(elf.ET_CORE):
		str = append(str, "Type :"+elf.ET_CORE.String())
	default:
		str = append(str, "Type :"+"none")
	}

	switch e.ElfHeader.EMachine {
	case uint16(elf.EM_386):
		str = append(str, "Machine :"+elf.EM_386.String())
	case uint16(elf.EM_MIPS):
		str = append(str, "Machine :"+elf.EM_MIPS.String())
	case uint16(elf.EM_ARM):
		str = append(str, "Machine :"+elf.EM_ARM.String())
	case uint16(elf.EM_X86_64):
		str = append(str, "Machine :"+elf.EM_X86_64.String())
	default:
		str = append(str, "Machine :"+"none")
	}

	switch e.ElfHeader.EVersion {
	case uint32(elf.EV_CURRENT):
		str = append(str, "Version :"+elf.EV_CURRENT.String())
	default:
		str = append(str, "Version :"+elf.EV_NONE.String())
	}

	str = append(str, "Entry :"+strconv.FormatUint(e.ElfHeader.EEntry, 16))
	str = append(str, "Phoff :"+strconv.FormatUint(e.ElfHeader.EPhoff, 16))
	str = append(str, "Shoff :"+strconv.FormatUint(e.ElfHeader.EShoff, 16))
	str = append(str, "Phnum :"+strconv.FormatUint(uint64(e.ElfHeader.EPhnum), 10))
	str = append(str, "Shnum :"+strconv.FormatUint(uint64(e.ElfHeader.EShnum), 10))
	str = append(str, "Shstrndx :"+strconv.FormatUint(uint64(e.ElfHeader.EShstrndx), 10))
	return str
}

// return section header table in an array of string arrays
/*
func (e *ElfFile) GetSectionHeaders() [][]string {
	var str [][]string
	for i := 0; i < len(e.Shdr); i++ {
		str[i] = append(str[i], GetSectionName(e.Shdr[i].ShName))
	}
}
*/

// return program header table in an array of string arrays
func (e *ElfFile) GetProgHeaders() [][]string {
	var hdrtab [][]string
	for i := 0; i < int(e.ElfHeader.EPhnum); i++ {
		str := make([]string, PHDR_TABLE_ENTRY_COUNT)
		switch e.Phdr[i].PType {
		case uint32(elf.PT_LOAD):
			str[0] = elf.PT_LOAD.String()
		case uint32(elf.PT_DYNAMIC):
			str[0] = elf.PT_DYNAMIC.String()
		case uint32(elf.PT_INTERP):
			str[0] = elf.PT_INTERP.String()
		case uint32(elf.PT_NULL):
			str[0] = elf.PT_NULL.String()
		case uint32(elf.PT_NOTE):
			str[0] = elf.PT_NOTE.String()
		case uint32(elf.PT_SHLIB):
			str[0] = elf.PT_SHLIB.String()
		case uint32(elf.PT_PHDR):
			str[0] = elf.PT_PHDR.String()
		case uint32(elf.PT_TLS):
			str[0] = elf.PT_TLS.String()
		default:
			str[0] = "NONE"
		}
		str[1] = strconv.FormatUint(e.Phdr[i].POffset, 16)
		str[2] = strconv.FormatUint(e.Phdr[i].PVaddr, 16)
		str[3] = strconv.FormatUint(e.Phdr[i].PPaddr, 16)
		str[4] = strconv.FormatUint(uint64(e.Phdr[i].PFilesz), 16)
		str[5] = strconv.FormatUint(uint64(e.Phdr[i].PMemsz), 16)

		switch e.Phdr[i].PFlags {
		case uint32(elf.PF_R):
			str[6] = "PF_R"
		case uint32(elf.PF_W):
			str[6] = "PF_W"
		case uint32(elf.PF_X):
			str[6] = "PF_X"
		case uint32(elf.PF_R | elf.PF_X):
			str[6] = "PF_R | PF_X"
		case uint32(elf.PF_R | elf.PF_W):
			str[6] = "PF_R | PF_W"
		case uint32(elf.PF_W | elf.PF_X):
			str[6] = "PF_W | PF_X"
		case uint32(elf.PF_R | elf.PF_W | elf.PF_X):
			str[6] = "PF_R | PF_W | PF_X"
		}
		str[7] = strconv.FormatInt(int64(e.Phdr[i].PAlign), 16)
		hdrtab = append(hdrtab, str)
	}
	return hdrtab
}

// get first segment of a segment type
func (e *ElfFile) GetSegmentByType(ptype uint32) (*Phdr, error) {
	for i := 0; i < int(e.ElfHeader.EPhnum); i++ {
		if e.Phdr[i].PType == ptype {
			return e.Phdr[i], nil
		}
	}
	return nil, errors.New("Segment not found")
}

// get the n segment of a segment type
func (e *ElfFile) GetNSegmentByType(ptype uint32, n int) (*Phdr, error) {
	k := 0
	for i := 0; i < int(e.ElfHeader.EPhnum); i++ {
		if e.Phdr[i].PType == ptype {
			k++
			if k == n {
				return e.Phdr[i], nil
			}
		}
	}
	return nil, errors.New("Segment not found")
}

// parse the whole elf binary
func LoadElf(e *ElfFile, pathname string) error {
	e.pathname = pathname
	f, err := openFile(pathname)
	if err != nil {
		return err
	}

	// parsing elf header
	decoder := binstruct.NewDecoder(f, binary.LittleEndian)
	err = decoder.Decode(&e.ElfHeader)
	if err != nil {
		return err
	}

	// parsing program header table
	for i := 0; i < int(e.ElfHeader.EPhnum); i++ {
		f.Seek(int64(e.ElfHeader.EPhoff+uint64(i*int(e.ElfHeader.EPhentsize))), os.SEEK_SET)
		decoder = binstruct.NewDecoder(f, binary.LittleEndian)
		ph := new(Phdr)
		err = decoder.Decode(ph)
		if err != nil {
			return err
		}
		e.Phdr = append(e.Phdr, ph)
	}

	// parsing section header table
	for i := 0; i < int(e.ElfHeader.EShnum); i++ {
		f.Seek(int64(e.ElfHeader.EShoff+uint64(i*int(e.ElfHeader.EShentsize))), os.SEEK_SET)
		decoder = binstruct.NewDecoder(f, binary.LittleEndian)
		sh := new(Shdr)
		err = decoder.Decode(sh)
		if err != nil {
			return err
		}
		e.Shdr = append(e.Shdr, sh)
	}

	err = e.ParseSectionHeaderStringTable(f)
	if err != nil {
		return err
	}

	return nil
}
