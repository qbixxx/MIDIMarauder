package ui

import (
	"github.com/rivo/tview"
	"github.com/gdamore/tcell/v2"
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
	Root	*tview.Grid
	MidiStream *tview.TextView
	Menu	*tview.TreeNode
	Tree	*tview.TreeView
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
	gridMenu.SetRows(10,-1)	

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
		SetColumns(-4,62).
		SetRows(-2, 1).
		SetBorders(false).
		AddItem(midiStream, 0, 0, 1, 1, 0, 0, true).
		AddItem(gridMenu, 0, 1, 1, 1, 0, 0, true)

	
	ui.Root = rootGrid
	ui.MidiStream = midiStream
	ui.Menu = rootTree
	ui.Tree = tree
	ui.Tree.SetSelectedFunc(func(node *tview.TreeNode){		node.SetExpanded(!node.IsExpanded())})

	return ui
}

func (ui *UI) GetMIDIStream() *tview.TextView {
	return ui.MidiStream
}

func (ui *UI) GetMenu() *tview.TreeNode {
	return ui.Menu
}

func(ui *UI) AddDevice2Menu(dev midi.MidiDevice){//man, prod, sn string, path []int, port, bus int, s gousb.Speed){
	
	
	node := tview.NewTreeNode(" "+dev.Product).SetSelectable(true).SetColor(tcell.ColorRed)

	mnode := tview.NewTreeNode("Manufacturer: "+dev.Manufacturer).SetSelectable(false)
	pnode := tview.NewTreeNode("Product: "+dev.Product).SetSelectable(false)
	pidnode := tview.NewTreeNode("PID: 0x" + dev.PID.String()).SetSelectable(false)
	vidnode := tview.NewTreeNode("VID: 0x" + dev.VID.String()).SetSelectable(false)
	cnode := tview.NewTreeNode("Class: " + dev.Class).SetSelectable(false)
	scnode := tview.NewTreeNode("SubClass: " + dev.SubClass.String()).SetSelectable(false)
	protnode := tview.NewTreeNode("Protocol: " + dev.Protocol.String()).SetSelectable(false)
	snnode := tview.NewTreeNode("Serial Number: " +  dev.SerialNumber).SetSelectable(false)
	epnode := tview.NewTreeNode("IN Endpoint: "+ dev.EndpointIn.String()).SetSelectable(false)
	node.AddChild(mnode)
	node.AddChild(pnode)
	node.AddChild(vidnode)
	node.AddChild(pidnode)
	node.AddChild(cnode)
	node.AddChild(scnode)
	node.AddChild(protnode)
	node.AddChild(snnode)
	node.AddChild(epnode)

	node.Collapse()

	ui.Menu.AddChild(node)

	//ui.Menu.ScrollToEnd()
}


func CreateMidiStream() *tview.TextView {
	return tview.NewTextView().SetDynamicColors(true)
}
