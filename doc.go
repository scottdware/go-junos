/*
package junos

Package junos provides automation for Junos (Juniper Networks) devices.

Establishing A Session

To connect to a Junos device, the process is fairly straightforward.

    jnpr := junos.NewSession(host, user, password)
    defer jnpr.Close()

Compare Rollback Configurations

If you want to view the difference between the current configuration and a rollback
one, then you can use the RollbackDiff() function.

    diff, err := jnpr.RollbackDiff(3)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(diff)
    
This will output exactly how it does on the CLI when you "| compare."
*/
package junos
