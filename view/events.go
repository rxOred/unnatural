package view

import (
	"os"

	ui "github.com/gizak/termui/v3"
)

func (av *AnalysisView) Eventloop() {
	for e := range ui.PollEvents() {
		for e.Type == ui.KeyboardEvent {
			switch e.ID {
			case "<C-c>":
				os.Exit(1)
			case "<C-s>":
				av.Guagebar.Title = "Analysing"
				// start analysis
			}
		}
	}
}

func (ev *ErrorView) Eventloop() {
	for e := range ui.PollEvents() {
		for e.Type == ui.KeyboardEvent {
			os.Exit(1)
		}
	}
}
