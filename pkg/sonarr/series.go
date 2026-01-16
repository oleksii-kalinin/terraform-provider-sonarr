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

// GetAllSeries retrieves all series currently in the Sonarr library.
// Returns a slice of Series or an error if the API call fails.
func (c *Client) GetAllSeries() ([]Series, error) {
	var series []Series

	url := fmt.Sprintf("%s/api/v3/series", c.BaseURL)
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
		return nil, fmt.Errorf("API error: status code %d", resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&series)
	if err != nil {
		return nil, err
	}
	return series, nil
}

func (c *Client) GetSeries(id int) (*Series, error) {
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

func (c *Client) CreateSeries(show *Series) (*Series, error) {
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
	u, err := url2.Parse(c.BaseURL)
	if err != nil {
		return err
	}

	u = u.JoinPath("api", "v3", "series", strconv.Itoa(int(id)))

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
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(res.Body)

	switch res.StatusCode {
	case http.StatusNotFound:
		return nil
	case http.StatusNoContent: // 204 is common success for DELETE
		return nil
	case http.StatusOK: // 200 also OK
		return nil
	default:
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("DELETE failed: %d %s - %s",
			res.StatusCode, res.Status, string(body))
	}
}

func (c *Client) UpdateSeries(show *Series) (*Series, error) {
	if show == nil {
		return nil, fmt.Errorf("series can't be found: %v", show)
	}
	jsonBytes, err := json.Marshal(show)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/api/v3/series/%d", c.BaseURL, show.Id)

	reqReader := bytes.NewBuffer(jsonBytes)

	req, err := http.NewRequest("PUT", url, reqReader)
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
	case http.StatusAccepted, http.StatusOK:
		var resSeries Series
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		if len(bodyBytes) == 0 {
			return show, nil
		}
		err = json.Unmarshal(bodyBytes, &resSeries)
		if err != nil {
			return nil, err
		}
		return &resSeries, nil
	default:
		bodyBytes, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("API Error: %d - %s", res.StatusCode, string(bodyBytes))
	}
}

func (s Series) String() string {
	return fmt.Sprint(s.Title)
}

// LookupSeries searches for series on TVDB via Sonarr's lookup endpoint.
// The term parameter is the search query (e.g., series title).
// Returns a slice of matching SeriesLookup results or an error if the API call fails.
func (c *Client) LookupSeries(term string) ([]SeriesLookup, error) {
	var results []SeriesLookup

	u, err := url2.Parse(c.BaseURL)
	if err != nil {
		return nil, err
	}

	u = u.JoinPath("api", "v3", "series", "lookup")
	q := u.Query()
	q.Set("term", term)
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
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
		return nil, fmt.Errorf("API error: status code %d", resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&results)
	if err != nil {
		return nil, err
	}
	return results, nil
}
