package urbanarea

import (
    "context"

    "fmt"

    "github.com/pkg/errors"
)

// Not persisted to database, used to get additional location data
type LocationDataType struct {
    FullName string
    Population int64 `json:"population"`
    TeleportGeoNameID int64
    TeleportGeoHash string
    TeleportName string
    TeleportFullName string
    TeleportLatitude float64
    TeleportLongitude float64
}

func SearchOrNew(
    ctx context.Context,
    city string,
    region string,
    country string,
) (_ Model, recordIsFound bool, _ LocationDataType, _ error) {
    var empty Model
    var emptyRIF bool
    var emptyLD LocationDataType

    citiesFoundByAPI, err := search(city)
    if err != nil {
        return empty, emptyRIF, emptyLD, errors.WithStack(err)
    }

    mr, err := match(city, region, country, citiesFoundByAPI)
    if err != nil {
        return empty, emptyRIF, emptyLD, errors.WithStack(err)
    }
    if !mr.IsFound {
        return empty, emptyRIF, emptyLD, nil
    }

    gd, err := geoData(mr.GeoNameIDURL)
    if err != nil {
        return empty, emptyRIF, emptyLD, errors.WithStack(err)
    }

    uaModel := Model{Name: gd.UrbanAreaName}
    if uaModel.IsEmpty() {
        return empty, false, LocationDataType{
            FullName: mr.FullName,
            Population: gd.Population,
            TeleportName: gd.Name,
            TeleportFullName: gd.FullName,
            TeleportGeoNameID: gd.GeoNameID,
            TeleportGeoHash: gd.GeoHash,
            TeleportLatitude: gd.Latitude,
            TeleportLongitude: gd.Longitude,
        }, nil
    }

    uniq, err := uaModel.Unique(ctx)
    if err != nil {
        return empty, emptyRIF, emptyLD, errors.WithStack(err)
    }
    if uniq.IsFound {
        return uniq.Record, true, LocationDataType{
            FullName: mr.FullName,
            Population: gd.Population,
            TeleportName: gd.Name,
            TeleportFullName: gd.FullName,
            TeleportGeoNameID: gd.GeoNameID,
            TeleportGeoHash: gd.GeoHash,
            TeleportLatitude: gd.Latitude,
            TeleportLongitude: gd.Longitude,
        },nil
    } else {
        if uniq.AreRecordsFound {
            return empty, emptyRIF, emptyLD, errors.New(fmt.Sprintf(
                "multiple urban area records of name %#v are found", gd.Name,
            ))
        }
    }

    ua, err := urbanArea(gd.UrbanAreaLink)
    if err != nil {
        return empty, emptyRIF, emptyLD, errors.WithStack(err)
    }

    return Model{
        Name: gd.UrbanAreaName,
        Link: gd.UrbanAreaLink,
        Scores: ua.UrbanAreaScores,
        Summary: ua.UrbanAreaSummary,
    }, false, LocationDataType{
        FullName: mr.FullName,
        Population: gd.Population,
        TeleportName: gd.Name,
        TeleportFullName: gd.FullName,
        TeleportGeoNameID: gd.GeoNameID,
        TeleportGeoHash: gd.GeoHash,
        TeleportLatitude: gd.Latitude,
        TeleportLongitude: gd.Longitude,
    }, nil
}
