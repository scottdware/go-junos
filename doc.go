/*
Package junos allows you to run commands on and configure Junos devices.

Establishing a session
	jnpr := junos.NewSession(host, user, password)

Lock config
	err := jnpr.Lock()
	if err != nil {
		log.Fatal(err)
	}

Unlock config
	err = jnpr.Unlock()
	if err != nil {
		log.Fatal(err)
	}

Compare the current configuration to a rollback config.
	diff, err := jnpr.RollbackDiff(3)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	fmt.Println(diff)

Rollback to an older configuration.
	err := jnpr.RollbackConfig(2)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}

Run operational mode commands, such as "show."
	output, err := jnpr.Command("show version", "text")
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	fmt.Println(output)

When you specify "text," the output will be just like it is on the CLI:

    node0:
    --------------------------------------------------------------------------
    Hostname: firewall-1
    Model: srx240h2
    JUNOS Software Release [12.1X47-D10.4]

    node1:
    --------------------------------------------------------------------------
    Hostname: firewall-2
    Model: srx240h2
    JUNOS Software Release [12.1X47-D10.4]

*/
package junos