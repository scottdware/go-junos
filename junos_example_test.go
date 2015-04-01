package junos

// To View the entire configuration, use the keyword "full" for the first
// argument. If anything else outside of "full" is specified, it will return
// the configuration of the specified top-level stanza only. So "security"
// would return everything under the "security" stanza.
func ExampleJunos_viewConfiguration() {
	// Establish our session first.
	jnpr, err := junos.NewSession(host, user, password)
	if err != nil {
		log.Fatal(err)
	}
	defer jnpr.Close()

	// Output format can be "text" or "xml".
	config, err := jnpr.GetConfig("full", "text")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(config)
}

// Comparing and working with rollback configurations.
func ExampleJunos_rollbackConfigurations() {
	// Establish our session first.
	jnpr, err := junos.NewSession(host, user, password)
	if err != nil {
		log.Fatal(err)
	}
	defer jnpr.Close()

	// If you want to view the difference between the current configuration and a rollback
	// one, then you can use the ConfigDiff() function to specify a previous config:
	diff, err := jnpr.ConfigDiff(3)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(diff)

	// You can rollback to a previous state, or the rescue configuration by using
	// the RollbackConfig() function:
	err := jnpr.RollbackConfig(3)
	if err != nil {
		fmt.Println(err)
	}

	// Create a rescue config from the active configuration.
	jnpr.Rescue("save")

	// You can also delete a rescue config.
	jnpr.Rescue("delete")

	// Rollback to the "rescue" configuration.
	err := jnpr.RollbackConfig("rescue")
	if err != nil {
		fmt.Println(err)
	}
}

// Configuring devices.
func ExampleJunos_configuringDevices() {
	// Use the Config() function to configure a Junos device.

	// When configuring a device, it is good practice to lock the configuration database,
	// load the config, commit the configuration, and then unlock the configuration database.
	// You can do this with the following functions: Lock(), Commit(), Unlock().

	// Multiple ways to commit a configuration.

	// Commit the configuration as normal.
	Commit()

	// Check the configuration for any syntax errors (NOTE: you must still issue a
	// Commit() afterwards).
	CommitCheck()

	// Commit at a later time, i.e. 4:30 PM.
	CommitAt("16:30:00")

	// Rollback configuration if a Commit() is not issued within the given <minutes>.
	CommitConfirm(15)

	// You can configure the Junos device by uploading a local file, or pulling from an
	// FTP/HTTP server. The LoadConfig() function takes three arguments:

	// filename or URL, format, and a boolean (true/false) "commit-on-load".

	// If you specify a URL, it must be in the following format:

	// ftp://<username>:<password>@hostname/pathname/file-name
	// http://<username>:<password>@hostname/pathname/file-name

	// Note: The default value for the FTP path variable is the user’s home directory. Thus,
	// by default the file path to the configuration file is relative to the user directory.
	// To specify an absolute path when using FTP, start the path with the characters %2F;
	// for example: ftp://username:password@hostname/%2Fpath/filename.

	// The format of the commands within the file must be one of the following types:

	// set
	// system name-server 1.1.1.1

	// text
	// system {
	//     name-server 1.1.1.1;
	// }

	// xml
	// <system>
	//     <name-server>
	//         <name>1.1.1.1</name>
	//     </name-server>
	// </system>

	// Establish our session first.
	jnpr, err := junos.NewSession(host, user, password)
	if err != nil {
		log.Fatal(err)
	}
	defer jnpr.Close()

	// Configure a device from a file:

	// If the third option is "true" then after the configuration is loaded, a commit
	// will be issued. If set to "false," you will have to commit the configuration
	// using one of the Commit() functions.
	jnpr.Lock()
	err = jnpr.Config("path-to-file.txt", "set", true)
	if err != nil {
		fmt.Println(err)
	}
	jnpr.Unlock()

	// Configure a device using commands specified in a []string:

	// You can also use the following format also:
	//
	// var configCommands = `
	//     set interfaces ge-0/0/5.0 family inet address 192.168.0.1/24
	//     set security zones security-zone trust interfaces ge-0/0/5.0
	// `
	//
	configCommands := []string{
		"set interfaces ge-0/0/5.0 family inet address 192.168.0.1/24",
		"set security zones security-zone trust interfaces ge-0/0/5.0",
	}
	err = jnpr.Config(configCommands, "set", true)
	if err != nil {
		fmt.Println(err)
	}
}

