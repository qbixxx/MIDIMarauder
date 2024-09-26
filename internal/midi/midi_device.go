package midi

import (
	"fmt"
	"github.com/google/gousb"
	"github.com/rivo/tview"
	"time"
)

type MidiDevice struct {
	Device			*gousb.Device
	Manufacturer	string
	Product			string
	VID				gousb.ID
	PID				gousb.ID
	EndpointIn		*gousb.InEndpoint
	Port			int
	Class			gousb.Class
	SubClass		gousb.Class
	Protocol		gousb.Protocol
	Speed			gousb.Speed
	MaxPacketSize	int
	DeviceConfig	string
	SerialNumber	string
}


func getNotesList() []string {
	return []string{"C ", "C#", "D ", "D#", "E ", "F ", "F#", "G ", "G#", "A ", "A#", "B "}
}
func (d *MidiDevice) GetProductInfo() (string, string, gousb.ID, gousb.ID, string, string, string, string, string){
	return d.Manufacturer, d.Product, d.VID, d.PID, d.SerialNumber, d.Class.String(), d.SubClass.String(), d.Protocol.String(), d.Speed.String()
}

func (d *MidiDevice) Read(midiStream *tview.TextView, app *tview.Application) bool {
	interval := time.Duration(1250000)
	ticker := time.NewTicker(interval)
	//ticker := time.NewTicker(d.PollInterval)
	defer ticker.Stop()

	buff := make([]byte, d.MaxPacketSize)

	for {
		select {
		case <-ticker.C:
			n, err := d.EndpointIn.Read(buff)
			if err != nil {
				fmt.Printf("Error: %s: %s - %s\n", err, d.Manufacturer, d.Product)
				return false
			}

			data := buff[:n]
			formattedMessage := formatMessage(data, d)

			app.QueueUpdateDraw(func() {
				fmt.Fprintln(midiStream, formattedMessage)
				midiStream.ScrollToEnd()
			})
		}
	}
}

 
func formatMessage(data []byte, d *MidiDevice) string {
	switch data[0] {
	case 0xA: // aftertouch
		note, octave := getNoteAndOctave(data[2])
		return fmt.Sprintf("[%s-%s]\t| After touch: %s%d|\tVelocity: %d\t\t|\tPacket Size %d\t|\tRAW DATA: % X", d.Manufacturer, d.Product, note, octave, data[3], len(data), data)
	case 0xB:// CC, 0xE:
		return fmt.Sprintf("[%s-%s]\t| CC:%d\t\t\t|\tValue: %d \t\t|\tPacket Size %d\t|\tRAW DATA: % X", d.Manufacturer, d.Product, data[2], data[3], len(data), data)
	case 0xE:
		return fmt.Sprintf("[%s-%s]\t| Pitch Bend:%d\t\t\t|\tValue: %d \t\t|\tPacket Size %d\t|\tRAW DATA: % X", d.Manufacturer, d.Product, data[2], data[3], len(data), data)
	case 0x8:
		note, octave := getNoteAndOctave(data[2])
		return fmt.Sprintf("[%s-%s]\t| Note OFF: %s%d \t|\tVelocity: %d \t\t|\tPacket Size %d\t|\tRAW DATA: % X", d.Manufacturer, d.Product, note, octave, data[3], len(data), data)
	case 0x9:
		note, octave := getNoteAndOctave(data[2])
		return fmt.Sprintf("[%s-%s]\t| Note ON: %s%d \t|\tVelocity: %d \t\t|\tPacket Size %d\t|\tRAW DATA: % X", d.Manufacturer, d.Product, note, octave, data[3], len(data), data)
	default:
		return fmt.Sprintf("[%s-%s]\t| UNKNOWN MESSAGE \t|\tRAW DATA: % X", d.Manufacturer, d.Product, data)
	}
}

func getNoteAndOctave(n byte) (string, int) {
	notes := getNotesList()
	note := notes[n % 12]      // Calcula la nota (dentro de una octava)
	octave := int(n / 12) - 1  // Calcula la octava, ajustada para MIDI
	return note, octave
}