package junos

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
)

// JunosSpace holds all of our information that we use for our server
// connection.
type JunosSpace struct {
	Host      string
	User      string
	Password  string
	Transport *http.Transport
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

// contentType holds all of the HTTP Content-Types that our Junos Space requests will use.
var contentType = map[string]string{
	"discover-devices": "application/vnd.net.juniper.space.device-management.discover-devices+xml;version=2;charset=UTF-8",
	"exec-deploy":      "application/vnd.net.juniper.space.software-management.exec-deploy+xml;version=1;charset=UTF-8",
	"exec-remove":      "application/vnd.net.juniper.space.software-management.exec-remove+xml;version=1;charset=UTF-8",
	"exec-stage":       "application/vnd.net.juniper.space.software-management.exec-stage+xml;version=1;charset=UTF-8",
}

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

// APIDelete builds our DELETE request to the server.
func (s *JunosSpace) APIDelete(uri string) error {
	var req *http.Request
	client := &http.Client{Transport: s.Transport}
	url := fmt.Sprintf("https://%s/api/space/%s", s.Host, uri)
	req, _ = http.NewRequest("DELETE", url, nil)
	req.SetBasicAuth(s.User, s.Password)
	res, err := client.Do(req)
	defer res.Body.Close()

	if err != nil {
		return err
	}

	return nil
}

// APIPost builds our POST request to the server.
func (s *JunosSpace) APIPost(uri, body, ct string) ([]byte, error) {
	var req *http.Request
	b := bytes.NewReader([]byte(body))
	client := &http.Client{Transport: s.Transport}
	url := fmt.Sprintf("https://%s/api/space/%s", s.Host, uri)
	req, _ = http.NewRequest("POST", url, b)
	req.Header.Set("Content-Type", contentType[ct])
	req.SetBasicAuth(s.User, s.Password)
	res, err := client.Do(req)
	defer res.Body.Close()

	if err != nil {
		return nil, err
	}

	data, _ := ioutil.ReadAll(res.Body)

	return data, nil
}

// APIRequest builds our GET request to the server.
func (s *JunosSpace) APIRequest(uri string) ([]byte, error) {
	var req *http.Request
	client := &http.Client{Transport: s.Transport}
	url := fmt.Sprintf("https://%s/api/space/%s", s.Host, uri)
	req, _ = http.NewRequest("GET", url, nil)
	req.SetBasicAuth(s.User, s.Password)
	res, err := client.Do(req)
	defer res.Body.Close()

	if err != nil {
		return nil, err
	}

	data, _ := ioutil.ReadAll(res.Body)

	return data, nil
}
