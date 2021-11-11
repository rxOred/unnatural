// +build arm 386

package parser

type Phdr struct {
	p_type   uint32
	p_offset uint32
	p_vaddr  uint32
	p_paddr  uint32
	p_filesz uint32
	p_memsz  uint32
	p_flags  uint32
	p_align  uint32
}

type Shdr struct {
	sh_name      uint32
	sh_type      uint32
	sh_flags     uint32
	sh_addr      uint32
	sh_offset    uint32
	sh_size      uint32
	sh_link      uint32
	sh_info      uint32
	sh_addralign uint32
	sh_entsize   uint32
}

type Sym struct {
	st_name  uint32
	st_value uint32
	st_size  uint32
	st_info  byte
	st_other byte
	st_shndx uint16
}
