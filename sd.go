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
	Hostname    string `xml:"host-name"`
}

// GroupMembers contains a list of all the members within a address or service
// group.
type GroupMembers struct {
	Members []Member `xml:"members>member"`
}

// Member contains information about each individual group member.
type Member struct {
	ID   int    `xml:"id"`
	Name string `xml:"name"`
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

// VariableManagement contains our session state when updating a polymorphic (variable) object.
type VariableManagement struct {
	Devices []SecurityDevice
	Space   *Space
}

// existingVariable contains all of our information in regards to said polymorphic (variable) object.
type existingVariable struct {
	XMLName            xml.Name         `xml:"variable-definition"`
	Name               string           `xml:"name"`
	Description        string           `xml:"description"`
	Type               string           `xml:"type"`
	Version            int              `xml:"edit-version"`
	DefaultName        string           `xml:"default-name"`
	DefaultValue       string           `xml:"default-value-detail>default-value"`
	VariableValuesList []variableValues `xml:"variable-values-list>variable-values"`
}

// variableValues contains the information for each device/object tied to the polymorphic (variable) object.
type variableValues struct {
	XMLName       xml.Name `xml:"variable-values"`
	DeviceMOID    string   `xml:"device>moid"`
	DeviceName    string   `xml:"device>name"`
	VariableValue string   `xml:"variable-value-detail>variable-value"`
	VariableName  string   `xml:"variable-value-detail>name"`
}

// existingAddress contains information about an address object before modification.
type existingAddress struct {
	Name        string `xml:"name"`
	EditVersion int    `xml:"edit-version"`
	Description string `xml:"description"`
}

// XML for creating an address object.
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

// XML for creating a dns-host address object.
var dnsXML = `
<address>
    <name>%s</name>
    <address-type>%s</address-type>
    <host-name>%s</host-name>
    <edit-version/>
    <members/>
    <address-version>IPV4</address-version>
    <definition-type>CUSTOM</definition-type>
    <ip-address/>
    <description>%s</description>
</address>
`

var modifyAddressXML = `
<address>
    <name>%s</name>
    <address-type>%s</address-type>
    <host-name/>
    <edit-version>%d</edit-version>
    <members/>
    <address-version>IPV4</address-version>
    <definition-type>CUSTOM</definition-type>
    <ip-address>%s</ip-address>
    <description>%s</description>
</address>
`

var modifyDNSXML = `
<address>
    <name>%s</name>
    <address-type>%s</address-type>
    <host-name>%s</host-name>
    <edit-version>%d</edit-version>
    <members/>
    <address-version>IPV4</address-version>
    <definition-type>CUSTOM</definition-type>
    <ip-address/>
    <description>%s</description>
</address>
`

// XML for creating a service object.
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

// XML for adding an address group.
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

// XML for adding a service group.
var serviceGroupXML = `
<service>
    <name>%s</name>
    <is-group>true</is-group>
    <description>%s</description>
</service>
`

// XML for removing an address or service from a group.
var removeXML = `
<diff>
    <remove sel="%s/members/member[name='%s']"/>
</diff>
`

// XML for adding addresses or services to a group.
var addGroupMemberXML = `
<diff>
    <add sel="%s/members">
        <member>
            <name>%s</name>
        </member>
    </add>
</diff>
`

// XML for renaming an address or service object.
var renameXML = `
<diff>
    <replace sel="%s/name">
        <name>%s</name>
    </replace>
</diff>
`

// XML for updating a security device.
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

// XML for publishing a changed policy.
var publishPolicyXML = `
<publish>
    <policy-ids>
        <policy-id>%d</policy-id>
    </policy-ids>
</publish>
`

// XML for adding a new variable object.
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

// XML for modifying variable objects.
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

// getDeviceID returns the ID of a managed device.
func (s *Space) getSDDeviceID(device interface{}) (int, error) {
	var err error
	var deviceID int
	ipRegex := regexp.MustCompile(`(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})`)
	devices, err := s.SecurityDevices()
	if err != nil {
		return 0, err
	}

	switch device.(type) {
	case int:
		deviceID = device.(int)
	case string:
		if ipRegex.MatchString(device.(string)) {
			for _, d := range devices.Devices {
				if d.IPAddress == device.(string) {
					deviceID = d.ID
				}
			}
		}
		for _, d := range devices.Devices {
			if d.Name == device.(string) {
				deviceID = d.ID
			}
		}
	}

	return deviceID, nil
}

// getObjectID returns the ID of the address or service object.
func (s *Space) getObjectID(object interface{}, otype string) (int, error) {
	var err error
	var objectID int
	var services *Services
	ipRegex := regexp.MustCompile(`(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\/\d+)`)
	if otype == "service" {
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
		if otype == "service" {
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
func (s *Space) getPolicyID(object string) (int, error) {
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
func (s *Space) getVariableID(variable string) (int, error) {
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

// getAddrTypeIP returns the address type and IP address of the given address object.
func (s *Space) getAddrTypeIP(address string) []string {
	var addrType string
	var ipaddr string
	r := regexp.MustCompile(`(\d+\.\d+\.\d+\.\d+)(\/\d+)?`)
	rDNS := regexp.MustCompile(`[-\w\.]*\.(com|net|org|us)$`)
	match := r.FindStringSubmatch(address)

	if rDNS.MatchString(address) {
		addrType = "DNS"
		ipaddr = address

		return []string{addrType, ipaddr}
	}

	switch match[2] {
	case "", "/32":
		addrType = "IPADDRESS"
		ipaddr = match[1]
	default:
		addrType = "NETWORK"
		ipaddr = address
	}

	return []string{addrType, ipaddr}
}

// modifyVariableContent creates the XML we use when modifying an existing polymorphic (variable) object.
func (s *Space) modifyVariableContent(data *existingVariable, moid, firewall, address string, vid int) string {
	var varValuesList string
	for _, d := range data.VariableValuesList {
		varValuesList += fmt.Sprintf("<variable-values><device><moid>%s</moid><name>%s</name></device>", d.DeviceMOID, d.DeviceName)
		varValuesList += fmt.Sprintf("<variable-value-detail><variable-value>%s</variable-value><name>%s</name></variable-value-detail></variable-values>", d.VariableValue, d.VariableName)
	}
	varValuesList += fmt.Sprintf("<variable-values><device><moid>%s</moid><name>%s</name></device>", moid, firewall)
	varValuesList += fmt.Sprintf("<variable-value-detail><variable-value>net.juniper.jnap.sm.om.jpa.AddressEntity:%d</variable-value><name>%s</name></variable-value-detail></variable-values>", vid, address)

	return varValuesList
}

// Addresses queries the Junos Space server and returns all of the information
// about each address that is managed by Space. Filter is optional, but if specified
// can help reduce the amount of objects returned.
func (s *Space) Addresses(filter ...string) (*Addresses, error) {
	var addresses Addresses
	p := url.Values{}
	p.Set("filter", "(global eq '')")

	if len(filter) > 0 {
		p.Set("filter", fmt.Sprintf("(global eq '%s')", filter[0]))
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

// AddAddress creates a new address object in Junos Space. Description is optional.
func (s *Space) AddAddress(name, ip string, description ...string) error {
	desc := ""
	re := regexp.MustCompile(`[-\w\.]*\.(com|net|org|us)$`)
	addrInfo := s.getAddrTypeIP(ip)

	if len(description) > 0 {
		desc = description[0]
	}

	address := fmt.Sprintf(addressesXML, name, addrInfo[0], addrInfo[1], desc)

	if re.MatchString(ip) {
		address = fmt.Sprintf(dnsXML, name, addrInfo[0], addrInfo[1], desc)
	}

	req := &APIRequest{
		Method:      "post",
		URL:         "/api/juniper/sd/address-management/addresses",
		Body:        address,
		ContentType: contentAddress,
	}
	_, err := s.APICall(req)
	if err != nil {
		return err
	}

	return nil
}

// EditAddress changes the IP/Network/FQDN of the given address object name.
func (s *Space) EditAddress(name, newip string) error {
	var existing existingAddress
	addrInfo := s.getAddrTypeIP(newip)
	re := regexp.MustCompile(`[-\w\.]*\.(com|net|org|us)$`)

	objectID, err := s.getObjectID(name, "address")
	if err != nil {
		return err
	}

	req := &APIRequest{
		URL:         fmt.Sprintf("/api/juniper/sd/address-management/addresses/%d", objectID),
		Method:      "get",
		ContentType: contentAddress,
	}

	data, err := s.APICall(req)
	if err != nil {
		return err
	}

	err = xml.Unmarshal(data, &existing)
	if err != nil {
		return err
	}

	updateContent := fmt.Sprintf(modifyAddressXML, existing.Name, addrInfo[0], existing.EditVersion, addrInfo[1], existing.Description)

	if re.MatchString(name) {
		updateContent = fmt.Sprintf(modifyDNSXML, existing.Name, addrInfo[0], existing.EditVersion, addrInfo[1], existing.Description)
	}

	modifyReq := &APIRequest{
		Method:      "put",
		URL:         fmt.Sprintf("/api/juniper/sd/address-management/addresses/%d", objectID),
		Body:        updateContent,
		ContentType: contentAddress,
	}

	_, err = s.APICall(modifyReq)
	if err != nil {
		return err
	}

	return nil
}

// AddService creates a new service object to Junos Space. For a single port, just enter in
// the number. For a range of ports, enter the low-high range in quotes like so: "10000-10002".
func (s *Space) AddService(protocol, name string, ports interface{}, description string, timeout int) error {
	var port string
	var protoNumber int
	var inactivity string
	var secs string
	ptype := fmt.Sprintf("PROTOCOL_%s", strings.ToUpper(protocol))
	protocol = strings.ToUpper(protocol)

	protoNumber = 6
	if protocol == "UDP" {
		protoNumber = 17
	}

	switch ports.(type) {
	case int:
		port = strconv.Itoa(ports.(int))
	case string:
		p := strings.Split(ports.(string), "-")
		port = fmt.Sprintf("%s-%s", p[0], p[1])
	}

	inactivity = "false"
	secs = fmt.Sprintf("<inactivity-timeout>%d</inactivity-timeout>", timeout)
	if timeout == 0 {
		inactivity = "true"
		secs = "<inactivity-timeout/>"
	}

	service := fmt.Sprintf(serviceXML, name, description, name, port, protocol, protocol, protoNumber, ptype, inactivity, secs)
	req := &APIRequest{
		Method:      "post",
		URL:         "/api/juniper/sd/service-management/services",
		Body:        service,
		ContentType: contentService,
	}
	_, err := s.APICall(req)
	if err != nil {
		return err
	}

	return nil
}

// AddGroup creates a new address or service group in Junos Space. Objecttype must be "address" or "service".
func (s *Space) AddGroup(grouptype, name string, description ...string) error {
	desc := ""
	uri := "/api/juniper/sd/address-management/addresses"
	addGroupXML := addressGroupXML
	content := contentAddress

	if len(description) > 0 {
		desc = description[0]
	}

	if grouptype == "service" {
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

// EditGroup adds or removes objects to/from an existing address or service group. Grouptype must be
// "address" or "service." Action must be add or remove.
func (s *Space) EditGroup(grouptype, action, object, name string) error {
	var err error
	var uri string
	var content string
	var rel string
	objectID, err := s.getObjectID(name, grouptype)
	if err != nil {
		return err
	}

	if objectID != 0 {
		var req *APIRequest
		uri = fmt.Sprintf("/api/juniper/sd/address-management/addresses/%d", objectID)
		content = contentAddressPatch
		rel = "address"

		if grouptype == "service" {
			uri = fmt.Sprintf("/api/juniper/sd/service-management/services/%d", objectID)
			content = contentServicePatch
			rel = "service"
		}

		switch action {
		case "add":
			req = &APIRequest{
				Method:      "patch",
				URL:         uri,
				Body:        fmt.Sprintf(addGroupMemberXML, rel, object),
				ContentType: content,
			}
		case "remove":
			req = &APIRequest{
				Method:      "patch",
				URL:         uri,
				Body:        fmt.Sprintf(removeXML, rel, object),
				ContentType: content,
			}
		}

		_, err = s.APICall(req)
		if err != nil {
			return err
		}
	}

	return nil
}

// RenameObject renames an address or service object to the given new name. Grouptype
// must be "address" or "service"
func (s *Space) RenameObject(grouptype, name, newname string) error {
	var err error
	var uri string
	var content string
	var rel string
	objectID, err := s.getObjectID(name, grouptype)
	if err != nil {
		return err
	}

	if objectID != 0 {
		var req *APIRequest
		uri = fmt.Sprintf("/api/juniper/sd/address-management/addresses/%d", objectID)
		content = contentAddressPatch
		rel = "address"

		if grouptype == "service" {
			uri = fmt.Sprintf("/api/juniper/sd/service-management/services/%d", objectID)
			content = contentServicePatch
			rel = "service"
		}

		req = &APIRequest{
			Method:      "patch",
			URL:         uri,
			Body:        fmt.Sprintf(renameXML, rel, newname),
			ContentType: content,
		}

		_, err = s.APICall(req)
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteObject removes an address or service object from Junos Space. Grouptype
// must be "address" or "service"
func (s *Space) DeleteObject(grouptype, name string) error {
	var err error
	var uri string
	objectID, err := s.getObjectID(name, grouptype)
	if err != nil {
		return err
	}

	if objectID != 0 {
		var req *APIRequest
		uri = fmt.Sprintf("/api/juniper/sd/address-management/addresses/%d", objectID)

		if grouptype == "service" {
			uri = fmt.Sprintf("/api/juniper/sd/service-management/services/%d", objectID)
		}

		req = &APIRequest{
			Method: "delete",
			URL:    uri,
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
func (s *Space) Services(filter ...string) (*Services, error) {
	var services Services
	p := url.Values{}
	p.Set("filter", "(global eq '')")

	if len(filter) > 0 {
		p.Set("filter", fmt.Sprintf("(global eq '%s')", filter[0]))
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

// GroupMembers lists all of the address or service objects within the
// given group. Grouptype must be "address" or "service".
func (s *Space) GroupMembers(grouptype, name string) (*GroupMembers, error) {
	var members GroupMembers
	objectID, err := s.getObjectID(name, grouptype)
	url := fmt.Sprintf("/api/juniper/sd/address-management/addresses/%d", objectID)

	if grouptype == "service" {
		url = fmt.Sprintf("/api/juniper/sd/service-management/services/%d", objectID)
	}

	req := &APIRequest{
		Method: "get",
		URL:    url,
	}
	data, err := s.APICall(req)
	if err != nil {
		return nil, err
	}

	err = xml.Unmarshal(data, &members)
	if err != nil {
		return nil, err
	}

	return &members, nil
}

// SecurityDevices queries the Junos Space server and returns all of the information
// about each security device that is managed by Space.
func (s *Space) SecurityDevices() (*SecurityDevices, error) {
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
func (s *Space) Policies() (*Policies, error) {
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
// update, then Junos Space will also update the device.
func (s *Space) PublishPolicy(object interface{}, update bool) (int, error) {
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
func (s *Space) UpdateDevice(device interface{}) (int, error) {
	var job jobID
	deviceID, err := s.getDeviceID(device)
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
func (s *Space) Variables() (*Variables, error) {
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
// The address option is a default address object that will be used. This address object must
// already exist on the server.
func (s *Space) AddVariable(name, address string, description ...string) error {
	desc := ""
	objectID, err := s.getObjectID(address, "address")
	if err != nil {
		return err
	}

	if len(description) > 0 {
		desc = description[0]
	}

	varBody := fmt.Sprintf(createVariableXML, name, "ADDRESS", desc, address, objectID)
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

// DeleteVariable removes the polymorphic (variable) object from Junos Space.
// If the variable object is in use by a policy, then it will not be deleted
// until you remove it from the policy.
func (s *Space) DeleteVariable(name string) error {
	var req *APIRequest
	varID, err := s.getVariableID(name)
	if err != nil {
		return err
	}

	req = &APIRequest{
		Method:      "delete",
		URL:         fmt.Sprintf("/api/juniper/sd/variable-management/variable-definitions/%d", varID),
		ContentType: contentVariable,
	}

	_, err = s.APICall(req)
	if err != nil {
		return err
	}

	return nil
}

// EditVariable creates a new state when adding/removing addresses to
// a polymorphic (variable) object. We do this to only get the list of
// security devices (SecurityDevices()) once, instead of call the function
// each time we want to modify a variable.
func (s *Space) EditVariable() (*VariableManagement, error) {
	devices, err := s.SecurityDevices()
	if err != nil {
		return nil, err
	}

	return &VariableManagement{
		Devices: devices.Devices,
		Space:   s,
	}, nil
}

// Add appends an address object to the given polymorphic (variable) object.
// Address is the address object you want to add, and name needs to be the variable
// object you wish to add the object to. You also must specify the device (firewall) that you
// want to associate the variable object to.
func (v *VariableManagement) Add(address, name, firewall string) error {
	var req *APIRequest
	var varData existingVariable
	var deviceID int

	varID, err := v.Space.getVariableID(name)
	if err != nil {
		return err
	}

	for _, d := range v.Devices {
		if d.Name == firewall {
			deviceID = d.ID
		}
	}
	moid := fmt.Sprintf("net.juniper.jnap.sm.om.jpa.SecurityDeviceEntity:%d", deviceID)

	vid, err := v.Space.getObjectID(address, "address")
	if err != nil {
		return err
	}

	existing := &APIRequest{
		Method: "get",
		URL:    fmt.Sprintf("/api/juniper/sd/variable-management/variable-definitions/%d", varID),
	}
	data, err := v.Space.APICall(existing)
	if err != nil {
		return err
	}

	err = xml.Unmarshal(data, &varData)
	if err != nil {
		return err
	}

	varContent := v.Space.modifyVariableContent(&varData, moid, firewall, address, vid)
	modifyVariable := fmt.Sprintf(modifyVariableXML, varData.Name, varData.Type, varData.Description, varData.Version, varData.DefaultName, varData.DefaultValue, varContent)

	req = &APIRequest{
		Method:      "put",
		URL:         fmt.Sprintf("/api/juniper/sd/variable-management/variable-definitions/%d", varID),
		Body:        modifyVariable,
		ContentType: contentVariable,
	}

	_, err = v.Space.APICall(req)
	if err != nil {
		return err
	}

	return nil
}
