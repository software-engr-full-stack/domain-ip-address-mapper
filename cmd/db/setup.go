package db

import (
    "fmt"

    "github.com/pkg/errors"

    "demo/config"
    "demo/db/manager"
)

func Setup(op string) error {
    setup, err := config.Setup()
    if err != nil {
        return errors.WithStack(err)
    }

    ctx := setup.Context

    switch op {
    case "db/setup/create":
        err := manager.CreateDatabase(ctx)
        if err != nil {
            return errors.WithStack(err)
        }

    case "db/setup/drop":
        err := manager.DropDatabase(ctx)
        if err != nil {
            return errors.WithStack(err)
        }

    case "db/setup/reset":
        err := manager.ResetDatabase(ctx)
        if err != nil {
            return errors.WithStack(err)
        }

    case "db/setup/migrate":
        err := manager.Migrate(ctx)
        if err != nil {
            return errors.WithStack(err)
        }

    case "db/setup/reset-and-migrate":
        err := manager.ResetDatabase(ctx)
        if err != nil {
            return errors.WithStack(err)
        }

        err = manager.Migrate(ctx)
        if err != nil {
            return errors.WithStack(err)
        }

    default:
        return errors.New(fmt.Sprintf("invalid op %#v", op))
    }

    return nil
}
