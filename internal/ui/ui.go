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
	//menu.SetTextAlign(tview.AlignLeft).SetDynamicColors(true)	

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
		SetColumns(-4, 54).
		SetRows(-2, 1).
		SetBorders(false).
		AddItem(midiStream, 0, 0, 1, 1, 0, 0, true).
		AddItem(gridMenu, 0, 1, 1, 1, 0, 0, true)

	
	ui.Root = rootGrid
	ui.MidiStream = midiStream
	ui.Menu = rootTree
	ui.Tree = tree
	ui.Tree.SetSelectedFunc(func(node *tview.TreeNode){
		node.SetExpanded(!node.IsExpanded())})

	return ui
}

func (ui *UI) GetMIDIStream() *tview.TextView {
	return ui.MidiStream
}

func (ui *UI) GetMenu() *tview.TreeNode {
	return ui.Menu
}

func(ui *UI) AddDevice2Menu(dev midi.MidiDevice){//man, prod, sn string, path []int, port, bus int, s gousb.Speed){
	
	
	node := tview.NewTreeNode(" " + dev.Manufacturer +" "+ dev.Product).SetSelectable(true).SetColor(tcell.ColorRed)

		mnode := tview.NewTreeNode(dev.Manufacturer).SetSelectable(false)
		pnode := tview.NewTreeNode(dev.Product).SetSelectable(false)
		vidnode := tview.NewTreeNode(dev.VID.String()).SetSelectable(false)
		pidnode := tview.NewTreeNode(dev.PID.String()).SetSelectable(false)
		cnode := tview.NewTreeNode(dev.Class.String()).SetSelectable(false)
		scnode := tview.NewTreeNode(dev.SubClass.String()).SetSelectable(false)
		protnode := tview.NewTreeNode(dev.Protocol.String()).SetSelectable(false)
		snnode := tview.NewTreeNode(dev.SerialNumber).SetSelectable(false)


		node.AddChild(mnode)
		node.AddChild(pnode)
		node.AddChild(vidnode)
		node.AddChild(pidnode)
		node.AddChild(cnode)
		node.AddChild(scnode)
		node.AddChild(protnode)
		node.AddChild(snnode)
	node.Collapse()

	ui.Menu.AddChild(node)

	//ui.Menu.ScrollToEnd()
}


func CreateMidiStream() *tview.TextView {
	return tview.NewTextView().SetDynamicColors(true)
}
