package junos

// Establishing a session
func Example() {
	jnpr, err := junos.NewSession(host, user, password)
	if err != nil {
		log.Fatal(err)
	}
	defer jnpr.Close()
}

// To View the entire configuration, use the keyword "full" for the first
// argument. If anything else outside of "full" is specified, it will return
// the configuration of the specified top-level stanza only. So "security" would return everything
// under the "security" stanza.
func Example_viewConfiguration() {
	// Output format can be "text" or "xml"
	config, err := jnpr.GetConfig("full", "text")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(config)
}

// Comparing and Rolling Back Configurations
func Example_rollbackConfigurations() {
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

// Device Configuration
func Example_configuringDevices() {
	// Use the LoadConfig() function to load the configuration from a file.

	// When configuring a device, it is good practice to lock the configuration database,
	// load the config, commit the configuration, and then unlock the configuration database.
	// You can do this with the following functions: Lock(), Commit(), Unlock()

	// Multiple ways to commit a configuration

	// Commit the configuration as normal
	Commit()

	// Check the configuration for any syntax errors (NOTE: you must still issue a Commit() afterwards)
	CommitCheck()

	// Commit at a later time, i.e. 4:30 PM
	CommitAt("16:30:00")

	// Rollback configuration if a Commit() is not issued within the given <minutes>.
	CommitConfirm(15)

	// You can configure the Junos device by uploading a local file, or pulling from an
	// FTP/HTTP server. The LoadConfig() function takes three arguments:

	// filename or URL, format, and a boolean (true/false) "commit-on-load"

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

	// If the third option is "true" then after the configuration is loaded, a commit
	// will be issued. If set to "false," you will have to commit the configuration
	// using one of the Commit() functions.
	jnpr.Lock()
	err := jnpr.LoadConfig("path-to-file.txt", "set", true)
	if err != nil {
		fmt.Println(err)
	}
	jnpr.Unlock()
}

// Running Commands on a Device
func Example_runCommands() {
	// You can run operational mode commands such as "show" and "request" by using the
	// Command() function. Output formats can be "text" or "xml":

	// Results returned in text format
	txtOutput, err := jnpr.Command("show chassis hardware", "text")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(txtOutput)

	// Results returned in XML format
	xmlOutput, err := jnpr.Command("show chassis hardware", "xml")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(xmlOutput)
}

func Example_displayPlatform_Information() {
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
