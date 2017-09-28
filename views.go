package junos

import (
	"encoding/xml"
	"errors"
	"strings"

	"github.com/Juniper/go-netconf/netconf"
)

// ArpTable contains the ARP table on the device.
type ArpTable struct {
	Count   int        `xml:"arp-entry-count"`
	Entries []ArpEntry `xml:"arp-table-entry"`
}

// ArpEntry holds each individual ARP entry.
type ArpEntry struct {
	MACAddress string `xml:"mac-address"`
	IPAddress  string `xml:"ip-address"`
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

// EthernetSwitchingTable contains the ethernet-switching table on the device.
type EthernetSwitchingTable struct {
	Entries []L2MACEntry `xml:"l2ng-l2ald-mac-entry-vlan"`
}

// L2MACEntry contains information about every MAC address on each VLAN.
type L2MACEntry struct {
	GlobalMACCount  int        `xml:"mac-count-global"`
	LearnedMACCount int        `xml:"learnt-mac-count"`
	RoutingInstance string     `xml:"l2ng-l2-mac-routing-instance"`
	VlanID          int        `xml:"l2ng-l2-vlan-id"`
	MACEntries      []MACEntry `xml:"l2ng-mac-entry"`
}

// MACEntry contains information about each individual MAC address. Flags are: S - static MAC, D - dynamic MAC,
// L - locally learned, P - persistent static, SE - statistics enabled, NM - non configured MAC, R - remote PE MAC,
// O - ovsdb MAC.
type MACEntry struct {
	VlanName         string `xml:"l2ng-l2-mac-vlan-name"`
	MACAddress       string `xml:"l2ng-l2-mac-address"`
	Age              string `xml:"l2ng-l2-mac-age"`
	Flags            string `xml:"l2ng-l2-mac-flags"`
	LogicalInterface string `xml:"l2ng-l2-mac-logical-interface"`
}

// HardwareInventory contains all the hardware information about the device.
type HardwareInventory struct {
	Chassis []Chassis `xml:"chassis"`
}

type srxHardwareInventory struct {
	Chassis []Chassis `xml:"multi-routing-engine-item>chassis-inventory>chassis"`
}

// Chassis contains all of the hardware information for each chassis, such as a clustered pair of SRX's or a
// virtual-chassis configuration.
type Chassis struct {
	Name         string   `xml:"name"`
	SerialNumber string   `xml:"serial-number"`
	Description  string   `xml:"description"`
	Modules      []Module `xml:"chassis-module"`
}

// Module contains information about each individual module.
type Module struct {
	Name         string      `xml:"name"`
	Version      string      `xml:"version,omitempty"`
	PartNumber   string      `xml:"part-number"`
	SerialNumber string      `xml:"serial-number"`
	Description  string      `xml:"description"`
	CLEICode     string      `xml:"clei-code"`
	ModuleNumber string      `xml:"module-number"`
	SubModules   []SubModule `xml:"chassis-sub-module"`
}

// SubModule contains information about each individual sub-module.
type SubModule struct {
	Name          string         `xml:"name"`
	Version       string         `xml:"version,omitempty"`
	PartNumber    string         `xml:"part-number"`
	SerialNumber  string         `xml:"serial-number"`
	Description   string         `xml:"description"`
	CLEICode      string         `xml:"clei-code"`
	ModuleNumber  string         `xml:"module-number"`
	SubSubModules []SubSubModule `xml:"chassis-sub-sub-module"`
}

// SubSubModule contains information about each sub-sub module, such as SFP's.
type SubSubModule struct {
	Name             string            `xml:"name"`
	Version          string            `xml:"version,omitempty"`
	PartNumber       string            `xml:"part-number"`
	SerialNumber     string            `xml:"serial-number"`
	Description      string            `xml:"description"`
	SubSubSubModules []SubSubSubModule `xml:"chassis-sub-sub-sub-module"`
}

// SubSubSubModule contains information about each sub-sub-sub module, such as SFP's on a
// PIC, which is tied to a MIC on an MX.
type SubSubSubModule struct {
	Name         string `xml:"name"`
	Version      string `xml:"version,omitempty"`
	PartNumber   string `xml:"part-number"`
	SerialNumber string `xml:"serial-number"`
	Description  string `xml:"description"`
}

// VirtualChassis contains information regarding the virtual-chassis setup for the device.
type VirtualChassis struct {
	PreProvisionedVCID   string     `xml:"preprovisioned-virtual-chassis-information>virtual-chassis-id"`
	PreProvisionedVCMode string     `xml:"preprovisioned-virtual-chassis-information>virtual-chassis-mode"`
	Members              []VCMember `xml:"member-list>member"`
}

// VCMember contains information about each individual virtual-chassis member.
type VCMember struct {
	Status       string             `xml:"member-status"`
	ID           int                `xml:"member-id"`
	FPCSlot      string             `xml:"fpc-slot"`
	SerialNumber string             `xml:"member-serial-number"`
	Model        string             `xml:"member-model"`
	Priority     int                `xml:"member-priority"`
	MixedMode    string             `xml:"member-mixed-mode"`
	RouteMode    string             `xml:"member-route-mode"`
	Role         string             `xml:"member-role"`
	Neighbors    []VCMemberNeighbor `xml:"neighbor-list>neighbor"`
}

// VCMemberNeighbor contains information about each virtual-chassis member neighbor.
type VCMemberNeighbor struct {
	ID        int    `xml:"neighbor-id"`
	Interface string `xml:"neighbor-interface"`
}

// BGPTable contains information about every BGP peer configured on the device.
type BGPTable struct {
	TotalGroups int       `xml:"group-count"`
	TotalPeers  int       `xml:"peer-count"`
	DownPeers   int       `xml:"down-peer-count"`
	Entries     []BGPPeer `xml:"bgp-peer"`
}

// BGPPeer contains information about each individual BGP peer.
type BGPPeer struct {
	Address            string `xml:"peer-address"`
	ASN                int    `xml:"peer-as"`
	InputMessages      int    `xml:"input-messages"`
	OutputMessages     int    `xml:"output-messages"`
	QueuedRoutes       int    `xml:"route-queue-count"`
	Flaps              int    `xml:"flap-count"`
	ElapsedTime        string `xml:"elapsed-time"`
	State              string `xml:"peer-state"`
	RoutingTable       string `xml:"bgp-rib>name"`
	ActivePrefixes     int    `xml:"bgp-rib>active-prefix-count"`
	ReceivedPrefixes   int    `xml:"bgp-rib>received-prefix-count"`
	AcceptedPrefixes   int    `xml:"bgp-rib>accepted-prefix-count"`
	SuppressedPrefixes int    `xml:"bgp-rib>suppressed-prefix-count"`
}

// StaticNats contains the static NATs configured on the device.
type StaticNats struct {
	Count   int              `xml:"total-static-nat-rules>total-rules"`
	Entries []StaticNatEntry `xml:"static-nat-rule-entry"`
}

// srxStaticNats contains static NATs configured across a clustered-mode SRX
type srxStaticNats struct {
	Count   int              `xml:"multi-routing-engine-item>static-nat-rule-information>total-static-nat-rules>total-rules"`
	Entries []StaticNatEntry `xml:"multi-routing-engine-item>static-nat-rule-information"`
}

// StaticNatEntry holds each individual static NAT entry.
type StaticNatEntry struct {
	Name       string `xml:"rule-name"`
	SetName    string `xml:"rule-set-name"`
	ID         string `xml:"rule-id"`
	FromZone   string `xml:"rule-from-context-name"`
	FakePrefix string `xml:"rule-destination-address-prefix"`
	RealPrefix string `xml:"rule-host-address-prefix"`
	Mask       string `xml:"rule-address-netmask"`
	Hits       string `xml:"succ-hits"`
}

// Views contains the information for the specific views. Note that some views aren't available for specific
// hardware platforms, such as the "VirtualChassis" view on an SRX.
type Views struct {
	Arp            ArpTable
	Route          RoutingTable
	Interface      Interfaces
	Vlan           Vlans
	EthernetSwitch EthernetSwitchingTable
	Inventory      HardwareInventory
	VirtualChassis VirtualChassis
	BGP            BGPTable
	StaticNat      StaticNats
}

var (
	viewCategories = map[string]string{
		"arp":            "<get-arp-table-information><no-resolve/></get-arp-table-information>",
		"route":          "<get-route-information/>",
		"interface":      "<get-interface-information/>",
		"vlan":           "<get-vlan-information/>",
		"ethernetswitch": "<get-ethernet-switching-table-information/>",
		"inventory":      "<get-chassis-inventory/>",
		"virtualchassis": "<get-virtual-chassis-information/>",
		"bgp":            "<get-bgp-summary-information/>",
		"staticnat":      "<get-static-nat-rule-information><all/></get-static-nat-rule-information>",
	}
)

func validatePlatform(j *Junos, v string) error {
	switch v {
	case "ethernetswitch":
		if strings.Contains(j.Platform[0].Model, "SRX") || strings.Contains(j.Platform[0].Model, "MX") {
			return errors.New("ethernet-switching information is not available on this platform")
		}
	case "virtualchassis":
		if strings.Contains(j.Platform[0].Model, "SRX") || strings.Contains(j.Platform[0].Model, "MX") {
			return errors.New("virtual-chassis information is not available on this platform")
		}
	}

	return nil
}

// Views gathers information on the device given the "view" specified. These views can be interrated/looped over to view the
// data (i.e. ARP table entries, interface details/statistics, routing tables, etc.). Supported views are:
// arp, route, interface, vlan, ethernetswitch, inventory, staticnat.
func (j *Junos) Views(view string) (*Views, error) {
	var results Views

	if strings.Contains(j.Platform[0].Model, "SRX") || strings.Contains(j.Platform[0].Model, "MX") {
		err := validatePlatform(j, view)

		if err != nil {
			return nil, err
		}
	}

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
	case "ethernetswitch":
		var ethtable EthernetSwitchingTable
		formatted := strings.Replace(reply.Data, "\n", "", -1)

		if err := xml.Unmarshal([]byte(formatted), &ethtable); err != nil {
			return nil, err
		}

		results.EthernetSwitch = ethtable
	case "inventory":
		var inventory HardwareInventory
		formatted := strings.Replace(reply.Data, "\n", "", -1)

		if strings.Contains(reply.Data, "multi-routing-engine-results") {
			var srxinventory srxHardwareInventory

			if err := xml.Unmarshal([]byte(formatted), &srxinventory); err != nil {
				return nil, err
			}

			for _, c := range srxinventory.Chassis {
				inventory.Chassis = append(inventory.Chassis, c)
			}

			results.Inventory = inventory
		} else {
			if err := xml.Unmarshal([]byte(formatted), &inventory); err != nil {
				return nil, err
			}

			results.Inventory = inventory
		}
	case "virtualchassis":
		var vc VirtualChassis
		formatted := strings.Replace(reply.Data, "\n", "", -1)

		if err := xml.Unmarshal([]byte(formatted), &vc); err != nil {
			return nil, err
		}

		results.VirtualChassis = vc
	case "bgp":
		var bgpTable BGPTable
		formatted := strings.Replace(reply.Data, "\n", "", -1)

		if err := xml.Unmarshal([]byte(formatted), &bgpTable); err != nil {
			return nil, err
		}

		results.BGP = bgpTable
	case "staticnat":
		var staticnats StaticNats
		formatted := strings.Replace(reply.Data, "\n", "", -1)

		if strings.Contains(reply.Data, "multi-routing-engine-results") {
			var srxstaticnats srxStaticNats

			if err := xml.Unmarshal([]byte(formatted), &srxstaticnats); err != nil {
				return nil, err
			}

			for _, c := range srxstaticnats.Entries {
				staticnats.Entries = append(staticnats.Entries, c)
			}

			results.StaticNat = staticnats
		} else {
			if err := xml.Unmarshal([]byte(formatted), &staticnats); err != nil {
				return nil, err
			}

			results.StaticNat = staticnats
		}
	}

	return &results, nil
}
