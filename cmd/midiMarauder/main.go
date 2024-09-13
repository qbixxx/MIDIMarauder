package main

import (
	//"fmt"
	"github.com/google/gousb"
	"github.com/rivo/tview"
	"log"
	"midiMarauder/internal/midi"
	"midiMarauder/internal/ui"
	"midiMarauder/internal/usb"
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
	

	for _, dev := range devices{
		man, prod, _, _ := dev.GetProductInfo()
		uiManager.AddDevice2Menu(man, prod)
	}

	// Iniciar lectura de dispositivos MIDI
	var wg sync.WaitGroup
	for _, device := range devices {
		wg.Add(1)
		go func(dev *midi.MidiDevice) {
			defer wg.Done()
			dev.Read(uiManager.GetMIDIStream(), app)
		}(device)
	}

	// Correr la aplicación
	if err := app.SetRoot(uiManager.Root, true).SetFocus(uiManager.Root).Run(); err != nil {
		log.Fatalf("Failed to run application: %v", err)
	}

	wg.Wait() // Esperar a que todas las lecturas terminen antes de salir
}
