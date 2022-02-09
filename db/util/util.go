package util

import (
    "context"

    "github.com/pkg/errors"
    "gorm.io/gorm"

    "demo/config"
    "demo/db"
)

func FindOpenConnection(ctx context.Context) (*gorm.DB, *db.Type, error) {
    cfg := ctx.Value("config").(config.Type)
    txif := ctx.Value("tx")

    var dborm *gorm.DB

    dbobj := cfg.DBObj
    var empty0 *gorm.DB
    var empty1 *db.Type
    switch true {
    case txif != nil:
        dborm = txif.(*gorm.DB)
    case dbobj.Instance != nil:
        dborm = dbobj.Instance
    default:
        return empty0, dbobj, nil
    }

    sqldb, err := dborm.DB()
    if err != nil {
        return empty0, empty1, errors.WithStack(err)
    }
    if err = sqldb.Ping(); err != nil {
        return nil, dbobj, nil
    }

    return dborm, dbobj, nil
}
