package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/magzhan/geotracker/internal/model"
)

func (s *server) GetLocationDetails(location *model.Location) error {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(
		"https://nominatim.openstreetmap.org/reverse?lat=%v&lon=%v&format=json",
		location.Latitude, location.Longitude), nil)
	if err != nil {
		s.logger.Error("Error creating request in GetLocationDetails", "error", err)
		return err
	}
	req.Header.Set("User-Agent", "geotracker/1.0 (em261777@gmail.com)")

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		s.logger.Error("The fetch is failed from GetLocationDetails", "error", err)
		return err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(location)
	if err != nil {
		s.logger.Error("The fetch is failed from GetLocationDetails")
		return err
	}

	return nil
}
