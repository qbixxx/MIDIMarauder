package main

import (
	"fmt"
	"github.com/google/gousb"
	"github.com/rivo/tview"
	
//	"os"
	 "sync"
	 "time"
)

const asciiTitle = "[cyan]  \\  | _)      | _)\n" +
	" |\\/ |  |   _` |  |                                 \n" +
	" |   |  |  (   |  |                                 \n" +
	"_|\\ _| _| \\__,_| _|                   [turquoise]|             \n" +
	" |\\/ |   _` |   __|  _` |  |   |   _` |   _ \\   __| \n" +
	" |   |  (   |  |    (   |  |   |  (   |   __/  |    \n" +
	"_|  _| \\__,_| _|   \\__,_| \\__,_| \\__,_| \\___| _| \n\n\n"

//const midiStream = "/---------------------------------------- MIDI STREAM ----------------------------------------/"

const (
	clrReset  = "\033[0m"
	Black     = "\033[30m"
	Red       = "\033[31m"
	Green     = "\033[32m"
	Yellow    = "\033[33m"
	Blue      = "\033[34m"
	Magenta   = "\033[35m"
	Cyan      = "\033[36m"
	White     = "\033[37m"
	Purple    = "\033[35m"
	YellowBG  = "\033[43m"
	BlueBG    = "\033[44m"
	MagentaBG = "\033[45m"
	CyanBG    = "\033[46m"
	WhiteBG   = "\033[47m"
	BlackBG   = "\033[40m"
	RedBG     = "\033[41m"
	GreenBG   = "\033[42m"
	Bold      = "\033[1m"
	// Add more colors if needed
)

type midiDev struct {
	// Fields for interacting with the USB connection
	device *gousb.Device
	man    string
	prod   string
	vid		gousb.ID
	pid	    gousb.ID
	//context  *gousb.Context

	//intf     *gousb.Interface
	endpoint *gousb.InEndpoint
}

func getNotesList() []string {
	return []string{"C ", "C#", "D ", "D#", "E ", "F ", "F#", "G ", "G#", "A ", "A#", "B "}
}

func scanForMIDIDevices(ctx *gousb.Context, menu *tview.TextView, app *tview.Application) []*midiDev { //[]gousb.ID {

	//defer ctx.Close()

	//store every device connected
	devices, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		return true

	})

	if err != nil {
		fmt.Printf("Failed to open devices: %v\n", err)
		return nil
	}
	defer func() {
		for _, dev := range devices {
			dev.Close()
		}
	}()

	count := 0
	midiDevices := make([]*midiDev, 0)
	// Iterate over the found devices.
	for _, dev := range devices {

		for _, cfg := range dev.Desc.Configs {

			for _, intf := range cfg.Interfaces {

				for _, ifSetting := range intf.AltSettings {

					for _, endpoint := range ifSetting.Endpoints {

						// Check if the device is a MIDI device (class 1, subclass 3).
						if ifSetting.Class == gousb.ClassAudio && ifSetting.SubClass == 3 && endpoint.Direction == gousb.EndpointDirectionIn {
							man, _ := dev.Manufacturer()
							prod, _ := dev.Product()
							vid := dev.Desc.Vendor
							pid := dev.Desc.Product
							//fmt.Printf("MIDI Device [%d]: %s - %s\n", count, man, pr)
							count++
							d := &midiDev{
								device: dev,
								man:    man,
								prod:   prod,
								vid:	vid,
								pid:	pid,
								//context:  ctx,
								//intf:     intf,
								//endpoint: endpoint,
							}
							//[]gousb.ID{dev.Desc.Product, dev.Desc.Vendor}
							midiDevices = append(midiDevices, d)
						}
					}
				}
			}
		}
	}
//s := fmt.Sprintf("[red]MIDI Devices Found: %d", len(midiDevices))

	var devMsg string

	if len(midiDevices) == 0{
		devMsg = fmt.Sprintf("[red]MIDI Devices Found: %d", len(midiDevices))
	}else{
		devMsg = fmt.Sprintf("[green]MIDI Devices Found: %d", len(midiDevices))
	}
	

	fmt.Fprintln(menu, devMsg)

	if len(midiDevices) > 0{
		for i, mdev := range midiDevices{
			//app.QueueUpdateDraw(func() {
				//fmt.Println(i)
				strDevice := fmt.Sprintf("Device [%d]: %s - %s",i, mdev.man, mdev.prod)
				fmt.Fprintln(menu, "[white]"+strDevice)
				menu.ScrollToEnd()
			//})
		}
	}
	

	//fmt.Println(Green+"MIDI devices found:", count, clrReset)

	
	return midiDevices
}

var ctx *gousb.Context

func main() {

	ctx := gousb.NewContext()	

	app := tview.NewApplication()
	midiStream := tview.NewTextView().SetDynamicColors(true)
	midiStream.Box.SetBorder(true).SetTitle(" Midi Stream ")

	menu := tview.NewTextView()
	menu.Box.SetBorder(true).SetTitle(" Menu ")
	menu.SetTextAlign(tview.AlignLeft).SetDynamicColors(true).
		SetText(asciiTitle)

	grid := tview.NewGrid().
		SetColumns(-4, 54).
		SetRows(-2, 2).
		SetBorders(false).
		AddItem(midiStream, 0, 0, 1, 1, 0, 0, true).
		AddItem(menu, 0, 1, 1, 1, 0, 0, true)

	
	//fmt.Println(Red + asciiTitle + clrReset)

	midiDevices := scanForMIDIDevices(ctx, menu, app)
	//fmt.Println(midiDevices[0].man,midiDevices[0].prod)
	//scanForMIDIDevices(ctx, menu, app)
	
	//if len(midiDevices) == 0 {
	//	fmt.Println(Red + "Exiting" + clrReset)
	//	os.Exit(0)
	//}
	//fmt.Println(midiStream)

	
		var wg sync.WaitGroup

		// Launch a goroutine for each MIDI device
		for _, dev := range midiDevices {
			wg.Add(1)
			go func(device midiDev) {
				defer wg.Done()

				readDevice(device, ctx, midiStream, &wg, app)
			}(*dev)
		}
	
	if err := app.SetRoot(grid, true).SetFocus(grid).Run(); err != nil {
		panic(err)
	}

	wg.Wait() // Ensure all goroutines finish before exiting
}


