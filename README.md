## go-junos
[![GoDoc](https://godoc.org/github.com/scottdware/go-junos?status.svg)](https://godoc.org/github.com/scottdware/go-junos) [![Travis-CI](https://travis-ci.org/scottdware/go-junos.svg?branch=master)](https://travis-ci.org/scottdware/go-junos)
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/scottdware/go-junos/master/LICENSE)

A Go package that interacts with Junos devices, as well as Junos Space, and allows you to do the following:

* Run operational mode commands, such as `show`, `request`, etc..
* Compare the active configuration to a rollback configuration (diff).
* Rollback the configuration to a given state or a "rescue" config.
* Configure devices by submitting commands, uploading a local file or from a remote FTP/HTTP server.
* Commit operations: lock, unlock, commit, commit-at, commit-confirmed, etc.
* [Device "views"][views] - This will allow you to quickly get all the information on the device given a specific view.
* [SRX] Convert from a zone-based address book to a global one.

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

### Installation
`go get -u github.com/scottdware/go-junos`

> **Note:** This package makes all of it's calls over [Netconf][netconf-rfc] using the [go-netconf][go-netconf] package from
 [Juniper Networks][juniper]. Please make sure you allow Netconf communication to your devices:
```
set system services netconf ssh
set security zones security-zone <xxx> interfaces <xxx> host-inbound-traffic system-services netconf
```

### Examples & Documentation
Visit the [GoDoc][godoc-go-junos] page for package documentation and examples.

Connect to a device, and view the current config to rollback 1.
```Go
jnpr, err := junos.NewSession("qfx-switch.company.com", "admin", "Juniper123!")
if err != nil {
    fmt.Println(err)
}

defer jnpr.Close()

diff, err := jnpr.Diff(1)
if err != nil {
    fmt.Println(err)
}

fmt.Println(diff)

// Will output the following

[edit vlans]
-   zzz-Test {
-       vlan-id 999;
-   }
-   zzz-Test2 {
-       vlan-id 1000;
-   }
```

View the routing-instance configuration.
```Go
jnpr, err := junos.NewSession("srx.company.com", "admin", "Juniper123!")
if err != nil {
    fmt.Println(err)
}

defer jnpr.Close()

riConfig, err := jnpr.GetConfig("text", "routing-instances")
if err != nil {
    fmt.Println(err)
}

fmt.Println(riConfig)

// Will output the following

## Last changed: 2017-03-24 12:26:58 EDT
routing-instances {
    default-ri {
        instance-type virtual-router;
        interface lo0.0;
        interface reth1.0;
        routing-options {
            static {
                route 0.0.0.0/0 next-hop 10.1.1.1;
            }
        }
    }
}
```

### Views
Device views allow you to quickly gather information regarding a specific "view", so that you may use that information
however you wish. A good example, is using the "interface" view to gather all of the interface information on the device,
then iterate over that view to see statistics, interface settings, etc.

> **Note:** Some of the views aren't available for all platforms, such as the `ethernetswitch` and `virtualchassis` on an SRX or MX.

Current out-of-the-box built-in views are:
* `arp`: Gathers the ARP table on the device.
    * Same as `show arp`
* `route`: Gathers the routing table (for all routing-instances) on the device.
    * Same as `show route`
* `interface`: Gathers the interface listing as well as statistics and settings on the device.
    * Same as `show interfaces`
* `vlan`: Gathers the VLAN table on the device.
    * Same as `show vlans`
* `ethernetswitch`: Gathers the ethernet-switching table on the device.
    * Same as `show ethernet-switching table`
* `inventory`: Gathers the physical inventory on the device.
    * Same as `show chassis hardware`
* `virtualchassis`: Gathers the virtual-chassis configuration on the device.
    * Same as `show virtual-chassis status`

You can even create your own views by creating a `struct` that models the XML output from using the `GetConfig()` function. Granted,
this is a little more work, and requires you to know a bit more about the Go language (such as unmarshalling XML), but if there's a custom
view that you want to see, it's possible to do this for anything you want.

I will be adding more views over time, but feel free to request ones you'd like to see by [emailing](mailto:scottdware@gmail.com) me, or drop
me a line on [Twitter](https://twitter.com/scottdware).

[netconf-rfc]: https://tools.ietf.org/html/rfc6241
[go-netconf]: https://github.com/Juniper/go-netconf
[juniper]: http://www.juniper.net
[godoc-go-junos]: https://godoc.org/github.com/scottdware/go-junos
[license]: https://github.com/scottdware/go-junos/blob/master/LICENSE
[views]: https://github.com/scottdware/go-junos#views
