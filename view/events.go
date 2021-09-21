package view

import (
	"os"
	"os/signal"
	"syscall"

	ui "github.com/gizak/termui/v3"
)

func (av *AnalysisView) Eventloop(ev *ErrorView) {
	sig_term := make(chan os.Signal, 2)
	signal.Notify(sig_term, os.Interrupt, syscall.SIGTERM)

	var highlight int8 = 0

	for e := range ui.PollEvents() {
		if e.Type == ui.KeyboardEvent {
			switch e.ID {
			case "<C-c>":
				os.Exit(1)

			case "<C-s>":
				if err := av.StartAnalysis(); err != nil {
					ShowErrorView(ev, err.Error())
				}

			case "<Left>":
				highlight--
				if highlight == -1 {
					highlight = 3
				}
				av.HighLight(highlight)

			case "<Right>":
				highlight++
				if highlight == 4 {
					highlight = 0
				}
				av.HighLight(highlight)

			case "<Down>":
				av.a_selected_list.ScrollDown()
				ui.Render(av.a_grid)

			case "<Up>":
				av.a_selected_list.ScrollUp()
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
