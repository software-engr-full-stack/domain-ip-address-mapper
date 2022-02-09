package config

import (
    "context"
    "os"
    "strings"

    "github.com/pkg/errors"
    "gorm.io/gorm"

    "demo/db"
)

const (
    // Limit IP addresses to save location API usage
    defaultIPCountLimit = 5
    defaultSubDomainCountLimit = 5
    defaultSecretsFilePath = "/temp/secrets.yml"
)

type Type struct {
    DBObj *db.Type
    DB *gorm.DB
    Env EnvType
    IPCountLimit int
    SubDomainCountLimit int
    SecretsFilePath string
}

func New(in Type) (Type, error) {
    // var empty Type

    ipCountLimit := in.IPCountLimit
    if ipCountLimit == 0 {
        ipCountLimit = defaultIPCountLimit
    }

    subDomainCountLimit := in.SubDomainCountLimit
    if subDomainCountLimit == 0 {
        subDomainCountLimit = defaultSubDomainCountLimit
    }

    secretsFilePath := strings.TrimSpace(os.Getenv("SECRETS_FILE_PATH"))
    if secretsFilePath == "" {
        secretsFilePath = defaultSecretsFilePath
    }

    return Type{
        DBObj: in.DBObj,
        DB: in.DB,
        IPCountLimit: ipCountLimit,
        SubDomainCountLimit: subDomainCountLimit,
        Env: in.Env,
        SecretsFilePath: secretsFilePath,
    }, nil
}

type SetupType struct {
    DBObj *db.Type
    Context context.Context
}

func Setup() (SetupType, error) {
    env, err := NewEnv()
    var empty SetupType
    if err != nil {
        return empty, errors.WithStack(err)
    }
    dbobj, err := db.New(db.InputType{Env: env.LowerCaseBaseName})
    if err != nil {
        return empty, errors.WithStack(err)
    }

    cfg, err := New(Type{DBObj: dbobj})
    if err != nil {
        return empty, errors.WithStack(err)
    }

    return SetupType{
        Context: context.WithValue(context.Background(), "config", cfg),
        DBObj: dbobj,
    }, nil
}

func (stp *SetupType) ContextWithDBInstance(instance *gorm.DB) (context.Context, error) {
    cfg, err := New(Type{DB: instance, DBObj: stp.DBObj})
    var empty context.Context
    if err != nil {
        return empty, errors.WithStack(err)
    }

    ctx := context.WithValue(context.Background(), "config", cfg)
    stp.Context = ctx

    return ctx, nil
}
