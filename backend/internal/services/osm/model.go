package osm

type Box struct {
	MinLon float64
	MinLat float64
	MaxLon float64
	MaxLat float64
}

type Type struct {
	Lat   float64
	Lon   float64
	Value string
}
