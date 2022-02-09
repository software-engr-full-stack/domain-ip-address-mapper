package app

import (
    "context"
    "fmt"

    "github.com/pkg/errors"
    "gorm.io/gorm"

    "demo/config"
    "demo/entities/data"
    "demo/entities/subdomain"
    "demo/entities/crobat"
)

func createSubdomains(
    ctx context.Context, tx *gorm.DB, dobj data.Model, in crobat.InputType,
) ([]subdomain.Model, error) {
    var empty []subdomain.Model
    sdms, err := crobat.FetchSubDomains(ctx, in)
    if err != nil {
        return empty, errors.WithStack(err)
    }

    subDomainCountLimit := ctx.Value("config").(config.Type).SubDomainCountLimit
    if subDomainCountLimit > 0 {
        if subDomainCountLimit < len(sdms) {
            sdms = sdms[:subDomainCountLimit]
        }
    }

    var sdmodels []subdomain.Model
    for _, sdmstr := range sdms {
        sdm := subdomain.Model{Domain: sdmstr}
        uniq, err := sdm.Unique(ctx)
        if err != nil {
            return empty, errors.WithStack(err)
        }
        if uniq.IsFound {
            sdmodels = append(sdmodels, uniq.Record)
        } else {
            if uniq.AreRecordsFound {
                return empty, errors.New(fmt.Sprintf(
                    "multiple subdomain records of domain %#v are found", sdmstr,
                ))
            }
            sdm.ParentID = dobj.ID
            if err = tx.Create(&sdm).Error; err != nil {
                return empty, errors.WithStack(err)
            }
            sdmodels = append(sdmodels, sdm)
        }
    }
    return sdmodels, nil
}
