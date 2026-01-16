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

// SeriesLookup represents a series returned from Sonarr's TVDB lookup endpoint.
type SeriesLookup struct {
	Title      string `json:"title"`
	SortTitle  string `json:"sortTitle"`
	Status     string `json:"status"`
	Overview   string `json:"overview"`
	Network    string `json:"network"`
	Year       int32  `json:"year"`
	TvdbId     int32  `json:"tvdbId"`
	ImdbId     string `json:"imdbId"`
	Runtime    int32  `json:"runtime"`
	SeasonCount int32 `json:"seasonCount"`
}