func readDevice(mdev midiDev, ctx *gousb.Context, midiStream *tview.TextView, wg *sync.WaitGroup, app *tview.Application) {

	//defer wg.Done()

	// Open each MIDI device with a given VID/PID
//	dev, err := ctx.OpenDeviceWithVIDPID(mdev.device.Desc.Vendor, mdev.device.Desc.Product)

	dev, err := ctx.OpenDeviceWithVIDPID(mdev.vid, mdev.pid)

	if err != nil {
		fmt.Println("Error >>>", err)
	}
	//if dev == nil{
	//	os.Exit(1)
	//}
	//fmt.Println("Error >>>", err)
	//fmt.Println("dev >>>", dev)
	//defer dev.Close()

	dev.SetAutoDetach(true)

	// Iterate through devices and endpoints
	for num := range dev.Desc.Configs {
		config, _ := dev.Config(num)

		defer config.Close()

		for _, desc := range config.Desc.Interfaces {
			intf, _ := config.Interface(desc.Number, 0)

			for _, endpointDesc := range intf.Setting.Endpoints {

				if endpointDesc.Direction == gousb.EndpointDirectionIn {

					endpoint, _ := intf.InEndpoint(endpointDesc.Number)

					mdev.endpoint = endpoint

					if !mdev.read(endpointDesc.MaxPacketSize, midiStream, app) {
						wg.Done()
					}

				}
			}
		}

	}

}



func (mdev *midiDev) read(maxSize int, midiStream *tview.TextView, app *tview.Application) bool {

	interval := time.Duration(12500000) //hardcoded, idkw it appears to be 0ms according to the endpoint poll description
	//fmt.Println(interval)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	list := getNotesList()

	//man, _ := mdev.device.Manufacturer()
	//pr, _ := mdev.device.Product()

	//fmt.Println(man, pr, "<<<")

	buff := make([]byte, maxSize)
	var note byte
	for {
		var formattedMessage string
		select {
		case <-ticker.C:
			//buff := make([]byte, maxSize)
			n, err := mdev.endpoint.Read(buff) 	//////////////
			if err != nil {
				fmt.Printf("Error!!!!!!!!!!!!!!!!!!!: %s: %s - %s\n", err, mdev.man, mdev.prod)
				return false
			}

			data := buff[:n]

			
			switch data[0] {

			case 10:
				note = getNotePosition(&data[2])
				//formattedMessage = fmt.Sprintf("mps: %d [%s-%s] >>> After touch: %s\tVelocity: %d",maxSize, mdev.man, mdev.prod, list[note], data[3])
				formattedMessage = fmt.Sprintf("[%s-%s]\t| After touch: %s|\tVelocity: %d\t\t|\tMax packet size: %d\t|\tRAW DATA: % X", mdev.man, mdev.prod, list[note], data[3], maxSize, data)
 

			case 11, 14:
				
				formattedMessage = styleText(fmt.Sprintf("[%s-%s]\t| CC:%d\t\t\t|\tValue: %d \t\t|\tMax packet size: %d\t|\tRAW DATA: % X", mdev.man, mdev.prod ,data[2], data[3], maxSize, data), "red", "black", true, true)
				//formattedMessage = fmt.Sprintf("CC_ x%",data)
				//fmt.Println(formattedMessage)
			case 8:
				note = getNotePosition(&data[2])
				formattedMessage = fmt.Sprintf("[%s-%s]\t| Note OFF: %s \t|\tVelocity: %d \t\t|\tMax packet size: %d\t|\tRAW DATA: % X", mdev.man, mdev.prod, list[note], data[3], maxSize, data)
			case 9:
				note = getNotePosition(&data[2])
				formattedMessage = styleText(fmt.Sprintf("[%s-%s]\t| Note ON: %s \t|\tVelocity: %d \t\t|\tMax packet size: %d\t|\tRAW DATA: % X", mdev.man, mdev.prod, list[note], data[3] ,maxSize, data), "white", "green", false, false)
			default:
				// Handle other MIDI message types (optional)
				formattedMessage = fmt.Sprintf("[%s-%s]\t| UNKNOWN MESSAGE \t|\tRAW DATA: %X", mdev.man, mdev.prod, data)

				
			}

			// Update the MIDI Stream text view in a thread-safe manner
			app.QueueUpdateDraw(func() {
				fmt.Fprintln(midiStream, formattedMessage)
				midiStream.ScrollToEnd()
			})
		}
	}
}

func styleText(text, color, background string, bold, underline bool) string {
    style := fmt.Sprintf("[%s:%s]", color, background)
    if bold {
        style += "[::b]"
    }
    if underline {
        style += "[::u]"
    }
    // Asegúrate de restablecer el color de fondo al final del texto
    return fmt.Sprintf("%s%s[white:black]", style, text)
}




func getNotePosition(n *byte) byte {

	for *n > 11 {
		*n = *n - 12
	}
	return *n
}
