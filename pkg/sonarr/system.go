package sonarr

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

func (c *Client) GetSystemStatus() (*SystemStatus, error) {
	status := SystemStatus{}

	url := c.BaseURL + "/api/v3/system/status"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Print(err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("got an error with code %v", resp.StatusCode))
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&status)
	if err != nil {
		return nil, err
	}

	return &status, nil
}

func (s SystemStatus) String() string {
	return fmt.Sprintf("Platform: %s, version: %v", s.AppName, s.Version)
}
