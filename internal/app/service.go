package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/magzhan/geotracker/internal/model"
)

func (s *server) GetLocationDetails(location model.Location) error {
	response, err := http.Get(fmt.Sprintf("https://nominatim.openstreetmap.org/reverse?lat=%v&lon=%v&format=json", location.Latitude, location.Longitude))
	if err != nil {
		s.logger.Error("The fetch is failed from GetLocationDetails")
		return err
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&location)
	if err != nil {
		s.logger.Error("The fetch is failed from GetLocationDetails")
		return err
	}

	return nil
}
