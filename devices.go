package junos

import (
	"encoding/xml"
	"fmt"
	"regexp"
)

// DeviceList holds a []Device slice of every device within Space.
type DeviceList struct {
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

// getDeviceID returns the given devices ID, which will be used for REST
// calls against it.
func (s *JunosSpace) getDeviceID(device interface{}) (int, error) {
	var err error
	var deviceID int
	ipRegex := regexp.MustCompile(`(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})`)
	devices, err := s.Devices()
	if err != nil {
		return 0, err
	}

	switch device.(type) {
	case int:
		deviceID = device.(int)
	case string:
		if ipRegex.MatchString(device.(string)) {
			for _, d := range devices.Devices {
				if d.IPAddress == device {
					deviceID = d.ID
				}
			}
		} else {
			for _, d := range devices.Devices {
				if d.Name == device {
					deviceID = d.ID
				}
			}
		}
	}

	return deviceID, nil
}

// AddDevice adds a new managed device to Junos Space, and returns the Job ID.
func (s *JunosSpace) AddDevice(host, user, password string) (int, error) {
	var job jobID
	ipRegex := regexp.MustCompile(`(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})`)
	inputXML := "<discover-devices>"

	if ipRegex.MatchString(host) {
		inputXML += fmt.Sprintf("<ipAddressDiscoveryTarget><ipAddress>%s</ipAddress></ipAddressDiscoveryTarget>", host)
	}

	inputXML += fmt.Sprintf("<hostNameDiscoveryTarget><hostName>%s</hostName></hostNameDiscoveryTarget>", host)
	inputXML += fmt.Sprintf("<sshCredential><userName>%s</userName><password>%s</password></sshCredential>", user, password)
	inputXML += "<manageDiscoveredSystemsFlag>true</manageDiscoveredSystemsFlag><usePing>true</usePing></discover-devices>"

	data, err := s.APIPost("space/device-management/discover-devices", inputXML, "discover-devices")
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
func (s *JunosSpace) Devices() (*DeviceList, error) {
	var devices DeviceList
	data, err := s.APIRequest("space/device-management/devices")
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
	deviceID, err := s.getDeviceID(device)
	if err != nil {
		return err
	}

	if deviceID != 0 {
		err = s.APIDelete(fmt.Sprintf("space/device-management/devices/%d", deviceID), "")
		if err != nil {
			return err
		}
	}

	return nil
}
