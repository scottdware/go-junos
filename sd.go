package junos

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// Addresses contains a list of address objects.
type Addresses struct {
	Addresses []Address `xml:"address"`
}

// An Address contains information about each individual address object.
type Address struct {
	ID          int    `xml:"id"`
	Name        string `xml:"name"`
	AddressType string `xml:"address-type"`
	Description string `xml:"description"`
	IPAddress   string `xml:"ip-address"`
}

// Services contains a list of service objects.
type Services struct {
	Services []Service `xml:"service"`
}

// A Service contains information about each individual service object.
type Service struct {
	ID          int    `xml:"id"`
	Name        string `xml:"name"`
	IsGroup     bool   `xml:"is-group"`
	Description string `xml:"description"`
}

// A Policy contains information about each individual firewall policy.
type Policy struct {
	ID          int    `xml:"id"`
	Name        string `xml:"name"`
	Description string `xml:"description"`
}

// Policies contains a list of firewall policies.
type Policies struct {
	Policies []Policy `xml:"firewall-policy"`
}

// SecurityDevices contains a list of security devices.
type SecurityDevices struct {
	XMLName xml.Name         `xml:"devices"`
	Devices []SecurityDevice `xml:"device"`
}

// A SecurityDevice contains information about each individual security device.
type SecurityDevice struct {
	ID        int    `xml:"id"`
	Family    string `xml:"device-family"`
	Platform  string `xml:"platform"`
	IPAddress string `xml:"device-ip"`
	Name      string `xml:"name"`
}

// Variables contains a list of all polymorphic (variable) objects.
type Variables struct {
	Variables []Variable `xml:"variable-definition"`
}

// A Variable contains information about each individual polymorphic (variable) object.
type Variable struct {
	ID          int    `xml:"id"`
	Name        string `xml:"name"`
	Description string `xml:"description"`
}

type existingVariable struct {
	XMLName	xml.Name `xml:"variable-definition"`
	Name               string           `xml:"name"`
	Description        string           `xml:"description"`
	Type               string           `xml:"type"`
	Version            int              `xml:"edit-version"`
	DefaultName        string           `xml:"default-name"`
	DefaultValue       string              `xml:"default-value-detail>default-value"`
	VariableValuesList []variableValues `xml:"variable-values-list>variable-values"`
}

type variableValues struct {
	XMLName xml.Name `xml:"variable-values"`
	DeviceMOID    string `xml:"device>moid"`
	DeviceName    string `xml:"device>name"`
	VariableValue string    `xml:"variable-value-detail>variable-value"`
	VariableName  string `xml:"variable-value-detail>name"`
}

// addressesXML is XML we send (POST) for creating an address object.
var addressesXML = `
<address>
    <name>%s</name>
    <address-type>%s</address-type>
    <host-name/>
    <edit-version/>
    <members/>
    <address-version>IPV4</address-version>
    <definition-type>CUSTOM</definition-type>
    <ip-address>%s</ip-address>
    <description>%s</description>
</address>
`

// serviceXML is XML we send (POST) for creating a service object.
var serviceXML = `
<service>
    <name>%s</name>
    <description>%s</description>
    <is-group>false</is-group>
    <protocols>
        <protocol>
            <name>%s</name>
            <dst-port>%s</dst-port>
            <sunrpc-protocol-type>%s</sunrpc-protocol-type>
            <msrpc-protocol-type>%s</msrpc-protocol-type>
            <protocol-number>%d</protocol-number>
            <protocol-type>%s</protocol-type>
            <disable-timeout>%s</disable-timeout>
            %s
        </protocol>
    </protocols>
</service>
`

// addressGroupXML is XML we send (POST) for adding an address group.
var addressGroupXML = `
<address>
    <name>%s</name>
    <address-type>GROUP</address-type>
    <host-name/>
    <edit-version/>
    <address-version>IPV4</address-version>
    <definition-type>CUSTOM</definition-type>
    <description>%s</description>
</address>
`

// serviceGroupXML is XML we send (POST) for adding a service group.
var serviceGroupXML = `
<service>
    <name>%s</name>
    <is-group>true</is-group>
    <description>%s</description>
</service>
`

// removeXML is XML we send (POST) for removing an address or service from a group.
var removeXML = `
<diff>
    <remove sel="%s/members/member[name='%s']"/>
</diff>
`

