package junos

import (
	"fmt"
	"strings"
)

// FirewallPolicy contains all of the rules that will be created for the policy.
type FirewallPolicy struct {
	Rules []Rule
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

var (
	// These are the Junos default applications in an SRX
	junosDefaultApps = map[string]string{
		"aol":                          "junos-aol",
		"bgp":                          "junos-bgp",
		"biff":                         "junos-biff",
		"bootpc":                       "junos-bootpc",
		"bootps":                       "junos-bootps",
		"chargen":                      "junos-chargen",
		"cifs":                         "junos-cifs",
		"cvspserver":                   "junos-cvspserver",
		"dhcp-client":                  "junos-dhcp-client",
		"dhcp-relay":                   "junos-dhcp-relay",
		"dhcp-server":                  "junos-dhcp-server",
		"discard":                      "junos-discard",
		"dns-tcp":                      "junos-dns-tcp",
		"dns-udp":                      "junos-dns-udp",
		"echo":                         "junos-echo",
		"finger":                       "junos-finger",
		"ftp":                          "junos-ftp",
		"gnutella":                     "junos-gnutella",
		"gopher":                       "junos-gopher",
		"gre":                          "junos-gre",
		"gtp":                          "junos-gtp",
		"h323":                         "junos-h323",
		"http":                         "junos-http",
		"http-ext":                     "junos-http-ext",
		"https":                        "junos-https",
		"icmp-all":                     "junos-icmp-all",
		"icmp-ping":                    "junos-icmp-ping",
		"icmp6-all":                    "junos-icmp6-all",
		"icmp6-dst-unreach-addr":       "junos-icmp6-dst-unreach-addr",
		"icmp6-dst-unreach-admin":      "junos-icmp6-dst-unreach-admin",
		"icmp6-dst-unreach-beyond":     "junos-icmp6-dst-unreach-beyond",
		"icmp6-dst-unreach-port":       "junos-icmp6-dst-unreach-port",
		"icmp6-dst-unreach-route":      "junos-icmp6-dst-unreach-route",
		"icmp6-echo-reply":             "junos-icmp6-echo-reply",
		"icmp6-echo-request":           "junos-icmp6-echo-request",
		"icmp6-packet-too-big":         "junos-icmp6-packet-too-big",
		"icmp6-param-prob-header":      "junos-icmp6-param-prob-header",
		"icmp6-param-prob-nexthdr":     "junos-icmp6-param-prob-nexthdr",
		"icmp6-param-prob-option":      "junos-icmp6-param-prob-option",
		"icmp6-time-exceed-reassembly": "junos-icmp6-time-exceed-reassembly",
		"icmp6-time-exceed-transit":    "junos-icmp6-time-exceed-transit",
		"ident":                        "junos-ident",
		"ike":                          "junos-ike",
		"ike-nat":                      "junos-ike-nat",
		"imap":                         "junos-imap",
		"imaps":                        "junos-imaps",
		"internet-locator-service":     "junos-internet-locator-service",
		"irc":                                  "junos-irc",
		"l2tp":                                 "junos-l2tp",
		"ldap":                                 "junos-ldap",
		"ldp-tcp":                              "junos-ldp-tcp",
		"ldp-udp":                              "junos-ldp-udp",
		"lpr":                                  "junos-lpr",
		"mail":                                 "junos-mail",
		"mgcp":                                 "junos-mgcp",
		"mgcp-ca":                              "junos-mgcp-ca",
		"mgcp-ua":                              "junos-mgcp-ua",
		"ms-rpc":                               "junos-ms-rpc",
		"ms-rpc-any":                           "junos-ms-rpc-any",
		"ms-rpc-epm":                           "junos-ms-rpc-epm",
		"ms-rpc-iis-com":                       "junos-ms-rpc-iis-com",
		"ms-rpc-iis-com-1":                     "junos-ms-rpc-iis-com-1",
		"ms-rpc-iis-com-adminbase":             "junos-ms-rpc-iis-com-adminbase",
		"ms-rpc-msexchange":                    "junos-ms-rpc-msexchange",
		"ms-rpc-msexchange-directory-nsp":      "junos-ms-rpc-msexchange-directory-nsp",
		"ms-rpc-msexchange-directory-rfr":      "junos-ms-rpc-msexchange-directory-rfr",
		"ms-rpc-msexchange-info-store":         "junos-ms-rpc-msexchange-info-store",
		"ms-rpc-tcp":                           "junos-ms-rpc-tcp",
		"ms-rpc-udp":                           "junos-ms-rpc-udp",
		"ms-rpc-uuid-any-tcp":                  "junos-ms-rpc-uuid-any-tcp",
		"ms-rpc-uuid-any-udp":                  "junos-ms-rpc-uuid-any-udp",
		"ms-rpc-wmic":                          "junos-ms-rpc-wmic",
		"ms-rpc-wmic-admin":                    "junos-ms-rpc-wmic-admin",
		"ms-rpc-wmic-admin2":                   "junos-ms-rpc-wmic-admin2",
		"ms-rpc-wmic-mgmt":                     "junos-ms-rpc-wmic-mgmt",
		"ms-rpc-wmic-webm-callresult":          "junos-ms-rpc-wmic-webm-callresult",
		"ms-rpc-wmic-webm-classobject":         "junos-ms-rpc-wmic-webm-classobject",
		"ms-rpc-wmic-webm-level1login":         "junos-ms-rpc-wmic-webm-level1login",
		"ms-rpc-wmic-webm-login-clientid":      "junos-ms-rpc-wmic-webm-login-clientid",
		"ms-rpc-wmic-webm-login-helper":        "junos-ms-rpc-wmic-webm-login-helper",
		"ms-rpc-wmic-webm-objectsink":          "junos-ms-rpc-wmic-webm-objectsink",
		"ms-rpc-wmic-webm-refreshing-services": "junos-ms-rpc-wmic-webm-refreshing-services",
		"ms-rpc-wmic-webm-remote-refresher":    "junos-ms-rpc-wmic-webm-remote-refresher",
		"ms-rpc-wmic-webm-services":            "junos-ms-rpc-wmic-webm-services",
		"ms-rpc-wmic-webm-shutdown":            "junos-ms-rpc-wmic-webm-shutdown",
		"ms-sql":                               "junos-ms-sql",
		"msn":                                  "junos-msn",
		"nbds":                                 "junos-nbds",
		"nbname":                               "junos-nbname",
		"netbios-session":                      "junos-netbios-session",
		"nfs":                                  "junos-nfs",
		"nfsd-tcp":                             "junos-nfsd-tcp",
		"nfsd-udp":                             "junos-nfsd-udp",
		"nntp":                                 "junos-nntp",
		"ns-global":                            "junos-ns-global",
		"ns-global-pro":                        "junos-ns-global-pro",
		"nsm":                                  "junos-nsm",
		"ntalk":                                "junos-ntalk",
		"ntp":                                  "junos-ntp",
		"ospf":                                 "junos-ospf",
		"pc-anywhere":                          "junos-pc-anywhere",
		"persistent-nat":                       "junos-persistent-nat",
		"ping":                                 "junos-ping",
		"pingv6":                               "junos-pingv6",
		"pop3":                                 "junos-pop3",
		"pptp":                                 "junos-pptp",
		"printer":                              "junos-printer",
		"r2cp":                                 "junos-r2cp",
		"radacct":                              "junos-radacct",
		"radius":                               "junos-radius",
		"realaudio":                            "junos-realaudio",
		"rip":                                  "junos-rip",
		"routing-inbound":                      "junos-routing-inbound",
		"rsh":                                  "junos-rsh",
		"rtsp":                                 "junos-rtsp",
		"sccp":                                 "junos-sccp",
		"sctp-any":                             "junos-sctp-any",
		"sip":                                  "junos-sip",
		"smb":                                  "junos-smb",
		"smb-session":                          "junos-smb-session",
		"smtp":                                 "junos-smtp",
		"snmp-agentx":                          "junos-snmp-agentx",
		"snpp":                                 "junos-snpp",
		"sql-monitor":                          "junos-sql-monitor",
		"sqlnet-v1":                            "junos-sqlnet-v1",
		"sqlnet-v2":                            "junos-sqlnet-v2",
		"ssh":                                  "junos-ssh",
		"stun":                                 "junos-stun",
		"sun-rpc":                              "junos-sun-rpc",
		"sun-rpc-any":                          "junos-sun-rpc-any",
		"sun-rpc-any-tcp":                      "junos-sun-rpc-any-tcp",
		"sun-rpc-any-udp":                      "junos-sun-rpc-any-udp",
		"sun-rpc-mountd":                       "junos-sun-rpc-mountd",
		"sun-rpc-mountd-tcp":                   "junos-sun-rpc-mountd-tcp",
		"sun-rpc-mountd-udp":                   "junos-sun-rpc-mountd-udp",
		"sun-rpc-nfs":                          "junos-sun-rpc-nfs",
		"sun-rpc-nfs-access":                   "junos-sun-rpc-nfs-access",
		"sun-rpc-nfs-tcp":                      "junos-sun-rpc-nfs-tcp",
		"sun-rpc-nfs-udp":                      "junos-sun-rpc-nfs-udp",
		"sun-rpc-nlockmgr":                     "junos-sun-rpc-nlockmgr",
		"sun-rpc-nlockmgr-tcp":                 "junos-sun-rpc-nlockmgr-tcp",
		"sun-rpc-nlockmgr-udp":                 "junos-sun-rpc-nlockmgr-udp",
		"sun-rpc-portmap":                      "junos-sun-rpc-portmap",
		"sun-rpc-portmap-tcp":                  "junos-sun-rpc-portmap-tcp",
		"sun-rpc-portmap-udp":                  "junos-sun-rpc-portmap-udp",
		"sun-rpc-rquotad":                      "junos-sun-rpc-rquotad",
		"sun-rpc-rquotad-tcp":                  "junos-sun-rpc-rquotad-tcp",
		"sun-rpc-rquotad-udp":                  "junos-sun-rpc-rquotad-udp",
		"sun-rpc-ruserd":                       "junos-sun-rpc-ruserd",
		"sun-rpc-ruserd-tcp":                   "junos-sun-rpc-ruserd-tcp",
		"sun-rpc-ruserd-udp":                   "junos-sun-rpc-ruserd-udp",
		"sun-rpc-sadmind":                      "junos-sun-rpc-sadmind",
		"sun-rpc-sadmind-tcp":                  "junos-sun-rpc-sadmind-tcp",
		"sun-rpc-sadmind-udp":                  "junos-sun-rpc-sadmind-udp",
		"sun-rpc-sprayd":                       "junos-sun-rpc-sprayd",
		"sun-rpc-sprayd-tcp":                   "junos-sun-rpc-sprayd-tcp",
		"sun-rpc-sprayd-udp":                   "junos-sun-rpc-sprayd-udp",
		"sun-rpc-status":                       "junos-sun-rpc-status",
		"sun-rpc-status-tcp":                   "junos-sun-rpc-status-tcp",
		"sun-rpc-status-udp":                   "junos-sun-rpc-status-udp",
		"sun-rpc-tcp":                          "junos-sun-rpc-tcp",
		"sun-rpc-udp":                          "junos-sun-rpc-udp",
		"sun-rpc-walld":                        "junos-sun-rpc-walld",
		"sun-rpc-walld-tcp":                    "junos-sun-rpc-walld-tcp",
		"sun-rpc-walld-udp":                    "junos-sun-rpc-walld-udp",
		"sun-rpc-ypbind":                       "junos-sun-rpc-ypbind",
		"sun-rpc-ypbind-tcp":                   "junos-sun-rpc-ypbind-tcp",
		"sun-rpc-ypbind-udp":                   "junos-sun-rpc-ypbind-udp",
		"sun-rpc-ypserv":                       "junos-sun-rpc-ypserv",
		"sun-rpc-ypserv-tcp":                   "junos-sun-rpc-ypserv-tcp",
		"sun-rpc-ypserv-udp":                   "junos-sun-rpc-ypserv-udp",
		"syslog":                               "junos-syslog",
		"tacacs":                               "junos-tacacs",
		"tacacs-ds":                            "junos-tacacs-ds",
		"talk":                                 "junos-talk",
		"tcp-any":                              "junos-tcp-any",
		"telnet":                               "junos-telnet",
		"tftp":                                 "junos-tftp",
		"udp-any":                              "junos-udp-any",
		"uucp":                                 "junos-uucp",
		"vdo-live":                             "junos-vdo-live",
		"vnc":                                  "junos-vnc",
		"wais":                                 "junos-wais",
		"who":                                  "junos-who",
		"whois":                                "junos-whois",
		"winframe":                             "junos-winframe",
		"wxcontrol":                            "junos-wxcontrol",
		"x-windows":                            "junos-x-windows",
		"xnm-clear-text":                       "junos-xnm-clear-text",
		"xnm-ssl":                              "junos-xnm-ssl",
		"ymsg":                                 "junos-ymsg",
	}
)

// CreatePolicy establishes a blank policy that will hold any newly created rules.
func CreatePolicy() *FirewallPolicy {
	return &FirewallPolicy{}
}

// AddRule creates a single rule and adds it to the policy.
func (p *FirewallPolicy) AddRule(name interface{}, srczone, src, dstzone, dst, application, action string) {
	var applications []string
	srcAddrs := strings.Split(src, ",")
	dstAddrs := strings.Split(dst, ",")
	apps := strings.Split(application, ",")

	for _, a := range apps {
		if junosDefaultApps[strings.TrimSpace(a)] != "" {
			applications = append(applications, junosDefaultApps[strings.TrimSpace(a)])
		}

		if junosDefaultApps[strings.TrimSpace(a)] == "" {
			applications = append(applications, strings.TrimSpace(a))
		}
	}

	rule := Rule{
		Name:          name,
		SourceZone:    srczone,
		SourceAddress: srcAddrs,
		DestZone:      dstzone,
		DestAddress:   dstAddrs,
		Application:   applications,
		Action:        action,
	}

	p.Rules = append(p.Rules, rule)
}

// BuildPolicy creates our SRX security policy configuration.
func (p *FirewallPolicy) BuildPolicy() []string {
	var policy []string
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
