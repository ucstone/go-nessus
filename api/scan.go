package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (n *nessusImpl) NewScan(
	TmplUUID string,
	Name string,
	TextTargets string,
) (*Scan, error) {
	data := NewScanRequest{
		UUID: TmplUUID,
		Settings: Settings{
			Name:        Name,
			Description: "",
			LaunchNow:   true,
			FolderID:    3,
			Enabled:     false,
			ScannerID:   "1",
			TextTargets: TextTargets,
		},
	}
	return n.CreateScan(data)
}

func (n *nessusImpl) CreateScan(newScanRequest NewScanRequest) (*Scan, error) {
	if n.verbose {
		log.Println("Creating a new scan...")
	}

	resp, err := n.Request("POST", "/scans", newScanRequest, []int{http.StatusOK})
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	reply := struct {
		Scan Scan `json:"scan"`
	}{}

	if err = json.NewDecoder(resp.Body).Decode(&reply); err != nil {
		return nil, err
	}
	return &reply.Scan, nil
}


// StartScan starts the given scan and returns its UUID.
func (n *nessusImpl) StartScan(scanID int64) (string, error) {
	if n.verbose {
		log.Println("Starting scan...")
	}

	resp, err := n.Request("POST", fmt.Sprintf("/scans/%d/launch", scanID), nil, []int{http.StatusOK})
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	reply := &startScanResp{}
	if err = json.NewDecoder(resp.Body).Decode(&reply); err != nil {
		return "", err
	}
	return reply.UUID, nil
}

func (n *nessusImpl) PauseScan(scanID int64) error {
	if n.verbose {
		log.Println("Pausing scan...")
	}

	_, err := n.Request("POST", fmt.Sprintf("/scans/%d/pause", scanID), nil, []int{http.StatusOK})
	return err
}

func (n *nessusImpl) ResumeScan(scanID int64) error {
	if n.verbose {
		log.Println("Resume scan...")
	}

	_, err := n.Request("POST", fmt.Sprintf("/scans/%d/resume", scanID), nil, []int{http.StatusOK})
	return err
}

func (n *nessusImpl) StopScan(scanID int64) error {
	if n.verbose {
		log.Println("Stop scan...")
	}

	_, err := n.Request("POST", fmt.Sprintf("/scans/%d/stop", scanID), nil, []int{http.StatusOK})
	return err
}

func (n *nessusImpl) DeleteScan(scanID int64) error {
	if n.verbose {
		log.Println("Deleting scan...")
	}

	_, err := n.Request("DELETE", fmt.Sprintf("/scans/%d", scanID), nil, []int{http.StatusOK})
	return err
}

func (n *nessusImpl) ScanDetails(scanID int64) (*ScanDetailsResp, error) {
	if n.verbose {
		log.Println("Getting details about a scan...")
	}

	resp, err := n.Request("GET", fmt.Sprintf("/scans/%d", scanID), nil, []int{http.StatusOK})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	reply := &ScanDetailsResp{}
	if err = json.NewDecoder(resp.Body).Decode(&reply); err != nil {
		return nil, err
	}
	return reply, nil
}

func (n *nessusImpl) Scans() (*ListScansResponse, error) {
	if n.verbose {
		log.Println("Getting scans list...")
	}

	resp, err := n.Request("GET", "/scans", nil, []int{http.StatusOK})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	reply := &ListScansResponse{}
	if err = json.NewDecoder(resp.Body).Decode(&reply); err != nil {
		return nil, err
	}
	return reply, nil
}
