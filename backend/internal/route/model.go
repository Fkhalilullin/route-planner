package route

const (
	FromLat = "from_lat"
	FromLon = "from_lon"
	ToLat   = "to_lat"
	ToLon   = "to_lon"
)

type Point struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type Route []Point
