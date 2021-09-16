package view

import (
	"fmt"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	parser "github.com/rxOred/unnatural/parser"
)

// Analysis view
type AnalysisView struct {
	e             *parser.Elf
	grid          *ui.Grid
	Header        *widgets.Paragraph
	Guagebar      *widgets.Gauge
	SectionList   *widgets.List
	ElfheaderList *widgets.List
	SymbolList    *widgets.List
	OutData       *widgets.List
}

// Error view
type ErrorView struct {
	grid     *ui.Grid
	ErrorBox *widgets.Paragraph
}

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
func (av *AnalysisView) SetupAnalysisGrid() error {
	av.grid = ui.NewGrid()

	top := ui.NewRow(1.0/8, av.Header)
	belowtop := ui.NewRow(1.0/8, av.Guagebar)
	mid := ui.NewRow(2.0/7,
		ui.NewCol(1.0/3, av.SectionList),
		ui.NewCol(1.0/3, av.ElfheaderList),
		ui.NewCol(1.0/3, av.SymbolList),
	)
	bottom := ui.NewRow(4.0/8, av.OutData)

	av.grid.Set(top, belowtop, mid, bottom)

	termwidth, termheight := ui.TerminalDimensions()
	av.grid.SetRect(0, 0, termwidth, termheight-1)
	return nil
}

func (av *AnalysisView) StartAnalysis() error {
	increasePercent(10, av.Guagebar)

	CheckTextPaddingInfection(av.e)
}

func InitAnalysisView(av *AnalysisView, ev *ErrorView, file string) {
	text := fmt.Sprintf("\t\tunnatural - Elf anomaly detector and disinfector\t\t\nTarget: %s", file)
	av.Header = CreateParagraph("", text, true, ui.ColorYellow, ui.ColorCyan)
	av.Guagebar = CreateGuage("Analysing", 0, ui.ColorYellow, ui.ColorCyan)
	av.ElfheaderList = CreateList("Elf header", true, ui.ColorYellow, ui.ColorCyan)
	av.SectionList = CreateList("Sections", true, ui.ColorYellow, ui.ColorCyan)
	av.SymbolList = CreateList("Symbols", true, ui.ColorYellow, ui.ColorCyan)
	av.OutData = CreateList("Analysis report", true, ui.ColorMagenta, ui.ColorRed)
	err := av.SetupAnalysisGrid()
	if err != nil {
		ShowErrorView(ev, err.Error())
	}
	if av.e, err = parser.InitElf(file); err != nil {
		ShowErrorView(ev, err.Error())
	}

	p = av.e
	header := p.GetElfHeader()
	sections := p.GetSectionHeaders()
	symbols := p.GetSymbols()

	increasePercent(10, av.Guagebar)

	av.SymbolList.Rows = symbols
	av.ElfheaderList.Rows = header
	av.SectionList.Rows = sections

	increasePercent(10, av.Guagebar)
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
