package view

import (
	"fmt"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/rxOred/unnatural/analyser"
	disinfect "github.com/rxOred/unnatural/disinfect"
	parser "github.com/rxOred/unnatural/parser"
)

// Analysis view
type AnalysisView struct {
	a_elf             *parser.Elf
	a_grid            *ui.Grid
	a_header          *widgets.Paragraph
	a_guage           *widgets.Gauge
	a_section_list    *widgets.List
	a_elf_header_list *widgets.List
	a_symbol_list     *widgets.List
	a_analysis_report *widgets.List

	// currently selected list
	a_selected_list *widgets.List
}

// Error view
type ErrorView struct {
	e_grid     *ui.Grid
	e_errorbox *widgets.Paragraph
}

const (
	ANA_ELFHEADER_HI    int8 = 0
	ANA_SECTIONS_HI     int8 = 1
	ANA_SYMBOLS_HI      int8 = 2
	ANA_ANALYSIS_REPORT int8 = 3
)

var (
	p parser.Parser
)

// creates and return a new list
func CreateList(title string, border bool, borderfg ui.Color, titlefg ui.Color) *widgets.List {
	list := widgets.NewList()

	list.Title = title
	list.Border = border
	list.BorderStyle.Fg = borderfg
	list.TitleStyle.Fg = titlefg

	return list
}

// creates and return a new paragraph
func CreateParagraph(title string, text string, border bool, borderfg ui.Color, titlefg ui.Color) *widgets.Paragraph {
	para := widgets.NewParagraph()

	para.Title = title
	para.Text = text
	para.Border = border
	para.BorderStyle.Fg = borderfg
	para.TitleStyle.Fg = titlefg

	return para
}

// creates and return a new guage
func CreateGuage(title string, percent int, borderfg ui.Color, titlefg ui.Color) *widgets.Gauge {
	g := widgets.NewGauge()

	g.Title = title
	g.Percent = percent
	g.BorderStyle.Fg = borderfg
	g.TitleStyle.Fg = titlefg

	return g
}

func increasePercent(val int, g *widgets.Gauge) {
	if g.Percent+val > 100 {
		g.Percent += (100 - g.Percent)
	} else {
		g.Percent += val
	}
}

// Analysis view
func (av *AnalysisView) HighLight(highlight int8) {
	switch highlight {
	case ANA_ELFHEADER_HI:
		av.a_elf_header_list.BorderStyle.Fg = ui.ColorRed
		av.a_section_list.BorderStyle.Fg = ui.ColorYellow
		av.a_symbol_list.BorderStyle.Fg = ui.ColorYellow
		av.a_analysis_report.BorderStyle.Fg = ui.ColorYellow

		av.a_selected_list = av.a_elf_header_list

	case ANA_SECTIONS_HI:
		av.a_section_list.BorderStyle.Fg = ui.ColorRed
		av.a_elf_header_list.BorderStyle.Fg = ui.ColorYellow
		av.a_symbol_list.BorderStyle.Fg = ui.ColorYellow
		av.a_analysis_report.BorderStyle.Fg = ui.ColorYellow

		av.a_selected_list = av.a_section_list

	case ANA_SYMBOLS_HI:
		av.a_symbol_list.BorderStyle.Fg = ui.ColorRed
		av.a_elf_header_list.BorderStyle.Fg = ui.ColorYellow
		av.a_section_list.BorderStyle.Fg = ui.ColorYellow
		av.a_analysis_report.BorderStyle.Fg = ui.ColorYellow

		av.a_selected_list = av.a_symbol_list

	case ANA_ANALYSIS_REPORT:
		av.a_analysis_report.BorderStyle.Fg = ui.ColorRed
		av.a_elf_header_list.BorderStyle.Fg = ui.ColorYellow
		av.a_section_list.BorderStyle.Fg = ui.ColorYellow
		av.a_symbol_list.BorderStyle.Fg = ui.ColorYellow

		av.a_selected_list = av.a_analysis_report
	}

	ui.Render(av.a_grid)
}

