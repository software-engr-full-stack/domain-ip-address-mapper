package ip

import (
    "context"
    "net"
    "strings"

    "github.com/pkg/errors"

    "demo/config"
)

func LookupByDomain(ctx context.Context, domain string) (models []Model, areFound bool, err error) {
    var empty []Model
    var emptyAreFound bool

    ips, err := net.LookupIP(domain)
    if err != nil {
        errMsg := err.Error()
        switch true {
        case strings.Contains(errMsg, "no such host"):
            return empty, false, nil
        }
        return empty, emptyAreFound, errors.WithStack(err)
    }

    ipCountLimit := ctx.Value("config").(config.Type).IPCountLimit
    if ipCountLimit > 0 {
        if ipCountLimit < len(ips) {
            ips = ips[:ipCountLimit]
        }
    }
    for _, netip := range ips {
        if netip.IsLoopback() ||
            netip.IsMulticast() ||
            netip.IsPrivate() ||
            netip.IsUnspecified() {
            continue
        }

        ipm := Model{Address: netip, Domain: domain}

        models = append(models, ipm)
    }

    return models, true, nil
}
