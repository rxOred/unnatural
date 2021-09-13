package view

import (
	"os"

	ui "github.com/gizak/termui/v3"
	parser "github.com/rxOred/unnatural/parser"
)

func (v *View) GotoExitState() {
	for e := range ui.PollEvents() {
		if e.Type == ui.KeyboardEvent {
			os.Exit(1)
		}
	}
}

func (v *View) GotoEventLoopState() {
	for e := range ui.PollEvents() {
		if e.Type == ui.KeyboardEvent {
			switch e.ID {
			case "<C-s>":
				parser.StartAnalysis()
			case "<C-c>":
				os.Exit(1)
			case "<C-i>":
				parser.ShowInfo()
			}
		}
	}
}