// addGroupMemberXML is XML we send (POST) for adding addresses or services to a group.
var addGroupMemberXML = `
<diff>
    <add sel="%s/members">
        <member>
            <name>%s</name>
        </member>
    </add>
</diff>
`

// renameXML is XML we send (POST) for renaming an address or service object.
var renameXML = `
<diff>
    <replace sel="%s/name">
        <name>%s</name>
    </replace>
</diff>
`

// updateDeviceXML is XML we send (POST) for updating a security device.
var updateDeviceXML = `
<update-devices>
    <sd-ids>
        <id>%d</id>
    </sd-ids>
    <service-types>
        <service-type>POLICY</service-type>
    </service-types>
    <update-options>
        <enable-policy-rematch-srx-only>false</enable-policy-rematch-srx-only>
    </update-options>
</update-devices>
`

// publishPolicyXML is XML we send (POST) for publishing a changed policy.
var publishPolicyXML = `
<publish>
    <policy-ids>
        <policy-id>%d</policy-id>
    </policy-ids>
</publish>
`

// createVariableXML is the XML we send (POST) for adding a new variable object.
var createVariableXML = `
<variable-definition>
    <name>%s</name>
    <type>%s</type>
	<description>%s</description>
    <context>DEVICE</context>
    <default-name>%s</default-name>
    <default-value-detail>
        <default-value>%d</default-value>
    </default-value-detail>
</variable-definition>
`

var modifyVariableXML = `
<variable-definition>
    <name>%s</name>
    <type>%s</type>
	<description>%s</description>
	<edit-version>%d</edit-version>
    <context>DEVICE</context>
    <default-name>%s</default-name>
    <default-value-detail>
        <default-value>%s</default-value>
    </default-value-detail>
	<variable-values-list>
	%s
	</variable-values-list>
</variable-definition>
`

// getObjectID returns the ID of the address or service object.
func (s *JunosSpace) getObjectID(object interface{}, otype bool) (int, error) {
	var err error
	var objectID int
	var services *Services
	ipRegex := regexp.MustCompile(`(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\/\d+)`)
	if !otype {
		services, err = s.Services(object.(string))
	}
	objects, err := s.Addresses(object.(string))
	if err != nil {
		return 0, err
	}

	switch object.(type) {
	case int:
		objectID = object.(int)
	case string:
		if !otype {
			for _, o := range services.Services {
				if o.Name == object {
					objectID = o.ID
				}
			}
		}
		if ipRegex.MatchString(object.(string)) {
			for _, o := range objects.Addresses {
				if o.IPAddress == object {
					objectID = o.ID
				}
			}
		}
		for _, o := range objects.Addresses {
			if o.Name == object {
				objectID = o.ID
			}
		}
	}

	return objectID, nil
}

// getPolicyID returns the ID of a firewall policy.
func (s *JunosSpace) getPolicyID(object string) (int, error) {
	var err error
	var objectID int
	objects, err := s.Policies()
	if err != nil {
		return 0, err
	}

	for _, o := range objects.Policies {
		if o.Name == object {
			objectID = o.ID
		}
	}

	return objectID, nil
}

// getVariableID returns the ID of a polymorphic (variable) object.
func (s *JunosSpace) getVariableID(variable string) (int, error) {
	var err error
	var variableID int
	vars, err := s.Variables()
	if err != nil {
		return 0, err
	}

	for _, v := range vars.Variables {
		if v.Name == variable {
			variableID = v.ID
		}
	}

	return variableID, nil
}

func (s *JunosSpace) modifyVariableContent(data *existingVariable, moid, firewall, obj string, vid int) string {
	// var varValuesList = "<variable-values-list>"
	var varValuesList string
	for _, d := range data.VariableValuesList {
		varValuesList += fmt.Sprintf("<variable-values><device><moid>%s</moid><name>%s</name></device>", d.DeviceMOID, d.DeviceName)
		varValuesList += fmt.Sprintf("<variable-value-detail><variable-value>%s</variable-value><name>%s</name></variable-value-detail></variable-values>", d.VariableValue, d.VariableName)
	}
	varValuesList += fmt.Sprintf("<variable-values><device><moid>%s</moid><name>%s</name></device>", moid, firewall)
	varValuesList += fmt.Sprintf("<variable-value-detail><variable-value>net.juniper.jnap.sm.om.jpa.AddressEntity:%d</variable-value><name>%s</name></variable-value-detail></variable-values>", vid, obj)
	// varValuesList += "</variable-values-list>"
	
	return varValuesList
}

