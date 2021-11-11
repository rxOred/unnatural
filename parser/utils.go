package parser

import (
	"os"
)

func openFile(pathname string) (*os.File, error) {
	f, err := os.OpenFile(pathname, os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func verifyElf(ident []byte) bool {
	if ident[0] == 0x7f && ident[1] == 0x45 && ident[2] == 0x4c && ident[3] == 0x46 {
		return true
	}
	return false
}
