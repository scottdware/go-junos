/*
Package junos allows you to run commands on and configure Junos devices.

	// Establish our session
	jnpr := junos.NewSession(host, user, password)

	// Lock config
	err := jnpr.Lock()
	if err != nil {
		log.Fatal(err)
	}

	// Unlock config
	err = jnpr.Unlock()
	if err != nil {
		log.Fatal(err)
	}

	// Rollback diff compare
	diff, err := jnpr.RollbackDiff(3)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	fmt.Println(diff)

	// Rollback config
	err := jnpr.RollbackConfig(2)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}s

	// Show command
	output, err := jnpr.Command("show version", "text")
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	fmt.Println(output)

	// Close the connection
	jnpr.Close()
*/
package junos