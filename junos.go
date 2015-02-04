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

// All of our RPC calls we use.
var (
	rpcCommand             = "<rpc><command format=\"text\">%s</command></rpc>"
	rpcCommandXML          = "<rpc><command format=\"xml\">%s</command></rpc>"
	rpcCommit              = "<rpc><commit-configuration/></rpc>"
	rpcCommitAt            = "<rpc><commit-configuration><at-time>%s</at-time></commit-configuration></rpc>"
	rpcCommitCheck         = "<rpc><commit-configuration><check/></commit-configuration></rpc>"
	rpcCommitConfirm       = "<rpc><commit-configuration><confirmed/><confirm-timeout>%d</confirm-timeout></commit-configuration></rpc>"
	rpcFactsRE             = "<rpc><get-route-engine-information/></rpc>"
	rpcFactsChassis        = "<rpc><get-chassis-inventory/></rpc>"
	rpcConfigFileSet       = "<rpc><load-configuration action=\"set\" format=\"text\"><configuration-set>%s</configuration-set></load-configuration></rpc>"
	rpcConfigFileText      = "<rpc><load-configuration format=\"text\"><configuration-text>%s</configuration-text></load-configuration></rpc>"
	rpcConfigFileXML       = "<rpc><load-configuration format=\"xml\"><configuration>%s</configuration></load-configuration></rpc>"
	rpcConfigURLSet        = "<rpc><load-configuration action=\"set\" format=\"text\" url=\"%s\"/></rpc>"
	rpcConfigURLText       = "<rpc><load-configuration format=\"text\" url=\"%s\"/></rpc>"
	rpcConfigURLXML        = "<rpc><load-configuration format=\"xml\" url=\"%s\"/></rpc>"
	rpcGetRescue           = "<rpc><get-rescue-information><format>text</format></get-rescue-information></rpc>"
	rpcGetRollback         = "<rpc><get-rollback-information><rollback>%d</rollback><format>text</format></get-rollback-information></rpc>"
	rpcGetRollbackCompare  = "<rpc><get-rollback-information><rollback>0</rollback><compare>%d</compare><format>text</format></get-rollback-information></rpc>"
	rpcHardware            = "<rpc><get-chassis-inventory/></rpc>"
	rpcLock                = "<rpc><lock><target><candidate/></target></lock></rpc>"
	rpcRescueConfig        = "<rpc><load-configuration rescue=\"rescue\"/></rpc>"
	rpcRescueDelete        = "<rpc><request-delete-rescue-configuration/></rpc>"
	rpcRescueSave          = "<rpc><request-save-rescue-configuration/></rpc>"
	rpcRollbackConfig      = "<rpc><load-configuration rollback=\"%d\"/></rpc>"
	rpcRoute               = "<rpc><get-route-engine-information/></rpc>"
	rpcSoftware            = "<rpc><get-software-information/></rpc>"
	rpcUnlock              = "<rpc><unlock><target><candidate/></target></unlock></rpc>"
	rpcVersion             = "<rpc><get-software-information/></rpc>"
)

// Junos contains our session state.
type Junos struct {
	Session      *netconf.Session
	Hostname     string
	RouteEngines int
	Platform     []routingEngine
}

type commandXML struct {
	Config string `xml:",innerxml"`
}

type commitError struct {
	Path    string `xml:"error-path"`
	Element string `xml:"error-info>bad-element"`
	Message string `xml:"error-message"`
}

type commitResults struct {
	XMLName xml.Name      `xml:"commit-results"`
	Errors  []commitError `xml:"rpc-error"`
}

type diffXML struct {
	XMLName xml.Name `xml:"rollback-information"`
	Config  string   `xml:"configuration-information>configuration-output"`
}

type routingEngine struct {
	Model   string
	Version string
}

type hardwareRouteEngines struct {
	XMLName xml.Name              `xml:"multi-routing-engine-results"`
	RE      []hardwareRouteEngine `xml:"multi-routing-engine-item>chassis-inventory"`
}

type hardwareRouteEngine struct {
	XMLName     xml.Name `xml:"chassis-inventory"`
	Serial      string   `xml:"chassis>serial-number"`
	Description string   `xml:"chassis>description"`
}

type versionRouteEngines struct {
	XMLName xml.Name             `xml:"multi-routing-engine-results"`
	RE      []versionRouteEngine `xml:"multi-routing-engine-item>software-information"`
}

type versionRouteEngine struct {
	XMLName     xml.Name             `xml:"software-information"`
	Hostname    string               `xml:"host-name"`
	Platform    string               `xml:"product-model"`
	PackageInfo []versionPackageInfo `xml:"package-information"`
}

type versionPackageInfo struct {
	XMLName         xml.Name `xml:"package-information"`
	PackageName     []string `xml:"name"`
	SoftwareVersion []string `xml:"comment"`
}

// Close disconnects our session to the device.
func (j *Junos) Close() {
	j.Session.Transport.Close()
}

