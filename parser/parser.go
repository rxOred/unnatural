package parser

import (
	"debug/elf"
	"encoding/binary"
	"errors"
	"os"
	"strconv"

	"github.com/ghostiam/binstruct"
)

type ElfFile struct {
	File      *os.File // file
	ElfHeader Ehdr     // elf header
	Phdr      []*Phdr  // program header table
	Shdr      []*Shdr  // section header table
	Symtab    []*Sym   // symbol table
	Strtab    []byte   // string table
	shstrtab  []byte   // section header string table. pretty much useless for user of this module
}

// return name of a section from the index
func (e *ElfFile) GetSectionNameByIndex(index uint32) (string, error) {
	if e.ElfHeader.EShstrndx == 0 {
		return "", errors.New("shstrndx not found")
	}

	if uint32(e.Shdr[e.ElfHeader.EShstrndx].ShSize) < index {
		return "", errors.New("index out of range")
	}

	return string(e.shstrtab[index]), nil
}

func (e *ElfFile) GetSectionIndexByName(name string) (int, error) {
	for i := 0; i < int(e.ElfHeader.EShnum); i++ {
		str, err := e.GetSectionNameByIndex(e.Shdr[i].ShName)
		if err != nil {
			return -1, err
		}
		if str == name {
			return i, nil
		}
	}

	return -1, nil
}

// parse section header string table if it exists
func (e *ElfFile) ParseSectionHeaderStringTable() error {
	if e.ElfHeader.EShstrndx <= 0 {
		return errors.New("Failed to find section header string table")
	}
	e.File.Seek(int64(e.Shdr[e.ElfHeader.EShstrndx].ShOffset), os.SEEK_SET)
	e.shstrtab = make([]byte, e.Shdr[e.ElfHeader.EShstrndx].ShSize)
	e.File.Read(e.shstrtab)

	return nil
}

// return elf header in a string array
func (e *ElfFile) GetElfHeader() []string {
	var str []string
	switch e.ElfHeader.EType {
	case uint16(elf.ET_REL):
		str = append(str, elf.ET_REL.String())
	case uint16(elf.ET_EXEC):
		str = append(str, elf.ET_EXEC.String())
	case uint16(elf.ET_DYN):
		str = append(str, elf.ET_DYN.String())
	case uint16(elf.ET_CORE):
		str = append(str, elf.ET_CORE.String())
	default:
		str = append(str, "NONE")
	}

	switch e.ElfHeader.EMachine {
	case uint16(elf.EM_386):
		str = append(str, elf.EM_386.String())
	case uint16(elf.EM_MIPS):
		str = append(str, elf.EM_MIPS.String())
	case uint16(elf.EM_ARM):
		str = append(str, elf.EM_ARM.String())
	case uint16(elf.EM_X86_64):
		str = append(str, elf.EM_X86_64.String())
	default:
		str = append(str, "NONE")
	}

	switch e.ElfHeader.EVersion {
	case uint32(elf.EV_CURRENT):
		str = append(str, elf.EV_CURRENT.String())
	default:
		str = append(str, elf.EV_NONE.String())
	}

	str = append(str, strconv.FormatUint(e.ElfHeader.EEntry, 16))
	str = append(str, strconv.FormatUint(e.ElfHeader.EPhoff, 16))
	str = append(str, strconv.FormatUint(e.ElfHeader.EShoff, 16))
	str = append(str, strconv.FormatUint(uint64(e.ElfHeader.EPhnum), 10))
	str = append(str, strconv.FormatUint(uint64(e.ElfHeader.EShnum), 10))
	str = append(str, strconv.FormatUint(uint64(e.ElfHeader.EShstrndx), 10))

	return str
}

