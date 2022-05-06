package manager

import (
    "context"
    "fmt"
    "github.com/pkg/errors"

    "demo/config"

    "demo/entities/user"
    "demo/entities/data"
    "demo/entities/subdomain"
    "demo/entities/urbanarea"
    "demo/entities/location"
    "demo/entities/ip"
)

func Migrate(ctx context.Context) (deferr error) {
    db := ctx.Value("config").(config.Type).DBObj

    _, err := db.Open()
    if err != nil {
        return errors.WithStack(err)
    }
    defer func() {
        if tempErr := db.Close(); tempErr != nil {
            deferr = errors.WithStack(tempErr)
        }
    }()

    models := []interface{}{
        &user.Model{},
        &data.Model{},
        &subdomain.Model{},
        &urbanarea.Model{},
        &urbanarea.ScoreModel{},
        &location.Model{},
        &ip.Model{},
    }

    for _, mdl := range models {
        fmt.Printf("... migrating %#T...\n", mdl)
        err := db.Instance.Migrator().CreateTable(mdl)
        if err != nil {
            return errors.WithStack(err)
        }
    }
    return deferr
}
