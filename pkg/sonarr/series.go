package sonarr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	url2 "net/url"
	"strconv"
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

func (c Client) CreateSeries(show *Series) (*Series, error) {
	jsonBytes, err := json.Marshal(show)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/api/v3/series", c.BaseURL)

	reqReader := bytes.NewBuffer(jsonBytes)

	req, err := http.NewRequest("POST", url, reqReader)
	if err != nil {
		return nil, err
	}

	res, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(res.Body)

	switch res.StatusCode {
	case http.StatusCreated:
		break
	default:
		bodyBytes, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("API error: %d - %s", res.StatusCode, string(bodyBytes))
	}

	var result Series
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&result)
	return &result, nil
}

func (c *Client) DeleteSeries(id int, deleteFiles bool) error {
	//url := fmt.Sprintf("%s/api/v3/series/%d?deleteFiles=%s", c.BaseURL, id, deleteFiles)
	u, _ := url2.Parse(c.BaseURL)
	u = u.JoinPath("api", "v3", "series", strconv.Itoa(id))

	q := u.Query()
	q.Set("deleteFiles", strconv.FormatBool(deleteFiles))
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("DELETE", u.String(), nil)
	if err != nil {
		return err
	}

	res, err := c.doRequest(req)
	if err != nil {
		return err
	}

	switch res.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return nil
	default:
		return fmt.Errorf("API error: %d", res.StatusCode)
	}
}

func (s Series) String() string {
	return fmt.Sprint(s.Title)
}
