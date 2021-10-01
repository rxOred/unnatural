package parser

type Ehdr struct {
	EIdent     [16]byte
	EType      uint16
	EMachine   uint16
	EVersion   uint32
	EEntry     uint64
	EPhoff     uint64
	EShoff     uint64
	EFlags     uint32
	EEhsize    uint16
	EPhentsize uint16
	EPhnum     uint16
	EShentsize uint16
	EShnum     uint16
	EShstrndx  uint16
}

type Phdr struct {
	PType   uint32
	POffset uint64
	PVaddr  uint64
	PPaddr  uint64
	PFilesz uint32
	PMemsz  uint32
	PFlags  uint32
	PAlign  uint32
}

type Shdr struct {
	sh_name      uint32
	sh_type      uint32
	sh_flags     uint32
	sh_addr      uint64
	sh_offset    uint64
	sh_size      uint32
	sh_link      uint32
	sh_info      uint32
	sh_addralign uint32
	sh_entsize   uint32
}

type Sym struct {
	st_name  uint32
	st_value uint64
	st_size  uint64
	st_info  byte
	st_other byte
	st_shndx uint16
}
