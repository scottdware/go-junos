## go-junos
[![GoDoc](https://godoc.org/github.com/scottdware/go-junos?status.svg)](https://godoc.org/github.com/scottdware/go-junos) [![Travis-CI](https://travis-ci.org/scottdware/go-junos.svg?branch=master)](https://travis-ci.org/scottdware/go-junos)
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/scottdware/go-junos/master/LICENSE)

A Go package that interacts with Junos devices, as well as Junos Space, and allows you to do the following:

* Run operational mode commands, such as `show`, `request`, etc..
* Compare the active configuration to a given rollback config.
* Rollback the configuration to a given state or a "rescue" config.
* Configure devices by uploading a local file or from an FTP/HTTP server.

Junos Space

* Get information from Junos Space managed devices.
* Add/remove devices from Junos Space.
* List all software image packages that are in Junos Space.
* Stage and deploy software images to devices from Junos Space.
* Add/remove address and service objects/groups.
* Modify objects (e.g. add/remove objects from groups, and delete objects).
* View all policies managed by Junos Space.
* Publish policies and update devices.
* Add/modify polymorphic (variable) objects.

### Examples & Documentation
Visit the [GoDoc][godoc-go-junos] page for package documentation and examples.

> **Note:** This package makes all of it's calls over [Netconf][netconf-rfc] using the [go-netconf][go-netconf] package from
 [Juniper Networks][juniper]

[netconf-rfc]: https://tools.ietf.org/html/rfc6241
[go-netconf]: https://github.com/Juniper/go-netconf
[juniper]: http://www.juniper.net
[godoc-go-junos]: https://godoc.org/github.com/scottdware/go-junos
[license]: https://github.com/scottdware/go-junos/blob/master/LICENSE