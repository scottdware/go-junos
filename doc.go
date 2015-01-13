/*
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

Device Configuration

You can configure the Junos device by uploading a local file, or pulling from an
FTP/HTTP server. The commands within the config file can be any of the following types:
    set, text, or xml

*/
package junos
