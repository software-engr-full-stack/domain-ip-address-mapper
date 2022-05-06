package app

import (
    "context"
    "strings"

    "github.com/pkg/errors"
    "gorm.io/gorm"

    "demo/entities/ip"
)

type createIPsInputType struct {
    ParentDataID uint
    ParentSubDomainID uint
    Domain string
}

func createIPs(ctx context.Context, tx *gorm.DB, in createIPsInputType) ([]ip.Model, error) {
    var empty []ip.Model
    if in.ParentDataID == 0 && in.ParentSubDomainID == 0 {
        return empty, errors.New("must pass owner id")
    }
    if strings.TrimSpace(in.Domain) == "" {
        return empty, errors.New("must pass domain")
    }

    ips, areFound, err := ip.LookupByDomain(ctx, in.Domain)
    if err != nil {
        return empty, errors.WithStack(err)
    }
    if !areFound {
        return empty, nil
    }

    var models []ip.Model
    for _, ipaddr := range ips {
        switch true {
        case in.ParentDataID != 0:
            ipaddr.ParentDataID = in.ParentDataID
        case in.ParentSubDomainID != 0:
            ipaddr.ParentSubDomainID = in.ParentSubDomainID
        default:
            return empty, errors.New("should be unreachable")
        }

        if err = tx.Create(&ipaddr).Error; err != nil {
            return empty, errors.WithStack(err)
        }
        models = append(models, ipaddr)
    }

    return models, nil
}
