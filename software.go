package junos

import (
	"encoding/xml"
	"fmt"
)

// SoftwarePackages holds a list of every software image managed within Space.
type SoftwarePackages struct {
	Packages []SoftwarePackage `xml:"package"`
}

// SoftwarePackage holds all the information about each software image within Space.
type SoftwarePackage struct {
	ID       int    `xml:"key,attr"`
	Name     string `xml:"fileName"`
	Version  string `xml:"version"`
	Platform string `xml:"platformType"`
}

// SoftwareUpgrade holds all of the options available for deploying/upgrading
// the software on a device through Junos Space.
type SoftwareUpgrade struct {
    // Use an image already staged on the device.
	UseDownloaded bool
    // Check/don't check compatibility with current configuration.
	Validate      bool
    // Reboot system after adding package.
	Reboot        bool
    // Reboot the system after "x" minutes.
	RebootAfter   int
    // Remove any pre-existing packages on the device.
	Cleanup       bool
    // Remove the package after successful installation.
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

// getSoftwareID returns the given software image ID, which will be used for REST
// calls against it.
func (s *JunosSpace) getSoftwareID(image string) (int, error) {
	var err error
	var softwareID int
	images, err := s.Software()
	if err != nil {
		return 0, err
	}

	for _, sw := range images.Packages {
		if sw.Name == image {
			softwareID = sw.ID
		}
	}

	return softwareID, nil
}

// DeploySoftware starts the upgrade process on the device, using the given image along
// with the options specified.
func (s *JunosSpace) DeploySoftware(device, image string, options *SoftwareUpgrade) (int, error) {
	var job jobID
	deviceID, _ := s.getDeviceID(device, false)
	softwareID, _ := s.getSoftwareID(image)
	deploy := fmt.Sprintf(deployXML, deviceID, options.UseDownloaded, options.Validate, options.Reboot, options.RebootAfter, options.Cleanup, options.RemoveAfter)
	req := &APIRequest{
		Method:      "post",
		URL:         fmt.Sprintf("/api/space/software-management/packages/%d/exec-deploy", softwareID),
		Body:        deploy,
		ContentType: ContentExecDeploy,
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

// RemoveStagedSoftware will delete the staged software image on the device.
func (s *JunosSpace) RemoveStagedSoftware(device, image string) (int, error) {
	var job jobID
	deviceID, _ := s.getDeviceID(device, false)
	softwareID, _ := s.getSoftwareID(image)
	remove := fmt.Sprintf(removeStagedXML, deviceID)
	req := &APIRequest{
		Method:      "post",
		URL:         fmt.Sprintf("/api/space/software-management/packages/%d/exec-remove", softwareID),
		Body:        remove,
		ContentType: ContentExecRemove,
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

// Software queries the Junos Space server and returns all of the information
// about each software image that Space manages.
func (s *JunosSpace) Software() (*SoftwarePackages, error) {
	var software SoftwarePackages
	req := &APIRequest{
		Method: "get",
		URL:    "/api/space/software-management/packages",
	}
	data, err := s.APICall(req)
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
	deviceID, _ := s.getDeviceID(device, false)
	softwareID, _ := s.getSoftwareID(image)
	stage := fmt.Sprintf(stageXML, deviceID, cleanup)
	req := &APIRequest{
		Method:      "post",
		URL:         fmt.Sprintf("/api/space/software-management/packages/%d/exec-stage", softwareID),
		Body:        stage,
		ContentType: ContentExecStage,
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
