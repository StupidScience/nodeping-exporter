package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/prometheus/common/log"
)

// NodePing is nodeping.com API endpoint access config
type NodePing struct {
	APIURL string
	Token  string
}

// Check describe nodeping response about it
type Check struct {
	ID         string `json:"_id,omitempty"`
	Label      string `json:"label"`
	Type       string `json:"type"`
	Parameters struct {
		Target    string `json:"target"`
		Follow    bool   `json:"follow,omitempty"`
		Threshold int    `json:"threshold"`
		Sens      int    `json:"sens"`
	} `json:"parameters"`
	State int `json:"state"`
}

// CheckStats describe nodeping response about it
type CheckStats struct {
	Type     string `json:"t"`
	Target   string `json:"tg"`
	Result   string `json:"sc"`
	Message  string `json:"m"`
	Success  bool   `json:"su"`
	Duration int    `json:"rt"`
}

// NodePingError custom error
type NodePingError struct {
	e          string
	statusCode int
}

func (npe NodePingError) Error() string {
	return fmt.Sprintf("%s, got status code %d", npe.e, npe.statusCode)
}

func (np NodePing) request(m string) (*http.Response, error) {
	reqURL := fmt.Sprintf("%s/%s", np.APIURL, m)
	client := &http.Client{}
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		log.Errorf("Cannot create new request: %v", err)
		return nil, err
	}
	req.SetBasicAuth(np.Token, "")

	r, err := client.Do(req)
	if err != nil {
		log.Errorf("Cannot make new request to %s: %v", reqURL, err)
		return nil, err
	}
	switch r.StatusCode {
	case http.StatusOK:
		return r, nil
	case http.StatusForbidden:
		return nil, fmt.Errorf("Can't access with provided token")
	default:
		return nil, NodePingError{
			e:          fmt.Sprintf("Error occurred on %s", reqURL),
			statusCode: r.StatusCode,
		}
	}
}

// GetAllChecks return list of all available checks
func (np NodePing) GetAllChecks() (map[string]Check, error) {
	r, err := np.request("checks")
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	checks := map[string]Check{}

	err = json.NewDecoder(r.Body).Decode(&checks)
	if err != nil {
		log.Errorf("Can't decode received body of all checks: %v", err)
		return nil, fmt.Errorf("Can't decode received body of all checks: %v", err)
	}

	return checks, nil
}

// GetCheckStats return info and status about check
func (np NodePing) GetCheckStats(checkID string) (CheckStats, error) {
	r, err := np.request(fmt.Sprintf("results/%s?limit=1", checkID))
	if err != nil {
		return CheckStats{}, err
	}
	defer r.Body.Close()
	cs := []CheckStats{}

	err = json.NewDecoder(r.Body).Decode(&cs)
	if err != nil {
		log.Errorf("Can't decode received body of check stats: %v", err)
		return CheckStats{}, fmt.Errorf("Can't decode received body of check stats: %v", err)
	}

	return cs[0], nil
}

// CheckAccess simple availability check
func (np NodePing) CheckAccess() error {
	_, err := np.request("info/probe")
	if err != nil {
		return err
	}

	return nil
}

// NewNodePing create new NodePing structure
func NewNodePing(apiURL, token string) (NodePing, error) {
	if token == "" {
		return NodePing{}, fmt.Errorf("Token should be specified")
	}
	np := NodePing{
		APIURL: apiURL,
		Token:  token,
	}
	err := np.CheckAccess()
	if err != nil {
		return NodePing{}, fmt.Errorf("Can't proceed initial check: %v", err)
	}

	return np, nil
}
