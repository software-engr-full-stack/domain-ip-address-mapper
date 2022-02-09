package app

import (
    "context"

    "github.com/pkg/errors"
    "gorm.io/gorm"

    "demo/entities/urbanarea"
)

func createUrbanArea(
    ctx context.Context, tx *gorm.DB,
    city string,
    region string,
    country string,
) (_ urbanarea.Model, _ urbanarea.LocationDataType, _ error) {
    ua, recordIsFound, locationData, err := urbanarea.SearchOrNew(ctx, city, region, country)
    var empty urbanarea.Model
    var emptyLD urbanarea.LocationDataType
    if err != nil {
        return empty, emptyLD, errors.WithStack(err)
    }
    if recordIsFound {
        return ua, locationData, nil
    }

    if ua.IsEmpty() {
        return empty, locationData, nil
    }

    if err = tx.Create(&ua).Error; err != nil {
        return empty, emptyLD, errors.WithStack(err)
    }

    return ua, locationData, nil
}
