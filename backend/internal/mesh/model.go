package mesh

type Request struct {
	TopLeftPoint struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	} `json:"topLeftPoint"`
	BotRightPoint struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	} `json:"botRightPoint"`
}

type Response struct {
	Points []Points `json:"points"`
}

type Points struct {
	Lat       float64 `json:"lat"`
	Lon       float64 `json:"lon"`
	Elevation float64 `json:"elevation"`
}
