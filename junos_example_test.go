package junos

// Establishing a session
func Example_session() {
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
