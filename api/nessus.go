package api

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// Nessus exposes the resources offered via the Tenable Nessus RESTful API.
type Nessus interface {
	SetVerbose(bool)
	AuthCookie() string
	Request(method string, resource string, js interface{}, wantStatus []int) (resp *http.Response, err error)
	Login(username, password string) error
	Logout() error
	Session() (Session, error)

	ServerProperties() (*ServerProperties, error)
	ServerStatus() (*ServerStatus, error)

	PluginFamilies() ([]PluginFamily, error)
	FamilyDetails(ID int64) (*FamilyDetails, error)
	PluginDetails(ID int64) (*PluginDetails, error)
	AllPlugins() (chan PluginDetails, error)

	NewScan(editorTmplUUID, settingsName string, outputFolderID, policyID, scannerID int64, launch string, targets []string) (*Scan, error)
	CreateScan(newScanRequest NewScanRequest) (*Scan, error)
	Scans() (*ListScansResponse, error)
	ScanTemplates() ([]Template, error)
	PolicyTemplates() ([]Template, error)
	StartScan(scanID int64) (string, error)
	PauseScan(scanID int64) error
	ResumeScan(scanID int64) error
	StopScan(scanID int64) error
	DeleteScan(scanID int64) error
	ScanDetails(scanID int64) (*ScanDetailsResp, error)
	ConfigureScan(scanID int64, scanSetting NewScanRequest) (*Scan, error)

	Timezones() ([]TimeZone, error)

	Folders() ([]Folder, error)
	CreateFolder(name string) error
	EditFolder(folderID int64, newName string) error
	DeleteFolder(folderID int64) error

	ExportScan(scanID int64, format string) (int64, error)
	ExportFinished(scanID, exportID int64) (bool, error)
	DownloadExport(scanID, exportID int64) ([]byte, error)

	Permissions(objectType string, objectID int64) ([]Permission, error)
}

type nessusImpl struct {
	client *http.Client
	authCookie string
	apiURL     string
	verbose bool
}

func NewInsecureNessus(apiURL string) (*nessusImpl, error) {
	return newNessus(apiURL)
}

func newNessus(apiURL string) (*nessusImpl, error) {
	timeout := time.Duration(10 * time.Second) //超时时间50ms
	client := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}}

	return &nessusImpl{apiURL: apiURL, client: client}, nil
}

func (n *nessusImpl) Folders() ([]Folder, error) {
	if n.verbose {
		log.Println("Getting list of folders...")
	}

	resp, err := n.Request("GET", "/folders", "", []int{http.StatusOK})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	reply := &listFoldersResp{}
	if err = json.NewDecoder(resp.Body).Decode(&reply); err != nil {
		return nil, err
	}
	return reply.Folders, nil
}
