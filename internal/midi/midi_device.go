package midi

import (
	"fmt"
	"github.com/google/gousb"
	"github.com/rivo/tview"
	"time"
)


// MIDIReader interface for reading MIDI messages and providing device info
type MIDIReader interface {
	ReadMIDI(midiStream *tview.TextView, app *tview.Application) error
	GetDeviceInfo() string
	GetDeviceDetails() [][2]string
}

// MidiDevice implements MIDIReader
type MidiDevice struct {
	Device        *gousb.Device
	Manufacturer  string
	Product       string
	VID           gousb.ID
	PID           gousb.ID
	EndpointIn    *gousb.InEndpoint
	MaxPacketSize int
	Class         string
	SubClass      string
	Protocol      string
	SerialNumber  string
	MaxPower      string
}


const (
	afterTouch    = 0xA
	controlChange = 0xB
	pitchBend     = 0xE
	noteOff       = 0x8
	noteOn        = 0x9
)

func getNotesList() []string {
	return []string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}
}


func (d *MidiDevice) ReadMIDI(midiStream *tview.TextView, app *tview.Application) error {
	interval := time.Duration(125000)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	buff := make([]byte, d.MaxPacketSize)

	for {
		select {
		case <-ticker.C:
			n, err := d.EndpointIn.Read(buff)
			if err != nil {
				return fmt.Errorf("error reading from %s: %v", d.Product, err)
			}

			data := buff[:n]
			formattedMessage := styleText(formatMessage(data, d), "white", false, true)

			app.QueueUpdateDraw(func() {
				fmt.Fprintln(midiStream, formattedMessage)
				midiStream.ScrollToEnd()
			})
		}
	}
}

func (d *MidiDevice) GetDeviceInfo() string {
	return fmt.Sprintf("MIDI Device: %s (%s) [%s:%s]", d.Manufacturer, d.Product, d.VID.String(), d.PID.String())
}

// GetDeviceDetails provides a detailed ordered slice of key-value pairs of the device's attributes
func (d *MidiDevice) GetDeviceDetails() [][2]string {
	return [][2]string{
		{"Manufacturer", d.Manufacturer},
		{"Product", d.Product},
		{"VID", "0x" + d.VID.String()},
		{"PID", "0x" + d.PID.String()},
		{"Class", d.Class},
		{"SubClass", d.SubClass},
		{"Protocol", d.Protocol},
		{"Serial Number", d.SerialNumber},
		{"Max Current", d.MaxPower},
		{"IN Endpoint", d.EndpointIn.String()},
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
	style := fmt.Sprintf("[%s]", color)
	if bold {
		style += "[::b]"
	}
	if underline {
		style += "[::u]"
	}
	return fmt.Sprintf("%s%s", style, text)
}

func getNoteAndOctave(n byte) (string, int) {
	notes := getNotesList()
	note := notes[n%12]     // calculate note inside octave
	octave := int(n/12) - 1 // calculates the octave
	return note, octave
}
