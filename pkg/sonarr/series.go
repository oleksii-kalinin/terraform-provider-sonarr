package sonarr

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func (c Client) GetSeries(id int32) (*Series, error) {
	series := Series{}

	url := fmt.Sprintf("%s/api/v3/series/%d", c.BaseURL, id)
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
	switch resp.StatusCode {
	case http.StatusNotFound:
		return nil, nil // No such resource
	case http.StatusOK:
		break
	default:
		return nil, fmt.Errorf("API error: status code %d", resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&series)
	if err != nil {
		return nil, err
	}
	return &series, nil
}

func (s Series) String() string {
	return fmt.Sprint(s.Title)
}
