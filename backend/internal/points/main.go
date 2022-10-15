package points

import (
	"fmt"
	"net/http"

	"github.com/Fkhalilullin/route-planner/internal/services/topodata"
)

func GetPoints(_ http.ResponseWriter, _ *http.Request) {
	evelationService := topodata.NewService()

	// TODO Отправлять массив точек
	elevation, err := evelationService.GetElevationPoints(-43.5, 172.5)
	if err != nil {
		return
	}

	fmt.Println(elevation)
}
