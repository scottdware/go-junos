package junos

import (
	"fmt"
	"github.com/scottdware/go-junos"
)

func Example_JunosSpace() {
	// Establish a connection to a Junos Space server.
	space := junos.NewServer("space.company.com", "admin", "juniper123")

	// Get the list of devices.
	devices, err := space.Devices()
	if err != nil {
		fmt.Println(err)
	}

	// Iterate over our device list and display some information about them.
	for _, device := range devices.Devices {
		fmt.Printf("Name: %s, IP Address: %s, Platform: %s\n", device.Name, device.IP, device.Platform)
	}
}