// Addresses queries the Junos Space server and returns all of the information
// about each address that is managed by Space.
func (s *JunosSpace) Addresses(filter string) (*Addresses, error) {
	var addresses Addresses
	p := url.Values{}
	p.Set("filter", "(global eq '')")

	if filter != "all" {
		p.Set("filter", fmt.Sprintf("(global eq '%s')", filter))
	}

	req := &APIRequest{
		Method: "get",
		URL:    fmt.Sprintf("/api/juniper/sd/address-management/addresses?%s", p.Encode()),
	}
	data, err := s.APICall(req)
	if err != nil {
		return nil, err
	}

	err = xml.Unmarshal(data, &addresses)
	if err != nil {
		return nil, err
	}

	return &addresses, nil
}

// AddAddress adds a new address object to Junos Space, and returns the Job ID.
func (s *JunosSpace) AddAddress(name, ip, desc string) (int, error) {
	var job jobID
	var addrType = "IPADDRESS"

	if strings.Contains(ip, "/") {
		addrType = "NETWORK"
	}

	address := fmt.Sprintf(addressesXML, name, addrType, ip, desc)
	req := &APIRequest{
		Method:      "post",
		URL:         "/api/juniper/sd/address-management/addresses",
		Body:        address,
		ContentType: contentAddress,
	}
	data, err := s.APICall(req)
	if err != nil {
		return 0, err
	}

	err = xml.Unmarshal(data, &job)
	if err != nil {
		return 0, err
	}

	return job.ID, nil
}

// AddService adds a new service object to Junos Space, and returns the Job ID. If adding just
// a single port, then enter in the same number in both the "low" and "high" parameters. For a
// range of ports, enter the starting port in "low" and the uppper limit in "high."
func (s *JunosSpace) AddService(proto, name string, low, high int, desc string, timeout int) (int, error) {
	var job jobID
	var port string
	var protoNumber int
	var inactivity string
	var secs string
	ptype := fmt.Sprintf("PROTOCOL_%s", strings.ToUpper(proto))
	protocol := strings.ToUpper(proto)

	protoNumber = 6
	if proto == "udp" {
		protoNumber = 17
	}

	port = strconv.Itoa(low)
	if low < high {
		port = fmt.Sprintf("%s-%s", strconv.Itoa(low), strconv.Itoa(high))
	}

	inactivity = "false"
	secs = fmt.Sprintf("<inactivity-timeout>%d</inactivity-timeout>", timeout)
	if timeout == 0 {
		inactivity = "true"
		secs = "<inactivity-timeout/>"
	}

	service := fmt.Sprintf(serviceXML, name, desc, name, port, protocol, protocol, protoNumber, ptype, inactivity, secs)
	req := &APIRequest{
		Method:      "post",
		URL:         "/api/juniper/sd/service-management/services",
		Body:        service,
		ContentType: contentService,
	}
	data, err := s.APICall(req)
	if err != nil {
		return 0, err
	}

	err = xml.Unmarshal(data, &job)
	if err != nil {
		return 0, err
	}

	return job.ID, nil
}

// AddGroup adds a new address or service group to Junos Space, and returns the Job ID.
func (s *JunosSpace) AddGroup(otype bool, name, desc string) error {
	uri := "/api/juniper/sd/address-management/addresses"
	addGroupXML := addressGroupXML
	content := contentAddress

	if !otype {
		uri = "/api/juniper/sd/service-management/services"
		addGroupXML = serviceGroupXML
		content = contentService
	}

	groupXML := fmt.Sprintf(addGroupXML, name, desc)
	req := &APIRequest{
		Method:      "post",
		URL:         uri,
		Body:        groupXML,
		ContentType: content,
	}
	_, err := s.APICall(req)

	if err != nil {
		return err
	}

	return nil
}

