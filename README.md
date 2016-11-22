## go-junos
[![GoDoc](https://godoc.org/github.com/scottdware/go-junos?status.svg)](https://godoc.org/github.com/scottdware/go-junos) [![Travis-CI](https://travis-ci.org/scottdware/go-junos.svg?branch=master)](https://travis-ci.org/scottdware/go-junos)
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/scottdware/go-junos/master/LICENSE)

A Go package that interacts with Junos devices, as well as Junos Space, and allows you to do the following:

* Run operational mode commands, such as `show`, `request`, etc..
* Compare the active configuration to a rollback configuration (diff).
* Rollback the configuration to a given state or a "rescue" config.
* Configure devices by submitting commands, uploading a local file or from a remote FTP/HTTP server.
* List files on devices.
* Commit operations: lock, unlock, commit, commit-at, commit-confirmed, etc.
* [SRX] Convert from a zone-based address book to a global one.
* [SRX] Create a site-to-site IPsec VPN.

Junos Space

* Get information from Junos Space managed devices.
* Add/remove devices from Junos Space.
* List all software image packages that are in Junos Space.
* Stage and deploy software images to devices from Junos Space.
* Create, edit and delete address and service objects/groups.
* Edit address and service groups by adding or removing objects to them.
* View all policies managed by Junos Space.
* Publish policies and update devices.
* Add/modify polymorphic (variable) objects.

### Examples & Documentation
Visit the [GoDoc][godoc-go-junos] page for package documentation and examples.

> **Note:** This package makes all of it's calls over [Netconf][netconf-rfc] using the [go-netconf][go-netconf] package from
 [Juniper Networks][juniper]. Please make sure you allow Netconf communication to your devices:
```
set system services netconf ssh
set security zones security-zone <xxx> interfaces <xxx> host-inbound-traffic system-services netconf
```

[netconf-rfc]: https://tools.ietf.org/html/rfc6241
[go-netconf]: https://github.com/Juniper/go-netconf
[juniper]: http://www.juniper.net
[godoc-go-junos]: https://godoc.org/github.com/scottdware/go-junos
[license]: https://github.com/scottdware/go-junos/blob/master/LICENSE