package maps

import (
	"errors"
	"fmt"
	"regexp"
	"os"

	g "github.com/serpapi/google-search-results-golang"
)

type Company struct {
	Title string
	Website string
	Phone string
	Address string
	PlaceId string
	Category string
	Rating float64
}	

//
// Query the Google Maps API to retrive a list of companies near a specific location
//
// Params:
//		latLon: The latitude and longitude of the location, example: "@45.123456,12.345678,15z"
//		searchQuery: The query to search, example: "Pizza" or "Restaurant"
//		lang: The language of the results, default is "it"
//
// Returns:
//		[]Company: A list of companies near the location
//		error: An error if the request fails
//
func GetMaps(latLon string, searchQuery string, lang ...string) ([]Company, error) {
	matched, err := regexp.MatchString(`^@-?\d+\.\d+,-?\d+\.\d+,(\d+)z$`, latLon)
	if err != nil {
		msg := fmt.Sprintf("Error compiling regex: %s", err)
		return nil, errors.New(msg) 
	}
	if !matched {
		msg := fmt.Sprintf("Invalid latLon format: %s", latLon)
		return nil, errors.New(msg)
	}

	language := "it"
	if len(lang) > 0 {
		language = lang[0]
	}

	parameter := map[string]string{
		"q": searchQuery,
		"ll": latLon,
		"hl": language,
	}

	mapsApiKey := os.Getenv("MAPS_API_KEY")

	if mapsApiKey == "" {
		msg := "MAPS_API_KEY is not set"
		return nil, errors.New(msg)
	}

	search := g.NewSearch("google_maps", parameter, mapsApiKey)
	rsp, err := search.GetJSON()

	if err != nil {
		msg := fmt.Sprintf("Error getting JSON: %s", err)
		return nil, errors.New(msg)
	}

	results := rsp["local_results"].([]interface{})
	var companies []Company

	for _, result := range results {
		result := result.(map[string]interface{})
		company := Company{
			Title: checkKeyValueString("title", result),
			Website: checkKeyValueString("website", result),
			Phone: checkKeyValueString("phone", result),
			Address: checkKeyValueString("address", result),
			PlaceId: checkKeyValueString("place_id", result),
			Category: checkKeyValueString("category", result),
			Rating: checkKeyValueFloat64("rating", result),
		}

		if company.Title != "" {
			companies = append(companies, company)
		} else {
			fmt.Println(result)
		}
	}
	
	return companies, nil
}

//
// Check if a key exists in a map and return the value as a string
//
// Params:
//		key: The key to search
//		data: The map to search
//
// Returns:
//		string: The value of the key as a string
//
func checkKeyValueString(key string, data map[string]interface{}) string {
	value, exists := data[key].(string)
	if !exists {
		value = ""
	}
	return value
}

//
// Check if a key exists in a map and return the value as a float64
//
// Params:
//		key: The key to search
//		data: The map to search
//
// Returns:
//		float64: The value of the key as a float64
//
func checkKeyValueFloat64(key string, data map[string]interface{}) float64 {
	value, exists := data[key].(float64)
	if !exists {
		value = 0
	}
	return value
}