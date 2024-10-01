package usb

import (
	//"fmt"
	"github.com/google/gousb"
	"midiMarauder/internal/midi"
	//"strconv"
)

// DeviceManager interface for managing devices
type DeviceManager interface {
	ScanDevices() ([]midi.MIDIReader, error)
	Close() error
}

// USBMIDIDeviceManager implements DeviceManager for USB MIDI devices
type USBMIDIDeviceManager struct {
	ctx *gousb.Context
}

func NewUSBMIDIDeviceManager() *USBMIDIDeviceManager {
	return &USBMIDIDeviceManager{ctx: gousb.NewContext()}
}

// ScanDevices scans for USB MIDI devices and returns a list of MIDIReader
func (m *USBMIDIDeviceManager) ScanDevices() ([]midi.MIDIReader, error) {
	devices, err := m.ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		// Logic for identifying MIDI devices
		return true
	})

	if err != nil {
		return nil, err
	}

	defer func() {
		for _, dev := range devices {
			dev.Close()
		}
	}()

	var midiDevices []midi.MIDIReader
	for _, dev := range devices {
		for num := range dev.Desc.Configs {
			config, _ := dev.Config(num)

			for _, intfDesc := range config.Desc.Interfaces {
				interFace, err := config.Interface(intfDesc.Number, 0)
				if err != nil {
					continue
				}

				for _, interFaceSetting := range intfDesc.AltSettings {
					for _, endpointDesc := range interFaceSetting.Endpoints {
						if interFaceSetting.Class == gousb.ClassAudio && interFaceSetting.SubClass == 3 && endpointDesc.Direction == gousb.EndpointDirectionIn {

							err = dev.SetAutoDetach(true)
							if err != nil {
								continue
							}

							endpoint, err := interFace.InEndpoint(endpointDesc.Number)
							if err != nil {
								continue
							}

							man, _ := dev.Manufacturer()
							prod, _ := dev.Product()
							vid := dev.Desc.Vendor
							pid := dev.Desc.Product
							mpSize := endpointDesc.MaxPacketSize

							midiDev := &midi.MidiDevice{
								Device:        dev,
								Manufacturer:  man,
								Product:       prod,
								VID:           vid,
								PID:           pid,
								EndpointIn:    endpoint,
								MaxPacketSize: mpSize,
							}

							midiDevices = append(midiDevices, midiDev)
						}
					}
				}
			}
		}
	}
	return midiDevices, nil
}

func (m *USBMIDIDeviceManager) Close() error {
	m.ctx.Close()
	return nil
}
