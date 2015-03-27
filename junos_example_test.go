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
