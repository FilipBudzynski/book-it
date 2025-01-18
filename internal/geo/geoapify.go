package geo

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

type (
	Timezone struct {
		Name             string `json:"name"`
		OffsetSTD        string `json:"offset_STD"`
		OffsetSTDSeconds int    `json:"offset_STD_seconds"`
		OffsetDST        string `json:"offset_DST"`
		OffsetDSTSeconds int    `json:"offset_DST_seconds"`
		AbbreviationSTD  string `json:"abbreviation_STD"`
		AbbreviationDST  string `json:"abbreviation_DST"`
	}

	Rank struct {
		Importance            float64 `json:"importance"`
		Confidence            float64 `json:"confidence"`
		ConfidenceCityLevel   float64 `json:"confidence_city_level"`
		ConfidenceStreetLevel float64 `json:"confidence_street_level"`
		MatchType             string  `json:"match_type"`
	}

	Datasource struct {
		Sourcename  string `json:"sourcename"`
		Attribution string `json:"attribution"`
		License     string `json:"license"`
		URL         string `json:"url"`
	}

	Result struct {
		Datasource    Datasource `json:"datasource"`
		Name          string     `json:"name"`
		Country       string     `json:"country"`
		CountryCode   string     `json:"country_code"`
		State         string     `json:"state"`
		City          string     `json:"city"`
		Postcode      string     `json:"postcode"`
		District      string     `json:"district,omitempty"`
		Suburb        string     `json:"suburb"`
		Quarter       string     `json:"quarter"`
		Street        string     `json:"street"`
		Lon           float64    `json:"lon"`
		Lat           float64    `json:"lat"`
		ResultType    string     `json:"result_type"`
		Formatted     string     `json:"formatted"`
		AddressLine1  string     `json:"address_line1"`
		AddressLine2  string     `json:"address_line2"`
		Timezone      Timezone   `json:"timezone"`
		PlusCode      string     `json:"plus_code"`
		PlusCodeShort string     `json:"plus_code_short"`
		Rank          Rank       `json:"rank"`
		PlaceID       string     `json:"place_id"`
	}
)

type GeoapifyAutocompleteResponse struct {
	Results []Result `json:"results"`
}

var GeoApifyKey string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}
	GeoApifyKey = os.Getenv("GEOAPIFY_KEY")
}

func GetLocalizationAutocomplete(query string) ([]Result, error) {
	params := url.Values{}
	params.Add("text", query)
	params.Add("format", "json")
	params.Add("apiKey", GeoApifyKey)

	queryUrl := "https://api.geoapify.com/v1/geocode/autocomplete?" + params.Encode()

	resp, err := http.Get(queryUrl)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("non-200 status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var geoapifyResponse GeoapifyAutocompleteResponse
	err = json.Unmarshal(body, &geoapifyResponse)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON response: %v", err)
	}

	if len(geoapifyResponse.Results) > 0 {
		for _, result := range geoapifyResponse.Results {
			fmt.Printf("Formatted Address: %s\n", result.Formatted)
			fmt.Printf("Latitude: %f, Longitude: %f\n", result.Lat, result.Lon)
			fmt.Printf("Country: %s, City: %s, State: %s\n", result.Country, result.City, result.State)
			fmt.Printf("Place ID: %s\n", result.PlaceID)
			fmt.Println("===================================")
		}
	} else {
		return []Result{}, nil
	}
	return geoapifyResponse.Results, nil
}
