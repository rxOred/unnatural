package parser

import (
	"debug/elf"
	"encoding/binary"
	"errors"
	"os"
)

func (e *ElfFile) WriteElfHeader() error {
	e.File.Seek(0, os.SEEK_SET)
	err := binary.Write(e.File, binary.LittleEndian, e.ElfHeader)
	return err
}

func (e *ElfFile) WritePhdrTable() error {
	e.File.Seek(int64(e.ElfHeader.EPhoff), os.SEEK_SET)
	err := binary.Write(e.File, binary.LittleEndian, e.Phdr)
	return err
}

func (e *ElfFile) WriteSegment(buffer []byte, index uint32) error {
	e.File.Seek(int64(e.Phdr[index].POffset), os.SEEK_SET)
	err := binary.Write(e.File, binary.LittleEndian, buffer)
	return err
}

func (e *ElfFile) WriteSection(buffer []byte, index int) error {
	i, err := e.GetSectionIndexByName(".text")
	if index == i {
		if e.ElfHeader.EIdent[elf.EI_DATA] == byte(elf.ELFDATA2LSB) {
			e.File.Seek(int64(e.Shdr[index].ShOffset), os.SEEK_SET)
			err := binary.Write(e.File, binary.LittleEndian, buffer)
			if err != nil {
				return nil
			}
		} else if e.ElfHeader.EIdent[elf.EI_DATA] == byte(elf.ELFDATA2MSB) {
			e.File.Seek(int64(e.Shdr[index].ShOffset), os.SEEK_SET)
			err := binary.Write(e.File, binary.BigEndian, buffer)
			if err != nil {
				return nil
			}
		} else {
			return errors.New("Unknown data format")
		}
	} else {
		e.File.Seek(int64(e.Shdr[index].ShOffset), os.SEEK_SET)
		err := binary.Write(e.File, binary.LittleEndian, buffer)
		if err != nil {
			return nil
		}
	}
	return err
}