// return section header table in an array of string arrays
func (e *ElfFile) GetSectionHeaders() ([][]string, error) {
	var hdrtab [][]string

	// parsing section names
	err := e.ParseSectionHeaderStringTable()
	if err != nil {

	}

	for i := 0; i < int(e.ElfHeader.EShnum); i++ {
		str := make([]string, SHDR_TABLE_ENTRY_COUNT)
		name, err := e.GetSectionNameByIndex(e.Shdr[i].ShName)
		// if we could not
		if err != nil {
			str[0] = strconv.FormatUint(uint64(e.Shdr[i].ShName), 10)
		}
		str[0] = name

		switch e.Shdr[i].ShType {
		case uint32(elf.SHT_NULL):
			str[1] = elf.SHT_NULL.String()
		case uint32(elf.SHT_PROGBITS):
			str[1] = elf.SHT_PROGBITS.String()
		case uint32(elf.SHT_SYMTAB):
			str[1] = elf.SHT_SYMTAB.String()
		case uint32(elf.SHT_STRTAB):
			str[1] = elf.SHT_STRTAB.String()
		case uint32(elf.SHT_RELA):
			str[1] = elf.SHT_RELA.String()
		case uint32(elf.SHT_REL):
			str[1] = elf.SHT_REL.String()
		case uint32(elf.SHT_HASH):
			str[1] = elf.SHT_HASH.String()
		case uint32(elf.SHT_DYNAMIC):
			str[1] = elf.SHT_NOTE.String()
		case uint32(elf.SHT_NOBITS):
			str[1] = elf.SHT_NOBITS.String()
		case uint32(elf.SHT_SHLIB):
			str[1] = elf.SHT_DYNAMIC.String()
		case uint32(elf.SHT_DYNSYM):
			str[1] = elf.SHT_DYNSYM.String()
		case uint32(elf.SHT_LOPROC):
			str[1] = elf.SHT_LOPROC.String()
		case uint32(elf.SHT_HIPROC):
			str[1] = elf.SHT_HIPROC.String()
		case uint32(elf.SHT_LOUSER):
			str[1] = elf.SHT_LOPROC.String()
		case uint32(elf.SHT_HIUSER):
			str[1] = elf.SHT_HIUSER.String()
		default:
			str[1] = "NONE"
		}

		switch e.Shdr[i].ShFlags {
		case uint64(elf.SHF_ALLOC):
			str[2] = elf.SHF_ALLOC.String()
		case uint64(elf.SHF_WRITE):
			str[2] = elf.SHF_WRITE.String()
		case uint64(elf.SHF_EXECINSTR):
			str[2] = elf.SHF_EXECINSTR.String()
		case uint64(elf.SHF_MASKPROC):
			str[2] = elf.SHF_MASKPROC.String()
		default:
			str[2] = "NONE"
		}

		str[3] = strconv.FormatUint(e.Shdr[i].ShAddr, 16)
		str[4] = strconv.FormatUint(e.Shdr[i].ShOffset, 16)
		str[5] = strconv.FormatUint(e.Shdr[i].ShSize, 16)
		str[6] = strconv.FormatUint(uint64(e.Shdr[i].ShLink), 16)
		str[7] = strconv.FormatUint(uint64(e.Shdr[i].ShInfo), 16)
		str[8] = strconv.FormatUint(e.Shdr[i].ShAddralign, 16)
		str[9] = strconv.FormatUint(e.Shdr[i].ShEntsize, 16)

		hdrtab = append(hdrtab, str)
	}

	return hdrtab, nil
}

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

func (e *ElfFile) ParseSymbolTable() error {
	var type_index, name_index int
	for i := 0; i < int(e.ElfHeader.EShnum); i++ {
		if e.Shdr[i].ShType == uint32(elf.SHT_SYMTAB) {
			type_index = i
		}
	}

	name_index, err := e.GetSectionIndexByName(".symtab")
	if err != nil {

	}

	if name_index != type_index {
		// indexes for symtab does not match, which means we are fucked
	}
}

// load the elf, parse main header tables
func LoadElf(e *ElfFile, pathname string) error {
	f, err := openFile(pathname)
	if err != nil {
		return err
	}
	e.File = f

	// parsing elf header
	decoder := binstruct.NewDecoder(e.File, binary.LittleEndian)
	err = decoder.Decode(&e.ElfHeader)
	if err != nil {
		return err
	}

	// parsing program header table
	for i := 0; i < int(e.ElfHeader.EPhnum); i++ {
		e.File.Seek(int64(e.ElfHeader.EPhoff+uint64(i*int(e.ElfHeader.EPhentsize))), os.SEEK_SET)
		decoder = binstruct.NewDecoder(e.File, binary.LittleEndian)
		ph := new(Phdr)
		err = decoder.Decode(ph)
		if err != nil {
			return err
		}
		e.Phdr = append(e.Phdr, ph)
	}

	// parsing section header table
	for i := 0; i < int(e.ElfHeader.EShnum); i++ {
		e.File.Seek(int64(e.ElfHeader.EShoff+uint64(i*int(e.ElfHeader.EShentsize))), os.SEEK_SET)
		decoder = binstruct.NewDecoder(e.File, binary.LittleEndian)
		sh := new(Shdr)
		err = decoder.Decode(sh)
		if err != nil {
			return err
		}
		e.Shdr = append(e.Shdr, sh)
	}

	return nil
}

func (e *ElfFile) UnloadFile() {
	e.File.Close()
}
