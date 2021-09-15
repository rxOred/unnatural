package view

import (
	"fmt"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

// Analysis view
type AnalysisView struct {
	grid        *ui.Grid
	Header      *widgets.Paragraph
	Guagebar    *widgets.Gauge
	SectionList *widgets.List
	SegmentList *widgets.List
	SymbolList  *widgets.List
	OutData     *widgets.List
}

// Error view
type ErrorView struct {
	grid     *ui.Grid
	ErrorBox *widgets.Paragraph
}

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

// Analysis view
func (av *AnalysisView) SetupAnalysisGrid() {
	av.grid = ui.NewGrid()

	top := ui.NewRow(1.0/8, av.Header)
	belowtop := ui.NewRow(1.0/8, av.Guagebar)
	mid := ui.NewRow(2.0/7,
		ui.NewCol(1.0/3, av.SectionList),
		ui.NewCol(1.0/3, av.SegmentList),
		ui.NewCol(1.0/3, av.SymbolList),
	)
	bottom := ui.NewRow(4.0/8, av.OutData)

	av.grid.Set(top, belowtop, mid, bottom)

	termwidth, termheight := ui.TerminalDimensions()
	av.grid.SetRect(0, 0, termwidth, termheight-1)
}

func InitAnalysisView(av *AnalysisView, target string) {
	text := fmt.Sprintf("\t\tunnatural - Elf anomaly detector and disinfector\t\t\nTarget: %s", target)
	av.Header = CreateParagraph("", text, true, ui.ColorYellow, ui.ColorCyan)
	av.Guagebar = CreateGuage("Analysing", 0, ui.ColorYellow, ui.ColorCyan)
	av.SectionList = CreateList("elf header", true, ui.ColorYellow, ui.ColorCyan)
	av.SegmentList = CreateList("Sections", true, ui.ColorYellow, ui.ColorCyan)
	av.SymbolList = CreateList("Strings", true, ui.ColorYellow, ui.ColorCyan)
	av.OutData = CreateList("Analysis report", true, ui.ColorMagenta, ui.ColorRed)
	av.SetupAnalysisGrid()
}

func ShowAnalysisView(av *AnalysisView) {
	ui.Clear()
	ui.Render(av.grid)
	av.Eventloop()
}

// error view
func (ev *ErrorView) SetupErrorGrid() {
	ev.grid = ui.NewGrid()
	body := ui.NewRow(1.0, ev.ErrorBox)

	ev.grid.Set(body)
	termwidth, termheight := ui.TerminalDimensions()
	ev.grid.SetRect(0, 0, termwidth, termheight-1)
}

func InitErrorView(ev *ErrorView) {
	ev.ErrorBox = CreateParagraph("Error", "", true, ui.ColorYellow, ui.ColorWhite)
	ev.SetupErrorGrid()
}

func ShowErrorView(ev *ErrorView, errmsg string) {
	ev.ErrorBox.Text = errmsg
	ui.Clear()
	ui.Render(ev.grid)
	ev.Eventloop()
}
