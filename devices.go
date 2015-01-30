package junos

import (
	"encoding/xml"
	"fmt"
	"regexp"
)

// Devices holds a []Device slice of every device within Space.
type Devices struct {
	XMLName xml.Name `xml:"devices"`
	Devices []Device `xml:"device"`
}

// Device holds all the information about each device within Space.
type Device struct {
	ID        int    `xml:"key,attr"`
	Family    string `xml:"deviceFamily"`
	Version   string `xml:"OSVersion"`
	Platform  string `xml:"platform"`
	Serial    string `xml:"serialNumber"`
	IPAddress string `xml:"ipAddr"`
	Name      string `xml:"name"`
}

// addDeviceIPXML is the XML used to add a device by IP address.
var addDeviceIPXML = `
<discover-devices>
    <ipAddressDiscoveryTarget>
        <ipAddress>%s</ipAddress>
    </ipAddressDiscoveryTarget>
    <sshCredential>
        <userName>%s</userName>
        <password>%s</password>
    </sshCredential>
    <manageDiscoveredSystemsFlag>true</manageDiscoveredSystemsFlag>
    <usePing>true</usePing>
</discover-devices>
`

// addDeviceHostXML is the XML used to add a device by hostname.
var addDeviceHostXML = `
<discover-devices>
    <hostNameDiscoveryTarget>
        <hostName>%s</hostName>
    </hostNameDiscoveryTarget>
    <sshCredential>
        <userName>%s</userName>
        <password>%s</password>
    </sshCredential>
    <manageDiscoveredSystemsFlag>true</manageDiscoveredSystemsFlag>
    <usePing>true</usePing>
</discover-devices>
`

// getDeviceID returns the given devices ID, which will be used for REST
// calls against it.
func (s *JunosSpace) getDeviceID(device interface{}, sd bool) (int, error) {
	var err error
	var deviceID int
	var sds *SecurityDevices
	ipRegex := regexp.MustCompile(`(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})`)
	if sd {
		sds, err = s.SecurityDevices()
	}
	devices, err := s.Devices()
	if err != nil {
		return 0, err
	}

	switch device.(type) {
	case int:
		deviceID = device.(int)
	case string:
		if sd {
			if ipRegex.MatchString(device.(string)) {
				for _, d := range sds.Devices {
					if d.IPAddress == device {
						deviceID = d.ID
					}
				}
			}
			for _, d := range sds.Devices {
				if d.Name == device {
					deviceID = d.ID
				}
			}
		}
		if ipRegex.MatchString(device.(string)) {
			for _, d := range devices.Devices {
				if d.IPAddress == device {
					deviceID = d.ID
				}
			}
		}
		for _, d := range devices.Devices {
			if d.Name == device {
				deviceID = d.ID
			}
		}
	}

	return deviceID, nil
}

// AddDevice adds a new managed device to Junos Space, and returns the Job ID.
func (s *JunosSpace) AddDevice(host, user, password string) (int, error) {
	var job jobID
	var addDevice string
	ipRegex := regexp.MustCompile(`(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})`)

	if ipRegex.MatchString(host) {
		addDevice = addDeviceIPXML
	}

	addDevice = addDeviceHostXML

	req := &APIRequest{
		Method:      "post",
		URL:         "space/device-management/discover-devices",
		Body:        fmt.Sprintf(addDevice, host, user, password),
		ContentType: ContentDiscoverDevices,
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

// Devices queries the Junos Space server and returns all of the information
// about each device that is managed by Space.
func (s *JunosSpace) Devices() (*Devices, error) {
	var devices Devices
	req := &APIRequest{
		Method: "get",
		URL:    "space/device-management/devices",
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

// RemoveDevice removes a device from Junos Space. You can specify the device ID, name
// or IP address.
func (s *JunosSpace) RemoveDevice(device interface{}) error {
	var err error
	deviceID, err := s.getDeviceID(device, false)
	if err != nil {
		return err
	}

	if deviceID != 0 {
		req := &APIRequest{
			Method: "delete",
			URL:    fmt.Sprintf("space/device-management/devices/%d", deviceID),
		}
		_, err = s.APICall(req)
		if err != nil {
			return err
		}
	}

	return nil
}
