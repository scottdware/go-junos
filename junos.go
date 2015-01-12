package junos

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/Juniper/go-netconf/netconf"
	"io/ioutil"
	"log"
	"strings"
)

// Junos holds the connection information to our Junos device.
type Junos struct {
	*netconf.Session
}

// rollbackXML parses our rollback diff configuration.
type rollbackXML struct {
	XMLName xml.Name `xml:"rollback-information"`
	Config  string   `xml:"configuration-information>configuration-output"`
}

// CommandXML parses our operational command responses.
type commandXML struct {
	Config string `xml:",innerxml"`
}

type software struct {
	RES []routingEngine `xml:"multi-routing-engine-item"`
}

type routingEngine struct {
	Name    string `xml:"re-name"`
	Model   string `xml:"software-information>product-model"`
	Type    string `xml:"software-information>package-information>name"`
	Version string `xml:"software-information>package-information>comment"`
}

// NewSession establishes a new connection to a Junos device that we will use
// to run our commands against.
func NewSession(host, user, password string) *Junos {
	s, err := netconf.DialSSH(host, netconf.SSHConfigPassword(user, password))
	if err != nil {
		log.Fatal(err)
	}

	return &Junos{
		s,
	}
}

// Commit commits the configuration.
func (j *Junos) Commit() error {
	reply, err := j.Exec(rpcCommand["commit"])
	if err != nil {
		return err
	}

	if reply.Ok == false {
		for _, m := range reply.Errors {
			return errors.New(m.Message)
		}
	}

	return nil
}

// Lock locks the candidate configuration.
func (j *Junos) Lock() error {
	reply, err := j.Exec(rpcCommand["lock"])
	if err != nil {
		return err
	}

	if reply.Ok == false {
		for _, m := range reply.Errors {
			return errors.New(m.Message)
		}
	}

	return nil
}

// Unlock unlocks the candidate configuration.
func (j *Junos) Unlock() error {
	resp, err := j.Exec(rpcCommand["unlock"])
	if err != nil {
		return err
	}

	if resp.Ok == false {
		for _, m := range resp.Errors {
			return errors.New(m.Message)
		}
	}

	return nil
}

// LoadConfig loads a given configuration file locally or from
// an FTP or HTTP server. Format is either "set" or "text."
func (j *Junos) LoadConfig(path, format string) error {
	var command string
	switch format {
	case "set":
		if strings.Contains(path, "tp://") {
			command = fmt.Sprintf(rpcCommand["load-config-url-set"], path)
		} else {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			command = fmt.Sprintf(rpcCommand["load-config-local-set"], string(data))
		}
	case "text":
		if strings.Contains(path, "tp://") {
			command = fmt.Sprintf(rpcCommand["load-config-url-text"], path)
		} else {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			command = fmt.Sprintf(rpcCommand["load-config-local-text"], string(data))
		}
	case "xml":
		if strings.Contains(path, "tp://") {
			command = fmt.Sprintf(rpcCommand["load-config-url-xml"], path)
		} else {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			command = fmt.Sprintf(rpcCommand["load-config-local-xml"], string(data))
		}
	}

	reply, err := j.Exec(command)
	if err != nil {
		return err
	}

	if reply.Ok == false {
		for _, m := range reply.Errors {
			return errors.New(m.Message)
		}
	}

	return nil
}

// RollbackConfig loads and commits the configuration of a given rollback or rescue state.
func (j *Junos) RollbackConfig(option interface{}) error {
	var command string
	switch option.(type) {
	case int:
		command = fmt.Sprintf(rpcCommand["rollback-config"], option)
	case string:
		command = fmt.Sprintf(rpcCommand["rescue-config"])
	}

	reply, err := j.Exec(command)
	if err != nil {
		return err
	}

	err = j.Commit()
	if err != nil {
		return err
	}

	if reply.Ok == false {
		for _, m := range reply.Errors {
			return errors.New(m.Message)
		}
	}

	return nil
}

// RollbackDiff compares the current active configuration to a given rollback configuration.
func (j *Junos) RollbackDiff(compare int) (string, error) {
	rb := &rollbackXML{}
	command := fmt.Sprintf(rpcCommand["get-rollback-information-compare"], compare)
	reply, err := j.Exec(command)
	if err != nil {
		return "", err
	}

	if reply.Ok == false {
		for _, m := range reply.Errors {
			return "", errors.New(m.Message)
		}
	}

	err = xml.Unmarshal([]byte(reply.Data), rb)
	if err != nil {
		return "", err
	}

	return rb.Config, nil
}

// Command runs any operational mode command, such as "show" or "request."
// Format is either "text" or "xml".
func (j *Junos) Command(cmd, format string) (string, error) {
	c := &commandXML{}
	var command string

	switch format {
	case "xml":
		command = fmt.Sprintf(rpcCommand["command-xml"], cmd)
	default:
		command = fmt.Sprintf(rpcCommand["command"], cmd)
	}
	reply, err := j.Exec(command)
	if err != nil {
		return err
	}

	if reply.Ok == false {
		for _, m := range reply.Errors {
			return "", errors.New(m.Message)
		}
	}

	err = xml.Unmarshal([]byte(reply.Data), &c)
	if err != nil {
		return "", err
	}

	if c.Config == "" {
		return "No output available.", nil
	}

	return c.Config, nil
}

// Close disconnects our session to the device.
func (j *Junos) Close() {
	j.Transport.Close()
}

// Software displays basic information about the device, such as software, hardware, etc.
func (j *Junos) Software() (*software, error) {
	data := &software{}
	reply, err := j.Exec(rpcCommand["software"])
	if err != nil {
		return nil, err
	}

	if reply.Ok == false {
		for _, m := range reply.Errors {
			return nil, errors.New(m.Message)
		}
	}

	err = xml.Unmarshal([]byte(reply.Data), &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
