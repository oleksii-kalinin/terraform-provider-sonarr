package sonarr

type SystemStatus struct {
	AppName string `json:"appName"`
	Version string `json:"version"`
}

type ApiInfo struct {
}

type Series struct {
	Id               int32  `json:"id"`
	Title            string `json:"title"`
	RootFolderPath   string `json:"rootFolderPath"`
	QualityProfileId int    `json:"qualityProfileId"`
	Monitored        bool   `json:"monitored"`
	SeasonFolder     bool   `json:"seasonFolder"`
	TvdbID           int32  `json:"tvdbID"`
}
