package handler

import (
    "context"
    "net/http"
    "encoding/json"
    "io"

    "github.com/pkg/errors"

    "demo/entities/crobat"
    "demo/services/app"
    "fmt"
)

func Index(ctx context.Context, wri http.ResponseWriter, req *http.Request) error {
    domain, err := parseDomain(req)
    if err != nil {
        return err
    }

    in := crobat.InputType{
        DomainSub:  domain,
        UniqueSort: true,
    }

    data, err := app.Service(ctx, in)
    if err != nil {
        return StatusError{500, errors.WithStack(err)}
    }

    JSONResponse(wri, data, nil)

    return nil
}

func parseDomain(req *http.Request) (string, error) {
    var temp map[string]interface{}
    err := json.NewDecoder(req.Body).Decode(&temp)
    var empty string
    switch {
    case err == io.EOF:
        return empty, StatusError{500, errors.New(fmt.Sprintf("empty request body"))}
    case err != nil:
        return empty, StatusError{500, errors.WithStack(err)}
    }
    defer req.Body.Close()

    domainKey := "domain"
    if domain, ok := temp[domainKey]; ok {
        return domain.(string), nil
    }

    return empty, StatusError{
        500, errors.New(fmt.Sprintf("%#v not found in request body", domainKey)),
    }
}
