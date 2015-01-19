package junos

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/Juniper/go-netconf/netconf"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
)

// Junos holds the connection information to our Junos device, as well
// as the platform and software version it is running.
type Junos struct {
	Session  *netconf.Session
	Hostname string
	MultiRE  bool
	Platform []routingEngine
}

// CommandXML parses our operational command responses.
type commandXML struct {
	Config string `xml:",innerxml"`
}

// commitError parses any errors during the commit process.
type commitError struct {
	Path    string `xml:"error-path"`
	Element string `xml:"error-info>bad-element"`
	Message string `xml:"error-message"`
}

// commitResults stores our errors if we have any.
type commitResults struct {
	XMLName xml.Name      `xml:"commit-results"`
	Errors  []commitError `xml:"rpc-error"`
}

// rollbackXML parses our rollback diff configuration.
type diffXML struct {
	XMLName xml.Name `xml:"rollback-information"`
	Config  string   `xml:"configuration-information>configuration-output"`
}

// routingEngine holds the device hardware and software information.
type routingEngine struct {
	Model   string
	Version string
}

// multiRE parses our XML if we have multiple routing engines.
type softwareMultiRE struct {
	XMLName xml.Name              `xml:"multi-routing-engine-results"`
	RE      []softwareRouteEngine `xml:"multi-routing-engine-item>software-information"`
}

// routeEngine holds all of our routing engine information.
type softwareRouteEngine struct {
	XMLName     xml.Name              `xml:"software-information"`
	Hostname    string                `xml:"host-name"`
	Platform    string                `xml:"product-model"`
	PackageInfo []softwarePackageInfo `xml:"package-information"`
}

// packageInfo holds our software information per routing engine.
type softwarePackageInfo struct {
	XMLName         xml.Name `xml:"package-information"`
	PackageName     []string `xml:"name"`
	SoftwareVersion []string `xml:"comment"`
}

// singleRE parses our XML if we only have one routing engine.
type softwareSingleRE struct {
	XMLName     xml.Name              `xml:"software-information"`
	Hostname    string                `xml:"host-name"`
	Platform    string                `xml:"product-model"`
	PackageInfo []softwarePackageInfo `xml:"package-information"`
}

// Close disconnects our session to the device.
func (j *Junos) Close() {
	j.Session.Transport.Close()
}

// Command runs any operational mode command, such as "show" or "request."
// Format is either "text" or "xml".
func (j *Junos) Command(cmd, format string) (string, error) {
	c := &commandXML{}
	var command string
	errMessage := "No output available. Please check the syntax of your command."

	switch format {
	case "xml":
		command = fmt.Sprintf(rpcCommand["command-xml"], cmd)
	default:
		command = fmt.Sprintf(rpcCommand["command"], cmd)
	}
	reply, err := j.Session.Exec(command)
	if err != nil {
		return errMessage, err
	}

	if reply.Errors != nil {
		for _, m := range reply.Errors {
			return errMessage, errors.New(m.Message)
		}
	}

	err = xml.Unmarshal([]byte(reply.Data), &c)
	if err != nil {
		return errMessage, err
	}

	if c.Config == "" {
		return errMessage, nil
	}

	return c.Config, nil
}

// Commit commits the configuration.
func (j *Junos) Commit() error {
	errs := &commitResults{}
	reply, err := j.Session.Exec(rpcCommand["commit"])
	if err != nil {
		return err
	}

	if reply.Errors != nil {
		for _, m := range reply.Errors {
			return errors.New(m.Message)
		}
	}

	err = xml.Unmarshal([]byte(reply.Data), &errs)
	if err != nil {
		return err
	}

	if errs.Errors != nil {
		for _, m := range errs.Errors {
			message := fmt.Sprintf("[%s]\n    %s\nError: %s", strings.Trim(m.Path, "[\r\n]"), strings.Trim(m.Element, "[\r\n]"), strings.Trim(m.Message, "[\r\n]"))
			return errors.New(message)
		}
	}

	return nil
}

// CommitAt commits the configuration at the specified <time>.
func (j *Junos) CommitAt(time string) error {
	errs := &commitResults{}
	command := fmt.Sprintf(rpcCommand["commit-at"], time)
	reply, err := j.Session.Exec(command)
	if err != nil {
		return err
	}

	if reply.Errors != nil {
		for _, m := range reply.Errors {
			return errors.New(m.Message)
		}
	}

	err = xml.Unmarshal([]byte(reply.Data), &errs)
	if err != nil {
		return err
	}

	if errs.Errors != nil {
		for _, m := range errs.Errors {
			message := fmt.Sprintf("[%s]\n    %s\nError: %s", strings.Trim(m.Path, "[\r\n]"), strings.Trim(m.Element, "[\r\n]"), strings.Trim(m.Message, "[\r\n]"))
			return errors.New(message)
		}
	}

	return nil
}

