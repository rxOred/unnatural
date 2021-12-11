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
	//done := false

	prevkey := ""

	for e := range ui.PollEvents() {
		if e.Type == ui.KeyboardEvent {
			switch e.ID {
			case "<C-c>":
				ui.Clear()
				os.Exit(1)

			case "<C-s>":
				if err := av.StartAnalysis(); err != nil {
					ShowErrorView(ev, err.Error())
				}
				//			done = true
				/*
					case "<C-d>":
						// dis infect
						if done {
							if err := av.StartDisInfection(); err != nil {
								ShowErrorView(ev, err.Error())
							}
						}
				*/
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

			case "g":
				if prevkey == "g" {
					av.a_selected_list.ScrollTop()
					ui.Render(av.a_grid)
				}

			case "G":
				av.a_selected_list.ScrollBottom()
				ui.Render(av.a_grid)

			default:
				av.a_analysis_report.Rows = append(av.a_analysis_report.Rows, prevkey)
				av.a_analysis_report.Rows = append(av.a_analysis_report.Rows, e.ID)
			}

			prevkey = e.ID
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
