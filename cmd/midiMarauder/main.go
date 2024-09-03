package main

import (
	"log"
	"midiMarauder/internal/midi"
	"midiMarauder/internal/ui"
	"midiMarauder/internal/usb"
	"github.com/google/gousb"
	"github.com/rivo/tview"
	"sync"
)

func main() {
	ctx := gousb.NewContext()
	defer ctx.Close()

	uiManager := ui.SetupUI()
	app := tview.NewApplication()

	// Escanear dispositivos MIDI
	devices, err := usb.ScanForMIDIDevices(ctx)
	if err != nil {
		log.Fatalf("Failed to scan for MIDI devices: %v", err)
	}

	// Iniciar lectura de dispositivos MIDI
	var wg sync.WaitGroup
	for _, device := range devices {
		wg.Add(1)
		go func(dev midi.MidiDevice) {
			defer wg.Done()
			dev.Read(ui.CreateMidiStream(), app)
		}(device)
	}

	if err := app.SetRoot(uiManager, true).SetFocus(uiManager).Run(); err != nil {
		log.Fatalf("Failed to run application: %v", err)
	}

	wg.Wait() // Esperar a que todas las lecturas terminen antes de salir
}
