package handler

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"

	here_api "github.com/rirachii/golivechat/external"
	fragments_template "github.com/rirachii/golivechat/templates/fragments"

)

type LocationJSON struct {
	JSONData LocationData `json:"locationData"`
}
type LocationData struct {
	Latitude  float32 `json:"userLatitude"`
	Longitude float32 `json:"userLongitude"`
	Accuracy  float32 `json:"locAccuracy"`
}

// POST "/locate-user"
func HandleLocateUser(c echo.Context) error {

	var locationJSON LocationJSON

	err := c.Bind(&locationJSON)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "could not bind location data to json.")
	}
	// c.Validate()

	var locationDataFragment fragments_template.FragmentLocationData

	// check cookie here to see if we should use api
	locCookie, err := c.Cookie("location")
	// log.Print("COOKIES:", locCookie, "COOKIE ERROR:", err)
	if err == nil {
		// cookie exists

		var locationCookie LocationCookieJSON
		// log.Print("asdasd", locCookie)

		decodeCookie, err := base64.StdEncoding.DecodeString(locCookie.Value)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "could not decode cookie")
		}

		jsonDecodeErr := json.Unmarshal(decodeCookie, &locationCookie)
		if jsonDecodeErr != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Could not parse location cookie")
		}

		locationDataFragment.County = locationCookie.County
		locationDataFragment.City = locationCookie.City
		locationDataFragment.Country = locationCookie.Country
		
		// log.Print(locationCookie)

	} else {
		// cookie does not exist
		// 

		apiKey := os.Getenv(envHereAPIKEY)
		res, err := here_api.GetReverseGeocode(apiKey, locationJSON.JSONData.Latitude, locationJSON.JSONData.Longitude)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadGateway, "could not reverse geocode")
		}

		data := res.Items[0]
		addr := data.Address

		// log.Printf("%+v", addr)

		locationData := LocationCookieJSON{
			Country:     addr.CountryName,
			CountryCode: addr.CountryCode,
			City:        addr.City,
			County:      addr.County,
			PostalCode:  addr.PostalCode,
		}

		// log.Printf("for cookie: %+v", locationData)

		locationCookie, err := createLocationCookie(locationData)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "could not create location cookie")
		}

		c.SetCookie(locationCookie)

	}	

	return c.Render(http.StatusOK, 
		fragments_template.LocationDataFragment.TemplateName, 
		locationDataFragment,
	)
}

func createLocationCookie(locationData LocationCookieJSON) (*http.Cookie, error) {

	domainName := os.Getenv(envDomainName)

	locationJSON, err := json.Marshal(locationData)
	if err != nil {
		return nil, errors.New("could not create location cookie")
	}

	b64_LocationJSON := base64.StdEncoding.EncodeToString(locationJSON)

	log.Println(locationJSON, string(locationJSON))

	cookie := &http.Cookie{
		Name:     "location",
		Value:    b64_LocationJSON,
		MaxAge:   3600,
		Path:     "/",
		Domain:   domainName,
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteDefaultMode,
	}

	return cookie, nil
}

type LocationCookieJSON struct {
	CountryCode string `json:"countryCode"`
	Country     string `json:"country"`
	City        string `json:"city"`
	County      string `json:"county"`
	PostalCode  string `json:"postalCode"`
}
