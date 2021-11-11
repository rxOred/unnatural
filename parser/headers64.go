// +build arm64 amd64

package parser

type Phdr struct {
	PType   uint32
	PFlags  uint32
	POffset uint64
	PVaddr  uint64
	PPaddr  uint64
	PFilesz uint64
	PMemsz  uint64
	PAlign  uint64
}

type Shdr struct {
	ShName      uint32
	ShType      uint32
	ShFlags     uint64
	ShAddr      uint64
	ShOffset    uint64
	ShSize      uint64
	ShLink      uint32
	ShInfo      uint32
	ShAddralign uint64
	ShEntsize   uint64
}

type Sym struct {
	StName  uint32
	StInfo  byte
	StOther byte
	StShndx uint16
	StValue uint64
	StSize  uint64
}
