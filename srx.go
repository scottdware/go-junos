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
)

// CreatePolicy establishes a blank policy that will hold any newly created rules.
func (j *Junos) CreatePolicy() *FirewallPolicy {
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
