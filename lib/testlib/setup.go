package testlib

import (
    "context"

    "gorm.io/gorm"

    "demo/config"
    "demo/db"
)

type SetupType struct {
    TestDB TestDbType
    DBObj *db.Type
}

func Setup() SetupType {
    tdb := NewTestDbType()

    env, err := config.NewEnv()
    if err != nil {
        panic(err)
    }
    dbobj, err := db.New(db.InputType{
        Env: env.LowerCaseBaseName,
        DBName: tdb.Name,
    })
    if err != nil {
        panic(err)
    }

    return SetupType{TestDB: tdb, DBObj: dbobj}
}

func (stp *SetupType) ContextWithDBInstance(instance *gorm.DB) context.Context {
    cfg, err := config.New(config.Type{DB: instance, DBObj: stp.DBObj})
    if err != nil {
        panic(err)
    }

    return context.WithValue(context.Background(), "config", cfg)
}
