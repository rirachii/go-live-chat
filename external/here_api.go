package here_api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

func GetReverseGeocode(apiKey string, latitude float32, longitude float32) (ReverseGeocodeResponse, error) {

	// https://www.here.com/docs/bundle/geocoding-and-search-api-developer-guide/page/topics/endpoint-reverse-geocode-brief.html

	formatRoute := func(apiKey string, latitude float32, longitude float32) string {

		// latitude, longitude, apikey

		hereAPI := `https://revgeocode.search.hereapi.com`
		apiRoute := `/v1/revgeocode?at=%[1]s&lang=en-US&apiKey=%[2]s`

		// latitude, longitude
		latLongQuery := fmt.Sprintf(`%[1]f,%[2]f`, latitude, longitude)

		return hereAPI + fmt.Sprintf(apiRoute, url.QueryEscape(latLongQuery), apiKey)
	}

	route := formatRoute(apiKey, latitude, longitude)

	log.Print("getting reverse geocode")

	res, err := http.Get(route)
	if err != nil {
		log.Print(err)
		return ReverseGeocodeResponse{}, err
	}

	var geocodeResponse ReverseGeocodeResponse
	decodeErr := json.NewDecoder(res.Body).Decode(&geocodeResponse)
	if decodeErr != nil {
		log.Print(decodeErr)
		return ReverseGeocodeResponse{}, decodeErr
	}

	// log.Printf("%+v", geocodeResponse)

	return geocodeResponse, nil

}

type ReverseGeocodeResponse struct {
	// https://www.here.com/docs/bundle/geocoding-and-search-api-developer-guide/page/topics/endpoint-reverse-geocode-brief.html
	Items []GeocodeItem `json:"items"`
}

type GeocodeItem struct {
	Title           string          `json:"title"`
	Id              string          `json:"id"`
	ResultType      string          `json:"resultType"`
	HouseNumberType string          `json:"houseNumberType"`
	Address         GeocodeAddress  `json:"address"`
	Position        GeocodePosition `json:"position"`
	Distance        int             `json:"distance"`
	MapView         GeocodeMapView  `json:"mapView"`
}

type GeocodeAddress struct {
	Label       string `json:"label"`
	CountryCode string `json:"countryCode"`
	CountryName string `json:"countryName"`
	State       string `json:"state"`
	County      string `json:"county"`
	City        string `json:"city"`
	District    string `json:"district"`
	Street      string `json:"street"`
	PostalCode  string `json:"postalCode"`
	HouseNumber string `json:"houseNumber"`
}

type GeocodePosition struct {
	Latitude  float32 `json:"lat"`
	Longitude float32 `json:"lng"`
}

type GeocodeMapView struct {
	North float32 `json:"north"`
	East  float32 `json:"east"`
	South float32 `json:"south"`
	West  float32 `json:"west"`
}
