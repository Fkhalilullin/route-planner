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
	MinElevation float64  `json:"minElevation"`
	MaxElevation float64  `json:"maxElevation"`
	RowCount     int      `json:"rowCount"`
	ColumnCount  int      `json:"columnCount"`
	Points       []Points `json:"points"`
}

type Points struct {
	Lat       float64 `json:"lat"`
	Lon       float64 `json:"lon"`
	Elevation float64 `json:"elevation"`
}
