package main

import (
	"fmt"
	"github.com/google/gousb"
	"github.com/google/gousb/usbid"
	"log"
	"time"
	// tea "github.com/charmbracelet/bubbletea"
)

const asciiTitle = " _ _   .   _/  .  /|,/  _   _  _        _/  _   _\n" +
	"/ / / /  /_/  /  /  /  /_| /  /_| /_/ /_/  /_' / \n" +
	"                                                 "

const mBegin = "---------------------------------------- MIDI BEGIN ----------------------------------------"

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

func main() {

	fmt.Println(asciiTitle)

	readDevice()

}

func readDevice() {

	// Initialize a new Context.
	ctx := gousb.NewContext()
	defer ctx.Close()

	fmt.Println("Context: \n", ctx)

	// Open any device with a given VID/PID using a convenience function.
	dev, err := ctx.OpenDeviceWithVIDPID(0x2467, 0x2034)
	if err != nil {
		log.Fatalf("Could not open a device: %v", err)
	}
	defer dev.Close()

	dev.SetAutoDetach(true)

	config := dev.String()

	strDesc3, _ := dev.GetStringDescriptor(3)
	fmt.Println("strDesc3: ", strDesc3)

	strDesc2, _ := dev.GetStringDescriptor(2)
	fmt.Println("strDesc2: ", strDesc2)

	strDesc1, _ := dev.GetStringDescriptor(1)
	fmt.Println("strDesc1: ", strDesc1)

	strDesc0, _ := dev.GetStringDescriptor(0)
	fmt.Println("strDesc0: ", strDesc0)

	manu, _ := dev.Manufacturer()

	midi_b, _ := dev.Product()
	fmt.Println("Config ", config)
	fmt.Println("product: ", midi_b)
	fmt.Println("Manufacturer: ", manu)

	// Iterate through configurations
	for num := range dev.Desc.Configs {
		config, _ := dev.Config(num)

		defer config.Close()

		for _, desc := range config.Desc.Interfaces {
			intf, _ := config.Interface(desc.Number, 0)

			classy := usbid.Classify(dev.Desc)

			fmt.Println("Classify: ", classy)

			fmt.Println("Interface: ", intf)

			for _, endpointDesc := range intf.Setting.Endpoints {

				if endpointDesc.Direction == gousb.EndpointDirectionIn {

					endpoint, _ := intf.InEndpoint(endpointDesc.Number)

					fmt.Println(endpoint)

					fmt.Println("Endpoint poll interval: ", endpointDesc.PollInterval, "\nEndpoint max packet size", endpointDesc.MaxPacketSize, "Bytes")

					mdev := &MIDIDEV{
						context:  ctx,
						device:   dev,
						intf:     intf,
						endpoint: endpoint,
					}

					fmt.Println("\n", mdev.context, "\n", mdev.device, "\n", mdev.intf, "\n", mdev.endpoint)
					mdev.read(endpointDesc.PollInterval, endpointDesc.MaxPacketSize)

				}
			}
		}

	}

}

func (mdev *MIDIDEV) read(interval time.Duration, maxSize int) {

	interval = 125000 //hardcoded, idkw it appears to be 0ms according to the endpoint poll description
	fmt.Println("time duration: ", interval, "max size: ", maxSize)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	list := getNotesList()

	fmt.Println(mBegin)

	for {

		select {
		case <-ticker.C:
			buff := make([]byte, maxSize)
			n, err := mdev.endpoint.Read(buff)
			if err != nil {
				fmt.Println("Error", err)
			}

			data := buff[:n]

			switch data[0] {
			case 11, 14:
				fmt.Println("CC:", data[2], " Value: ", data[3])

			case 8:

				note := getNotePosition(&data[2])

				fmt.Println("Note OFF: ", list[note], " Velocity: ", data[3])
			case 9:
				note := getNotePosition(&data[2])

				fmt.Println("Note ON:  ", list[note], " Velocity: ", data[3])
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
