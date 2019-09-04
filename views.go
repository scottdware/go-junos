package junos

import (
	"encoding/xml"
	"errors"
	"fmt"
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
	InterfaceType           string             `xml:"if-type"`
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
	CIDR               string `xml:"address-family>interface-address>ifa-destination"`
	IPAddress          string `xml:"address-family>interface-address>ifa-local"`
	LocalIndex         int    `xml:"local-index"`
	SNMPIndex          int    `xml:"snmp-index"`
	Encapsulation      string `xml:"encapsulation"`
	LAGInputPackets    uint64 `xml:"lag-traffic-statistics>lag-bundle>input-packets"`
	LAGInputPps        int    `xml:"lag-traffic-statistics>lag-bundle>input-pps"`
	LAGInputBytes      int    `xml:"lag-traffic-statistics>lag-bundle>input-bytes"`
	LAGInputBps        int    `xml:"lag-traffic-statistics>lag-bundle>input-bps"`
	LAGOutputPackets   uint64 `xml:"lag-traffic-statistics>lag-bundle>output-packets"`
	LAGOutputPps       int    `xml:"lag-traffic-statistics>lag-bundle>output-pps"`
	LAGOutputBytes     int    `xml:"lag-traffic-statistics>lag-bundle>output-bytes"`
	LAGOutputBps       int    `xml:"lag-traffic-statistics>lag-bundle>output-bps"`
	ZoneName           string `xml:"logical-interface-zone-name"`
	InputPackets       uint64 `xml:"traffic-statistics>input-packets"`
	OutputPackets      uint64 `xml:"traffic-statistics>output-packets"`
	AddressFamily      string `xml:"address-family>address-family-name"`
	AggregatedEthernet string `xml:"address-family>ae-bundle-name,omitempty"`
	LinkAddress        string `xml:"link-address,omitempty"`
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

type LLDPNeighbors struct {
	Entries []LLDPNeighbor `xml:"lldp-neighbor-information"`
}

type LLDPNeighbor struct {
	LocalPortId              string `xml:"lldp-local-port-id"`
	LocalParentInterfaceName string `xml:"lldp-local-parent-interface-name"`
	RemoteChassisIdSubtype   string `xml:"lldp-remote-chassis-id-subtype"`
	RemoteChassisId          string `xml:"lldp-remote-chassis-id"`
	RemotePortDescription    string `xml:"lldp-remote-port-description"`
	RemotePortId             string `xml:"lldp-remote-port-id"`
	RemoteSystemName         string `xml:"lldp-remote-system-name"`
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

// Storage contains information about all of the file systems on the device.
type Storage struct {
	Entries []SystemStorage `xml:"system-storage-information"`
}

type multiStorage struct {
	Entries []SystemStorage `xml:"multi-routing-engine-item>system-storage-information"`
}

// SystemStorage stores the file system information for each node, routing-engine, etc. on the device.
type SystemStorage struct {
	FileSystems []FileSystem `xml:"filesystem"`
}

// FileSystem contains the information for each partition.
type FileSystem struct {
	Name            string `xml:"filesystem-name"`
	TotalBlocks     int    `xml:"total-blocks"`
	UsedBlocks      int    `xml:"used-blocks"`
	AvailableBlocks int    `xml:"available-blocks"`
	UsedPercent     string `xml:"used-percent"`
	MountedOn       string `xml:"mounted-on"`
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
	Count   int
	Entries []StaticNatEntry `xml:"static-nat-rule-entry"`
}

// srxStaticNats contains static NATs configured across a clustered-mode SRX
type srxStaticNats struct {
	Entries []StaticNatEntry `xml:"multi-routing-engine-item>static-nat-rule-information>static-nat-rule-entry"`
}

// StaticNatEntry holds each individual static NAT entry.
type StaticNatEntry struct {
	Name                    string `xml:"rule-name"`
	SetName                 string `xml:"rule-set-name"`
	ID                      string `xml:"rule-id"`
	RuleMatchingPosition    int    `xml:"rule-matching-position"`
	FromContext             string `xml:"rule-from-context"`
	FromZone                string `xml:"rule-from-context-name"`
	SourceAddressLowRange   string `xml:"static-source-address-range-entry>rule-source-address-low-range"`
	SourceAddressHighRange  string `xml:"static-source-address-range-entry>rule-source-address-high-range"`
	DestinaionAddressPrefix string `xml:"rule-destination-address-prefix"`
	DestinationPortLow      int    `xml:"rule-destination-port-low"`
	DestinationPortHigh     int    `xml:"rule-destination-port-high"`
	HostAddressPrefix       string `xml:"rule-host-address-prefix"`
	HostPortLow             int    `xml:"rule-host-port-low"`
	HostPortHigh            int    `xml:"rule-host-port-high"`
	Netmask                 string `xml:"rule-address-netmask"`
	RoutingInstance         string `xml:"rule-host-routing-instance"`
	TranslationHits         int    `xml:"rule-translation-hits"`
	SuccessfulSessions      int    `xml:"succ-hits"`
	ConcurrentHits          int    `xml:"concurrent-hits"`
}

// SourceNats contains the source NATs configured on the device.
type SourceNats struct {
	Count   int
	Entries []SourceNatEntry `xml:"source-nat-rule-entry"`
}

type srxSourceNats struct {
	Entries []SourceNatEntry `xml:"multi-routing-engine-item>source-nat-rule-detail-information>source-nat-rule-entry"`
}

// SourceNatEntry holds each individual source NAT entry.
type SourceNatEntry struct {
	Name                     string   `xml:"rule-name"`
	SetName                  string   `xml:"rule-set-name"`
	ID                       string   `xml:"rule-id"`
	RuleMatchingPosition     int      `xml:"rule-matching-position"`
	FromContext              string   `xml:"rule-from-context"`
	FromZone                 string   `xml:"rule-from-context-name"`
	ToContext                string   `xml:"rule-to-context"`
	ToZone                   string   `xml:"rule-to-context-name"`
	SourceAddressLowRange    string   `xml:"source-address-range-entry>rule-source-address-low-range"`
	SourceAddressHighRange   string   `xml:"source-address-range-entryrule-source-address-high-range"`
	SourceAddresses          []string `xml:"source-address-range-entry>rule-source-address"`
	DestinationAddresses     []string `xml:"destination-address-range-entry>rule-destination-address"`
	DestinationPortLow       int      `xml:"destination-port-entry>rule-destination-port-low"`
	DestinationPortHigh      int      `xml:"destination-port-entry>rule-destination-port-high"`
	SourcePortLow            int      `xml:"source-port-entry>rule-source-port-low"`
	SourcePortHigh           int      `xml:"source-port-entry>rule-source-port-high"`
	SourceNatProtocol        string   `xml:"src-nat-protocol-entry"`
	RuleAction               string   `xml:"source-nat-rule-action-entry>source-nat-rule-action"`
	PersistentNatType        string   `xml:"source-nat-rule-action-entry>persistent-nat-type"`
	PersistentNatMappingType string   `xml:"source-nat-rule-action-entry>persistent-nat-mapping-type"`
	PersistentNatTimeout     int      `xml:"source-nat-rule-action-entry>persistent-nat-timeout"`
	PersistentNatMaxSession  int      `xml:"source-nat-rule-action-entry>persistent-nat-max-session"`
	TranslationHits          int      `xml:"source-nat-rule-hits-entry>rule-translation-hits"`
	SuccessfulSessions       int      `xml:"source-nat-rule-hits-entry>succ-hits"`
	ConcurrentHits           int      `xml:"source-nat-rule-hits-entry>concurrent-hits"`
}

// FirewallPolicy contains the entire firewall policy for the device.
type FirewallPolicy struct {
	XMLName xml.Name          `xml:"security-policies"`
	Entries []SecurityContext `xml:"security-context"`
}

type srxFirewallPolicy struct {
	Entries []SecurityContext `xml:"multi-routing-engine-item>security-policies>security-context"`
}

// SecurityContext contains the policies for each context, such as rules from trust to untrust zones.
type SecurityContext struct {
	SourceZone      string `xml:"context-information>source-zone-name"`
	DestinationZone string `xml:"context-information>destination-zone-name"`
	Rules           []Rule `xml:"policies>policy-information"`
}

// Rule contains each individual element that makes up a security policy rule.
type Rule struct {
	Name                 string   `xml:"policy-name"`
	State                string   `xml:"policy-state"`
	Identifier           int      `xml:"policy-identifier"`
	ScopeIdentifier      int      `xml:"scope-policy-identifier"`
	SequenceNumber       int      `xml:"policy-sequence-number"`
	SourceAddresses      []string `xml:"source-addresses>source-address>address-name"`
	DestinationAddresses []string `xml:"destination-addresses>destination-address>address-name"`
	Applications         []string `xml:"applications>application>application-name"`
	SourceIdentities     []string `xml:"source-identities>source-identity>role-name"`
	PolicyAction         string   `xml:"policy-action>action-type"`
	PolicyTCPOptions     struct {
		SYNCheck      string `xml:"policy-tcp-options-syn-check"`
		SequenceCheck string `xml:"policy-tcp-options-sequence-check"`
	} `xml:"policy-action>policy-tcp-options"`
}

// Views contains the information for the specific views. Note that some views aren't available for specific
// hardware platforms, such as the "VirtualChassis" view on an SRX.
type Views struct {
	Arp            ArpTable
	BGP            BGPTable
	EthernetSwitch EthernetSwitchingTable
	FirewallPolicy FirewallPolicy
	Interface      Interfaces
	Inventory      HardwareInventory
	LLDPNeighbors  LLDPNeighbors
	Route          RoutingTable
	SourceNat      SourceNats
	StaticNat      StaticNats
	Storage        Storage
	VirtualChassis VirtualChassis
	Vlan           Vlans
}

var (
	viewCategories = map[string]string{
		"arp":            "<get-arp-table-information><no-resolve/></get-arp-table-information>",
		"route":          "<get-route-information/>",
		"interface":      "<get-interface-information/>",
		"vlan":           "<get-vlan-information/>",
		"lldp":           "<get-lldp-neighbors-information/>",
		"ethernetswitch": "<get-ethernet-switching-table-information/>",
		"inventory":      "<get-chassis-inventory/>",
		"virtualchassis": "<get-virtual-chassis-information/>",
		"bgp":            "<get-bgp-summary-information/>",
		"staticnat":      "<get-static-nat-rule-information><all/></get-static-nat-rule-information>",
		"sourcenat":      "<get-source-nat-rule-sets-information><all/></get-source-nat-rule-sets-information>",
		"storage":        "<get-system-storage/>",
		"firewallpolicy": "<get-firewall-policies/>",
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

// View gathers information on the device given the "view" specified. These views can be interrated/looped over to view the
// data (i.e. ARP table entries, interface details/statistics, routing tables, etc.). Supported views are:
//
// arp, route, bgp, interface, vlan, ethernetswitch, inventory, virtualchassis, staticnat, sourcenat, storage, fireawllpolicy
//
// By default, the interface view will return all interfaces on the device. If you wish to only see a particular physical interface,
// and all logical interfaces underneath it, you can use the option parameter to specify the name of the interface, e.g.:
//
// View("interface", "ge-0/0/0")
func (j *Junos) View(view string, option ...string) (*Views, error) {
	var results Views
	var reply *netconf.RPCReply
	var err error

	if strings.Contains(j.Platform[0].Model, "SRX") || strings.Contains(j.Platform[0].Model, "MX") {
		err := validatePlatform(j, view)

		if err != nil {
			return nil, err
		}
	}

	if view == "interface" && len(option) > 0 {
		rpcIntName := fmt.Sprintf("<get-interface-information><interface-name>%s</interface-name></get-interface-information>", option[0])
		reply, err = j.Session.Exec(netconf.RawMethod(rpcIntName))
		if err != nil {
			return nil, err
		}
	} else {
		reply, err = j.Session.Exec(netconf.RawMethod(viewCategories[view]))
		if err != nil {
			return nil, err
		}
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

	case "lldp":
		var lldpNeighbors LLDPNeighbors
		formatted := strings.Replace(reply.Data, "\n", "", -1)

		if err := xml.Unmarshal([]byte(formatted), &lldpNeighbors); err != nil {
			return nil, err
		}

		results.LLDPNeighbors = lldpNeighbors
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

			actualrules := len(srxstaticnats.Entries) / 2
			staticnats.Count = actualrules

			for i, c := range srxstaticnats.Entries {
				if i < actualrules {
					staticnats.Entries = append(staticnats.Entries, c)
				}

				i++
			}

			results.StaticNat = staticnats
		} else {
			var staticnatentry StaticNatEntry
			if err := xml.Unmarshal([]byte(formatted), &staticnats); err != nil {
				return nil, err
			}

			actualrules := len(staticnats.Entries)
			staticnats.Count = actualrules

			staticnats.Entries = append(staticnats.Entries, staticnatentry)

			results.StaticNat = staticnats
		}
	case "sourcenat":
		var sourcenats SourceNats
		formatted := strings.Replace(reply.Data, "\n", "", -1)

		if strings.Contains(reply.Data, "multi-routing-engine-results") {
			var srxsourcenats srxSourceNats

			if err := xml.Unmarshal([]byte(formatted), &srxsourcenats); err != nil {
				return nil, err
			}

			actualrules := len(srxsourcenats.Entries) / 2
			sourcenats.Count = actualrules

			for i, c := range srxsourcenats.Entries {
				if i < actualrules {
					sourcenats.Entries = append(sourcenats.Entries, c)
				}

				i++
			}

			results.SourceNat = sourcenats
		} else {
			var sourcenatentry SourceNatEntry
			if err := xml.Unmarshal([]byte(formatted), &sourcenats); err != nil {
				return nil, err
			}

			actualrules := len(sourcenats.Entries)
			sourcenats.Count = actualrules

			sourcenats.Entries = append(sourcenats.Entries, sourcenatentry)

			results.SourceNat = sourcenats
		}
	case "storage":
		var storage Storage
		formatted := strings.Replace(reply.Data, "\n", "", -1)

		if strings.Contains(reply.Data, "multi-routing-engine-results") {
			var multistorage multiStorage

			if err := xml.Unmarshal([]byte(formatted), &multistorage); err != nil {
				return nil, err
			}

			for _, s := range multistorage.Entries {
				storage.Entries = append(storage.Entries, s)
			}

			results.Storage = storage
		} else {
			var sysstorage SystemStorage
			if err := xml.Unmarshal([]byte(formatted), &sysstorage); err != nil {
				return nil, err
			}

			storage.Entries = append(storage.Entries, sysstorage)

			results.Storage = storage
		}
	case "firewallpolicy":
		var fwpolicy FirewallPolicy
		formatted := strings.Replace(reply.Data, "\n", "", -1)

		if strings.Contains(reply.Data, "multi-routing-engine-results") {
			var multifwpolicy srxFirewallPolicy

			if err := xml.Unmarshal([]byte(formatted), &multifwpolicy); err != nil {
				return nil, err
			}

			for _, s := range multifwpolicy.Entries {
				fwpolicy.Entries = append(fwpolicy.Entries, s)
			}

			results.FirewallPolicy = fwpolicy
		} else {
			if err := xml.Unmarshal([]byte(formatted), &fwpolicy); err != nil {
				return nil, err
			}

			results.FirewallPolicy = fwpolicy
		}
	}

	return &results, nil
}