func (av *AnalysisView) SetupAnalysisGrid() error {
	av.a_grid = ui.NewGrid()

	top := ui.NewRow(1.0/8, av.a_header)
	belowtop := ui.NewRow(1.0/8, av.a_guage)
	mid := ui.NewRow(2.0/7,
		ui.NewCol(1.0/3, av.a_elf_header_list),
		ui.NewCol(1.0/3, av.a_section_list),
		ui.NewCol(1.0/3, av.a_symbol_list),
	)
	bottom := ui.NewRow(4.0/8, av.a_analysis_report)

	av.a_grid.Set(top, belowtop, mid, bottom)

	termwidth, termheight := ui.TerminalDimensions()
	av.a_grid.SetRect(0, 0, termwidth, termheight-1)
	return nil
}

func (av *AnalysisView) StartAnalysis() error {
	increasePercent(3, av.a_guage)
	av.a_guage.Title = "Analysing"
	ui.Render(av.a_grid)

	// text padding infection
	r := analyser.CheckSegmentInfections(av.a_elf.E_file)
	if r.R_class == analyser.ELF_TEXT_PADDING {
		for i := 0; i < len(r.R_info); i++ {
			av.a_analysis_report.Rows = append(av.a_analysis_report.Rows, r.R_info[i])
		}
	}
	increasePercent(3, av.a_guage)
	ui.Render(av.a_grid)
	return nil
}

func (av *AnalysisView) StartDisInfection() error {
	if r := disinfect.DisinfectTextPaddingInfection(av.a_elf.E_file); r != nil {
		for i := 0; i < len(r); i++ {
			av.a_analysis_report.Rows = append(av.a_analysis_report.Rows, r[i])
		}
	}
	av.a_analysis_report.Rows = append(av.a_analysis_report.Rows, "")

	if r := disinfect.DisinfectDataSegmentInfection(av.a_elf.E_file); r != nil {
		for i := 0; i < len(r); i++ {
			av.a_analysis_report.Rows = append(av.a_analysis_report.Rows, r[i])
		}
	}
	av.a_analysis_report.Rows = append(av.a_analysis_report.Rows, "")
	return nil
}

func InitAnalysisWidgets(av *AnalysisView, ev *ErrorView, file string) {
	text := fmt.Sprintf("\t\tunnatural - Elf anomaly detector and disinfector\t\t\nTarget: %s", file)
	av.a_header = CreateParagraph("", text, true, ui.ColorYellow, ui.ColorCyan)
	av.a_guage = CreateGuage("", 0, ui.ColorYellow, ui.ColorCyan)

	// elf header is the initial highlight
	av.a_elf_header_list = CreateList("Elf header", true, ui.ColorRed, ui.ColorCyan)
	av.a_selected_list = av.a_elf_header_list

	av.a_section_list = CreateList("Sections", true, ui.ColorYellow, ui.ColorCyan)
	av.a_symbol_list = CreateList("Symbols", true, ui.ColorYellow, ui.ColorCyan)
	av.a_analysis_report = CreateList("Analysis report", true, ui.ColorMagenta, ui.ColorRed)
	err := av.SetupAnalysisGrid()
	if err != nil {
		ShowErrorView(ev, err.Error())
	}
	if av.a_elf, err = parser.InitElf(file); err != nil {
		ShowErrorView(ev, err.Error())
	}

	p = av.a_elf
	header := p.GetElfHeader()
	sections := p.GetSectionHeaders()
	symbols := p.GetSymbols()

	increasePercent(10, av.a_guage)

	av.a_symbol_list.Rows = symbols
	av.a_elf_header_list.Rows = header
	av.a_section_list.Rows = sections

	increasePercent(10, av.a_guage)
}

// clear screen, render the UI, start eventloop
func ShowAnalysisView(av *AnalysisView, ev *ErrorView) {
	ui.Clear()
	ui.Render(av.a_grid)
	av.Eventloop(ev)
}

// error view
func (ev *ErrorView) SetupErrorGrid() {
	ev.e_grid = ui.NewGrid()
	body := ui.NewRow(1.0, ev.e_errorbox)

	ev.e_grid.Set(body)
	termwidth, termheight := ui.TerminalDimensions()
	ev.e_grid.SetRect(0, 0, termwidth, termheight-1)
}

func InitErrorView(ev *ErrorView) {
	ev.e_errorbox = CreateParagraph("Error", "", true, ui.ColorYellow, ui.ColorWhite)
	ev.SetupErrorGrid()
}

func ShowErrorView(ev *ErrorView, errmsg string) {
	ev.e_errorbox.Text = errmsg
	ui.Clear()
	ui.Render(ev.e_grid)
	ev.Eventloop()
}