// Command runs any operational mode command, such as "show" or "request."
// Format can be one of "text" or "xml."
func (j *Junos) Command(cmd, format string) (string, error) {
	var c commandXML
	var command = fmt.Sprintf(rpcCommand, cmd)
	errMessage := "No output available. Please check the syntax of your command."

	if format == "xml" {
		command = fmt.Sprintf(rpcCommandXML, cmd)
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
	var errs commitResults
	reply, err := j.Session.Exec(rpcCommit)
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
	var errs commitResults
	command := fmt.Sprintf(rpcCommitAt, time)
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
	var errs commitResults
	reply, err := j.Session.Exec(rpcCommitCheck)
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
	var errs commitResults
	command := fmt.Sprintf(rpcCommitConfirm, delay)
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
	var rb diffXML
	command := fmt.Sprintf(rpcGetRollbackCompare, compare)
	reply, err := j.Session.Exec(command)
	if err != nil {
		return "", err
	}

	if reply.Errors != nil {
		for _, m := range reply.Errors {
			return "", errors.New(m.Message)
		}
	}

	err = xml.Unmarshal([]byte(reply.Data), &rb)
	if err != nil {
		return "", err
	}

	return rb.Config, nil
}

// PrintFacts prints information about the device, such as model and software.
func (j *Junos) PrintFacts() {
	var str string
	fpcRegex := regexp.MustCompile(`^(EX).*`)
	srxRegex := regexp.MustCompile(`^(SRX).*`)
	mRegex := regexp.MustCompile(`^(M[X]?).*`)
	str += fmt.Sprintf("Routing Engines/FPC's: %d\n\n", j.RouteEngines)
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

// GetConfig returns the full configuration, or configuration starting at <section>.
// Format can be one of "text" or "xml."
func (j *Junos) GetConfig(section, format string) (string, error) {
	command := fmt.Sprintf("<rpc><get-configuration format=\"%s\"><configuration>", format)
	if section == "full" {
		command += "</configuration></get-configuration></rpc>"
	}
	command += fmt.Sprintf("<%s/></configuration></get-configuration></rpc>", section)

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
		var output commandXML
		err = xml.Unmarshal([]byte(reply.Data), &output)
		if err != nil {
			return "", err
		}

		return output.Config, nil
	}

	return reply.Data, nil
}

// LoadConfig loads a given configuration file from your local machine or
// a remote (FTP or HTTP server) location. Format can be one of "set" "text" or "xml."
func (j *Junos) LoadConfig(path, format string, commit bool) error {
	var command string
	switch format {
	case "set":
		if strings.Contains(path, "tp://") {
			command = fmt.Sprintf(rpcConfigURLSet, path)
		}
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		command = fmt.Sprintf(rpcConfigFileSet, string(data))
	case "text":
		if strings.Contains(path, "tp://") {
			command = fmt.Sprintf(rpcConfigURLText, path)
		}
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		command = fmt.Sprintf(rpcConfigFileText, string(data))
	case "xml":
		if strings.Contains(path, "tp://") {
			command = fmt.Sprintf(rpcConfigURLXML, path)
		}
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		command = fmt.Sprintf(rpcConfigFileXML, string(data))
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
	reply, err := j.Session.Exec(rpcLock)
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
// about the device.
func NewSession(host, user, password string) (*Junos, error) {
	rex := regexp.MustCompile(`^.*\[(.*)\]`)
	s, err := netconf.DialSSH(host, netconf.SSHConfigPassword(user, password))
	if err != nil {
		log.Fatal(err)
	}

	reply, err := s.Exec(rpcVersion)
	if err != nil {
		return nil, err
	}

	if reply.Errors != nil {
		for _, m := range reply.Errors {
			return nil, errors.New(m.Message)
		}
	}

	if strings.Contains(reply.Data, "multi-routing-engine-results") {
		var facts versionRouteEngines
		err = xml.Unmarshal([]byte(reply.Data), &facts)
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
			Session:      s,
			Hostname:     hostname,
			RouteEngines: numRE,
			Platform:     res,
		}, nil
	}

	var facts versionRouteEngine
	err = xml.Unmarshal([]byte(reply.Data), &facts)
	if err != nil {
		return nil, err
	}

	res := make([]routingEngine, 0)
	hostname := facts.Hostname
	version := rex.FindStringSubmatch(facts.PackageInfo[0].SoftwareVersion[0])
	model := strings.ToUpper(facts.Platform)
	res = append(res, routingEngine{Model: model, Version: version[1]})

	return &Junos{
		Session:      s,
		Hostname:     hostname,
		RouteEngines: 1,
		Platform:     res,
	}, nil
}

// Rescue will create or delete the rescue configuration given "save" or "delete."
func (j *Junos) Rescue(action string) error {
	command := fmt.Sprintf(rpcRescueSave)

	if action == "delete" {
		command = fmt.Sprintf(rpcRescueDelete)
	}

	reply, err := j.Session.Exec(command)
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
	var command = fmt.Sprintf(rpcRollbackConfig, option)

	if option == "rescue" {
		command = fmt.Sprintf(rpcRescueConfig)
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
	reply, err := j.Session.Exec(rpcUnlock)
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
