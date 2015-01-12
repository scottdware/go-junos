go-junos
========

[![GoDoc](https://godoc.org/github.com/scottdware/go-junos?status.svg)](https://godoc.org/github.com/scottdware/go-junos)

A Go package that interacts with Junos devices and allows you to do the following:

* Run operational mode commands, such as `show`, `request`, etc..
* Compare the active configuration to a given rollback config.
* Rollback the configuration to a given state or a "rescue" config.
* Configure devices by uploading a file locally or using a FTP/HTTP url.

Visit the [GoDoc][4] page for package documentation.

> This package makes all of it's calls over [Netconf][1] using the [go-netconf][2] package from [Juniper Networks][3]

Example
-------
```Go
package main

import (
	"fmt"
	"github.com/scottdware/go-junos"
	"os"
)

func main() {
    // Establish a connection to a device.
	jnpr := junos.NewSession("srx-1", "admin", "juniper123")
    defer jnpr.Close()

    // Compare the current running config to "rollback 1."
	diff, _ := jnpr.RollbackDiff(1)
	fmt.Println(diff)

    // Load a configuration file with "set" commands, and commit it.
	err := jnpr.LoadConfig("C:/Configs/juniper.txt", "set", true)
	if err != nil {
		fmt.Println(err)
	}
    
    // Rollback to a previous config.
    err = jnpr.RollbackConfig(5)
    if err != nil {
        fmt.Println(err)
    }
    
    // Run a command and return the results in "text" format (similar to CLI).
    output, err := jnpr.Command("show security ipsec inactive-tunnels", "text")
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(output)
}
```

[1]: https://tools.ietf.org/html/rfc6241
[2]: https://github.com/Juniper/go-netconf
[3]: http://www.juniper.net
[4]: https://godoc.org/github.com/scottdware/go-junos