// CommitCheck checks the configuration for syntax errors.
func (j *Junos) CommitCheck() error {
	errs := &commitResults{}
	reply, err := j.Session.Exec(rpcCommand["commit-check"])
	if err != nil {
		return err
	}

	if reply.Errors != nil {
		for _, m := range reply.Errors {
			return errors.New(m.Message)
		}
	}

	err = xml.Unmarshal([]byte(reply.Data), &errs)
	if err != nil {
		return err
	}

	if errs.Errors != nil {
		for _, m := range errs.Errors {
			message := fmt.Sprintf("[%s]\n    %s\nError: %s", strings.Trim(m.Path, "[\r\n]"), strings.Trim(m.Element, "[\r\n]"), strings.Trim(m.Message, "[\r\n]"))
			return errors.New(message)
		}
	}

	return nil
}

// CommitConfirm rolls back the configuration after <delay> minutes.
func (j *Junos) CommitConfirm(delay int) error {
	errs := &commitResults{}
	command := fmt.Sprintf(rpcCommand["commit-confirm"], delay)
	reply, err := j.Session.Exec(command)
	if err != nil {
		return err
	}

	if reply.Errors != nil {
		for _, m := range reply.Errors {
			return errors.New(m.Message)
		}
	}

	err = xml.Unmarshal([]byte(reply.Data), &errs)
	if err != nil {
		return err
	}

	if errs.Errors != nil {
		for _, m := range errs.Errors {
			message := fmt.Sprintf("[%s]\n    %s\nError: %s", strings.Trim(m.Path, "[\r\n]"), strings.Trim(m.Element, "[\r\n]"), strings.Trim(m.Message, "[\r\n]"))
			return errors.New(message)
		}
	}

	return nil
}

