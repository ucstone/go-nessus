package api

type NewScanRequest struct {
	UUID     string   `json:"uuid"`
	Settings Settings `json:"settings"`
}

type Settings struct {
	LaunchNow   bool   `json:"launch_now"`
	Enabled     bool   `json:"enabled"`
	LiveResults string `json:"live_results"`
	Name        string `json:"name"`
	Description string `json:"description"`
	FolderID    int    `json:"folder_id"`
	ScannerID   string `json:"scanner_id"`
	TextTargets string `json:"text_targets"`
	FileTargets string `json:"file_targets"`
}
