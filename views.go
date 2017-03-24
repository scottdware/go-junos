package junos

import (
	"encoding/xml"
	"errors"
	"strings"

	"github.com/Juniper/go-netconf/netconf"
)

// ArpTable contains the ARP table on the device.
type ArpTable struct {
	Count   string     `xml:"arp-entry-count"`
	Entries []ArpEntry `xml:"arp-table-entry"`
}

// ArpEntry holds each individual ARP entry.
type ArpEntry struct {
	MACAddress string `xml:"mac-address"`
	IPAddress  string `xml:"ip-address"`
	Hostname   string `xml:"hostname"`
	Interface  string `xml:"interface-name"`
}

// RoutingTable contains every routing table on the device.
type RoutingTable struct {
	RouteTables []RouteTable `xml:"route-table"`
}

// RouteTable holds all the route information for each table.
type RouteTable struct {
	Name           string  `xml:"table-name"`
	TotalRoutes    int     `xml:"total-route-count"`
	ActiveRoutes   int     `xml:"active-route-count"`
	HolddownRoutes int     `xml:"holddown-route-count"`
	HiddenRoutes   int     `xml:"hidden-routes"`
	Entries        []Route `xml:"rt"`
}

// Route holds information about each individual route.
type Route struct {
	Destination           string `xml:"rt-destination"`
	Active                string `xml:"rt-entry>active-tag"`
	Protocol              string `xml:"rt-entry>protocol-name"`
	Preference            int    `xml:"rt-entry>preference"`
	Age                   string `xml:"rt-entry>age"`
	NextHop               string `xml:"rt-entry>nh>to,omitempty"`
	NextHopInterface      string `xml:"rt-entry>nh>via,omitempty"`
	NextHopTable          string `xml:"rt-entry>nh>nh-table,omitempty"`
	NextHopLocalInterface string `xml:"rt-entry>nh>nh-local-interface,omitempty"`
}

// Interfaces contains information about every interface on the device.
type Interfaces struct {
	Entries []PhysicalInterface `xml:"physical-interface"`
}

// PhysicalInterface contains information about each individual physical interface.
type PhysicalInterface struct {
	Name                    string             `xml:"name"`
	AdminStatus             string             `xml:"admin-status"`
	OperStatus              string             `xml:"oper-status"`
	LocalIndex              int                `xml:"local-index"`
	SNMPIndex               int                `xml:"snmp-index"`
	LinkLevelType           string             `xml:"link-level-type"`
	MTU                     string             `xml:"mtu"`
	LinkMode                string             `xml:"link-mode"`
	Speed                   string             `xml:"speed"`
	FlowControl             string             `xml:"if-flow-control"`
	AutoNegotiation         string             `xml:"if-auto-negotiation"`
	HardwarePhysicalAddress string             `xml:"hardware-physical-address"`
	Flapped                 string             `xml:"interface-flapped"`
	InputBps                int                `xml:"traffic-statistics>input-bps"`
	InputPps                int                `xml:"traffic-statistics>input-pps"`
	OutputBps               int                `xml:"traffic-statistics>output-bps"`
	OutputPps               int                `xml:"traffic-statistics>output-pps"`
	LogicalInterfaces       []LogicalInterface `xml:"logical-interface"`
}

// LogicalInterface contains information about the logical interfaces tied to a physical interface.
type LogicalInterface struct {
	Name               string `xml:"name"`
	MTU                string `xml:"address-family>mtu"`
	IPAddress          string `xml:"address-family>interface-address>ifa-local"`
	LocalIndex         int    `xml:"local-index"`
	SNMPIndex          int    `xml:"snmp-index"`
	Encapsulation      string `xml:"encapsulation"`
	LAGInputPackets    int    `xml:"lag-traffic-statistics>lag-bundle>input-packets"`
	LAGInputPps        int    `xml:"lag-traffic-statistics>lag-bundle>input-pps"`
	LAGInputBytes      int    `xml:"lag-traffic-statistics>lag-bundle>input-bytes"`
	LAGInputBps        int    `xml:"lag-traffic-statistics>lag-bundle>input-bps"`
	LAGOutputPackets   int    `xml:"lag-traffic-statistics>lag-bundle>output-packets"`
	LAGOutputPps       int    `xml:"lag-traffic-statistics>lag-bundle>output-pps"`
	LAGOutputBytes     int    `xml:"lag-traffic-statistics>lag-bundle>output-bytes"`
	LAGOutputBps       int    `xml:"lag-traffic-statistics>lag-bundle>output-bps"`
	ZoneName           string `xml:"logical-interface-zone-name"`
	InputPackets       int    `xml:"traffic-statistics>input-packets"`
	OutputPackets      int    `xml:"traffic-statistics>output-packets"`
	AddressFamily      string `xml:"address-family>address-family-name"`
	AggregatedEthernet string `xml:"address-family>ae-bundle-name,omitempty"`
}

// Vlans contains all of the VLAN information on the device.
type Vlans struct {
	Entries []Vlan `xml:"l2ng-l2ald-vlan-instance-group"`
}

// Vlan contains information about each individual VLAN.
type Vlan struct {
	Name             string   `xml:"l2ng-l2rtb-vlan-name"`
	Tag              int      `xml:"l2ng-l2rtb-vlan-tag"`
	MemberInterfaces []string `xml:"l2ng-l2rtb-vlan-member>l2ng-l2rtb-vlan-member-interface"`
}

// Views contains the information for the specific views.
type Views struct {
	Arp       ArpTable
	Route     RoutingTable
	Interface Interfaces
	Vlan      Vlans
}

var (
	viewCategories = map[string]string{
		"arp":       "<get-arp-table-information><no-resolve/></get-arp-table-information>",
		"route":     "<get-route-information/>",
		"interface": "<get-interface-information/>",
		"vlan":      "<get-vlan-information/>",
	}
)

// Views gathers information on the device given the "view" specified. These views can be interrated/looped over to view the
// information (i.e. ARP table entries, interface details/statistics, routing tables, etc.). Supported views are:
// arp, route, interface, vlan.
func (j *Junos) Views(view string) (*Views, error) {
	var results Views

	reply, err := j.Session.Exec(netconf.RawMethod(viewCategories[view]))
	if err != nil {
		return nil, err
	}

	if reply.Errors != nil {
		for _, m := range reply.Errors {
			return nil, errors.New(m.Message)
		}
	}

	if reply.Data == "" {
		return nil, errors.New("no output available - please check the syntax of your command")
	}

	switch view {
	case "arp":
		var arpTable ArpTable
		formatted := strings.Replace(reply.Data, "\n", "", -1)

		if err := xml.Unmarshal([]byte(formatted), &arpTable); err != nil {
			return nil, err
		}

		results.Arp = arpTable
	case "route":
		var routingTable RoutingTable
		formatted := strings.Replace(reply.Data, "\n", "", -1)

		if err := xml.Unmarshal([]byte(formatted), &routingTable); err != nil {
			return nil, err
		}

		results.Route = routingTable
	case "interface":
		var ints Interfaces
		formatted := strings.Replace(reply.Data, "\n", "", -1)

		if err := xml.Unmarshal([]byte(formatted), &ints); err != nil {
			return nil, err
		}

		results.Interface = ints
	case "vlan":
		var vlan Vlans
		formatted := strings.Replace(reply.Data, "\n", "", -1)

		if err := xml.Unmarshal([]byte(formatted), &vlan); err != nil {
			return nil, err
		}

		results.Vlan = vlan
	}

	return &results, nil
}