// ConfigDiff compares the current active configuration to a given rollback configuration.
func (j *Junos) ConfigDiff(compare int) (string, error) {
	rb := &diffXML{}
	command := fmt.Sprintf(rpcCommand["get-rollback-information-compare"], compare)
	reply, err := j.Session.Exec(command)
	if err != nil {
		return "", err
	}

	if reply.Errors != nil {
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

// Facts returns information about the device, such as model and software.
func (j *Junos) Facts() {
	var str string
	fpcRegex := regexp.MustCompile(`^(EX).*`)
	srxRegex := regexp.MustCompile(`^(SRX).*`)
	mRegex := regexp.MustCompile(`^(M[X]?).*`)
	str += fmt.Sprintf("Multiple RE: %t\n\n", j.MultiRE)
	for i, p := range j.Platform {
		model := p.Model
		version := p.Version
		switch model {
		case fpcRegex.FindString(model):
			str += fmt.Sprintf("fpc%d\n--------------------------------------------------------------------------\n", i)
			str += fmt.Sprintf("Hostname: %s\nModel: %s\nVersion: %s\n\n", j.Hostname, model, version)
		case srxRegex.FindString(model):
			str += fmt.Sprintf("node%d\n--------------------------------------------------------------------------\n", i)
			str += fmt.Sprintf("Hostname: %s\nModel: %s\nVersion: %s\n\n", j.Hostname, model, version)
		case mRegex.FindString(model):
			str += fmt.Sprintf("re%d\n--------------------------------------------------------------------------\n", i)
			str += fmt.Sprintf("Hostname: %s\nModel: %s\nVersion: %s\n\n", j.Hostname, model, version)
		}
	}

	fmt.Println(str)
}

// GetConfig returns the full configuration, or starting a given <section>.
// Format can either be "text" or "xml."
func (j *Junos) GetConfig(section, format string) (string, error) {
	command := fmt.Sprintf("<rpc><get-configuration format=\"%s\"><configuration>", format)

	if section == "full" {
		command += "</configuration></get-configuration></rpc>"
	} else {
		command += fmt.Sprintf("<%s/></configuration></get-configuration></rpc>", section)
	}

	reply, err := j.Session.Exec(command)
	if err != nil {
		return "", err
	}

	if reply.Errors != nil {
		for _, m := range reply.Errors {
			return "", errors.New(m.Message)
		}
	}

	if format == "text" {
		output := &commandXML{}
		err = xml.Unmarshal([]byte(reply.Data), output)
		if err != nil {
			return "", err
		}

		return output.Config, nil
	}

	return reply.Data, nil
}

// LoadConfig loads a given configuration file locally or from
// an FTP or HTTP server. Format is either "set" "text" or "xml."
func (j *Junos) LoadConfig(path, format string, commit bool) error {
	var command string
	switch format {
	case "set":
		if strings.Contains(path, "tp://") {
			command = fmt.Sprintf(rpcCommand["config-url-set"], path)
		} else {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			command = fmt.Sprintf(rpcCommand["config-file-set"], string(data))
		}
	case "text":
		if strings.Contains(path, "tp://") {
			command = fmt.Sprintf(rpcCommand["config-url-text"], path)
		} else {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			command = fmt.Sprintf(rpcCommand["config-file-text"], string(data))
		}
	case "xml":
		if strings.Contains(path, "tp://") {
			command = fmt.Sprintf(rpcCommand["config-url-xml"], path)
		} else {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			command = fmt.Sprintf(rpcCommand["config-file-xml"], string(data))
		}
	}

	reply, err := j.Session.Exec(command)
	if err != nil {
		return err
	}

	if commit {
		err = j.Commit()
		if err != nil {
			return err
		}
	}

	if reply.Errors != nil {
		for _, m := range reply.Errors {
			return errors.New(m.Message)
		}
	}

	return nil
}

// Lock locks the candidate configuration.
func (j *Junos) Lock() error {
	reply, err := j.Session.Exec(rpcCommand["lock"])
	if err != nil {
		return err
	}

	if reply.Errors != nil {
		for _, m := range reply.Errors {
			return errors.New(m.Message)
		}
	}

	return nil
}

// NewSession establishes a new connection to a Junos device that we will use
// to run our commands against. NewSession also gathers software information
// pertaining to the device.
func NewSession(host, user, password string) (*Junos, error) {
	rex := regexp.MustCompile(`^.*\[(.*)\]`)
	s, err := netconf.DialSSH(host, netconf.SSHConfigPassword(user, password))
	if err != nil {
		log.Fatal(err)
	}

	reply, err := s.Exec(rpcCommand["version"])
	if err != nil {
		return nil, err
	}

	if reply.Errors != nil {
		for _, m := range reply.Errors {
			return nil, errors.New(m.Message)
		}
	}

	if strings.Contains(reply.Data, "multi-routing-engine-results") {
		facts := &softwareMultiRE{}
		err = xml.Unmarshal([]byte(reply.Data), facts)
		if err != nil {
			return nil, err
		}

		numRE := len(facts.RE)
		hostname := facts.RE[0].Hostname
		res := make([]routingEngine, 0, numRE)

		for i := 0; i < numRE; i++ {
			version := rex.FindStringSubmatch(facts.RE[i].PackageInfo[0].SoftwareVersion[0])
			model := strings.ToUpper(facts.RE[i].Platform)
			res = append(res, routingEngine{Model: model, Version: version[1]})
		}

		return &Junos{
			Session:  s,
			Hostname: hostname,
			MultiRE:  true,
			Platform: res,
		}, nil
	} else {
		facts := &softwareSingleRE{}
		err = xml.Unmarshal([]byte(reply.Data), facts)
		if err != nil {
			return nil, err
		}

		res := make([]routingEngine, 1)
		hostname := facts.Hostname
		version := rex.FindStringSubmatch(facts.PackageInfo[0].SoftwareVersion[0])
		model := strings.ToUpper(facts.Platform)
		res = append(res, routingEngine{Model: model, Version: version[1]})

		return &Junos{
			Session:  s,
			Hostname: hostname,
			MultiRE:  false,
			Platform: res,
		}, nil
	}
}

// Rescue will create or delete the rescue configuration given "save" or "delete."
func (j *Junos) Rescue(action string) error {
	command := fmt.Sprintf("rescue-%s", action)
	reply, err := j.Session.Exec(rpcCommand[command])
	if err != nil {
		return err
	}

	if reply.Errors != nil {
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
		if option == "rescue" {
			command = fmt.Sprintf(rpcCommand["rescue-config"])
		}
	}

	reply, err := j.Session.Exec(command)
	if err != nil {
		return err
	}

	err = j.Commit()
	if err != nil {
		return err
	}

	if reply.Errors != nil {
		for _, m := range reply.Errors {
			return errors.New(m.Message)
		}
	}

	return nil
}

// Unlock unlocks the candidate configuration.
func (j *Junos) Unlock() error {
	resp, err := j.Session.Exec(rpcCommand["unlock"])
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
