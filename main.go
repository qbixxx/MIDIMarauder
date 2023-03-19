package main

import (
	"os"
	"fmt"
	"log"
	"time"
	"github.com/google/gousb"
	"github.com/google/gousb/usbid"
	
)

func main(){
	
	Example_simple()
	
}

func Example_simple() {
	// Initialize a new Context.
	ctx := gousb.NewContext()
	defer ctx.Close()

	fmt.Println("Context: \n", ctx)



	// Open any device with a given VID/PID using a convenience function.
	dev, err := ctx.OpenDeviceWithVIDPID(0x2467, 0x2034)
	if err != nil {
		log.Fatalf("Could not open a device: %v", err)
	}
	defer dev.Close()

	dev.SetAutoDetach(true)

	
	config := dev.String() 

	a, _ := dev.GetStringDescriptor(3)

	fmt.Println("a: ", a)
	b, _ := dev.GetStringDescriptor(2)

	fmt.Println("b: ", b)
	c, _ := dev.GetStringDescriptor(1)

	fmt.Println("c: ", c)

	manu, _ := dev.Manufacturer()
	
	midi_b, _ := dev.Product()
	fmt.Println("Config ",	config)
	fmt.Println("product: ",	midi_b)
	fmt.Println("Manufacturer: ", manu)

// Iterate through configurations
	for num := range dev.Desc.Configs {
		config, _ := dev.Config(num)
	
		

		// In a scenario where we have an error, we can continue
		// to the next config. Same is true for interfaces and
		// endpoints.
		defer config.Close()
	
		// Iterate through available interfaces for this configuration
		for _, desc := range config.Desc.Interfaces {
		  intf, _ := config.Interface(desc.Number, 0)	



		  classy := usbid.Classify(dev.Desc)


		  fmt.Println("Classify: \n", classy)


		  fmt.Println("Interface: ", intf)
		  // Iterate through endpoints available for this interface.
		  for _, endpointDesc := range intf.Setting.Endpoints {

			  if endpointDesc.Direction == gousb.EndpointDirectionIn {
			
			  
			  endpoint, _ := intf.InEndpoint(endpointDesc.Number)
			
			
			  fmt.Println(endpoint)


			  fmt.Println("endpoint poll interval", endpointDesc.PollInterval, "endpoint max packet size", endpointDesc.MaxPacketSize)

			  mdev := &MIDIDEV{
				context:   ctx,
				device:    dev,
				intf:      intf,
				endpoint:  endpoint,
			  }

			  fmt.Println("\n",mouse.context,"\n",mouse.device,"\n",mouse.intf,"\n",mouse.endpoint)
			  mdev .read(endpointDesc.PollInterval, endpointDesc.MaxPacketSize)


		  }
		}
	  }


	}

}


func (mdev *MIDIDEV) read(interval time.Duration, maxSize int) {

	
	interval = 2

	maxSize = 64
	fmt.Println("time duration: ", interval, "max size: ", maxSize)
	
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	

	for {
		
	  select {
	  case <-ticker.C:
		buff := make([]byte, maxSize)
		n, _ := mouse.endpoint.Read(buff)

		data := buff[:n]
	  
		if data[0] == 14{
			fmt.Println("CC:", data[2]," Value: ", data[3])
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


type MIDIDEV struct {
	// Fields for interacting with the USB connection
	context  *gousb.Context
	device   *gousb.Device
	intf     *gousb.Interface
	endpoint *gousb.InEndpoint
  
	// Additional fields we'll get to later
  }
