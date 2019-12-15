## API Client for ipgeolocation.io

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

	// Leave the IP empty if you want to check your host IP
	ip := ""
	data, err := c.GetGeolocation(ip)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(data)
}
```

#### Planned Additions

- Adding all the other features which ipgeolocation.io provides
- Unit tests
