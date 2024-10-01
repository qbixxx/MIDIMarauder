package ui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"midiMarauder/internal/midi"
)


const asciiTitle = "\n[cyan]    \\  | _)      | _)\n" +
	"   |\\/ |  |   _` |  |                                 \n" +
	"   |   |  |  (   |  |                                 \n" +
	"  _|\\ _| _| \\__,_| _|                   [turquoise]|             \n" +
	"   |\\/ |   _` |   __|  _` |  |   |   _` |   _ \\   __| \n" +
	"   |   |  (   |  |    (   |  |   |  (   |   __/  |    \n" +
	"  _|  _| \\__,_| _|   \\__,_| \\__,_| \\__,_| \\___| _| \n[-:-:-:-]"

type UI struct {
	Root       *tview.Grid
	MidiStream *tview.TextView
	Menu       *tview.TreeNode
	Tree       *tview.TreeView
}

func SetupUI() *UI {
	ui := new(UI)
	midiStream := tview.NewTextView().SetDynamicColors(true)
	midiStream.Box.SetBorder(true).SetTitle(" Midi Stream ")

	title := tview.NewTextView()
	title.Box.SetBorder(false)
	title.SetTextAlign(tview.AlignLeft).
		SetDynamicColors(true).
		SetText(asciiTitle).SetScrollable(false)

	menu := tview.NewGrid()
	menu.Box.SetBorder(true).SetTitle(" Menu ")

	gridMenu := tview.NewGrid()
	gridMenu.SetRows(10, -1)

	gridMenu.AddItem(title, 0, 0, 1, 1, 0, 0, false)
	gridMenu.AddItem(menu, 1, 0, 1, 1, 0, 0, true)

	rootMsg := "MIDI devices:\n"
	rootTree := tview.NewTreeNode(rootMsg).SetSelectable(false).
		SetColor(tcell.ColorGreen)
	tree := tview.NewTreeView().
		SetRoot(rootTree).
		SetCurrentNode(rootTree)

	menu.AddItem(tree, 0, 0, 1, 1, 0, 0, true)

	rootGrid := tview.NewGrid().
		SetColumns(-4, 62).
		SetRows(-2, 1).
		SetBorders(false).
		AddItem(midiStream, 0, 0, 1, 1, 0, 0, true).
		AddItem(gridMenu, 0, 1, 1, 1, 0, 0, true)

	ui.Root = rootGrid
	ui.MidiStream = midiStream
	ui.Menu = rootTree
	ui.Tree = tree
	ui.Tree.SetSelectedFunc(func(node *tview.TreeNode) { node.SetExpanded(!node.IsExpanded()) }).
		SetGraphicsColor(tcell.ColorDarkCyan)

	return ui
}

// AddDevice2Menu adds a MIDI device and its details to the tree view
func (ui *UI) AddDevice2Menu(device midi.MIDIReader) {
	deviceInfo := device.GetDeviceDetails()

	node := tview.NewTreeNode(deviceInfo[1][1]).SetSelectable(true).SetColor(tcell.ColorRed)

	// Add the dev details
	for _, detail := range deviceInfo {
		node.AddChild(createDeviceNode(detail[0], detail[1], false))
	}

	node.Collapse()
	ui.Menu.AddChild(node)
}

// Helper function to create a tree node for device details
func createDeviceNode(label string, value string, selectable bool) *tview.TreeNode {
	return tview.NewTreeNode(fmt.Sprintf("%s: %s", label, value)).SetSelectable(selectable)
}

func (ui *UI) GetMIDIStream() *tview.TextView {
	return ui.MidiStream
}