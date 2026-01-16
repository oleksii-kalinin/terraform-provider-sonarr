package sonarr

type SystemStatus struct {
	AppName string `json:"appName"`
	Version string `json:"version"`
	OsName  string `json:"osName"`
}

type ApiInfo struct {
}

type Series struct {
	Id               int         `json:"id"`
	Title            string      `json:"title"`
	RootFolderPath   string      `json:"rootFolderPath"`
	QualityProfileId int32       `json:"qualityProfileId"`
	Monitored        bool        `json:"monitored"`
	SeasonFolder     bool        `json:"seasonFolder"`
	TvdbID           int32       `json:"tvdbId"`
	AddOptions       *AddOptions `json:"addOptions"`
}

type AddOptions struct {
	Monitor string `json:"monitor"`
}
