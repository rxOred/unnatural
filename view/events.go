package view

import (
	"os"

	ui "github.com/gizak/termui/v3"
)

func (av *AnalysisView) Eventloop() {
	for e := range ui.PollEvents() {
		if e.Type == ui.KeyboardEvent {
			switch e.ID {
			case "<C-c>":
				os.Exit(1)
			case "<C-s>":
				_ = av.StartAnalysis()
			case "<Down>":
			case "<Up>":
			case "<Left>":
			case "<Right>":
			default:
				av.a_report.Rows = append(av.a_report.Rows, "key:"+e.ID)
				ui.Render(av.a_grid)
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
