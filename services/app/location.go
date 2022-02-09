package app

import (
    "context"
    "net"

    "github.com/pkg/errors"
    "gorm.io/gorm"

    "demo/entities/location"

    "fmt"
)

func createLocation(
    ctx context.Context, tx *gorm.DB, netip net.IP,
) (_ location.Model, _ location.IPDataType, recordIsFound bool, _ error) {
    var empty location.Model
    var emptyIPData location.IPDataType
    var emptyRIF bool

    loc, ipDataFromLoc, err := location.New(ctx, netip)
    if err != nil {
        return empty, emptyIPData, emptyRIF, errors.WithStack(err)
    }
    if loc.IsEmpty() {
        return empty, emptyIPData, emptyRIF, nil
    }

    // Because some locations have different latitudes and longitudes even if they
    //   have the same cities, regions, countries and zip codes.
    crcz := location.Model{
        City: loc.City, Region: loc.Region, Country: loc.Country, ZipCode: loc.ZipCode,
    }

    uniq, err := crcz.Unique(ctx)
    if err != nil {
        return empty, emptyIPData, emptyRIF, errors.WithStack(err)
    }
    if uniq.IsFound {
        return uniq.Record, emptyIPData, true, nil
    } else {
        // TODO: I don't know about this whole else section.
        //   Probably a tighter search than just using IP alone.
        //   An IP address can have multiple locations.
        //   This should probably be handled in the IP entity.
        if uniq.AreRecordsFound {
            return empty, emptyIPData, emptyRIF, errors.New(fmt.Sprintf(
                "multiple location records of IP %#v are found", netip.String(),
            ))
        }
        if err = tx.Create(&loc).Error; err != nil {
            return empty, emptyIPData, emptyRIF, errors.WithStack(err)
        }
    }

    return loc, ipDataFromLoc, false, nil
}
