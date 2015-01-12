go-junos
========

[![GoDoc](https://godoc.org/github.com/scottdware/go-junos?status.svg)](https://godoc.org/github.com/scottdware/go-junos)

A Go package that interacts with Junos devices and allows you to do the following:

* Run operational mode commands, such as `show`, `request`, etc..
* Compare the active configuration to a given rollback config.
* Rollback the configuration to a given state or a "rescue" config.
* Configure devices by uploading a file with commands in "set" format.

Visit the [GoDoc][4] page for example usage and documentation.

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

var (
	host     = os.Args[1]
	user     = os.Args[2]
	password = os.Args[3]
)

func main() {
	jnpr := junos.NewSession(host, user, password)
    defer jnpr.Close()

    // Compare the current running config to "rollback 1"
	diff, _ := jnpr.RollbackDiff(1)
	fmt.Println(diff)

    err := jnpr.Lock()
    if err != nil {
        fmt.Println(err)
    }
    
    // Load a configuration file with "set" commands
	err = jnpr.LoadConfig("C:/Configs/juniper.txt", "set", false)
	if err != nil {
		fmt.Println(err)
	}
    
    err = jnpr.CommitCheck()
    if err != nil {
        fmt.Println(err)
    }
    
    err = jnpr.Commit()
    if err != nil {
        fmt.Println(err)
    }
    
    err = jnpr.Unlock()
    if err != nil {
        fmt.Println(err)
    }
}
```

[1]: https://tools.ietf.org/html/rfc6241
[2]: https://github.com/Juniper/go-netconf
[3]: http://www.juniper.net
[4]: https://godoc.org/github.com/scottdware/go-junos