package junos

import (
	"encoding/xml"
)

// DeviceList holds a []Device slice of each individual device.
type DeviceList struct {
	XMLName xml.Name `xml:"devices"`
	Device  []Device `xml:"device"`
}

// Device holds all the information about each device.
type Device struct {
	ID        string `xml:"key,attr"`
	Family    string `xml:"deviceFamily"`
	Version   string `xml:"OSVersion"`
	Platform  string `xml:"platform"`
	Serial    string `xml:"serialNumber"`
	IPAddress string `xml:"ipAddr"`
	Name      string `xml:"name"`
}

// Devices
func (s *JunosSpace) Devices() (*DeviceList, error) {
	var devices DeviceList
	data, err := s.APICall("device-management/devices")
	if err != nil {
		return nil, err
	}

	err = xml.Unmarshal(data, &devices)
	if err != nil {
		return nil, err
	}

	return devices, nil
}
