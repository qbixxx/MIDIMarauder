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
	ui.Tree.SetSelectedFunc(func(node *tview.TreeNode) { node.SetExpanded(!node.IsExpanded()) })

	return ui
}

func (ui *UI) GetMIDIStream() *tview.TextView {
	return ui.MidiStream
}

func (ui *UI) GetMenu() *tview.TreeNode {
	return ui.Menu
}

func (ui *UI) AddDevice2Menu(dev midi.MidiDevice) {
	node := tview.NewTreeNode(dev.Product).SetSelectable(true).SetColor(tcell.ColorRed)

	node.AddChild(createDeviceNode("Manufacturer", dev.Manufacturer, false))
	node.AddChild(createDeviceNode("Product", dev.Product, false))
	node.AddChild(createDeviceNode("PID", "0x"+dev.PID.String(), false))
	node.AddChild(createDeviceNode("VID", "0x"+dev.VID.String(), false))
	node.AddChild(createDeviceNode("Class", dev.Class, false))
	node.AddChild(createDeviceNode("SubClass", dev.SubClass.String(), false))
	node.AddChild(createDeviceNode("Protocol", dev.Protocol.String(), false))
	node.AddChild(createDeviceNode("Serial Number", dev.SerialNumber, false))
	node.AddChild(createDeviceNode("IN Endpoint", dev.EndpointIn.String(), false))

	node.Collapse()
	ui.Menu.AddChild(node)
}

func createDeviceNode(label string, value string, selectable bool) *tview.TreeNode {
	node := tview.NewTreeNode(fmt.Sprintf("%s: %s", label, value)).SetSelectable(selectable)
	return node
}

func CreateMidiStream() *tview.TextView {
	return tview.NewTextView().SetDynamicColors(true)
}
