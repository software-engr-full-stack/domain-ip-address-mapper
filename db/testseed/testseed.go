package testseed

import (
    "context"

    "github.com/pkg/errors"

    "demo/entities/user"
    "demo/entities/crobat"

    "demo/services/app"
)

func TestSeed(ctx context.Context, in crobat.InputType) (error) {
    // err := cleanAll(ctx)
    // if err != nil {
    //     return errors.WithStack(err)
    // }

    _, err := app.Service(ctx, in)
    if err != nil {
        return errors.WithStack(err)
    }

    return nil
}

func cleanAll(ctx context.Context) error {
    // Enough to clean all because deletion cascades to child tables.
    _, err := user.DeleteAll(ctx)
    if err != nil {
        return errors.WithStack(err)
    }

    return nil
}
