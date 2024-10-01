package main

import (
	"github.com/rivo/tview"
	"log"
	"midiMarauder/internal/ui"
	"midiMarauder/internal/usb"
	"midiMarauder/internal/midi"
	"sync"
)

func main() {
	// Initialize USB Device Manager
	deviceManager := usb.NewUSBMIDIDeviceManager()
	defer deviceManager.Close()

	uiManager := ui.SetupUI()
	app := tview.NewApplication().EnableMouse(true)

	// Scan for MIDI devices
	devices, err := deviceManager.ScanDevices()
	if err != nil {
		log.Fatalf("Failed to scan for MIDI devices: %v", err)
	}

	// Add devices to UI. (Has to be implemented to use concurrency)
	for _, dev := range devices {
		uiManager.AddDevice2Menu(dev)
	}

	// Start reading every midi device simultaneously
	var wg sync.WaitGroup
	for _, device := range devices {
		wg.Add(1)
		go func(dev midi.MIDIReader) {
			defer wg.Done()
			err := dev.ReadMIDI(uiManager.GetMIDIStream(), app)
			if err != nil {
				log.Printf("Error reading MIDI: %v", err)
			}
		}(device)
	}

	// Run the application
	if err := app.SetRoot(uiManager.Root, true).SetFocus(uiManager.Tree).Run(); err != nil {
		log.Fatalf("Failed to run application: %v", err)
	}
}
