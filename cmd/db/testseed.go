package db

import (
    "os"
    "fmt"

    "github.com/pkg/errors"

    "demo/config"
    "demo/db/testseed"
    "demo/entities/crobat"
)

func TestSeed(op string) error {
    setup, err := config.Setup()
    if err != nil {
        return errors.WithStack(err)
    }

    ctx := setup.Context

    switch op {
    case "db/testseed/testseed":
        domain := "reddit.com"
        if len(os.Args) > 1 {
            domain = os.Args[1]
        }
        in := crobat.InputType{
            DomainSub:  domain,
            UniqueSort: true,
        }
        err = testseed.TestSeed(ctx, in)
        if err != nil {
            return errors.WithStack(err)
        }

    default:
        return errors.New(fmt.Sprintf("invalid op %#v", op))
    }

    return nil
}
