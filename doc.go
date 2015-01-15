/*
Package junos provides automation for Junos (Juniper Networks) devices.

Establishing A Session

To connect to a Junos device, the process is fairly straightforward.

    jnpr := junos.NewSession(host, user, password)
    defer jnpr.Close()

Viewing The Configuration

    // Output format can be "text" or "xml"
    config, err := jnpr.GetConfig("text", "full")
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(config)

    // Viewing only a certain part of the configuration
    routing, err := jnpr.GetConfig("text", "routing-instances")
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(routing)

Compare Rollback Configurations

If you want to view the difference between the current configuration and a rollback
one, then you can use the ConfigDiff() function.

    diff, err := jnpr.ConfigDiff(3)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(diff)

This will output exactly how it does on the CLI when you "| compare."

Rolling Back to a Previous State

You can also rollback to a previous state, or the "rescue" configuration by using
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
*/
package junos
