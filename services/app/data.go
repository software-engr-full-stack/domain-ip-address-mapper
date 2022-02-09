package app

import (
    "context"
    "fmt"

    "github.com/pkg/errors"
    "gorm.io/gorm"

    "demo/entities/user"
    "demo/entities/data"
    "demo/entities/crobat"
)

func createData(
    ctx context.Context, tx *gorm.DB, root user.Model, in crobat.InputType,
) (dmph data.Model, isDbRecordFound bool, eph error) {
    dobj := data.Model{Domain: in.DomainSub}

    uniq, err := dobj.Unique(ctx)
    var empty data.Model
    var empty2 bool
    if err != nil {
        return empty, empty2, errors.WithStack(err)
    }
    if uniq.IsFound {
        return uniq.Record, true, nil
    } else {
        if uniq.AreRecordsFound {
            return empty, empty2, errors.New(fmt.Sprintf(
                "multiple data records of domain %#v are found", in.DomainSub,
            ))
        }
        dobj.UserID = root.ID
        if err = tx.Create(&dobj).Error; err != nil {
            return empty, empty2, errors.WithStack(err)
        }
    }

    return dobj, false, nil
}
