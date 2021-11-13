package main

import (
	"flag"
	"fmt"
	"os"

	parser "github.com/rxOred/unnatural/parser"
	//	view "github.com/rxOred/unnatural/view"
)

var (
	binpathFlag = flag.String("path", "", "/path/to/binary")
	helpFlag    = flag.Bool("help", false, "print help and exit")
	versionFlag = flag.Bool("version", false, "print version and exit")

//	av view.AnalysisView
//	ev view.ErrorView
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: unnatural [options]\n\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *helpFlag {
		flag.Usage()
		os.Exit(1)
	}

	if *versionFlag {
		flag.Usage()
		os.Exit(0)
	}
}

func main() {
	//if err := ui.Init(); err != nil {
	//	log.Fatal(err)
	//}
	//defer ui.Close()

	// initialize TUIs
	//view.InitErrorView(&ev)
	//view.InitAnalysisWidgets(&av, &ev, *binpathFlag)

	// Show Analysis view to the user
	//view.ShowAnalysisView(&av, &ev)
	var e parser.ElfFile
	parser.LoadElf(&e, *binpathFlag)
	str := e.GetElfHeader()
	for i := 0; i < len(str); i++ {
		fmt.Println(str[i])
	}

	str2 := e.GetProgHeaders()

	for i := 0; i < int(e.ElfHeader.EPhnum); i++ {
		for j := 0; j < 8; j++ {
			fmt.Print(str2[i][j], "\t")
		}
		fmt.Println()
	}
	for i := 0; i < int(e.Shdr[e.ElfHeader.EShstrndx].ShSize); i++ {
		fmt.Print(string(e.Shstrtab[i]))
	}

	str3, err := e.GetSectionHeaders()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	for i := 0; i < int(e.ElfHeader.EPhnum); i++ {
		for j := 0; j < parser.SHDR_TABLE_ENTRY_COUNT; j++ {
			fmt.Print(str3[i][j], "\t")
		}
		fmt.Println()
	}
}
