package sonarr

type SystemStatus struct {
	AppName string `json:"appName"`
	Version string `json:"version"`
}

type ApiInfo struct {
}

type Series struct {
	Id               int    `json:"id"`
	Title            string `json:"title"`
	RootFolderPath   string `json:"rootFolderPath"`
	QualityProfileId int    `json:"qualityProfileId"`
	Monitored        bool   `json:"monitored"`
	SeasonFolder     bool   `json:"seasonFolder"`
	TvdbID           int    `json:"tvdbId"`
}
