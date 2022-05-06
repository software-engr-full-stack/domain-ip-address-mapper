package location

import (
    "context"
    "net"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "strings"

    "fmt"

    "github.com/pkg/errors"
)

// For API http://ip-api.com/json/1.2.3.4
type ipapicomDataType struct {
    Query string `json:"query"`
    Status string `json:"status"`
    Country string `json:"country"`
    CountryCode string `json:"countryCode"`
    Region string `json:"region"`
    RegionName string `json:"regionName"`
    City string `json:"city"`
    Zip string `json:"zip"`
    Lat float64 `json:"lat"`
    Lon float64 `json:"lon"`
    TimeZone string `json:"timezone"`
    ISP string `json:"isp"`
    Org string `json:"org"`
    AS string `json:"as"`
}

type ipapicomType struct {
    IsEmpty bool
    Model
    IPData IPDataType
}

// For API http://ip-api.com/json/1.2.3.4
func ipapicomNew(ctx context.Context, netip net.IP) (ipapicomType, error) {
    var empty ipapicomType

    link := fmt.Sprintf("http://ip-api.com/json/%s", netip.String())

    req, err := http.NewRequest(http.MethodGet, link, nil)
    // fmt.Printf("... DEBUG: fetching %#v\n", link)
    if err != nil {
        return empty, errors.WithStack(err)
    }

    req.Header.Set("User-Agent", "Mozilla/4.0 (compatible; MSIE 6.0; Windows 98)")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return empty, errors.WithStack(err)
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return empty, errors.WithStack(err)
    }

    var location ipapicomDataType
    if err := json.Unmarshal(body, &location); err != nil {
        return empty, errors.WithStack(err)
    }

    if strings.TrimSpace(location.City) == "" ||
       strings.TrimSpace(location.Region) == "" ||
       strings.TrimSpace(location.CountryCode) == "" {
        return ipapicomType{IsEmpty: true}, nil
    }

    return ipapicomType{
        IsEmpty: false,
        Model: Model{
            CountryCode: location.CountryCode,
            Country: location.Country,
            RegionCode: location.Region,
            Region: location.RegionName,
            City: location.City,
            ZipCode: location.Zip,
            TimeZone: location.TimeZone,
            Latitude: location.Lat,
            Longitude: location.Lon,
        },
        IPData: IPDataType{
            ISP: location.ISP,
            Org: location.Org,
            AS: location.AS,
        },
    }, nil
}
