package usb

import (
	"fmt"
	"github.com/google/gousb"
	"midiMarauder/internal/midi"
)

func ScanForMIDIDevices(ctx *gousb.Context) ([]*midi.MidiDevice, error) {
	devices, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
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

	var midiDevices []*midi.MidiDevice
	for _, dev := range devices {

		for num := range dev.Desc.Configs {

			config, _ := dev.Config(num)

			for _, intfDesc := range config.Desc.Interfaces {
				interFace, err := config.Interface(intfDesc.Number, 0)
				if err != nil {
					//fmt.Printf("Error initializing interface: %v\n", err)
					continue
				}

				for _, interFaceSetting := range intfDesc.AltSettings {
					for _, endpointDesc := range interFaceSetting.Endpoints {
						if interFaceSetting.Class == gousb.ClassAudio && interFaceSetting.SubClass == 3 && endpointDesc.Direction == gousb.EndpointDirectionIn {

							err = dev.SetAutoDetach(true)
							if err != nil {
								//fmt.Printf("Failed to detach kernel driver: %v\n", err)
								continue
							}

							endpoint, err := interFace.InEndpoint(endpointDesc.Number)
							if err != nil {
								fmt.Printf("Error accessing InEndpoint: %v\n", err)
								continue
							}

							man, _ := dev.Manufacturer()
							prod, _ := dev.Product()
							serialN, _ := dev.SerialNumber()
							vid := dev.Desc.Vendor
							pid := dev.Desc.Product
							mpSize := endpointDesc.MaxPacketSize
							class := dev.Desc.Class.String()
							subClass := dev.Desc.SubClass
							protocol := dev.Desc.Protocol
							speed := dev.Desc.Speed

							d := midi.MidiDevice{
								Device:        dev,
								Manufacturer:  man,
								Product:       prod,
								VID:           vid,
								PID:           pid,
								EndpointIn:    endpoint,
								MaxPacketSize: mpSize,
								DeviceConfig:  config.String(),
								SerialNumber:  serialN,
								Class:         class,
								SubClass:      subClass,
								Protocol:      protocol,
								Speed:         speed,
							}

							midiDevices = append(midiDevices, &d)

						}
					}
				}

			}

		}

	}

	return midiDevices, nil
}
