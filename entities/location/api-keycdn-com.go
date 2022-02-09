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

// https://tools.keycdn.com/geo.json?host=1.2.3.4
type keyCDNComWrapperType struct {
    Status string `json:"status"`
    Description string `json:"description"`
    Data keyCDNComDataType `json:"data"`
}

type keyCDNComDataType struct {
    Geo keyCDNComGeoType `json:"geo"`
}

type keyCDNComGeoType struct {
    Host string `json:"host"`
    IP string `json:"ip"`
    RDNS string `json:"rdns"`
    ASN int `json:"asn"`
    ISP string `json:"isp"`
    CountryName string `json:"country_name"`
    CountryCode string `json:"country_code"`
    RegionName string `json:"region_name"`
    RegionCode string `json:"region_code"`
    City string `json:"city"`
    PostalCode string `json:"postal_code"`
    ContinentName string `json:"continent_name"`
    ContinentCode string `json:"continent_code"`
    Latitude float64 `json:"latitude"`
    Longitude float64 `json:"longitude"`
    MetroCode int `json:"metro_code"`
    Timezone string `json:"timezone"`
    Datetime string `json:"datetime"`
}

type keyCDNComType struct {
    IsEmpty bool
    Model
    IPData IPDataType
}

func keyCDNComNew(ctx context.Context, netip net.IP) (keyCDNComType, error) {
    link := fmt.Sprintf("https://tools.keycdn.com/geo.json?host=%s", netip.String())
    req, err := http.NewRequest(http.MethodGet, link, nil)
    var empty keyCDNComType
    if err != nil {
        return empty, errors.WithStack(err)
    }

    req.Header.Set("User-Agent", "keycdn-tools:https://www.example.com")

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

    if yes, err := isEmpty(body); err != nil {
        return empty, errors.WithStack(err)
    } else {
        if yes {
            return keyCDNComType{IsEmpty: true}, nil
        }
    }

    var keyCDN keyCDNComWrapperType
    if err := json.Unmarshal(body, &keyCDN); err != nil {
        return empty, errors.WithStack(err)
    }
    if keyCDN.Status != "success" {
        return empty, errors.New(keyCDN.Description)
    }

    geo := keyCDN.Data.Geo

    if strings.TrimSpace(geo.City) == "" ||
       // strings.TrimSpace(geo.RegionName) == "" ||
       strings.TrimSpace(geo.CountryName) == "" {
        return keyCDNComType{IsEmpty: true}, nil
    }

    return keyCDNComType{
        IsEmpty: false,
        Model: Model{
            CountryCode: geo.CountryCode,
            Country: geo.CountryName,
            RegionCode: geo.RegionCode,
            Region: geo.RegionName,
            City: geo.City,
            ZipCode: geo.PostalCode,
            TimeZone: geo.Timezone,
            Latitude: geo.Latitude,
            Longitude: geo.Longitude,
        },
        IPData: IPDataType{
            ISP: geo.ISP,
            AS: fmt.Sprintf("%d", geo.ASN),
        },
    }, nil
}

// Hack because trash API result, empty strings, etc.
func isEmpty(body []byte) (bool, error) {
    var empty bool

    var temp map[string]interface{}
    if err := json.Unmarshal(body, &temp); err != nil {
        return empty, errors.WithStack(err)
    }

    geo := temp["data"].(map[string]interface{})["geo"].(map[string]interface{})

    keys := []string{"city", "country_name"}

    for _, key := range keys {
        if strings.TrimSpace(strings.TrimSpace(geo[key].(string))) == "" {
            return true, nil
        }
    }

    return false, nil
}
