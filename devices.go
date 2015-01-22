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
	ID        int `xml:"key,attr"`
	Family    string `xml:"deviceFamily"`
	Version   string `xml:"OSVersion"`
	Platform  string `xml:"platform"`
	Serial    string `xml:"serialNumber"`
	IPAddress string `xml:"ipAddr"`
	Name      string `xml:"name"`
}

// AddDevice adds a new managed device to Junos Space, and returns the Job ID.
func (s *JunosSpace) AddDevice(host, user, password string) error {
    ipRegex := regexp.MustCompile(`(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})`)
    inputXML := "<discover-devices>"
    
    if ipRegex.MatchString(host) {
        inputXML += fmt.Sprintf("<ipAddressDiscoveryTarget><ipAddress>%s</ipAddress></ipAddressDiscoveryTarget>", host)
    } else {
        inputXML += fmt.Sprintf("<hostNameDiscoveryTarget><hostName>%s</hostName></hostNameDiscoveryTarget>", host)
    }
    
    inputXML += fmt.Sprintf("<sshCredential><userName>%s</userName><password>%s</password></sshCredential>", user, password)
    inputXML += "<manageDiscoveredSystemsFlag>true</manageDiscoveredSystemsFlag><usePing>true</usePing></discover-devices>"
    
    err := s.APIPost("device-management/discover-devices", inputXML, "discover-devices")
    if err != nil {
        return err
    }
    
    return nil
}

// Devices queries the Junos Space server and returns all of the information
// about each device that is managed by Space.
func (s *JunosSpace) Devices() (*DeviceList, error) {
	var devices DeviceList
	data, err := s.APIRequest("device-management/devices")
	if err != nil {
		return nil, err
	}

	err = xml.Unmarshal(data, &devices)
	if err != nil {
		return nil, err
	}

	return &devices, nil
}

// RemoveDevice removes a device from Junos Space.
func (s *JunosSpace) RemoveDevice(id int) error {
    err := s.APIDelete(fmt.Sprintf("device-management/devices/%d", id))
    if err != nil {
        return err
    }
    
    return nil
}
