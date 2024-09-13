package ui

import (
	"github.com/rivo/tview"
	"fmt"
)

const asciiTitle = "[cyan]  \\  | _)      | _)\n" +
	" |\\/ |  |   _` |  |                                 \n" +
	" |   |  |  (   |  |                                 \n" +
	"_|\\ _| _| \\__,_| _|                   [turquoise]|             \n" +
	" |\\/ |   _` |   __|  _` |  |   |   _` |   _ \\   __| \n" +
	" |   |  (   |  |    (   |  |   |  (   |   __/  |    \n" +
	"_|  _| \\__,_| _|   \\__,_| \\__,_| \\__,_| \\___| _| \n\n\n[clear]"

type UI struct {
	Root	*tview.Grid
	MidiStream *tview.TextView
	Menu	*tview.TextView
}

func SetupUI() *UI {
	ui := new(UI)
	midiStream := tview.NewTextView().SetDynamicColors(true)
	midiStream.Box.SetBorder(true).SetTitle(" Midi Stream ")

	menu := tview.NewTextView()
	menu.Box.SetBorder(true).SetTitle(" Menu ")
	menu.SetTextAlign(tview.AlignLeft).SetDynamicColors(true)	
	menu.SetText(asciiTitle)


	grid := tview.NewGrid().
		SetColumns(-4, 54).
		SetRows(-2, 2).
		SetBorders(true).
		AddItem(midiStream, 0, 0, 1, 1, 0, 0, true).
		AddItem(menu, 0, 1, 1, 1, 0, 0, true)

	
		ui.Root = grid
	ui.MidiStream = midiStream
	ui.Menu = menu

	return ui
}

func (ui *UI) GetMIDIStream() *tview.TextView {
	return ui.MidiStream
}

func (ui *UI) GetMenu() *tview.TextView {
	return ui.Menu
}

func(ui *UI) AddDevice2Menu(man, prod string){
	fmt.Fprintln(ui.Menu, man + " - " + prod)
	ui.Menu.ScrollToEnd()
}


func CreateMidiStream() *tview.TextView {
	return tview.NewTextView().SetDynamicColors(true)
}