// ModifyObject modifies an existing address, service or group. The actions are as follows:
//
// "otype" is either "true" (for address) or "false" (for service)
//
// ModifyObject(otype, "add", "existing-group", "member-to-add")
// ModifyObject(otype, "remove", "existing-group", "member-to-remove")
// ModifyObject(otype, "rename", "old-name", "new-name")
// ModifyObject(otype, "delete", "object-to-delete")
func (s *JunosSpace) ModifyObject(otype bool, actions ...interface{}) error {
	var err error
	var uri string
	var content string
	var rel string
	objectID, err := s.getObjectID(actions[1], otype)
	if err != nil {
		return err
	}

	if objectID != 0 {
		var req *APIRequest
		uri = fmt.Sprintf("/api/juniper/sd/address-management/addresses/%d", objectID)
		content = contentAddressPatch
		rel = "address"

		if !otype {
			uri = fmt.Sprintf("/api/juniper/sd/service-management/services/%d", objectID)
			content = contentServicePatch
			rel = "service"
		}

		switch actions[0] {
		case "add":
			req = &APIRequest{
				Method:      "patch",
				URL:         uri,
				Body:        fmt.Sprintf(addGroupMemberXML, rel, actions[2]),
				ContentType: content,
			}
		case "remove":
			req = &APIRequest{
				Method:      "patch",
				URL:         uri,
				Body:        fmt.Sprintf(removeXML, rel, actions[2]),
				ContentType: content,
			}
		case "rename":
			req = &APIRequest{
				Method:      "patch",
				URL:         uri,
				Body:        fmt.Sprintf(renameXML, rel, actions[2]),
				ContentType: content,
			}
		case "delete":
			req = &APIRequest{
				Method: "delete",
				URL:    uri,
			}
		}

		_, err = s.APICall(req)
		if err != nil {
			return err
		}
	}

	return nil
}

// Services queries the Junos Space server and returns all of the information
// about each service that is managed by Space.
func (s *JunosSpace) Services(filter string) (*Services, error) {
	var services Services
	p := url.Values{}
	p.Set("filter", "(global eq '')")

	if filter != "all" {
		p.Set("filter", fmt.Sprintf("(global eq '%s')", filter))
	}

	req := &APIRequest{
		Method: "get",
		URL:    fmt.Sprintf("/api/juniper/sd/service-management/services?%s", p.Encode()),
	}
	data, err := s.APICall(req)
	if err != nil {
		return nil, err
	}

	err = xml.Unmarshal(data, &services)
	if err != nil {
		return nil, err
	}

	return &services, nil
}

// SecurityDevices queries the Junos Space server and returns all of the information
// about each security device that is managed by Space.
func (s *JunosSpace) SecurityDevices() (*SecurityDevices, error) {
	var devices SecurityDevices
	req := &APIRequest{
		Method: "get",
		URL:    "/api/juniper/sd/device-management/devices",
	}
	data, err := s.APICall(req)
	if err != nil {
		return nil, err
	}

	err = xml.Unmarshal(data, &devices)
	if err != nil {
		return nil, err
	}

	return &devices, nil
}

// Policies returns a list of all firewall policies managed by Junos Space.
func (s *JunosSpace) Policies() (*Policies, error) {
	var policies Policies
	req := &APIRequest{
		Method: "get",
		URL:    "/api/juniper/sd/fwpolicy-management/firewall-policies",
	}
	data, err := s.APICall(req)
	if err != nil {
		return nil, err
	}

	err = xml.Unmarshal(data, &policies)
	if err != nil {
		return nil, err
	}

	return &policies, nil
}

// PublishPolicy publishes a changed firewall policy. If "true" is specified for
// <update>, then Junos Space will also update the device.
func (s *JunosSpace) PublishPolicy(object interface{}, update bool) (int, error) {
	var err error
	var job jobID
	var id int
	var uri = "/api/juniper/sd/fwpolicy-management/publish"

	switch object.(type) {
	case int:
		id = object.(int)
	case string:
		id, err = s.getPolicyID(object.(string))
		if err != nil {
			return 0, err
		}
		if id == 0 {
			return 0, errors.New("no policy found")
		}
	}
	publish := fmt.Sprintf(publishPolicyXML, id)

	if update {
		uri = "/api/juniper/sd/fwpolicy-management/publish?update=true"
	}

	req := &APIRequest{
		Method:      "post",
		URL:         uri,
		Body:        publish,
		ContentType: contentPublish,
	}
	data, err := s.APICall(req)
	if err != nil {
		return 0, err
	}

	err = xml.Unmarshal(data, &job)
	if err != nil {
		return 0, errors.New("no policy changes to publish")
	}

	return job.ID, nil
}

