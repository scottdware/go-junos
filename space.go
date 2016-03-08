package junos

// All of our HTTP Content-Types we use.
var (
	contentDiscoverDevices = "application/vnd.net.juniper.space.device-management.discover-devices+xml;version=2;charset=UTF-8"
	contentExecDeploy      = "application/vnd.net.juniper.space.software-management.exec-deploy+xml;version=1;charset=UTF-8"
	contentExecRemove      = "application/vnd.net.juniper.space.software-management.exec-remove+xml;version=1;charset=UTF-8"
	contentExecStage       = "application/vnd.net.juniper.space.software-management.exec-stage+xml;version=1;charset=UTF-8"
	contentAddress         = "application/vnd.juniper.sd.address-management.address+xml;version=1;charset=UTF-8"
	contentUpdateDevices   = "application/vnd.juniper.sd.device-management.update-devices+xml;version=1;charset=UTF-8"
	contentPublish         = "application/vnd.juniper.sd.fwpolicy-management.publish+xml;version=1;charset=UTF-8"
	contentAddressPatch    = "application/vnd.juniper.sd.address-management.address_patch+xml;version=1;charset=UTF-8"
	contentService         = "application/vnd.juniper.sd.service-management.service+xml;version=1;charset=UTF-8"
	contentServicePatch    = "application/vnd.juniper.sd.service-management.service_patch+xml;version=1;charset=UTF-8"
	contentVariable        = "application/vnd.juniper.sd.variable-management.variable-definition+xml;version=1;charset=UTF-8"
	contentResync          = "application/vnd.net.juniper.space.device-management.exec-resync+xml;version=1"
)

// Space contains our session state.
type Space struct {
	Host     string
	User     string
	Password string
}

// APIRequest builds our request before sending it to the server.
type APIRequest struct {
	Method      string
	URL         string
	Body        string
	ContentType string
}

type jobID struct {
	ID int `xml:"id"`
}

type jobDetail struct {
	ID      jobID
	Name    string  `xml:"name"`
	State   string  `xml:"job-state"`
	Status  string  `xml:"job-status"`
	Percent float64 `xml:"percent-complete"`
}

// NewServer sets up our connection to the Junos Space server.
func NewServer(host, user, passwd string) *Space {
	return &Space{
		Host:     host,
		User:     user,
		Password: passwd,
	}
}
