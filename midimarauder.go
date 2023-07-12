package main

import (
	"fmt"
	"github.com/google/gousb"
	"os"
	"sync"
	"time"
)

const asciiTitle = "  \\  | _)      | _)\n" +
	" |\\/ |  |   _` |  |                                 \n" +
	" |   |  |  (   |  |                                 \n" +
	"_|\\ _| _| \\__,_| _|                   |             \n" +
	" |\\/ |   _` |   __|  _` |  |   |   _` |   _ \\   __| \n" +
	" |   |  (   |  |    (   |  |   |  (   |   __/  |    \n" +
	"_|  _| \\__,_| _|   \\__,_| \\__,_| \\__,_| \\___| _| \n"
const mBegin = "/---------------------------------------- MIDI STREAM ----------------------------------------/"

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

type MIDIDEV struct {
	// Fields for interacting with the USB connection
	context  *gousb.Context
	device   *gousb.Device
	intf     *gousb.Interface
	endpoint *gousb.InEndpoint
}

func getNotesList() []string {
	return []string{"C ", "C#", "D ", "D#", "E ", "F ", "F#", "G ", "G#", "A ", "A#", "B "}
}
func scanForMIDIDevices() [][]gousb.ID {

	ctx := gousb.NewContext()
	defer ctx.Close()

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
	out := make([][]gousb.ID, 0)

	// Iterate over the found devices.
	for _, dev := range devices {

		for _, cfg := range dev.Desc.Configs {

			for _, intf := range cfg.Interfaces {

				for _, ifSetting := range intf.AltSettings {

					for _, endpoint := range ifSetting.Endpoints {

						// Check if the device is a MIDI device (class 1, subclass 3).
						if ifSetting.Class == gousb.ClassAudio && ifSetting.SubClass == 3 && endpoint.Direction == gousb.EndpointDirectionIn {
							man, _ := dev.Manufacturer()
							pr, _ := dev.Product()
							fmt.Printf("MIDI Device [%d]: %s - %s\n", count, man, pr)
							count++
							s := []gousb.ID{dev.Desc.Product, dev.Desc.Vendor}
							out = append(out, s)
						}
					}
				}
			}
		}
	}

	fmt.Println(Green+"MIDI devices found:", count, clrReset)
	return out
}

func main() {

	fmt.Println(Red + asciiTitle + clrReset)

	midiDevices := scanForMIDIDevices()
	if len(midiDevices) == 0 {
		fmt.Println(Red + "Exiting" + clrReset)
		os.Exit(0)
	}
	fmt.Println(mBegin)
	var wg sync.WaitGroup
	resultChan := make(chan []byte)

	// Launch a goroutine for each MIDI device
	d := len(midiDevices)
	for _, device := range midiDevices {
		wg.Add(1)
		go func(dev []gousb.ID) {
			defer wg.Done()
			readDevice(dev, resultChan)
			d--
			if d == 0 {
				fmt.Println(Bold + Red + "Fatal: no devices left." + clrReset)
				os.Exit(0)
			}
		}(device)
	}

	// Wait for all goroutines to finish
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for resultChan != nil {
	}

	fmt.Println("Finished reading from all devices")

}

func readDevice(device []gousb.ID, resultChan chan<- []byte) {

	// Initialize a new Context.
	ctx := gousb.NewContext()
	defer ctx.Close()

	// Open each MIDI device with a given VID/PID
	dev, err := ctx.OpenDeviceWithVIDPID(device[1], device[0])
	if err != nil {
		fmt.Println("Error", err)
	}
	defer dev.Close()

	dev.SetAutoDetach(true)

	// Iterate through configurations
	for num := range dev.Desc.Configs {
		config, _ := dev.Config(num)

		defer config.Close()

		for _, desc := range config.Desc.Interfaces {
			intf, _ := config.Interface(desc.Number, 0)

			for _, endpointDesc := range intf.Setting.Endpoints {

				if endpointDesc.Direction == gousb.EndpointDirectionIn {

					endpoint, _ := intf.InEndpoint(endpointDesc.Number)

					mdev := &MIDIDEV{

						context:  ctx,
						device:   dev,
						intf:     intf,
						endpoint: endpoint,
					}

					if !mdev.read(endpointDesc.MaxPacketSize) {

						return
					}

				}
			}
		}

	}

}

func (mdev *MIDIDEV) read(maxSize int) bool {

	interval := time.Duration(1250000) //hardcoded, idkw it appears to be 0ms according to the endpoint poll description

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	list := getNotesList()

	man, _ := mdev.device.Manufacturer()
	pr, _ := mdev.device.Product()

	for {

		select {
		case <-ticker.C:
			buff := make([]byte, maxSize)
			n, err := mdev.endpoint.Read(buff)
			if err != nil {
				fmt.Printf("Error: %s: %s - %s\n", err, man, pr)
				return false
			}

			data := buff[:n]

			switch data[0] {
			case 11, 14:

				fmt.Println(Bold+Red+"["+man+"-"+pr+"] >>> "+clrReset+(Green+"CC:"), data[2], "Value: ", data[3], clrReset)

			case 8:

				note := getNotePosition(&data[2])
				fmt.Println(Bold+Red+"["+man+"-"+pr+"] >>> "+clrReset+(Cyan+"Note OFF: "), list[note], " Velocity: ", data[3], clrReset)

			case 9:
				note := getNotePosition(&data[2])

				fmt.Println(CyanBG+Bold+Red+"["+man+"-"+pr+"] >>> "+clrReset+CyanBG+Black+"Note ON:  ", list[note], " Velocity: ", data[3], clrReset)
			}

		}

	}

}

func getNotePosition(n *byte) byte {

	for *n > 11 {
		*n = *n - 12
	}
	return *n
}
