package location

import (
    "context"
    "net"
    "strings"

    "fmt"

    "github.com/pkg/errors"
)

// Not persisted to database, used to get additional location data
type IPDataType struct {
    ISP string `json:"isp"`
    Org string `json:"org"`
    AS  string `json:"as"`
}

func New(ctx context.Context, netip net.IP) (Model, IPDataType, error) {
    var empty Model
    var emptyIPData IPDataType
    fgip, err := freegeoipappNew(ctx, netip)
    if err != nil {
        return empty, emptyIPData, errors.WithStack(err)
    }
    if !fgip.IsEmpty {
        return injectNames(fgip.Model), fgip.IPData, nil
    }

    ipapi, err := ipapicomNew(ctx, netip)
    if err != nil {
        return empty, emptyIPData, errors.WithStack(err)
    }
    if !ipapi.IsEmpty {
        return injectNames(ipapi.Model), ipapi.IPData, nil
    }

    keyCDN, err := keyCDNComNew(ctx, netip)
    if err != nil {
        return empty, emptyIPData, errors.WithStack(err)
    }

    if !keyCDN.IsEmpty {
        return injectNames(keyCDN.Model), keyCDN.IPData, nil
    }

    return Model{
        Error: errors.New(fmt.Sprintf("no location data found for %#v", netip.String())),
    }, emptyIPData, nil
}

func injectNames(model Model) Model {
    city := strings.TrimSpace(model.City)
    comp := []string{city}
    region := strings.TrimSpace(model.Region)
    if region != "" {
        comp = append(comp, region)
    }

    country := strings.TrimSpace(model.Country)
    comp = append(comp, country)

    zc := strings.TrimSpace(model.ZipCode)
    if zc != "" {
        comp = append(comp, zc)
    }

    model.FullName = strings.Join(comp, ", ")

    comp = []string{city}
    regionCode := strings.TrimSpace(model.RegionCode)
    if regionCode != "" {
        comp = append(comp, regionCode)
    } else {
        if region != "" {
            comp = append(comp, region)
        }
    }

    cc := strings.TrimSpace(model.CountryCode)
    if cc != "" {
        if cc != "US" {
            comp = append(comp, cc)
        }
    } else {
        if country != "" {
            comp = append(comp, country)
        }
    }

    model.ViewFullName = strings.Join(comp, ", ")

    return model
}
