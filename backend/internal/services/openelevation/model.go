package openelevation

type Results struct {
	Longitude float64 `json:"longitude"`
	Elevation float64 `json:"elevation"`
	Latitude  float64 `json:"latitude"`
}

type ElevationResponse struct {
	Results []Results `json:"results"`
}

type Locations struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type ElevationRequest struct {
	Locations []Locations `json:"locations"`
}
