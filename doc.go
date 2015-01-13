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
FTP/HTTP server. The LoadConfig() function takes three arguments:

    filename or URL, format, and commit-on-load
    
If you specify a URL, it must be in the following format:

    ftp://user@password:path-to-file
    http://user@password/path-to-file
    
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
    
Using the LoadConfig() function, here's how you would do this.

    err := jnpr.LoadConfig("path-to-file.txt", "set", true)
    if err != nil {
        fmt.Println(err)
    }
    
If th
*/
package junos
