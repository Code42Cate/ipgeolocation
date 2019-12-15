// ipgeolocation is a small unofficial wrapper client for ipgeolocation.io with the goal to cover the complete API
package ipgeolocation

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

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

func addQuery(url *url.URL, key string, value string) {
	q := url.Query()
	q.Set(key, value)
	url.RawQuery = q.Encode()
}
