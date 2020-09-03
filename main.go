package linksys

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"
)

type jnapResponse struct {
	Result string `json:"result"`
	Error  string `json:"error"`
	Output interface{}
}

// Client represents a client to a Linksys router.
type Client struct {
	Endpoint      string
	Authorization string
	PasswordHint  string
}

// NewClient returns a new client initialized to http://192.168.1.1/JNAP/.
func NewClient() *Client {
	client := &Client{
		Endpoint: "http://192.168.1.1/JNAP/",
	}
	return client
}

// MakeRequest performs a request to the endpoint.
func (client Client) MakeRequest(action string, body, output interface{}) error {
	if body == nil {
		body = struct{}{}
	}
	marshalled, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", client.Endpoint, bytes.NewReader(marshalled))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-JNAP-Authorization", client.Authorization)
	req.Header.Set("X-JNAP-Action", "http://linksys.com/jnap/"+action)

	res, err := HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	// the status should be 200 even if there is an error
	if res.StatusCode != 200 {
		return ErrStatusCode
	}

	response := jnapResponse{
		Output: output,
	}
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return err
	}

	log.Println(response.Result)

	if response.Result != "OK" {
		if response.Error != "" {
			return errors.New(response.Error)
		}
		return errors.New(response.Result)
	}

	return nil
}

func InitClient(password string) (client *Client) {

	client = NewClient()

	err := client.Authorize(password)
	if err != nil {
		log.Println(err.Error())
	}

	return client
}

// Authorize logs in using the router's password (different from network password).
func (client *Client) Authorize(password string) error {
	client.Authorization = "Basic " + base64.StdEncoding.EncodeToString([]byte("admin:"+password))
	return client.MakeRequest("core/CheckAdminPassword", nil, nil)
}

func DetectDevicePresence(client *Client, macAddress string) {

	devices, err := client.GetDevices(0)
	if err != nil {
		log.Println(err.Error())
	}

	numDevices := len(devices.Devices)

	dev := false

	for i := 0; i < numDevices; i++ {

		numConnections := len(devices.Devices[i].Connections)

		for j := 0; j < numConnections; j++ {
			if devices.Devices[i].Connections[j].MAC == macAddress {
				dev = true
				log.Println(devices.Devices[i].Hostname)
				break
			} else {

			}
		}
	}
	if dev == true {
		log.Println("User is here!")
	} else {
		log.Println("User is not here...")
	}
}

func (client *Client) StartPolling(duration time.Duration, macAddress string) {
	for {
		time.Sleep(duration)
		go DetectDevicePresence(client, macAddress)
	}
}

// HTTPClient is the HTTP client that will be used for all requests.
var HTTPClient = http.DefaultClient

// ErrStatusCode is returned when an API response does not have a 200 OK status.
var ErrStatusCode = errors.New("non-200 status code")
