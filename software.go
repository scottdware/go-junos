package junos

import (
	"encoding/xml"
	"fmt"

	"github.com/scottdware/go-rested"
)

// SoftwarePackages contains a list of software packages managed by Junos Space.
type SoftwarePackages struct {
	Packages []SoftwarePackage `xml:"package"`
}

// A SoftwarePackage contains information about each individual software package.
type SoftwarePackage struct {
	ID       int    `xml:"key,attr"`
	Name     string `xml:"fileName"`
	Version  string `xml:"version"`
	Platform string `xml:"platformType"`
}

// SoftwareUpgrade consists of options available to use before issuing a software upgrade.
type SoftwareUpgrade struct {
	UseDownloaded bool // Use an image already staged on the device.
	Validate      bool // Check/don't check compatibility with current configuration.
	Reboot        bool // Reboot system after adding package.
	RebootAfter   int  // Reboot the system after "x" minutes.
	Cleanup       bool // Remove any pre-existing packages on the device.
	RemoveAfter   bool // Remove the package after successful installation.
}

// deployXML is XML we send (POST) for image deployment.
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

// removeStagedXML is XML we send (POST) for removing a staged image.
var removeStagedXML = `
<exec-remove>
    <devices>
        <device href="/api/space/device-management/devices/%d"/>
    </devices>
</exec-remove>
`

// stageXML is XML we send (POST) for staging an image on a device.
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

// getSoftwareID returns the ID of the software package.
func (s *Space) getSoftwareID(image string) (int, error) {
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
func (s *Space) DeploySoftware(device, image string, options *SoftwareUpgrade) (int, error) {
	r := rested.NewRequest()
	r.BasicAuth(s.User, s.Password)
	headers := map[string]string{
		"Content-Type": contentExecDeploy,
	}
	var job jobID
	deviceID, _ := s.getDeviceID(device)
	softwareID, _ := s.getSoftwareID(image)
	deploy := fmt.Sprintf(deployXML, deviceID, options.UseDownloaded, options.Validate, options.Reboot, options.RebootAfter, options.Cleanup, options.RemoveAfter)
	uri := fmt.Sprintf("https://%s/api/space/software-management/packages/%d/exec-deploy", s.Host, softwareID)

	resp := r.Send("post", uri, []byte(deploy), headers, nil)
	if resp.Error != nil {
		return 0, resp.Error
	}

	err := xml.Unmarshal(resp.Body, &job)
	if err != nil {
		return 0, err
	}

	return job.ID, nil
}

// RemoveStagedSoftware will delete the staged software image on the device.
func (s *Space) RemoveStagedSoftware(device, image string) (int, error) {
	r := rested.NewRequest()
	r.BasicAuth(s.User, s.Password)
	headers := map[string]string{
		"Content-Type": contentExecRemove,
	}
	var job jobID
	deviceID, _ := s.getDeviceID(device)
	softwareID, _ := s.getSoftwareID(image)
	remove := fmt.Sprintf(removeStagedXML, deviceID)
	uri := fmt.Sprintf("https://%s/api/space/software-management/packages/%d/exec-remove", s.Host, softwareID)

	resp := r.Send("post", uri, []byte(remove), headers, nil)
	if resp.Error != nil {
		return 0, resp.Error
	}

	err := xml.Unmarshal(resp.Body, &job)
	if err != nil {
		return 0, err
	}

	return job.ID, nil
}

// Software queries the Junos Space server and returns all of the information
// about each software image that Space manages.
func (s *Space) Software() (*SoftwarePackages, error) {
	r := rested.NewRequest()
	r.BasicAuth(s.User, s.Password)
	var software SoftwarePackages
	uri := fmt.Sprintf("https://%s/api/space/software-management/packages", s.Host)

	resp := r.Send("get", uri, nil, nil, nil)
	if resp.Error != nil {
		return nil, resp.Error
	}

	err := xml.Unmarshal(resp.Body, &software)
	if err != nil {
		return nil, err
	}

	return &software, nil
}

// StageSoftware loads the given software image onto the device but does not
// upgrade it. The package is placed in the /var/tmp directory.
func (s *Space) StageSoftware(device, image string, cleanup bool) (int, error) {
	r := rested.NewRequest()
	r.BasicAuth(s.User, s.Password)
	headers := map[string]string{
		"Content-Type": contentExecStage,
	}
	var job jobID
	deviceID, _ := s.getDeviceID(device)
	softwareID, _ := s.getSoftwareID(image)
	stage := fmt.Sprintf(stageXML, deviceID, cleanup)
	uri := fmt.Sprintf("https://%s/api/space/software-management/packages/%d/exec-stage", s.Host, softwareID)

	resp := r.Send("post", uri, []byte(stage), headers, nil)
	if resp.Error != nil {
		return 0, resp.Error
	}

	err := xml.Unmarshal(resp.Body, &job)
	if err != nil {
		return 0, err
	}

	return job.ID, nil
}
