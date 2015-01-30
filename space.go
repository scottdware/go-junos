package junos

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// JunosSpace holds all of our information that we use for our server
// connection.
type JunosSpace struct {
	Host      string
	User      string
	Password  string
	Transport *http.Transport
}

// APIRequest holds all of our options when building our API call.
type APIRequest struct {
	Method      string
	URL         string
	Body        string
	ContentType string
}

// jobID parses the job ID from the returned XML.
type jobID struct {
	ID int `xml:"id"`
}

// jobDetail holds the information about a given job.
type jobDetail struct {
	ID      jobID
	Name    string  `xml:"name"`
	State   string  `xml:"job-state"`
	Status  string  `xml:"job-status"`
	Percent float64 `xml:"percent-complete"`
}

// These are all of the HTTP Content-Type's that we use for POST and DELETE requests.
// You can also call an API yourself using a URL and Content-Type if one is not
// listed here.
var (
	ContentDiscoverDevices       = "application/vnd.net.juniper.space.device-management.discover-devices+xml;version=2;charset=UTF-8"
	ContentExecDeploy            = "application/vnd.net.juniper.space.software-management.exec-deploy+xml;version=1;charset=UTF-8"
	ContentExecRemove            = "application/vnd.net.juniper.space.software-management.exec-remove+xml;version=1;charset=UTF-8"
	ContentExecStage             = "application/vnd.net.juniper.space.software-management.exec-stage+xml;version=1;charset=UTF-8"
	ContentAddress               = "application/vnd.juniper.sd.address-management.address+xml;version=1;charset=UTF-8"
	ContentDeleteAddressResponse = "application/vnd.juniper.sd.address-management.delete-address-response+xml;version=1;q=0.01"
	ContentUpdateDevices         = "application/vnd.juniper.sd.device-management.update-devices+xml;version=1;charset=UTF-8"
	ContentPublish               = "application/vnd.juniper.sd.fwpolicy-management.publish+xml;version=1;charset=UTF-8"
)

// NewServer sets up our connection to the Junos Space server.
func NewServer(host, user, passwd string) *JunosSpace {
	return &JunosSpace{
		Host:     host,
		User:     user,
		Password: passwd,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
}

// APICall is used to query the Junos Space server API's.
func (s *JunosSpace) APICall(options *APIRequest) ([]byte, error) {
	var req *http.Request
	client := &http.Client{Transport: s.Transport}
	url := fmt.Sprintf("https://%s/api/%s", s.Host, options.URL)
	body := bytes.NewReader([]byte(options.Body))
	req, _ = http.NewRequest(strings.ToUpper(options.Method), url, body)
	req.SetBasicAuth(s.User, s.Password)

	if len(options.ContentType) > 0 {
		req.Header.Set("Content-Type", options.ContentType)
	}

	res, err := client.Do(req)
	defer res.Body.Close()

	if err != nil {
		return nil, err
	}

	data, _ := ioutil.ReadAll(res.Body)

	return data, nil
}
