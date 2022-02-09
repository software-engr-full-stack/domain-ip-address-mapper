package debug

import (
    "fmt"

    "github.com/pkg/errors"

    "demo/config"
    "demo/entities/user"
    "demo/entities/data"
    "demo/entities/ip"
)

func Debug(op string) (deferr error) {
    setup, err := config.Setup()
    if err != nil {
        return errors.WithStack(err)
    }

    instance, err := setup.DBObj.Open()
    if err != nil {
        return errors.WithStack(err)
    }
    defer func() {
        if tempErr := setup.DBObj.Close(); tempErr != nil {
            deferr = errors.WithStack(tempErr)
        }
    }()

    ctx, err := setup.ContextWithDBInstance(instance)
    if err != nil {
        return errors.WithStack(err)
    }

    root, isFound, err := user.Root(ctx)
    if err != nil {
        return errors.WithStack(err)
    }
    if !isFound {
        return errors.New("root user must exist")
    }

    switch op {
    case "debug/debug":
        showLoc := func(domain string, ipaddr ip.Model) {
            loc := ipaddr.Location
            ua := loc.UrbanArea
            fmt.Printf(
                "%#v %#v %#v %#v %#v %#v\n",
                domain,
                ipaddr.Address.String(),
                loc.FullName,
                loc.ViewFullName,
                loc.TeleportFullName,
                len(ua.Scores),
            )
        }
        for _, dobj := range root.Data {
            fmt.Printf("... %#v\n", len(dobj.IPs))
            for _, ipaddr := range dobj.IPs {
                showLoc(dobj.Domain, ipaddr)
            }

            for _, sdm := range dobj.SubDomains {
                for _, ipaddr := range sdm.IPs {
                    showLoc(sdm.Domain, ipaddr)
                }
            }
            fmt.Println()
        }

    case "debug/default":
        _, err := data.All(ctx)
        if err != nil {
            return errors.WithStack(err)
        }
        fmt.Println()
    default:
        return errors.New(fmt.Sprintf("invalid op %#v", op))
    }

    return deferr
}
