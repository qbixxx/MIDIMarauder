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
	Class			string
	SubClass		gousb.Class
	Protocol		gousb.Protocol
	Speed			gousb.Speed
	MaxPacketSize	int
	DeviceConfig	string
	SerialNumber	string
	Color			string
}

const (
    afterTouch   = 0xA
    controlChange = 0xB
    pitchBend    = 0xE
    noteOff      = 0x8
    noteOn       = 0x9
)

func getNotesList() []string {
	return []string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}
}
func (d *MidiDevice) GetProductInfo() (string, string, gousb.ID, gousb.ID, string, string, string, string, string){
	return d.Manufacturer, d.Product, d.VID, d.PID, d.SerialNumber, d.Class, d.SubClass.String(), d.Protocol.String(), d.Speed.String()
}

func (d *MidiDevice) Read(midiStream *tview.TextView, app *tview.Application){
	interval := time.Duration(125000)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	buff := make([]byte, d.MaxPacketSize)

	for {
		select {
		case <-ticker.C:
			n, err := d.EndpointIn.Read(buff)
			if err != nil {
				fmt.Printf("Error: %s: %s - %s\n", err, d.Manufacturer, d.Product) 
			}

			data := buff[:n]
			formattedMessage := styleText(formatMessage(data, d), "lightpink", true, true)

			app.QueueUpdateDraw(func() {
				fmt.Fprintln(midiStream, formattedMessage)
				midiStream.ScrollToEnd()
			})
		}
	}
}

 
func formatMessage(data []byte, d *MidiDevice) string {
	
	switch data[0] {
	
	case afterTouch:
		note, octave := getNoteAndOctave(data[2])
		return fmt.Sprintf("[%s-%s]\t| After touch: %s%d|\tVelocity: %d\t\t|\tPacket Size %d\t|\tRAW DATA: % X", d.Manufacturer, d.Product, note, octave, data[3], len(data), data)
	
	case controlChange:
		return fmt.Sprintf("[%s-%s]\t| CC:%d\t\t\t|\tValue: %d \t\t|\tPacket Size %d\t|\tRAW DATA: % X", d.Manufacturer, d.Product, data[2], data[3], len(data), data)
	
	case pitchBend:
		return fmt.Sprintf("[%s-%s]\t| Pitch Bend\t\t|\tValue: %d \t\t|\tPacket Size %d\t|\tRAW DATA: % X", d.Manufacturer, d.Product, data[3], len(data), data)

	case noteOff:
		note, octave := getNoteAndOctave(data[2])
		return fmt.Sprintf("[%s-%s]\t| Note OFF: %s%d \t\t|\tVelocity: %d \t|\tPacket Size %d\t|\tRAW DATA: % X", d.Manufacturer, d.Product, note, octave, data[3], len(data), data)
	
	case noteOn:
		note, octave := getNoteAndOctave(data[2])
		return fmt.Sprintf("[%s-%s]\t| Note ON: %s%d \t\t|\tVelocity: %d \t|\tPacket Size %d\t|\tRAW DATA: % X", d.Manufacturer, d.Product, note, octave, data[3], len(data), data)
	
	default:
		return fmt.Sprintf("[%s-%s]\t| UNKNOWN MESSAGE \t|\tRAW DATA: % X", d.Manufacturer, d.Product, data)
	}
}

func styleText(text, color string, bold, underline bool) string {
	style := fmt.Sprintf("[%s:]", color)
	if bold {
		style += "[::b]"
	}
	if underline {
		style += "[::u]"
	}
	// Asegúrate de restablecer el color de fondo al final del texto
	return fmt.Sprintf("%s%s[-:-:-:-]", style, text)
}

func getNoteAndOctave(n byte) (string, int) {
	notes := getNotesList()
	note := notes[n % 12]      // calculate note inside octave
	octave := int(n / 12) - 1  // calculates the octave
	return note, octave
}