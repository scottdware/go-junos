package junos

import (
	"encoding/xml"
	"fmt"
)

// SoftwarePackages holds a []Device slice of every device within Space.
type SoftwarePackages struct {
	Packages []SoftwarePackage `xml:"package"`
}

// SoftwarePackage holds all the information about each device within Space.
type SoftwarePackage struct {
	ID       int    `xml:"key,attr"`
	Name     string `xml:"fileName"`
	Version  string `xml:"version"`
	Platform string `xml:"platformType"`
}

// SoftwareOptions holds all of the options available for deploying/upgrading
// the software on a device.
type SoftwareOptions struct {
	UseDownloaded bool
	Validate      bool
	Reboot        bool
	RebootAfter   int
	Cleanup       bool
	RemoveAfter   bool
}

// deployXML is what we send to the server for image deployment.
var deployXML = `
<exec-deploy>
    <devices>
        <device href= "/api/space/device-management/devices/%d"/>
    </devices> 
    <deployOptions> 
        <useAlreadyDownloaded>%t</useAlreadyDownloaded>
        <validate>%t</validate>
        <bestEffortLoad>false</bestEffortLoad>
        <snapShotRequired>false</snapShotRequired>
        <rebootDevice>%t</rebootDevice>
        <rebootAfterXMinutes>%d</rebootAfterXMinutes>
        <cleanUpExistingOnDevice>%t</cleanUpExistingOnDevice>
        <removePkgAfterInstallation>%t</removePkgAfterInstallation>
    </deployOptions>
</exec-deploy>
`

// removeStagedXML is what we send to the server for removing a staged image.
var removeStagedXML = `
<exec-remove>
    <devices>
        <device href="/api/space/device-management/devices/%d"/>
    </devices>
</exec-remove>
`

// stageXML is what we send to the server for staging an image on a device.
var stageXML = `
<exec-stage>
    <devices>
        <device href="/api/space/device-management/devices/%d"/>
    </devices>
    <stageOptions>
        <cleanUpExistingOnDevice>%t</cleanUpExistingOnDevice>
    </stageOptions>
</exec-stage>
`

// getSoftwareID returns the given software images ID, which will be used for REST
// calls against it.
func (s *JunosSpace) getSoftwareID(image string) (int, error) {
	var err error
	var softwareID int
	images, err := s.Software()
	if err != nil {
		return -1, err
	}

	for _, sw := range images.Packages {
		if sw.Name == image {
			softwareID = sw.ID
		}
	}

	return softwareID, nil
}

// DeploySoftware starts the upgrade process on the device, using the given image.
func (s *JunosSpace) DeploySoftware(device, image string, options *SoftwareOptions) (int, error) {
	var job jobID
	deviceID, _ := s.getDeviceID(device)
	softwareID, _ := s.getSoftwareID(image)
	deploy := fmt.Sprintf(deployXML, deviceID, options.UseDownloaded, options.Validate, options.Reboot, options.RebootAfter, options.Cleanup, options.RemoveAfter)
	data, err := s.APIPost(fmt.Sprintf("software-management/packages/%d/exec-deploy", softwareID), deploy, "exec-deploy")
	if err != nil {
		return -1, err
	}

	err = xml.Unmarshal(data, &job)
	if err != nil {
		return -1, err
	}

	return job.ID, nil
}

// RemoveStagedSoftware will delete the staged software image on the device.
func (s *JunosSpace) RemoveStagedSoftware(device, image string) (int, error) {
	var job jobID
	deviceID, _ := s.getDeviceID(device)
	softwareID, _ := s.getSoftwareID(image)
	remove := fmt.Sprintf(removeStagedXML, deviceID)
	data, err := s.APIPost(fmt.Sprintf("software-management/packages/%d/exec-remove", softwareID), remove, "exec-remove")
	if err != nil {
		return -1, err
	}

	err = xml.Unmarshal(data, &job)
	if err != nil {
		return -1, err
	}

	return job.ID, nil
}

// Software queries the Junos Space server and returns all of the information
// about each software image that Space manages.
func (s *JunosSpace) Software() (*SoftwarePackages, error) {
	var software SoftwarePackages
	data, err := s.APIRequest("software-management/packages")
	if err != nil {
		return nil, err
	}

	err = xml.Unmarshal(data, &software)
	if err != nil {
		return nil, err
	}

	return &software, nil
}

// StageSoftware loads the given software image onto the device but does not
// upgrade it. The package is placed in the /var/tmp directory.
func (s *JunosSpace) StageSoftware(device, image string, cleanup bool) (int, error) {
	var job jobID
	deviceID, _ := s.getDeviceID(device)
	softwareID, _ := s.getSoftwareID(image)
	stage := fmt.Sprintf(stageXML, deviceID, cleanup)
	data, err := s.APIPost(fmt.Sprintf("software-management/packages/%d/exec-stage", softwareID), stage, "exec-stage")
	if err != nil {
		return -1, err
	}

	err = xml.Unmarshal(data, &job)
	if err != nil {
		return -1, err
	}

	return job.ID, nil
}
