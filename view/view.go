package view

import (
	"fmt"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type View struct {
	Grid *ui.Grid
}

type AnalysisView struct {
	v           *View
	Header      *widgets.Paragraph
	Guagebar    *widgets.Gauge
	SectionList *widgets.List
	SegmentList *widgets.List
	SymbolList  *widgets.List
	OutData     *widgets.List
}

type ErrorView struct {
	v        *View
	ErrorBox *widgets.Paragraph
}

type InfoView struct {
	v         *View
	Header    *widgets.Paragraph
	InfoTable *widgets.Table
}

func CreateList(border_name string) *widgets.List {
	list := widgets.NewList()

	list.Title = border_name
	list.Border = true

	return list
}

func CreateParagraph(title string, text string, border bool) *widgets.Paragraph {
	para := widgets.NewParagraph()

	para.Title = title
	para.Text = text
	para.Border = border

	return para
}

func CreateGuage(title string, percent int, borderfg ui.Color, titlefg ui.Color) *widgets.Gauge {
	g := widgets.NewGauge()

	g.Title = title
	g.Percent = percent
	g.BorderStyle.Fg = borderfg
	g.TitleStyle.Fg = titlefg

	return g
}

func CreateTable(title string) *widgets.Table {
	tab := widgets.NewTable()

	tab.Title = title

	return tab
}

// analysis view
func (av *AnalysisView) SetupAnalysisGrid() {
	av.v.Grid = ui.NewGrid()

	top := ui.NewRow(1.0/8, av.Header)
	belowtop := ui.NewRow(1.0/8, av.Guagebar)
	mid := ui.NewRow(2.0/7,
		ui.NewCol(1.0/3, av.SectionList),
		ui.NewCol(1.0/3, av.SegmentList),
		ui.NewCol(1.0/3, av.SymbolList),
	)
	bottom := ui.NewRow(4.0/8, av.OutData)

	av.v.Grid.Set(top, belowtop, mid, bottom)

	termwidth, termheight := ui.TerminalDimensions()
	av.v.Grid.SetRect(0, 0, termwidth, termheight-1)
}

func NewAnalysisView(av *AnalysisView, target string) {
	text := fmt.Sprintf("\t\tunnatural - Elf anomaly detector and disinfector\t\t\nTarget: %s", target)
	av.Header = CreateParagraph("", text, true)
	av.Guagebar = CreateGuage("Analysing", 0, ui.ColorWhite, ui.ColorCyan)
	av.SectionList = CreateList("Sections")
	av.SegmentList = CreateList("Segments")
	av.SymbolList = CreateList("Symbols")
	av.OutData = CreateList("Analysis report")
	av.SetupAnalysisGrid()
}

func (ev *ErrorView) SetupErrorGrid() {
	ev.v.Grid = ui.NewGrid()
	body := ui.NewRow(1.0, ev.ErrorBox)

	ev.v.Grid.Set(body)
	termwidth, termheight := ui.TerminalDimensions()
	ev.v.Grid.SetRect(0, 0, termwidth, termheight-1)
}

// error view
func NewErrorView(ev *ErrorView, errmsg string) {
	ev.ErrorBox = CreateParagraph("Error", errmsg, true)
	ev.SetupErrorGrid()
}

func (v *View) ShowErrorMsg(errmsg string) {
	ui.Clear()
	er := widgets.NewParagraph()

	errwindow := ui.NewRow(1.0, er)

	v.Grid.Set(errwindow)

	termwidth, termheight := ui.TerminalDimensions()
	v.Grid.SetRect(0, 0, termwidth, termheight)
	ui.Render(v.Grid)
	v.GotoExitState()
}

// info view
func (iv *InfoView) SetupInfoGrid() {
	iv.v.Grid = ui.NewGrid()
	header := ui.NewRow(1.0/3, iv.Header)
	body := ui.NewRow(2.0/3, iv.InfoTable)

	iv.v.Grid.Set(header, body)
	termwidth, termheight := ui.TerminalDimensions()
	iv.v.Grid.SetRect(0, 0, termwidth, termheight-1)
}

func NewInfoView(iv *InfoView) {
	iv.Header = CreateParagraph("Info View", "", false)
	iv.InfoTable = CreateTable()
}
