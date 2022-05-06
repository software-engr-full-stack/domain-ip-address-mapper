package app

import (
    "context"

    "github.com/pkg/errors"
    "gorm.io/gorm"

    "demo/config"
    "demo/entities/user"
    "demo/entities/data"
    "demo/entities/ip"
    "demo/entities/crobat"
)

func Service(ctx context.Context, in crobat.InputType) (dataobj data.Model, deferr error) {
    db := ctx.Value("config").(config.Type).DBObj

    dbinst, err := db.Open()
    var empty data.Model
    if err != nil {
        return empty, errors.WithStack(err)
    }
    defer func() {
        if tempErr := db.Close(); tempErr != nil {
            deferr = errors.WithStack(tempErr)
        }
    }()

    err = dbinst.Transaction(func(tx *gorm.DB) error {
        ctx = context.WithValue(ctx, "tx", tx)
        root, err := createRoot(ctx, tx)
        if err != nil {
            return errors.WithStack(err)
        }

        var isDbRecordFound bool
        dataobj, isDbRecordFound, err = createData(ctx, tx, root, in)
        if err != nil {
            return errors.WithStack(err)
        }
        if isDbRecordFound {
            return nil
        }

        ips, err := createIPs(ctx, tx, createIPsInputType{ParentDataID: dataobj.ID, Domain: in.DomainSub})
        if err != nil {
            return errors.WithStack(err)
        }

        var ipsWithAdditionalData []ip.Model
        for _, ipaddr := range ips {
            ipwad, err := createLocationAndUrbanArea(ctx, tx, ipaddr, "ParentSubDomainID")
            if err != nil {
                return errors.WithStack(err)
            }
            ipsWithAdditionalData = append(ipsWithAdditionalData, ipwad)
        }

        dataobj.IPs = ipsWithAdditionalData

        sdms, err := createSubdomains(ctx, tx, dataobj, in)
        if err != nil {
            return errors.WithStack(err)
        }
        for ix := range sdms {
            ips, err = createIPs(
                ctx, tx, createIPsInputType{ParentSubDomainID: sdms[ix].ID, Domain: sdms[ix].Domain},
            )
            if err != nil {
                return errors.WithStack(err)
            }

            ipsWithAdditionalData = []ip.Model{}
            for _, ipaddr := range ips {
                ipwad, err := createLocationAndUrbanArea(ctx, tx, ipaddr, "ParentDataID")
                if err != nil {
                    return errors.WithStack(err)
                }
                ipsWithAdditionalData = append(ipsWithAdditionalData, ipwad)
            }
            sdms[ix].IPs = ipsWithAdditionalData
        }

        dataobj.SubDomains = sdms

        return nil
    })
    if err != nil {
        return empty, errors.WithStack(err)
    }

    return dataobj, deferr
}

func cleanAll(ctx context.Context) error {
    // Enough to clean all because deletion cascades to child tables.
    _, err := user.DeleteAll(ctx)
    if err != nil {
        return errors.WithStack(err)
    }

    return nil
}

func createLocationAndUrbanArea(
    ctx context.Context, tx *gorm.DB, ipaddr ip.Model,
    ipFieldToOmit string,
) (ip.Model, error) {
    var empty ip.Model
    loc, ipDataFromLoc, recordIsFound, err := createLocation(ctx, tx, ipaddr.Address)
    if err != nil {
        return empty, errors.WithStack(err)
    }

    if loc.IsEmpty() {
        return ipaddr, nil
    }
    ipaddr.Location = loc
    ipaddr.ISP = ipDataFromLoc.ISP
    ipaddr.Org = ipDataFromLoc.Org
    ipaddr.AS = ipDataFromLoc.AS

    if err = tx.Omit(ipFieldToOmit).Save(&ipaddr).Error; err != nil {
        return empty, errors.WithStack(err)
    }
    if recordIsFound {
        return ipaddr, nil
    }

    // ua, locationData, err := createUrbanArea(ctx, tx, loc.City, loc.Region, loc.Country)
    // if err != nil {
    //     return empty, errors.WithStack(err)
    // }

    // uatx := tx
    // if ua.IsEmpty() {
    //     uatx = tx.Omit("UrbanArea", "UrbanAreaID")
    // } else {
    //     loc.UrbanArea = ua
    // }
    // loc.Population = locationData.Population
    // loc.TeleportGeoNameID = locationData.TeleportGeoNameID
    // loc.TeleportGeoHash = locationData.TeleportGeoHash
    // loc.TeleportName = locationData.TeleportName
    // loc.TeleportFullName = locationData.TeleportFullName
    // loc.TeleportLatitude = locationData.TeleportLatitude
    // loc.TeleportLongitude = locationData.TeleportLongitude

    // if err = uatx.Save(&loc).Error; err != nil {
    //     return empty, errors.WithStack(err)
    // }

    ipaddr.Location = loc
    return ipaddr, nil
}