// UpdateDevice will update a changed security device, synchronizing it with
// Junos Space.
func (s *JunosSpace) UpdateDevice(device interface{}) (int, error) {
	var job jobID
	deviceID, err := s.getDeviceID(device, true)
	if err != nil {
		return 0, err
	}

	update := fmt.Sprintf(updateDeviceXML, deviceID)
	req := &APIRequest{
		Method:      "post",
		URL:         "/api/juniper/sd/device-management/update-devices",
		Body:        update,
		ContentType: contentUpdateDevices,
	}
	data, err := s.APICall(req)
	if err != nil {
		return 0, err
	}

	err = xml.Unmarshal(data, &job)
	if err != nil {
		return 0, err
	}

	return job.ID, nil
}

// Variables returns a listing of all polymorphic (variable) objects.
func (s *JunosSpace) Variables() (*Variables, error) {
	var vars Variables
	req := &APIRequest{
		Method: "get",
		URL:    "/api/juniper/sd/variable-management/variable-definitions",
	}
	data, err := s.APICall(req)
	if err != nil {
		return nil, err
	}

	err = xml.Unmarshal(data, &vars)
	if err != nil {
		return nil, err
	}

	return &vars, nil
}

// AddVariable creates a new polymorphic object (variable) on the Junos Space server.
func (s *JunosSpace) AddVariable(name, vtype, desc, obj string) error {
	objID, err := s.getObjectID(obj, true)
	if err != nil {
		return err
	}

	varBody := fmt.Sprintf(createVariableXML, name, strings.ToUpper(vtype), desc, obj, objID)
	req := &APIRequest{
		Method:      "post",
		URL:         "/api/juniper/sd/variable-management/variable-definitions",
		Body:        varBody,
		ContentType: contentVariable,
	}
	_, err = s.APICall(req)
	if err != nil {
		return err
	}

	return nil
}

// ModifyVariable adds or deletes entries to the polymorphic (variable) object.
func (s *JunosSpace) ModifyVariable(actions ...interface{}) error {
	var err error
	var req *APIRequest
	var varData existingVariable
	var deviceID int
	var varID int
	var moid string
	var vid int
	var data []byte
	
	go func() {
		deviceID, _ = s.getDeviceID(actions[2].(string), true)
		// if err != nil {
			// return err
		// }
	}()
	
	go func() {
		moid = fmt.Sprintf("net.juniper.jnap.sm.om.jpa.SecurityDeviceEntity:%d", deviceID)
		varID, _ = s.getVariableID(actions[1].(string))
		// if err != nil {
			// return err
		// }
	}()
	
	go func() {
		vid, _ = s.getObjectID(actions[3].(string), true)
		// if err != nil {
			// return err
		// }

		existing := &APIRequest{
			Method: "get",
			URL:    fmt.Sprintf("/api/juniper/sd/variable-management/variable-definitions/%d", varID),
		}
		data, _ = s.APICall(existing)
		// if err != nil {
			// return err
		// }
	}()

	err = xml.Unmarshal(data, &varData)
	if err != nil {
		return err
	}
		
	varContent := s.modifyVariableContent(&varData, moid, actions[2].(string), actions[3].(string), vid)
	modifyVariable := fmt.Sprintf(modifyVariableXML, varData.Name, varData.Type, varData.Description, varData.Version, varData.DefaultName, varData.DefaultValue, varContent)

	fmt.Println(modifyVariable)
	if varID != 0 {
		switch actions[0].(string) {
		case "delete":
			req = &APIRequest{
				Method:      "delete",
				URL:         fmt.Sprintf("/api/juniper/sd/variable-management/variable-definitions/%d", varID),
				ContentType: contentVariable,
			}
		case "add":
			req = &APIRequest{
				Method:      "put",
				URL:         fmt.Sprintf("/api/juniper/sd/variable-management/variable-definitions/%d", varID),
				Body:        modifyVariable,
				ContentType: contentVariable,
			}
		}
	}

	_, err = s.APICall(req)
	if err != nil {
		return err
	}
	
	return nil
}
