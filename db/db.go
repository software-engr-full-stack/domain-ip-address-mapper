package db

import (
    "strings"
    "fmt"

    "github.com/pkg/errors"

    "gorm.io/gorm"
    "gorm.io/driver/postgres"
)

type Type struct {
    Config ConfigType
    Instance *gorm.DB
}

type InputType struct {
    Env string
    DBName string
}

func New(in InputType) (*Type, error) {
    config, err := NewConfig(in)
    var empty *Type
    if err != nil {
        return empty, errors.WithStack(err)
    }

    return &Type{Config: config}, nil
}

func (t *Type) Open() (*gorm.DB, error) {
    empty := &gorm.DB{}
    config := t.Config
    connectionSpec := []string{
        fmt.Sprintf("host=%s", config.Host),
        fmt.Sprintf("dbname=%s", config.Name),
        "sslmode=disable",
    }
    if config.User != "" {
        connectionSpec = append(connectionSpec, fmt.Sprintf("user=%s", config.User))
    }
    if config.Password != "" {
        connectionSpec = append(connectionSpec, fmt.Sprintf("password=%s", config.Password))
    }
    if config.Port != 0 {
        connectionSpec = append(connectionSpec, fmt.Sprintf("port=%d", config.Port))
    }
    instance, err := gorm.Open(postgres.Open(strings.Join(connectionSpec, " ")), &gorm.Config{})
    if err != nil {
        return empty, errors.WithStack(err)
    }

    t.Instance = instance

    return instance, nil
}

func (t *Type) Close() error {
    sqldb, err := t.Instance.DB()
    if err != nil {
        return errors.WithStack(err)
    }
    return sqldb.Close()
}