// Running operational mode commands on a device.
func ExampleJunos_runningCommands() {
	// Establish our session first.
	jnpr, err := junos.NewSession(host, user, password)
	if err != nil {
		log.Fatal(err)
	}
	defer jnpr.Close()

	// You can run operational mode commands such as "show" and "request" by using the
	// Command() function. Output formats can be "text" or "xml".

	// Results returned in text format.
	txtOutput, err := jnpr.RunCommand("show chassis hardware", "text")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(txtOutput)

	// Results returned in XML format.
	xmlOutput, err := jnpr.RunCommand("show chassis hardware", "xml")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(xmlOutput)

	// Reboot the device
	jnpr.Reboot()
}

// Viewing basic information about the device.
func ExampleJunos_deviceInformation() {
	// Establish our session first.
	jnpr, err := junos.NewSession(host, user, password)
	if err != nil {
		log.Fatal(err)
	}
	defer jnpr.Close()

	// When you call the PrintFacts() function, it just prints out the platform
	// and software information to the console.
	jnpr.PrintFacts()

	// You can also loop over the struct field that contains this information yourself:
	fmt.Printf("Hostname: %s", jnpr.Hostname)
	for _, data := range jnpr.Platform {
		fmt.Printf("Model: %s, Version: %s", data.Model, data.Version)
	}
	// Output: Model: SRX240H2, Version: 12.1X47-D10.4
}

// Establishing a connection to Junos Space and working with devices.
func ExampleJunosSpace_devices() {
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
	fmt.Println(job)
}

// Software upgrades using Junos Space.
func ExampleJunosSpace_softwareUpgrade() {
	// Establish a connection to a Junos Space server.
	space := junos.NewServer("space.company.com", "admin", "juniper123")

	// Staging software on a device. The last parameter is whether or not to remove any
	// existing images from the device; boolean.
	//
	// This will not upgrade the device, but only place the image there to be used at a later
	// time.
	jobID, err := space.StageSoftware("sdubs-fw", "junos-srxsme-12.1X46-D30.2-domestic.tgz", false)
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

	jobID, err := space.DeploySoftware("sdubs-fw", "junos-srxsme-12.1X46-D30.2-domestic.tgz", options)
	if err != nil {
		fmt.Println(err)
	}

	// Remove a staged image from the device.
	jobID, err := space.RemoveStagedSoftware("sdubs-fw", "junos-srxsme-12.1X46-D30.2-domestic.tgz")
	if err != nil {
		fmt.Println(err)
	}
}

// Viewing information about Security Director devices (SRX, J-series, etc.).
func ExampleJunosSpace_securityDirectorDevices() {
	// Establish a connection to a Junos Space server.
	space := junos.NewServer("space.company.com", "admin", "juniper123")

	// List all security devices:
	devices, err := space.SecurityDevices()
	if err != nil {
		fmt.Println(err)
	}

	for _, device := range devices.Devices {
		fmt.Printf("%+v\n", device)
	}
}

// Working with address and service objects.
func ExampleJunosSpace_addressObjects() {
	// Establish a connection to a Junos Space server.
	space := junos.NewServer("space.company.com", "admin", "juniper123")

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
}

// Working with polymorphic (variable) objects.
func ExampleJunosSpace_variables() {
	// Establish a connection to a Junos Space server.
	space := junos.NewServer("space.company.com", "admin", "juniper123")

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
}

// Working with policies.
func ExampleJunosSpace_policies() {
	// Establish a connection to a Junos Space server.
	space := junos.NewServer("space.company.com", "admin", "juniper123")

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
	job, err := space.PublishPolicy("Internet-Firewall-Policy", true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Job ID: %d\n", job)

	// Let's update a device knowing that we have some previously published services.
	job, err := space.UpdateDevice("firewall-1.company.com")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Job ID: %d\n", job)
}
