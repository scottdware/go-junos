package junos

import (
	"encoding/xml"
)

// DeviceList holds a []Device slice of every device within Space.
type DeviceList struct {
	XMLName xml.Name `xml:"devices"`
	Devices  []Device `xml:"device"`
}

// Device holds all the information about each device within Space.
type Device struct {
	ID        string `xml:"key,attr"`
	Family    string `xml:"deviceFamily"`
	Version   string `xml:"OSVersion"`
	Platform  string `xml:"platform"`
	Serial    string `xml:"serialNumber"`
	IPAddress string `xml:"ipAddr"`
	Name      string `xml:"name"`
}

// Devices queries the Junos Space server and returns all of the information
// about each device that is managed by Space.
func (s *JunosSpace) Devices() (*DeviceList, error) {
	var devices DeviceList
	data, err := s.APICall("get", "device-management/devices", nil)
	if err != nil {
		return nil, err
	}

	err = xml.Unmarshal(data, &devices)
	if err != nil {
		return nil, err
	}

	return &devices, nil
}
