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
