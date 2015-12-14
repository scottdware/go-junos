package junos

import (
	"encoding/xml"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// FirewallPolicy contains all of the rules that will be created for the policy.
type FirewallPolicy struct {
	Addresses     []string
	Applications  []string
	JunosDefaults []string
	Rules         []Rule
	AppConfig     []string
}

// Rule contains information about a single firewall/policy rule.
type Rule struct {
	Name          interface{}
	SourceZone    string
	SourceAddress []string
	DestZone      string
	DestAddress   []string
	Application   []string
	Action        string
}

// ExistingAddresses contains information about every global address-book entry.
type ExistingAddresses struct {
	XMLName             xml.Name             `xml:"configuration"`
	ExistingAddresses   []ExistingAddress    `xml:"security>address-book>address"`
	ExistingAddressSets []ExistingAddressSet `xml:"security>address-book>address-set"`
}

// ExistingAddress contains information about each individual address-book entry.
type ExistingAddress struct {
	Name     string `xml:"name"`
	IP       string `xml:"ip-prefix,omitempty"`
	DNSName  string `xml:"dns-name>name,omitempty"`
	Wildcard string `xml:"wildcard-address>name,omitempty"`
}

// ExistingAddressSet contains all of the address-sets (groups) in the address-book.
type ExistingAddressSet struct {
	Name              string            `xml:"name"`
	ExistingAddresses []ExistingAddress `xml:"address"`
}

// ExistingApplications contains information about every application entry.
type ExistingApplications struct {
	XMLName                 xml.Name                 `xml:"configuration"`
	ExistingApplications    []ExistingApplication    `xml:"applications>application"`
	ExistingApplicationSets []ExistingApplicationSet `xml:"applications>application-set"`
}

// ExistingApplication contains information about each individual application entry.
type ExistingApplication struct {
	Name string `xml:"name"`
}

// ExistingApplicationSet contains all of the application-sets (service groups) in the SRX.
type ExistingApplicationSet struct {
	Name                 string                `xml:"name"`
	ExistingApplications []ExistingApplication `xml:"application-set>application"`
}

// SecurityZones contains all of our security-zone information.
type SecurityZones struct {
	XMLName xml.Name `xml:"configuration"`
	Zones   []Zone   `xml:"security>zones>security-zone"`
}

// Zone contains information about each individual security-zone.
type Zone struct {
	Name           string          `xml:"name"`
	AddressEntries []AddressEntry  `xml:"address-book>address"`
	AddressSets    []AddressSet    `xml:"address-book>address-set"`
	ZoneInterfaces []ZoneInterface `xml:"interfaces"`
}

// AddressEntry contains information about each individual address-book entry.
type AddressEntry struct {
	Name     string `xml:"name"`
	IP       string `xml:"ip-prefix,omitempty"`
	DNSName  string `xml:"dns-name>name,omitempty"`
	Wildcard string `xml:"wildcard-address>name,omitempty"`
}

// AddressSet contains all of the address-sets (groups) in the address-book.
type AddressSet struct {
	Name           string         `xml:"name"`
	AddressEntries []AddressEntry `xml:"address"`
}

// ZoneInterface contains a list of all interfaces that belong to the zone.
type ZoneInterface struct {
	Name string `xml:"name"`
}

// IPsecVPN contains the necessary information when creating a new site-to-site VPN.
type IPsecVPN struct {
	Name              string
	Local             string
	Peer              string
	ExternalInterface string
	Zone              string
	Mode              string
	PSK               string
	PFS               int
	Establish         string
	St0               string
	Gateway           []string
	P1Proposals       []P1
	P2Proposals       []P2
	TrafficSelectors  []string
}

// P1 contains any IKE phase 1 proposal information.
type P1 struct {
	P1Name           string
	P1DiffeHellman   int
	P1Authentication string
	P1Encryption     string
	P1Seconds        int
}

// P2 contains any IKE phase 2 proposal information.
type P2 struct {
	P2Name           string
	P2Authentication string
	P2Encryption     string
	P2Seconds        int
	P2Protocol       string
}

// st0Interface holds all of the current st0 interfaces on the SRX.
type st0Interface struct {
	XMLName xml.Name `xml:"interface-information"`
	Units   []unit   `xml:"physical-interface>logical-interface"`
}

// unit contains each individual st0.<unit> name.
type unit struct {
	Name string `xml:"name"`
}

var (
	// These are the Junos default applications in an SRX
	junosDefaultApps = []string{
		"junos-aol",
		"junos-bgp",
		"junos-biff",
		"junos-bootpc",
		"junos-bootps",
		"junos-chargen",
		"junos-cifs",
		"junos-cvspserver",
		"junos-dhcp-client",
		"junos-dhcp-relay",
		"junos-dhcp-server",
		"junos-discard",
		"junos-dns-tcp",
		"junos-dns-udp",
		"junos-echo",
		"junos-finger",
		"junos-ftp",
		"junos-gnutella",
		"junos-gopher",
		"junos-gre",
		"junos-gtp",
		"junos-h323",
		"junos-http",
		"junos-http-ext",
		"junos-https",
		"junos-icmp-all",
		"junos-icmp-ping",
		"junos-icmp6-all",
		"junos-icmp6-dst-unreach-addr",
		"junos-icmp6-dst-unreach-admin",
		"junos-icmp6-dst-unreach-beyond",
		"junos-icmp6-dst-unreach-port",
		"junos-icmp6-dst-unreach-route",
		"junos-icmp6-echo-reply",
		"junos-icmp6-echo-request",
		"junos-icmp6-packet-too-big",
		"junos-icmp6-param-prob-header",
		"junos-icmp6-param-prob-nexthdr",
		"junos-icmp6-param-prob-option",
		"junos-icmp6-time-exceed-reassembly",
		"junos-icmp6-time-exceed-transit",
		"junos-ident",
		"junos-ike",
		"junos-ike-nat",
		"junos-imap",
		"junos-imaps",
		"junos-internet-locator-service",
		"junos-irc",
		"junos-l2tp",
		"junos-ldap",
		"junos-ldp-tcp",
		"junos-ldp-udp",
		"junos-lpr",
		"junos-mail",
		"junos-mgcp",
		"junos-mgcp-ca",
		"junos-mgcp-ua",
		"junos-ms-rpc",
		"junos-ms-rpc-any",
		"junos-ms-rpc-epm",
		"junos-ms-rpc-iis-com",
		"junos-ms-rpc-iis-com-1",
		"junos-ms-rpc-iis-com-adminbase",
		"junos-ms-rpc-msexchange",
		"junos-ms-rpc-msexchange-directory-nsp",
		"junos-ms-rpc-msexchange-directory-rfr",
		"junos-ms-rpc-msexchange-info-store",
		"junos-ms-rpc-tcp",
		"junos-ms-rpc-udp",
		"junos-ms-rpc-uuid-any-tcp",
		"junos-ms-rpc-uuid-any-udp",
		"junos-ms-rpc-wmic",
		"junos-ms-rpc-wmic-admin",
		"junos-ms-rpc-wmic-admin2",
		"junos-ms-rpc-wmic-mgmt",
		"junos-ms-rpc-wmic-webm-callresult",
		"junos-ms-rpc-wmic-webm-classobject",
		"junos-ms-rpc-wmic-webm-level1login",
		"junos-ms-rpc-wmic-webm-login-clientid",
		"junos-ms-rpc-wmic-webm-login-helper",
		"junos-ms-rpc-wmic-webm-objectsink",
		"junos-ms-rpc-wmic-webm-refreshing-services",
		"junos-ms-rpc-wmic-webm-remote-refresher",
		"junos-ms-rpc-wmic-webm-services",
		"junos-ms-rpc-wmic-webm-shutdown",
		"junos-ms-sql",
		"junos-msn",
		"junos-nbds",
		"junos-nbname",
		"junos-netbios-session",
		"junos-nfs",
		"junos-nfsd-tcp",
		"junos-nfsd-udp",
		"junos-nntp",
		"junos-ns-global",
		"junos-ns-global-pro",
		"junos-nsm",
		"junos-ntalk",
		"junos-ntp",
		"junos-ospf",
		"junos-pc-anywhere",
		"junos-persistent-nat",
		"junos-ping",
		"junos-pingv6",
		"junos-pop3",
		"junos-pptp",
		"junos-printer",
		"junos-r2cp",
		"junos-radacct",
		"junos-radius",
		"junos-realaudio",
		"junos-rip",
		"junos-routing-inbound",
		"junos-rsh",
		"junos-rtsp",
		"junos-sccp",
		"junos-sctp-any",
		"junos-sip",
		"junos-smb",
		"junos-smb-session",
		"junos-smtp",
		"junos-snmp-agentx",
		"junos-snpp",
		"junos-sql-monitor",
		"junos-sqlnet-v1",
		"junos-sqlnet-v2",
		"junos-ssh",
		"junos-stun",
		"junos-sun-rpc",
		"junos-sun-rpc-any",
		"junos-sun-rpc-any-tcp",
		"junos-sun-rpc-any-udp",
		"junos-sun-rpc-mountd",
		"junos-sun-rpc-mountd-tcp",
		"junos-sun-rpc-mountd-udp",
		"junos-sun-rpc-nfs",
		"junos-sun-rpc-nfs-access",
		"junos-sun-rpc-nfs-tcp",
		"junos-sun-rpc-nfs-udp",
		"junos-sun-rpc-nlockmgr",
		"junos-sun-rpc-nlockmgr-tcp",
		"junos-sun-rpc-nlockmgr-udp",
		"junos-sun-rpc-portmap",
		"junos-sun-rpc-portmap-tcp",
		"junos-sun-rpc-portmap-udp",
		"junos-sun-rpc-rquotad",
		"junos-sun-rpc-rquotad-tcp",
		"junos-sun-rpc-rquotad-udp",
		"junos-sun-rpc-ruserd",
		"junos-sun-rpc-ruserd-tcp",
		"junos-sun-rpc-ruserd-udp",
		"junos-sun-rpc-sadmind",
		"junos-sun-rpc-sadmind-tcp",
		"junos-sun-rpc-sadmind-udp",
		"junos-sun-rpc-sprayd",
		"junos-sun-rpc-sprayd-tcp",
		"junos-sun-rpc-sprayd-udp",
		"junos-sun-rpc-status",
		"junos-sun-rpc-status-tcp",
		"junos-sun-rpc-status-udp",
		"junos-sun-rpc-tcp",
		"junos-sun-rpc-udp",
		"junos-sun-rpc-walld",
		"junos-sun-rpc-walld-tcp",
		"junos-sun-rpc-walld-udp",
		"junos-sun-rpc-ypbind",
		"junos-sun-rpc-ypbind-tcp",
		"junos-sun-rpc-ypbind-udp",
		"junos-sun-rpc-ypserv",
		"junos-sun-rpc-ypserv-tcp",
		"junos-sun-rpc-ypserv-udp",
		"junos-syslog",
		"junos-tacacs",
		"junos-tacacs-ds",
		"junos-talk",
		"junos-tcp-any",
		"junos-telnet",
		"junos-tftp",
		"junos-udp-any",
		"junos-uucp",
		"junos-vdo-live",
		"junos-vnc",
		"junos-wais",
		"junos-who",
		"junos-whois",
		"junos-winframe",
		"junos-wxcontrol",
		"junos-x-windows",
		"junos-xnm-clear-text",
		"junos-xnm-ssl",
		"junos-ymsg",
	}
	groups = map[int]string{
		1:  "group1",
		2:  "group2",
		5:  "group5",
		14: "group14",
		19: "group19",
		20: "group20",
		24: "group24",
	}
	encrAlgorithm = map[string]string{
		"3des":    "3des-cbc",
		"aes-128": "aes-128-cbc",
		"aes-192": "aes-192-cbc",
		"aes-256": "aes-256-cbc",
		"des":     "des-cbc",
	}
)

// NewPolicy establishes a blank security policy that will hold any newly created rules.
func (j *Junos) NewPolicy() *FirewallPolicy {
	var addrs ExistingAddresses
	var apps ExistingApplications
	addresses := []string{}
	applications := []string{}
	appConfig := []string{}
	getAddresses, _ := j.GetConfig("security>address-book", "xml")
	getApplications, _ := j.GetConfig("applications", "xml")

	if err := xml.Unmarshal([]byte(getAddresses), &addrs); err != nil {
		fmt.Println(err)
	}

	if err := xml.Unmarshal([]byte(getApplications), &apps); err != nil {
		fmt.Println(err)
	}

	for _, addr := range addrs.ExistingAddresses {
		addresses = append(addresses, addr.Name)
	}

	for _, addrSet := range addrs.ExistingAddressSets {
		addresses = append(addresses, addrSet.Name)
	}

	for _, app := range apps.ExistingApplications {
		applications = append(applications, app.Name)
	}

	for _, appSet := range apps.ExistingApplicationSets {
		applications = append(applications, appSet.Name)
	}

	for _, d := range junosDefaultApps {
		applications = append(applications, d)
	}

	return &FirewallPolicy{
		Addresses:    addresses,
		Applications: applications,
		AppConfig:    appConfig,
	}
}

// CreateApplication creates a TCP or UDP service that is previously not defined on the SRX.
func (p *FirewallPolicy) CreateApplication(name, protocol, dstport string) {
	p.AppConfig = append(p.AppConfig, fmt.Sprintf("set applications application %s protocol %s destination-port %s\n", name, protocol, dstport))
}

// AddRule creates a single rule and adds it to the policy. For services/applications, you MUST
// provide an already existing application. If you wish to create one, use CreateApplication().
func (p *FirewallPolicy) AddRule(name interface{}, srczone, src, dstzone, dst, application, action string) {
	srcAddrs := strings.Split(src, ",")
	dstAddrs := strings.Split(dst, ",")
	apps := strings.Split(application, ",")

	rule := Rule{
		Name:          name,
		SourceZone:    srczone,
		SourceAddress: srcAddrs,
		DestZone:      dstzone,
		DestAddress:   dstAddrs,
		Application:   apps,
		Action:        action,
	}

	p.Rules = append(p.Rules, rule)
}

// BuildPolicy creates our SRX security policy configuration. You can use Config() to apply the changes.
func (p *FirewallPolicy) BuildPolicy() []string {
	var policy []string

	for _, a := range p.AppConfig {
		policy = append(policy, a)
	}

	for _, r := range p.Rules {
		srcAddrs := "[ "
		for _, s := range r.SourceAddress {
			srcAddrs += fmt.Sprintf("%s ", strings.TrimSpace(s))
		}
		srcAddrs += "]"

		dstAddrs := "[ "
		for _, d := range r.DestAddress {
			dstAddrs += fmt.Sprintf("%s ", strings.TrimSpace(d))
		}
		dstAddrs += "]"

		apps := "[ "
		for _, a := range r.Application {
			apps += fmt.Sprintf("%s ", strings.TrimSpace(a))
		}
		apps += "]"

		rule := fmt.Sprintf("set security policies from-zone %s to-zone %s policy %v match source-address %s destination-address %s application %s\n", r.SourceZone, r.DestZone, r.Name, srcAddrs, dstAddrs, apps)
		rule += fmt.Sprintf("set security policies from-zone %s to-zone %s policy %v then %s\n", r.SourceZone, r.DestZone, r.Name, r.Action)
		rule += fmt.Sprintf("set security policies from-zone %s to-zone %s policy %v then log session-init session-close\n", r.SourceZone, r.DestZone, r.Name)

		policy = append(policy, rule)
	}

	return policy
}

// ConvertAddressBook will generate the configuration needed to migrate from a zone-based address
// book to a global one. You can then use Config() to apply the changes if necessary.
func (j *Junos) ConvertAddressBook() []string {
	vrx := regexp.MustCompile(`(\d+)\.(\d+)([RBISX]{1})(\d+)(\.(\d+))?`)

	for _, d := range j.Platform {
		if strings.Contains(d.Model, "FIREFLY") {
			continue
		}

		if !strings.Contains(d.Model, "SRX") {
			fmt.Printf("This device doesn't look to be an SRX (%s). You can only run this script against an SRX.\n", d.Model)
			os.Exit(0)
		}
		versionBreak := vrx.FindStringSubmatch(d.Version)
		maj, _ := strconv.Atoi(versionBreak[1])
		min, _ := strconv.Atoi(versionBreak[2])
		// rel := versionBreak[3]
		// build, _ := strconv.Atoi(versionBreak[4])

		if maj <= 11 && min < 2 {
			fmt.Println("You must be running JUNOS version 11.2 or above in order to use this conversion tool.")
			os.Exit(0)
		}
	}

	var seczones SecurityZones
	globalAddressBook := []string{}

	zoneConfig, _ := j.GetConfig("security>zones", "xml")
	if err := xml.Unmarshal([]byte(zoneConfig), &seczones); err != nil {
		fmt.Println(err)
	}

	for _, z := range seczones.Zones {
		for _, a := range z.AddressEntries {
			if a.DNSName != "" {
				globalConfig := fmt.Sprintf("set security address-book global address %s dns-name %s\n", a.Name, a.DNSName)
				globalAddressBook = append(globalAddressBook, globalConfig)
			}
			if a.Wildcard != "" {
				globalConfig := fmt.Sprintf("set security address-book global address %s wildcard-address %s\n", a.Name, a.Wildcard)
				globalAddressBook = append(globalAddressBook, globalConfig)
			}
			if a.IP != "" {
				globalConfig := fmt.Sprintf("set security address-book global address %s %s\n", a.Name, a.IP)
				globalAddressBook = append(globalAddressBook, globalConfig)
			}
		}

		for _, as := range z.AddressSets {
			for _, addr := range as.AddressEntries {
				globalConfig := fmt.Sprintf("set security address-book global address-set %s address %s\n", as.Name, addr.Name)
				globalAddressBook = append(globalAddressBook, globalConfig)
			}
		}

		removeConfig := fmt.Sprintf("delete security zones security-zone %s address-book\n", z.Name)
		globalAddressBook = append(globalAddressBook, removeConfig)
	}

	return globalAddressBook
}

// NewIPsecVPN creates the initial template needed to bulid a new site-to-site VPN. Options are
// as follows:
//
// <name> - Name of the VPN.
//
// <local> - the public IP address of the SRX terminating the VPN.
//
// <peer> - the remote devices' IP address.
//
// <iface> - the external interface of the SRX (typically where the <local> IP is tied to).
//
// <zone> - the security-zone where the routed st0.<unit> interfaces reside.
//
// <pfs> - 1, 2, 5, 14, 19, 20, 24; use 0 to disable.
//
// <establish> - "traffic" or "immediately"; whether or not to establish the tunnel on-traffic or immediately.
//
// <mode> - "main" or "aggressive."
//
// <psk> - Pre-shared key.
func (j *Junos) NewIPsecVPN(name, local, peer, iface, zone string, pfs int, establish string, mode, psk string) *IPsecVPN {
	var ints st0Interface
	gateway := []string{}
	ontraffic := map[string]string{
		"traffic":     "on-traffic",
		"immediately": "immediately",
	}
	st0, _ := j.RunCommand("show interfaces st0", "xml")

	if err := xml.Unmarshal([]byte(st0), &ints); err != nil {
		fmt.Println(err)
	}

	stInts := map[int]int{}

	for n, i := range ints.Units {
		trimmed := strings.TrimPrefix(strings.TrimSpace(i.Name), "st0.")
		unit, _ := strconv.Atoi(trimmed)
		stInts[n] = unit
	}

	total := len(stInts)
	newst0 := fmt.Sprintf("st0.%d", stInts[total-1]+1)

	gateway = append(gateway, fmt.Sprintf("set security ike gateway %s address %s\n", name, peer))
	gateway = append(gateway, fmt.Sprintf("set security ike gateway %s external-interface %s\n", name, iface))
	gateway = append(gateway, fmt.Sprintf("set security ike gateway %s ike-policy %s\n", name, name))
	gateway = append(gateway, fmt.Sprintf("set security ike gateway %s local-address %s\n", name, local))

	return &IPsecVPN{
		Name:              name,
		Local:             local,
		Peer:              peer,
		ExternalInterface: iface,
		St0:               newst0,
		Zone:              zone,
		PFS:               pfs,
		Establish:         ontraffic[establish],
		Mode:              mode,
		PSK:               psk,
		Gateway:           gateway,
	}
}

// Phase1 creates the IKE proposal to use for the site-to-site VPN. Options are as follows:
//
// <name> - Name of the proposal.
//
// <dh> - 1, 2, 5, 14, 19, 20, 24
//
// <auth> - "md5" or "sha1"
//
// <encryption> - "3des", "aes-128", "aes-192", "aes-256", "des"
//
// <lifetime> - Lifetime in seconds.
func (i *IPsecVPN) Phase1(name string, dh int, auth, encryption string, lifetime int) {
	authAlgorithm := map[string]string{
		"md5":  "md5",
		"sha1": "sha1",
	}

	p1 := P1{
		P1Name:           name,
		P1DiffeHellman:   dh,
		P1Authentication: authAlgorithm[auth],
		P1Encryption:     encrAlgorithm[encryption],
		P1Seconds:        lifetime,
	}

	i.P1Proposals = append(i.P1Proposals, p1)
}

// Phase2 creates the IPsec proposal to use for the site-to-site VPN. Options are as follows:
//
// <name> - Name of the proposal.
//
// <auth> - "md5" or "sha1"
//
// <encryption> - "3des", "aes-128", "aes-192", "aes-256", "des"
//
// <lifetime> - Lifetime in seconds.
//
// <protocol> - "ah" or "esp"
func (i *IPsecVPN) Phase2(name string, auth, encryption string, lifetime int, protocol string) {
	authAlgorithm := map[string]string{
		"md5":  "hmac-md5-96",
		"sha1": "hmac-sha1-96",
	}

	p2 := P2{
		P2Name:           name,
		P2Authentication: authAlgorithm[auth],
		P2Encryption:     encrAlgorithm[encryption],
		P2Seconds:        lifetime,
		P2Protocol:       protocol,
	}

	i.P2Proposals = append(i.P2Proposals, p2)
}

// TrafficSelector creates the security-association (SA) configuration needed when building
// a site-to-site VPN. <local> and <remote> are a []string of IP addresses.
func (i *IPsecVPN) TrafficSelector(local, remote []string) {
	count := 1
	ts := []string{}

	for _, l := range local {
		for _, r := range remote {
			tsCfg := fmt.Sprintf("set security ipsec vpn %s traffic-selector ts%d local-ip %s remote-ip %s\n", i.Name, count, l, r)
			ts = append(ts, tsCfg)

			count++
		}
	}

	i.TrafficSelectors = ts
}

// BuildIPsecVPN creates the configuration of the site-to-site VPN to be commited.
func (i *IPsecVPN) BuildIPsecVPN() []string {
	config := []string{}

	config = append(config, fmt.Sprintf("set interfaces %s family inet\n", i.St0))
	config = append(config, fmt.Sprintf("set security zones security-zone %s interfaces %s\n", i.Zone, i.St0))

	for _, p1 := range i.P1Proposals {
		phase1 := []string{}
		phase1 = append(phase1, fmt.Sprintf("set security ike proposal %s authentication-method pre-shared-keys\n", p1.P1Name))
		phase1 = append(phase1, fmt.Sprintf("set security ike proposal %s dh-group %s\n", p1.P1Name, groups[p1.P1DiffeHellman]))
		phase1 = append(phase1, fmt.Sprintf("set security ike proposal %s authentication-algorithm %s\n", p1.P1Name, p1.P1Authentication))
		phase1 = append(phase1, fmt.Sprintf("set security ike proposal %s encryption-algorithm %s\n", p1.P1Name, p1.P1Encryption))
		phase1 = append(phase1, fmt.Sprintf("set security ike proposal %s lifetime-seconds %d\n", p1.P1Name, p1.P1Seconds))

		for _, p := range phase1 {
			config = append(config, p)
		}
	}

	config = append(config, fmt.Sprintf("set security ike policy %s mode %s\n", i.Name, i.Mode))
	config = append(config, fmt.Sprintf("set security ike policy %s pre-shared-key ascii-text \"%s\"\n", i.Name, i.PSK))

	for _, p1props := range i.P1Proposals {
		config = append(config, fmt.Sprintf("set security ike policy %s proposals %s\n", i.Name, p1props.P1Name))
	}

	for _, g := range i.Gateway {
		config = append(config, g)
	}

	for _, p2 := range i.P2Proposals {
		phase2 := []string{}
		phase2 = append(phase2, fmt.Sprintf("set security ipsec proposal %s protocol %s\n", p2.P2Name, p2.P2Protocol))
		phase2 = append(phase2, fmt.Sprintf("set security ipsec proposal %s authentication-algorithm %s\n", p2.P2Name, p2.P2Authentication))
		phase2 = append(phase2, fmt.Sprintf("set security ipsec proposal %s encryption-algorithm %s\n", p2.P2Name, p2.P2Encryption))
		phase2 = append(phase2, fmt.Sprintf("set security ipsec proposal %s lifetime-seconds %d\n", p2.P2Name, p2.P2Seconds))

		for _, p := range phase2 {
			config = append(config, p)
		}
	}

	for _, p2props := range i.P2Proposals {
		config = append(config, fmt.Sprintf("set security ipsec policy %s proposals %s\n", i.Name, p2props.P2Name))
	}

	if i.PFS != 0 {
		config = append(config, fmt.Sprintf("set security ipsec policy %s perfect-forward-secrecy keys %s\n", i.Name, groups[i.PFS]))
	}

	config = append(config, fmt.Sprintf("set security ipsec vpn %s bind-interface %s\n", i.Name, i.St0))
	config = append(config, fmt.Sprintf("set security ipsec vpn %s ike gateway %s\n", i.Name, i.Name))
	config = append(config, fmt.Sprintf("set security ipsec vpn %s ike idle-time 60\n", i.Name))
	config = append(config, fmt.Sprintf("set security ipsec vpn %s ike ipsec-policy %s\n", i.Name, i.Name))
	config = append(config, fmt.Sprintf("set security ipsec vpn %s establish-tunnels %s\n", i.Name, i.Establish))

	for _, ts := range i.TrafficSelectors {
		config = append(config, ts)
	}

	return config
}
