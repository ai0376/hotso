package cloud

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// A client represents a client connection to a {own|next}cloud
type Client struct {
	Url      *url.URL
	Username string
	Password string
}

// Error type encapsulates the returned error messages from the
// server.
type Error struct {
	// Exception contains the type of the exception returned by
	// the server.
	Exception string `xml:"exception"`

	// Message contains the error message string from the server.
	Message string `xml:"message"`
}

// Dial connects to an {own|next}Cloud instance at the specified
// address using the given credentials.
func Dial(host, username, password string) (*Client, error) {
	url, err := url.Parse(host)
	if err != nil {
		return nil, err
	}
	return &Client{
		Url:      url,
		Username: username,
		Password: password,
	}, nil
}

// Mkdir creates a new directory on the cloud with the specified name.
func (c *Client) Mkdir(path string) error {
	_, err := c.sendRequest("MKCOL", path)
	return err

}

// Delete removes the specified folder from the cloud.
func (c *Client) Delete(path string) error {
	_, err := c.sendRequest("DELETE", path)
	return err
}

// Upload uploads the specified source to the specified destination
// path on the cloud.
func (c *Client) Upload(src []byte, dest string) error {

	destUrl, err := url.Parse(dest)
	if err != nil {
		return err
	}

	// Create the https request

	client := &http.Client{}
	req, err := http.NewRequest("PUT", c.Url.ResolveReference(destUrl).String(), bytes.NewReader(src))
	if err != nil {
		return err
	}

	req.SetBasicAuth(c.Username, c.Password)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if len(body) > 0 {
		error := Error{}
		err = xml.Unmarshal(body, &error)
		if err != nil {
			return fmt.Errorf("Error during XML Unmarshal for response %s. The error is %s", body, err)
		}
		if error.Exception != "" {
			return fmt.Errorf("Exception: %s, Message: %s", error.Exception, error.Message)
		}

	}

	return nil
}

// Download downloads a file from the specified path.
func (c *Client) Download(path string) ([]byte, error) {

	pathUrl, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	// Create the https request

	client := &http.Client{}
	req, err := http.NewRequest("GET", c.Url.ResolveReference(pathUrl).String(), nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.Username, c.Password)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	error := Error{}
	err = xml.Unmarshal(body, &error)
	if err == nil {
		if error.Exception != "" {
			return nil, fmt.Errorf("Exception: %s, Message: %s", error.Exception, error.Message)
		}
	}

	return body, nil
}

func (c *Client) Exists(path string) bool {
	_, err := c.sendRequest("PROPFIND", path)
	return err == nil
}

func (c *Client) sendRequest(request string, path string) ([]byte, error) {
	// Create the https request

	folderUrl, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	req, err := http.NewRequest(request, c.Url.ResolveReference(folderUrl).String(), nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.Username, c.Password)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if len(body) > 0 {
		error := Error{}
		err = xml.Unmarshal(body, &error)
		if err != nil {
			return body, fmt.Errorf("Error during XML Unmarshal for response %s. The error was %s", body, err)
		}
		if error.Exception != "" {
			return nil, fmt.Errorf("Exception: %s, Message: %s", error.Exception, error.Message)
		}

	}

	return body, nil
}
