// ipgeolocation is a small unofficial wrapper client for ipgeolocation.io with the goal to cover the complete API
package ipgeolocation

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

// Options struct to build query in GetGeolocationWithOptions
type Options struct {
	IP       string
	Language string
	Exclude  []string
	Include  []string
}

// Helper function to validate include/exclude fields
func isDataField(field string) bool {
	val := reflect.ValueOf(GeolocationData{})
	for i := 0; i < val.Type().NumField(); i++ {
		if val.Type().Field(i).Tag.Get("json") == field {
			return true
		}
	}
	return false
}

// SetIncludes calls SetInclude for each element. Returns non-nil error if one of the elements failed
func (c *Options) SetIncludes(includes []string) error {
	for _, include := range includes {
		if err := c.SetInclude(include); err != nil {
			return err
		}
	}
	return nil
}

// SetInclude validates the field by checking if its a valid json field in the GeolocationData struct. Appends field to Options.Include if valid or returns non-nil error otherwise
func (c *Options) SetInclude(include string) error {
	if isDataField(include) {
		c.Include = append(c.Include, include)
		return nil
	}
	return errors.New("exclude is invalid")
}

// SetExcludes calls SetExclude for each element. Returns non-nil error if one of the elements failed
func (c *Options) SetExcludes(excludes []string) error {
	for _, exclude := range excludes {
		if err := c.SetExclude(exclude); err != nil {
			return err
		}
	}
	return nil
}

// SetExclude validates the field by checking if its a valid json field in the GeolocationData struct. Appends field to Options.Exclude if valid or returns non-nil error otherwise
func (c *Options) SetExclude(exclude string) error {
	if isDataField(exclude) {
		c.Exclude = append(c.Exclude, exclude)
		return nil
	}
	return errors.New("exclude is invalid")
}

// SetLanguage validates the given language and sets it as Options attribute if valid. Returns non-nil error otherwise
func (c *Options) SetLanguage(language string) error {
	possibleLanguages := []string{"en", "de", "ru", "ja", "fr", "cn", "es", "cs", "it"}

	for _, b := range possibleLanguages {
		if b == language {
			c.Language = language
			return nil
		}
	}
	return errors.New("language is not valid")
}

// Set the IP which you want to check. Validates the IP by trying to parse it with net.ParseIP
func (c *Options) SetIP(IP string) error {
	parsed := net.ParseIP(IP)
	if parsed != nil {
		c.IP = IP
		return nil
	}
	return errors.New("ip got invalid format")
}

// GeolocationData represents the full response which we get from the API.
// Full documentation: https://ipgeolocation.io/documentation/ip-geolocation-api.html
type GeolocationData struct {
	IP             string `json:"ip"`
	ContinentCode  string `json:"continent_code"`
	ContinentName  string `json:"continent_name"`
	CountryCode2   string `json:"country_code2"`
	CountryCode3   string `json:"country_code3"`
	CountryName    string `json:"country_name"`
	CountryCapital string `json:"country_capital"`
	StateProv      string `json:"state_prov"`
	District       string `json:"district"`
	City           string `json:"city"`
	ZipCode        string `json:"zipcode"`
	Latitude       string `json:"latitude"`
	Longitude      string `json:"longitude"`
	IsEU           bool   `json:"is_eu"`
	CallingCode    string `json:"calling_code"`
	CountryTLD     string `json:"country_tld"`
	Languages      string `json:"languages"`
	CountryFlag    string `json:"country_flag"`
	GeonameID      string `json:"geoname_id"`
	ISP            string `json:"isp"`
	ConnectionType string `json:"connection_type"`
	Organization   string `json:"organization"`
	Currency       struct {
		Code   string `json:"code"`
		Name   string `json:"name"`
		Symbol string `json:"symbol"`
	} `json:"currency"`
	Timezone struct {
		Name            string  `json:"name"`
		Offset          int     `json:"offset"`
		CurrentTime     string  `json:"current_time"`
		CurrentTimeUnix float64 `json:"current_time_unix"`
		IsDST           bool    `json:"is_dst"`
		DSTSavings      int     `json:"dst_savings"`
	} `json:"time_zone"`
}

// API base URL where we are going to append our queries
const baseURL = "https://api.ipgeolocation.io/ipgeo"

// Client is the IP-Geolocation API client
type Client struct {
	ClientToken string
}

// SetToken sets the client API token
func (c *Client) SetToken(token string) {
	c.ClientToken = token
}

// You can also pass an empty IP address, it will then use your client's IP.
//
// Please read the API documentation to know what fields to expect: https://ipgeolocation.io/documentation/ip-geolocation-api.html
//
// Will return non-nil error if making the request or parsing the response failed.
func (c *Client) GetGeolocation(ip string) (GeolocationData, error) {
	var data GeolocationData

	u, _ := url.Parse(baseURL)

	// API Key, required.
	addQuery(u, "apiKey", c.ClientToken)

	// IP
	if ip != "" {
		addQuery(u, "ip", ip)
	}

	res, err := http.Get(u.String())
	if err != nil {
		return data, err
	}

	if res.StatusCode != 200 {
		return data, errors.New("got non 200 status code")
	}

	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}
// GetGeolocationWithOptions, GetGeolocation but with more options!:D
func (c *Client) GetGeolocationWithOptions(options Options) (GeolocationData, error) {
	var data GeolocationData

	u, _ := url.Parse(baseURL)

	addQuery(u, "apiKey", c.ClientToken)

	// IP
	if options.IP != "" {
		addQuery(u, "ip", options.IP)
	}

	includes := strings.Join(options.Include, ",")
	if includes != "" {
		addQuery(u, "fields", includes)
	}
	fmt.Println(includes)

	excludes := strings.Join(options.Exclude, ",")
	if excludes != "" {
		addQuery(u, "excludes", excludes)
	}

	if options.Language != "" {
		addQuery(u, "lang", options.Language)
	}

	res, err := http.Get(u.String())
	if err != nil {
		return data, err
	}

	switch res.StatusCode {
	case 400:
		return data, errors.New("your subscription is paused from use")
	case 401:
		return data, errors.New("401 error code, permission denied")
	case 403:
		return data, errors.New("ip address or domain is not valid")
	case 404:
		return data, errors.New("ip not found in database")
	case 423:
		return data, errors.New("private ip address")
	}

	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return data, err
	}
	return data, nil

}

func addQuery(url *url.URL, key string, value string) {
	q := url.Query()
	q.Set(key, value)
	url.RawQuery = q.Encode()
}
