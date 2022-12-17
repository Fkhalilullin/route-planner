package route

type Request struct {
	TopLeftPoint struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	} `json:"topLeftPoint"`
	BotRightPoint struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	} `json:"botRightPoint"`
	Paths []Path `json:"paths"`
}

type Response struct {
	Paths []Path `json:"paths"`
}

type Path struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}
