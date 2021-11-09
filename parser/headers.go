package parser

const (
	ELFHEADER_ENTRY_COUNT  = 14
	PHDR_TABLE_ENTRY_COUNT = 8
	SHDR_TABLE_ENTRY_COUNT = 10
)

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
	ShName      uint32
	ShType      uint32
	ShFlags     uint32
	ShAddr      uint64
	ShOffset    uint64
	ShSize      uint32
	ShLink      uint32
	ShInfo      uint32
	ShAddralign uint32
	ShEntsize   uint32
}

type Sym struct {
	StName  uint32
	StValue uint64
	StSize  uint64
	StInfo  byte
	StOther byte
	StShndx uint16
}
