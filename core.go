package linksys

// RouterInfo represents information about the router.
type RouterInfo struct {
	Description     string   `json:"description"`
	FirmwareDate    string   `json:"firmwareDate"`
	FirmwareVersion string   `json:"firmwareVersion"`
	HardwareVersion string   `json:"hardwareVersion"`
	Manufacturer    string   `json:"manufacturer"`
	ModelNumber     string   `json:"modelNumber"`
	SerialNumber    string   `json:"serialNumber"`
	Services        []string `json:"services"`
}

//
// Unauthorized Requests
//

// GetAdminPasswordHint returns the admin password hint.
func (client Client) GetAdminPasswordHint() (string, error) {
	var res struct {
		Hint string `json:"passwordHint"`
	}
	err := client.MakeRequest("core/GetAdminPasswordHint", nil, &res)
	return res.Hint, err
}

// GetDeviceInfo returns information about the router.
func (client Client) GetDeviceInfo() (RouterInfo, error) {
	var info RouterInfo
	err := client.MakeRequest("core/GetDeviceInfo", nil, &info)
	return info, err
}

//
// Authorized Requests
//

// This action validates the password provided with the X-JNAP-Authorization.
func (client Client) CheckAdminPassword() error {
	err := client.MakeRequest("core/CheckAdminPassword", nil, nil)
	if err != nil {
		return err
	}
	return nil
}

// This action instructs the router to reboot.
func (client Client) Reboot() error {
	return client.MakeRequest("core/Reboot", nil, nil)
}

// This action sets the router admin password, NOT the network password.
// * adminPassword a string set to the password
// * passwordHint a hint viewable by anyone connected to the network
func (client Client) SetAdminPassword(password, hint string) error {
	err := client.MakeRequest("core/SetAdminPassword2", map[string]string{
		"adminPassword": password,
		"passwordHint":  hint,
	}, nil)
	if err != nil {
		return err
	}

	return client.Authorize(password)
}
