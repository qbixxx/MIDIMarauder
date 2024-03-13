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
const midiStream = "/---------------------------------------- MIDI STREAM ----------------------------------------/"

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
func scanForMIDIDevices(ctx *gousb.Context) [][]gousb.ID {

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

	ctx := gousb.NewContext()

	fmt.Println(Red + asciiTitle + clrReset)

	midiDevices := scanForMIDIDevices(ctx)
	if len(midiDevices) == 0 {
		fmt.Println(Red + "Exiting" + clrReset)
		os.Exit(0)
	}
	fmt.Println(midiStream)

	resultChan := make(chan string)
	var wg sync.WaitGroup

	// Launch a goroutine for each MIDI device

	for _, dev := range midiDevices {
		wg.Add(1)

		go readDevice(dev, ctx, resultChan, &wg)

	}

	// Wait for all goroutines to finish
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for data := range resultChan {
		// Aquí puedes realizar alguna acción con los datos recibidos del canal
		fmt.Println(data)
	}

	fmt.Println("Finished reading from all devices")

}

func readDevice(device []gousb.ID, ctx *gousb.Context, resultChan chan<- string, wg *sync.WaitGroup) {

	defer wg.Done()

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

					if !mdev.read(endpointDesc.MaxPacketSize, resultChan) {

						wg.Done()
					}

				}
			}
		}

	}

}

func (mdev *MIDIDEV) read(maxSize int, resultChan chan<- string) bool {

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
				resultChan <- fmt.Sprintf("[%s-%s] >>> CC:%d Value: %d", man, pr, data[2], data[3])
				//fmt.Sprintf(Bold+Red+"["+man+"-"+pr+"] >>> "+clrReset+(Green+"CC:%d"), data[2], "Value: %i", data[3], clrReset)

			case 8:

				note := getNotePosition(&data[2])
				//fmt.Println(Bold+Red+"["+man+"-"+pr+"] >>> "+clrReset+(Cyan+"Note OFF: "), list[note], " Velocity: ", data[3], clrReset)
				resultChan <- fmt.Sprintf("[%s-%s] >>> Note OFF: %s\tVelocity: %d", man, pr, list[note], data[3])

			case 9:
				note := getNotePosition(&data[2])
				//s, _ := fmt.Println(CyanBG+Bold+Red+"["+man+"-"+pr+"] >>> "+clrReset+CyanBG+Black+"Note ON:  ", list[note], " Velocity: ", data[3], clrReset)
				resultChan <- fmt.Sprintf("[%s-%s] >>> Note ON: %s\tVelocity: %d", man, pr, list[note], data[3])
				//fmt.Sprintf("%s",s)
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
