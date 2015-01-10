## go-junos

[![GoDoc](https://godoc.org/github.com/scottdware/go-junos?status.svg)](https://godoc.org/github.com/scottdware/go-junos)

A Go package that interacts with Junos devices and allows you to do the following:

* Run operational mode commands, such as `show`, `request`, etc..
* Compare active configuration to a given rollback config.
* Rollback the configuration to a given state or a "rescue" config.
* Configure devices by uploading a file with commands in "set" format.

Visit the [GoDoc][4] page for example usage and documentation.

This package makes all of it's calls over [Netconf][1] using the [go-netconf][2] package from 
[Juniper Networks][3]

[1]: https://tools.ietf.org/html/rfc6241
[2]: https://github.com/Juniper/go-netconf
[3]: http://www.juniper.net
[4]: https://godoc.org/github.com/scottdware/go-junos