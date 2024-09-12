package midi

import (
	"fmt"
	"github.com/google/gousb"
	"github.com/rivo/tview"
	"time"
)

type MidiDevice struct {
	Device       *gousb.Device
	Manufacturer string
	Product      string
	Vid          gousb.ID
	Pid          gousb.ID
	Endpoint     *gousb.InEndpoint
	MaxPacketSize int
}

func (d *MidiDevice) Read(midiStream *tview.TextView, app *tview.Application) bool {
	interval := time.Duration(12500000)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	list := getNotesList()
	buff := make([]byte, d.MaxPacketSize)

	for {
		select {
		case <-ticker.C:
			n, err := d.Endpoint.Read(buff)
			if err != nil {
				fmt.Printf("Error: %s: %s - %s\n", err, d.Manufacturer, d.Product)
				return false
			}

			data := buff[:n]
			formattedMessage := formatMessage(data, list, d)
			fmt.Println(formattedMessage)

			app.QueueUpdateDraw(func() {
				fmt.Fprintln(midiStream, formattedMessage)
				midiStream.ScrollToEnd()
			})
		}
	}
}

func getNotesList() []string {
	return []string{"C ", "C#", "D ", "D#", "E ", "F ", "F#", "G ", "G#", "A ", "A#", "B "}
}

func formatMessage(data []byte, list []string, d *MidiDevice) string {
	var note byte
	switch data[0] {
	case 10:
		note = getNotePosition(&data[2])
		fmt.Printf("[%s-%s]\t| After touch: %s|\tVelocity: %d\t\t|\tMax packet size: %d\t|\tRAW DATA: % X", d.Manufacturer, d.Product, list[note], data[3], len(data), data)
		return fmt.Sprintf("[%s-%s]\t| After touch: %s|\tVelocity: %d\t\t|\tMax packet size: %d\t|\tRAW DATA: % X", d.Manufacturer, d.Product, list[note], data[3], len(data), data)
	case 11, 14:
		fmt.Printf("[%s-%s]\t| CC:%d\t\t\t|\tValue: %d \t\t|\tMax packet size: %d\t|\tRAW DATA: % X", d.Manufacturer, d.Product, data[2], data[3], len(data), data)
		return fmt.Sprintf("[%s-%s]\t| CC:%d\t\t\t|\tValue: %d \t\t|\tMax packet size: %d\t|\tRAW DATA: % X", d.Manufacturer, d.Product, data[2], data[3], len(data), data)
	case 8:
		note = getNotePosition(&data[2])
		fmt.Printf("[%s-%s]\t| Note OFF: %s \t|\tVelocity: %d \t\t|\tMax packet size: %d\t|\tRAW DATA: % X", d.Manufacturer, d.Product, list[note], data[3], len(data), data)
		return fmt.Sprintf("[%s-%s]\t| Note OFF: %s \t|\tVelocity: %d \t\t|\tMax packet size: %d\t|\tRAW DATA: % X", d.Manufacturer, d.Product, list[note], data[3], len(data), data)
	case 9:
		note = getNotePosition(&data[2])
		return fmt.Sprintf("[%s-%s]\t| Note ON: %s \t|\tVelocity: %d \t\t|\tMax packet size: %d\t|\tRAW DATA: % X", d.Manufacturer, d.Product, list[note], data[3], len(data), data)
	default:
		return fmt.Sprintf("[%s-%s]\t| UNKNOWN MESSAGE \t|\tRAW DATA: %X", d.Manufacturer, d.Product, data)
	}
}

func getNotePosition(n *byte) byte {
	for *n > 11 {
		*n = *n - 12
	}
	return *n
}
