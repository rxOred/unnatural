package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rxOred/unnatural/parser"
	view "github.com/rxOred/unnatural/view"
)

var (
	binpathFlag = flag.String("path", "", "/path/to/binary")
	helpFlag    = flag.Bool("help", false, "print help and exit")
	versionFlag = flag.Bool("version", false, "print version and exit")

	av view.AnalysisView
	ev view.ErrorView
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
	/*
		if err := ui.Init(); err != nil {
			log.Fatal(err)
		}
		defer ui.Close()

		// initialize TUIs
		view.InitErrorView(&ev)
		view.InitAnalysisWidgets(&av, &ev, *binpathFlag)

		// Show Analysis view to the user
		view.ShowAnalysisView(&av, &ev)
	*/
	var e parser.ElfFile
	parser.LoadElf(&e, *binpathFlag)
	varia := e.GetSymbolNames()
	for i := 0; i < len(varia); i++ {
		fmt.Println(varia[i])
	}
	fmt.Println(len(e.Symtab))
}
