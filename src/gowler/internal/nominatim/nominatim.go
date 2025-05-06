package nominatim

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"fmt"
)

//
// Query the Nomination API to retrive tha latitude and longitude of a specific location
//
// Params:
//		query: The location to search, example: "New York" or "Eiffel Tower, France"
//		limit: The number of results to return
//
// Returns:
//		[]string: A list of strings with the latitude and longitude of the location
//		error: An error if the request fails
//
func GetLatLon(query string, limit uint) ([]string, error) {
	nominationUrl := "https://nominatim.openstreetmap.org/search?q=" + url.QueryEscape(query) + "&format=json&limit=" + fmt.Sprintf("%d", limit)
	resp, err := http.Get(nominationUrl)

	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var datas []map[string]interface{}
	err = json.Unmarshal(body, &datas)

	if err != nil {
		return nil, err
	}

	coordinates := make([]string, len(datas), limit)
	for i, data := range datas {
		lat := data["lat"]
		lon := data["lon"]
		coordinates[i] = "@" + lat.(string) + "," + lon.(string) + ",15z"
	}

	return coordinates, nil
}