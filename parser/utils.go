package parser

import (
	"os"

	mmap "github.com/edsrzf/mmap-go"
)

func openFile(pathname string) (mmap.MMap, error) {
	f, err := os.OpenFile(pathname, os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	mmap, err := mmap.Map(f, mmap.RDWR, 0)
	if err != nil {
		return nil, err
	}

	return mmap, nil
}

func verifyElf(ident []byte) bool {
	if ident[0] == 0x7f && ident[1] == 0x45 && ident[2] == 0x4c && ident[3] == 0x46 {
		return true
	}
	return false
}
