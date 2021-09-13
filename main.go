package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	view "github.com/rxOred/unnatural/view"
)

var binpathFlag = flag.String("path", "", "/path/to/binary")
var helpFlag = flag.Bool("help", false, "print help and exit")
var versionFlag = flag.Bool("version", false, "print version and exit")

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
	if err := ui.Init(); err != nil {
		log.Fatal(err)
	}
	defer ui.Close()

	if len(*binpathFlag) <= 0 {
		er := widgets.NewParagraph()
		er.Text = "[Error]\n/path/to/binary not specified\npress any key to exit..."
		er.SetRect(0, 0, 40, 5)
		ui.Render(er)
		for e := range ui.PollEvents() {
			if e.Type == ui.KeyboardEvent {
				break
			}
		}
	} else {
		var v view.View
		view.InitView(&v, *binpathFlag)
		v.GotoSigIntState()
	}

}
