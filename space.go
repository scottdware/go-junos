package junos

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// JunosSpace contains our session state.
type JunosSpace struct {
	Host      string
	User      string
	Password  string
	Transport *http.Transport
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

// APICall is used to query the Junos Space server API.
func (s *JunosSpace) APICall(options *APIRequest) ([]byte, error) {
	var req *http.Request
	client := &http.Client{Transport: s.Transport}
	url := fmt.Sprintf("https://%s%s", s.Host, options.URL)
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
