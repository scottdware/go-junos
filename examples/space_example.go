package junos

import (
	"fmt"
	"github.com/scottdware/go-junos"
	"log"
)

func main() {
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

	// Add a device to Junos Space.
	jobID, err = space.AddDevice("sdubs-fw", "admin", "juniper123")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(jobID)
	// Output: 1345283

	// Remove a device from Junos Space.
	err = space.RemoveDevice("sdubs-fw")
	if err != nil {
		fmt.Println(err)
	}

	// Resynchronize a device. A good option if you do a lot of configuration to a device
	// outside of Junos Space.
	job, err := space.Resync("firewall-A")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Job ID: %d\n", job)
	
	// Staging software on a device. The last parameter is whether or not to remove any
	// existing images from the device; boolean.
	//
	// This will not upgrade the device, but only place the image there to be used at a later
	// time.
	jobID, err = space.StageSoftware("sdubs-fw", "junos-srxsme-12.1X46-D30.2-domestic.tgz", false)
	if err != nil {
		fmt.Println(err)
	}

	// If you want to issue a software upgrade to the device, here's how:

	// Configure our options, such as whether or not to reboot the device, etc.
	options := &junos.SoftwareUpgrade{
		UseDownloaded: true,
		Validate:      false,
		Reboot:        false,
		RebootAfter:   0,
		Cleanup:       false,
		RemoveAfter:   false,
	}

	jobID, err = space.DeploySoftware("sdubs-fw", "junos-srxsme-12.1X46-D30.2-domestic.tgz", options)
	if err != nil {
		fmt.Println(err)
	}

	// Remove a staged image from the device.
	jobID, err = space.RemoveStagedSoftware("sdubs-fw", "junos-srxsme-12.1X46-D30.2-domestic.tgz")
	if err != nil {
		fmt.Println(err)
	}
	
	// List all security devices:
	devices, err := space.SecurityDevices()
	if err != nil {
		fmt.Println(err)
	}

	for _, device := range devices.Devices {
		fmt.Printf("%+v\n", device)
	}
	
	// To view the address and service objects, you use the Addresses() and Services() functions. Both of them
	// take a "filter" parameter, which lets you search for objects matching your filter.

	//If you leave the parameter blank (e.g. ""), or specify "all", then every object is returned.

	// Address objects
	addresses, err := space.Addresses("all")
	if err != nil {
		fmt.Println(err)
	}

	for _, address := range addresses.Addresses {
		fmt.Printf("%+v\n", address)
	}

	// Service objects
	services, err := space.Services("all")
	if err != nil {
		fmt.Println(err)
	}

	for _, service := range services.Services {
		fmt.Printf("%+v\n", service)
	}

	// Add an address group. "true" as the first parameter means that we assume the
	// group is going to be an address group.
	space.AddGroup(true, "Blacklist-IPs", "Blacklisted IP addresses")

	// Add a service group. We do this by specifying "false" as the first parameter.
	space.AddGroup(false, "Web-Protocols", "All web-based protocols and ports")

	// Add an address object
	space.AddAddress("my-laptop", "2.2.2.2", "My personal laptop")

	// Add a network
	space.AddAddress("corporate-users", "192.168.1.0/24", "People on campus")

	// Add a service object with an 1800 second inactivity timeout (using "0" disables this feature)
	space.AddService("udp", "udp-5000", 5000, 5000, "UDP port 5000", 1800)

	// Add a service object with a port range
	space.AddService("tcp", "high-port-range", 40000, 65000, "TCP high ports", 0)

	// If you want to modify an existing object group, you do this with the ModifyObject() function. The
	// first parameter is whether the object is an address group (true) or a service group (false).

	// Add a service to a group
	space.ModifyObject(false, "add", "service-group", "service-name")

	// Remove an address object from a group
	space.ModifyObject(true, "remove", "Whitelisted-Addresses", "bad-ip")

	// Rename an object
	space.ModifyObject(false, "rename", "Web-Services", "Web-Ports")

	// Delete an object
	space.ModifyObject(true, "delete", "my-laptop")
	
	// Add a variable
	// The parameters are as follows: variable-name, description, default-value
	space.AddVariable("test-variable", "Our test variable", "default-object")

	// Create our session state for modifying variables
	v, err := space.ModifyVariable()
	if err != nil {
		log.Fatal(err)
	}

	// Adding objects to the variable
	v.Add("test-variable", "srx-1", "user-pc")
	v.Add("test-variable", "corp-firewall", "db-server")

	// Delete a variable
	space.DeleteVariable("test-variable")
	
	// List all security policies Junos Space manages:
	policies, err := space.Policies()
	if err != nil {
		fmt.Println(err)
	}

	for _, policy := range policies.Policies {
		fmt.Printf("%s\n", policy.Name)
	}

	// For example, say we have been adding and removing objects in a group, and that group
	// is referenced in a firewall policy. Here's how to update the policy:

	// Update the policy. If "false" is specified, then the policy is only published, and the
	// device is not updated.
	job, err = space.PublishPolicy("Internet-Firewall-Policy", true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Job ID: %d\n", job)

	// Let's update a device knowing that we have some previously published services.
	job, err = space.UpdateDevice("firewall-1.company.com")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Job ID: %d\n", job)
}
