package location

import (
    "context"
    "net"
    "net/http"
    "net/url"
    "io/ioutil"
    "encoding/json"
    "strings"

    "fmt"

    "github.com/pkg/errors"

    "demo/config/secret"
)

// For API https://api.freegeoip.app/json/1.2.3.4?apikey=xxxxx
type freegeoipappDataType struct {
    IP string `json:"ip"`
    CountryCode string `json:"country_code"`
    CountryName string `json:"country_name"`
    RegionCode string `json:"region_code"`
    RegionName string `json:"region_name"`
    City string `json:"city"`
    ZipCode string `json:"zip_code"`
    TimeZone string `json:"time_zone"`
    Latitude float64 `json:"latitude"`
    Longitude float64 `json:"longitude"`
    MetroCode int `json:"metro_code"`
}

type freegeoipappType struct {
    IsEmpty bool
    Model
    IPData IPDataType
}

// For API https://api.freegeoip.app/json/1.2.3.4?apikey=xxxxx
func freegeoipappNew(ctx context.Context, netip net.IP) (freegeoipappType, error) {
    var empty freegeoipappType

    uobj, err := buildURL(netip)
    if err != nil {
        return empty, errors.WithStack(err)
    }

    if uobj.isEmpty {
        return freegeoipappType{IsEmpty: true}, nil
    }

    req, err := http.NewRequest(http.MethodGet, uobj.sec.String(), nil)
    if err != nil {
        return empty, errors.WithStack(err)
    }

    req.Header.Set("User-Agent", "Mozilla/4.0 (compatible; MSIE 6.0; Windows 98)")

    client := &http.Client{}

    // TODO: IP 151.101.193.140
    //   Warning: "Transport: unhandled response frame type *http.http2UnknownFrame"
    // https://github.com/golang/go/issues/40359
    // https://github.com/golang/net/blob/0dd24b26b47d4eb2d45eb3c7a4bcb809d7c1edb8/http2/transport.go#L2149
    resp, err := client.Do(req)

    if err != nil {
        return empty, errors.WithStack(err)
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return empty, errors.WithStack(err)
    }

    var location freegeoipappDataType
    isEmpty := freegeoipappType{IsEmpty: true}
    if err := json.Unmarshal(body, &location); err != nil {
        switch true {
        case strings.Contains(err.Error(), "invalid character '<' looking for beginning of value"):
            return isEmpty, nil
        }
        // return empty, errors.WithStack(err)
    }

    if strings.TrimSpace(location.City) == "" ||
       strings.TrimSpace(location.RegionCode) == "" ||
       strings.TrimSpace(location.CountryCode) == "" {
        return isEmpty, nil
    }

    return freegeoipappType{
        IsEmpty: false,
        Model: Model{
            CountryCode: location.CountryCode,
            Country: location.CountryName,
            RegionCode: location.RegionCode,
            Region: location.RegionName,
            City: location.City,
            ZipCode: location.ZipCode,
            TimeZone: location.TimeZone,
            Latitude: location.Latitude,
            Longitude: location.Longitude,
        },
    }, nil
}

type urlType struct {
    sec url.URL
    pub url.URL
    isEmpty bool
}

func buildURL(netip net.IP) (urlType, error) {
    baseURL := url.URL{
        Scheme: "https",
        Host:   "api.freegeoip.app",
        Path: fmt.Sprintf("/json/%s", netip.String()),
    }
    secURL := baseURL

    var empty urlType
    sec, err := secret.New()
    if err != nil {
        return empty, errors.WithStack(err)
    }

    if sec.FreeGeoIPAppKey == "" {
        return urlType{isEmpty: true}, nil
    }

    query := secURL.Query()
    query.Set("apikey", sec.FreeGeoIPAppKey)
    secURL.RawQuery = query.Encode()

    return urlType{
        sec: secURL,
        pub: baseURL,
        isEmpty: false,
    }, nil
}
