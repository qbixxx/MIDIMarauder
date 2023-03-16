package main

import (
	"fmt"
	"log"
	"time"
	"github.com/google/gousb"
	//"github.com/google/usbid"
)


type DEVICE struct {
	// Fields for interacting with the USB connection
	context  *gousb.Context
	device   *gousb.Device
	intf     *gousb.Interface
	endpoint *gousb.InEndpoint
  
  }


func (deviceSt *DEVICE) read(interval time.Duration, maxSize int) {

	
	interval = 2 // hardcoded elements

	maxSize = 64
	fmt.Println("time duration: ", interval, "max size: ", maxSize)
	
	ticker := time.NewTicker(interval)

	defer ticker.Stop()
	
	for {
		
	  select {
	  case <-ticker.C:
		buff := make([]byte, maxSize)
		n, _ := deviceSt.endpoint.Read(buff)

		data := buff[:n]

		if data[0] == 14{
			fmt.Println("Pitch Bend-> Parameter:", data[2]," Value: ", data[3])
		}
		if data[0] == 11{
			fmt.Println("CC:", data[2]," Value: ", data[3])
		}
		if data[0] == 8 || data[0] == 9{
			fmt.Println("Note: ", data[2]," Value: ", data[3])
		}
	}

	}
}
  

func Example_simple() {
	// Initialize strDesc new Context.
	ctx := gousb.NewContext()
	defer ctx.Close()

	fmt.Println("Context: \n", ctx)

	// Open any device with strDesc given VID/PID using strDesc convenience function.
	dev, err := ctx.OpenDeviceWithVIDPID(0x2467,0x2034) // MIDI Keyboard VID,PID (Nektar GX61)
	if err != nil {
		log.Fatalf("Could not open strDesc device: %v", err)
	}
	defer dev.Close()

	dev.SetAutoDetach(true)



	intDesc, _ := dev.InterfaceDescription(1,1,0)
	strDesc, _ := dev.GetStringDescriptor(3)
	fmt.Println("strDesc: ", strDesc)

	fmt.Println("ConfigDesc: ",*(&intDesc))
	manuf, _ := dev.Manufacturer()

	midi_b, _ := dev.Product()
	fmt.Println("product: ",	midi_b)
	fmt.Println("Manufacturer: ", manuf)


		
	defer config.Close()
	
		// Iterate through available interfaces for this configuration
	for _, desc := range config.Desc.Interfaces {
		intf, _ := config.Interface(desc.Number, 0)
		fmt.Println("Interface: ", intf)
	  // Iterate through endpoints available for this interface.
		for _, endpointDesc := range intf.Setting.Endpoints {
		// We only want to read, so we're looking for IN endpoints.
			if endpointDesc.Direction == gousb.EndpointDirectionIn {
			
			endpoint, _ := intf.InEndpoint(endpointDesc.Number)
			fmt.Println(endpoint)


			fmt.Println("endpoint poll interval", endpointDesc.PollInterval, "endpoint max packet size", endpointDesc.MaxPacketSize)

			deviceSt := &DEVICE{
				context:   ctx,
				device:    dev,
				intf:      intf,
				endpoint:  endpoint,
			}

			fmt.Println("\n",deviceSt.context,"\n",deviceSt.device,"\n",deviceSt.intf,"\n",deviceSt.endpoint)
			  
			 
			deviceSt.read(endpointDesc.PollInterval, endpointDesc.MaxPacketSize)



			}
		}
	}	


}
	

func main(){

	Example_simple()
	
}
