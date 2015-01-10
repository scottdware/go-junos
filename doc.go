/*
Package junos allows you to run commands on and configure Junos devices.

Establishing a session
	jnpr := junos.NewSession(host, user, password)

Locking the configuration
	err := jnpr.Lock()
	if err != nil {
		log.Fatal(err)
	}

Commiting the configuration
	err = jnpr.Commit()
	if err != nil {
		log.Fatal(err)
	}
    
Unlocking the configuration
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
    
The output will be exactly as it is running the "| compare" command on the CLI:

    [edit forwarding-options helpers bootp server 192.168.10.2]
    -     routing-instance srx-vr;
    [edit forwarding-options helpers bootp server 192.168.10.3]
    -     routing-instance srx-vr;
    [edit security address-book global]
         address server1 { ... }
    +    address dc-console 192.168.20.15/32;
    +    address dc-laptop 192.168.22.7/32;
    [edit security zones security-zone vendors interfaces]
          reth0.1000 { ... }
    +     reth0.520 {
    +         host-inbound-traffic {
    +             system-services {
    +                 dhcp;
    +                 ping;
    +             }
    +         }
    +     }

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