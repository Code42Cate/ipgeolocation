[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://godoc.org/github.com/Code42Cate/ipgeolocation)

## IP Geolocation API Golang SDK

This is a small, incomplete and unofficial library to provide easier access to the [ipgeolocation.io](https://ipgeolocation.io/documentation/ip-geolocation-api.html) API

#### Installation
`go get -u github.com/Code42Cate/ipgeolocation`

#### Example Usage

```go
package main

import (
	"fmt"
	"log"

	"github.com/Code42Cate/ipgeolocation"
)

func main() {

	c := ipgeolocation.Client{}
	c.SetToken("your token")

	// Get all the data for your IP
	d, err := c.GetGeolocation("") // Or change it to any other IP

	if err != nil {
		log.Fatal(err)
	}

	// You can also specify more options:
	options := Options{
		IP:       "", // No IP <=> your IP
		Language: "de",
		Exclude:  []string{},       // You can set excludes here, or with the setter functions
		Include:  []string{"city"}, // Same as for excludes
	}

	options.SetInclude("ISP")
	// Or you could also just do:
	// options.SetIncludes("ISP, "city")
	options.SetLanguage("en") // Will overwrite "de"

	d, err = c.GetGeolocationWithOptions(options)
	fmt.Println(d) // You will still get the complete struct, just with default values in the fields you did not include.

	options = Options{}
	options.SetExcludes([]string{"currency", "time_zone", "country_flag", "geoname_id"})

    // Premium features:
	options.SetSecurity(true)
	options.SetHostname(true)

	d, err = c.GetGeolocationWithOptions(options)
	fmt.Println(d) // You will still get the complete struct, just with default values in the fields you excluded.
}
```

#### Planned Additions

- Adding premium features which ipgeolocation.io provides
- Unit tests?
