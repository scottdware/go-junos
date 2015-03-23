/*
Package junos provides automation for Junos (Juniper Networks) devices, as
well as interaction with Junos Space.

Establishing A Session

To connect to a Junos device, the process is fairly straightforward:

    jnpr := junos.NewSession(host, user, password)
    defer jnpr.Close()

Viewing The Configuration

To View the entire configuration, use the keyword "full" for the second
argument. If anything else outside of "full" is specified, it will return
the configuration of that section only. So "security" would return everything
under the "security" stanza.

    // Output format can be "text" or "xml"
    config, err := jnpr.GetConfig("full", "text")
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(config)

    // Viewing only a certain part of the configuration
    routing, err := jnpr.GetConfig("routing-instances", "text")
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(routing)

Compare Rollback Configurations

If you want to view the difference between the current configuration and a rollback
one, then you can use the ConfigDiff() function:

    diff, err := jnpr.ConfigDiff(3)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(diff)

This will output exactly how it does on the CLI when you "| compare."

Rolling Back to a Previous State

You can also rollback to a previous state, or the rescue configuration by using
the RollbackConfig() function:

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

Device Configuration

When configuring a device, it is good practice to lock the configuration database,
load the config, commit the configuration, and then unlock the configuration database.

You can do this with the following functions:

    Lock(), Commit(), Unlock()

There are multiple ways to commit a configuration as well:

    // Commit the configuration as normal
    Commit()

    // Check the configuration for any syntax errors (NOTE: you must still issue a Commit())
    CommitCheck()

    // Commit at a later time, i.e. 4:30 PM
    CommitAt("16:30:00")

    // Rollback configuration if a Commit() is not issued within the given <minutes>.
    CommitConfirm(15)

You can configure the Junos device by uploading a local file, or pulling from an
FTP/HTTP server. The LoadConfig() function takes three arguments:

    filename or URL, format, and commit-on-load

If you specify a URL, it must be in the following format:

    ftp://<username>:<password>@hostname/pathname/file-name
    http://<username>:<password>@hostname/pathname/file-name

    Note: The default value for the FTP path variable is the user’s home directory. Thus,
    by default the file path to the configuration file is relative to the user directory.
    To specify an absolute path when using FTP, start the path with the characters %2F;
    for example: ftp://username:password@hostname/%2Fpath/filename.

The format of the commands within the file must be one of the following types:

    set
    // system name-server 1.1.1.1

    text
    // system {
    //     name-server 1.1.1.1;
    // }

    xml
    // <system>
    //     <name-server>
    //         <name>1.1.1.1</name>
    //     </name-server>
    // </system>

If the third option is "true" then after the configuration is loaded, a commit
will be issued. If set to "false," you will have to commit the configuration
using the Commit() function.

    jnpr.Lock()
    err := jnpr.LoadConfig("path-to-file.txt", "set", true)
    if err != nil {
        fmt.Println(err)
    }
    jnpr.Unlock()

You don't have to use Lock() and Unlock() if you wish, but if by chance someone
else tries to edit the device configuration at the same time, there can be conflics
and most likely an error will be returned.

Running Commands

You can run operational mode commands such as "show" and "request" by using the
Command() function. Output formats can be "text" or "xml":

    // Results returned in text format
    output, err := jnpr.Command("show chassis hardware", "text")

    // Results returned in XML format
    output, err := jnpr.Command("show chassis hardware", "xml")

Viewing Platform and Software Information

When you call the PrintFacts() function, it prints out the platform and software information:

    jnpr.PrintFacts()

    // Returns output similar to the following
    node0
    --------------------------------------------------------------------------
    Hostname: firewall-1
    Model: SRX240H2
    Version: 12.1X47-D10.4

    node1
    --------------------------------------------------------------------------
    Hostname: firewall-1
    Model: SRX240H2
    Version: 12.1X47-D10.4

You can also loop over the struct field that contains this information yourself:

    fmt.Printf("Hostname: %s", jnpr.Hostname)
    for _, data := range jnpr.Platform {
        fmt.Printf("Model: %s, Version: %s", data.Model, data.Version)
    }

Junos Space - Network Management Platform Functions

Here's an example of how to connect to a Junos Space server, and get information about
all of the managed devices:

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

How to add and remove devices:

    // Add a device to Junos Space.
    jobID, err = space.AddDevice("sdubs-fw", "admin", "juniper123")
    if err != nil {
        fmt.Println(err)
    }

    // Remove a device from Junos Space.
    err = space.RemoveDevice("sdubs-fw")
    if err != nil {
        fmt.Println(err)
    }

    // Here's a good way to loop through all devices, and find the one you want to delete:
    d, err := space.Devices()
    if err != nil {
        fmt.Println(err)
    }

    for _, device := range d.Devices {
        if device.Name == "sdubs-fw" {
            err = space.RemoveDevice(device.Name)
            if err != nil {
                fmt.Println(err)
            }

            fmt.Printf("Deleted device: %s\n", device.Name)
        }
    }

Stage a software image on a device. This basically just downloads it to the device, and does
not upgrade it:

    // The third parameter is whether or not to remove any existing images from the device - true or false.
    jobID, err := space.StageSoftware("sdubs-fw", "junos-srxsme-12.1X46-D30.2-domestic.tgz", false)
    if err != nil {
        fmt.Println(err)
    }

If you want to issue a software upgrade to the device, here's how:

    // Configure our options, such as whether or not to reboot the device, etc.
    options := &junos.SoftwareUpgrade{
        UseDownloaded: true,
        Validate: false,
        Reboot: false,
        RebootAfter: 0,
        Cleanup: false,
        RemoveAfter: false,
    }

    jobID, err := space.DeploySoftware("sdubs-fw", "junos-srxsme-12.1X46-D30.2-domestic.tgz", options)
    if err != nil {
        fmt.Println(err)
    }

Remove a staged image from a device:

    // Remove a staged image from the device.
    jobID, err := space.RemoveStagedSoftware("sdubs-fw", "junos-srxsme-12.1X46-D30.2-domestic.tgz")
    if err != nil {
        fmt.Println(err)
    }

Junos Space - Security Director Functions

List all security devices, and display the SecurityDevices struct information about them:

    devices, err := space.SecurityDevices()
    if err != nil {
        fmt.Println(err)
    }

    for _, device := range devices.Devices {
        fmt.Printf("%+v\n", device)
    }

To view the address and service objects, you use the Addresses() and Services() functions. Both of them
take a "filter" parameter, which lets you search for objects matching your filter. If you specify "all",
then the full list is returned.

If you leave the parameter blank (e.g. ""), all objects are returned.

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

To add objects to Junos Space, you can do the following:

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

Adding polymorphic (variable) objects works in a similar way:

	// Add a variable
	// The parameters are as follows: variable-name, description, default-value
	space.AddVariable("test-variable", "Our test variable", "default-object")

	// Modify a variable by assigning devices/objects to it
	// The parameters are as follows: action (add or delete), variable-name, SD device (firewall), address-object
	space.ModifyVariable("add", "test-variable", "srx-firewall1", "my-home-network")

	// Delete a variable
	space.ModifyVariable("delete", "test-variable")

If you want to modify an existing object, you can do that with the ModifyObject() function. The
first parameter is whether the object is an address (true) or a service object (false).

    // Add a service to a group
    space.ModifyObject(false, "add", "service-group", "service-name")

    // Remove an address object from a group
    space.ModifyObject(true, "remove", "Whitelisted-Addresses", "bad-ip")

    // Rename an object
    space.ModifyObject(false, "rename", "Web-Services", "Web-Ports")

    // Delete an object
    space.ModifyObject(true, "delete", "my-laptop")

Let's take a look at what security policies Space manages:

    policies, err := space.Policies()
    if err != nil {
        fmt.Println(err)
    }

    for _, policy := range policies.Policies {
        fmt.Printf("%s\n", policy.Name)
    }

For example, say we have been adding and removing objects to/from a group. That group
is referenced in a firewall policy, so let's update it:

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
*/
package junos
