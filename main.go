// Copyright 2017 the gousb Authors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"os"
	"fmt"
	"log"
	"time"
	"github.com/google/gousb"
	//"github.com/google/gousb/usbid"
	
)


type MOUSE struct {
	// Fields for interacting with the USB connection
	context  *gousb.Context
	device   *gousb.Device
	intf     *gousb.Interface
	endpoint *gousb.InEndpoint
  
	// Additional fields we'll get to later
  }

  type SubClass struct {
	Name     string
	Protocol map[uint8]string
}

// This examples demonstrates the use of a few convenience functions that
// can be used in simple situations and with simple devices.
// It opens a device with a given VID/PID,
// claims the default interface (use the same config as currently active,
// interface 0, alternate setting 0) and tries to write 5 bytes of data
// to endpoint number 7.


func (mouse *MOUSE) read(interval time.Duration, maxSize int) {

	
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
		// Do something with this data
		//fmt.Printf("DATA: %x\n", data[0])
	  
		if data[0] == 14{
			fmt.Println("CC:", data[2]," Value: ", data[3])
		}
		if data[0] == 11{
			fmt.Println("CC:", data[2]," Value: ", data[3])
		}
		if data[0] == 8 || data[0] == 9{
			fmt.Println("Note: ", data[2]," Value: ", data[3])
		}

		if data[0] == 1 && data[2] == 6{
			fmt.Println("Control C !!!")
			os.Exit(0)
		}
	  }

	}
	
  }
  

func Example_simple() {
	// Initialize a new Context.


	ctx := gousb.NewContext()
	defer ctx.Close()

	fmt.Println("Context: \n", ctx)



	// Open any device with a given VID/PID using a convenience function.
	dev, err := ctx.OpenDeviceWithVIDPID(0x2467,0x2034)
	if err != nil {
		log.Fatalf("Could not open a device: %v", err)
	}
	defer dev.Close()

	dev.SetAutoDetach(true)

	
	// Claim the default interface using a convenience function.
	// The default interface is always #0 alt #0 in the currently active
	// config.
	config := dev.String() 

	configDesc, _ := dev.InterfaceDescription(1,1,0)
	a, _ := dev.GetStringDescriptor(3)
	//class := a.Classify()

	fmt.Println("a: ", a)
		//if errror != nil{
		//	fmt.Println("Error != nil")
		//	fmt.Println("Error: ", errror)
//
		//}
		//if errror == nil{
		//	fmt.Println("Error == nil")
		//	fmt.Println("Error: ", errror)
		//}
		
	fmt.Println("ConfigDesc: ",*(&configDesc))
	manu, _ := dev.Manufacturer()
	
	midi_b, _ := dev.Product()
	//midi_b := dev.Class()
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



		  fmt.Println("Interface: ", intf)
		  // Iterate through endpoints available for this interface.
		  for _, endpointDesc := range intf.Setting.Endpoints {
			// We only want to read, so we're looking for IN endpoints.
			if endpointDesc.Direction == gousb.EndpointDirectionIn {
			
			  
			  endpoint, _ := intf.InEndpoint(endpointDesc.Number)
			
			//  buff := make([]byte, 30)
			
			  fmt.Println(endpoint)

	//		  fmt.Println("Class: ", desc.Class())

			  //fmt.Println("ACA")
			  fmt.Println("endpoint poll interval", endpointDesc.PollInterval, "endpoint max packet size", endpointDesc.MaxPacketSize)

			  mouse := &MOUSE{
				context:   ctx,
				device:    dev,
				intf:      intf,
				endpoint:  endpoint,
			  }

			  fmt.Println("\n",mouse.context,"\n",mouse.device,"\n",mouse.intf,"\n",mouse.endpoint)
			  mouse.read(endpointDesc.PollInterval, endpointDesc.MaxPacketSize)

			  fmt.Println("ACA222")

			  // When we get here, we have an endpoint where we can
			  // read data from the USB device
			
		  }
		}
	  }






func main(){
	Example_simple()
	
}
