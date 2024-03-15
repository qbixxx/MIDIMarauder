package main

import (
	"github.com/rivo/tview"
	"time"
	"fmt"
)


const asciiTitle ="[cyan]  \\  | _)      | _)\n" +
	" |\\/ |  |   _` |  |                                 \n" +
	" |   |  |  (   |  |                                 \n" +
	"_|\\ _| _| \\__,_| _|                   [turquoise]|             \n" +
	" |\\/ |   _` |   __|  _` |  |   |   _` |   _ \\   __| \n" +
	" |   |  (   |  |    (   |  |   |  (   |   __/  |    \n" +
	"_|  _| \\__,_| _|   \\__,_| \\__,_| \\__,_| \\___| _| \n\n\n"

func main() {
//	newPrimitive := func(text string) tview.Primitive {
//		return tview.NewTextView().
//			SetTextAlign(tview.AlignCenter).
//			SetText(text)
//	}
//	menu := newPrimitive("Menu")
//	main := newPrimitive("Main content")
//	sideBar := newPrimitive("Side Bar")

	app := tview.NewApplication()


	midiStream := tview.NewTextView()
	midiStream.Box.SetBorder(true).SetTitle(" Midi Stream ")

	menu := tview.NewTextView()
	menu.Box.SetBorder(true).SetTitle(" Menu ") 
	menu.SetTextAlign(tview.AlignLeft).SetDynamicColors(true).
	SetText(asciiTitle)



	grid := tview.NewGrid().
		SetColumns(-4,54).
		SetRows(-2,2).
		SetBorders(false).
		//AddItem(newPrimitive("Header\nHeader2"), 0, 0, 1, 1, 0, 0, true).
		AddItem(midiStream,0,0,1,1,0,0,true).
		AddItem(menu,0,1,1,1,0,0,true)
		//AddItem(newPrimitive("Footer"), 1, 0, 1, 1, 0, 0, true)



		go func(){
	
			app.QueueUpdateDraw(func(){
				time.Sleep(2000*time.Millisecond)
				fmt.Fprintf(menu, "Nektar Impakt")
				
			})
			
		}()

	
		go func(){
			
			for i := 0; i < 100; i++{
				app.QueueUpdateDraw(func(){
					fmt.Fprintf(midiStream, "----------- MIDI MESSAGE:___%d---------\n", i)
					midiStream.ScrollToEnd()			
				})
				time.Sleep(10*time.Millisecond)
			}


			

			
		}()

	
	if err := app.SetRoot(grid, true).SetFocus(grid).Run(); err != nil {
		panic(err)
	}

	


}