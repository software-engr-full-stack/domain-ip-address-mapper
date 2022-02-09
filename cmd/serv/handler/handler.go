package handler

import (
    "context"
    "net/http"
    "encoding/json"
    "log"
    "fmt"
)

type Error interface {
    error
    Status() int

    // ...
    JSON() string
}

type StatusError struct {
    Code int
    Err  error
}

func (se StatusError) Error() string {
    // With stack
    return fmt.Sprintf("%+v", se.Err)
    // // Terse
    // return se.Err.Error()
}

func (se StatusError) Status() int {
    return se.Code
}

func (se StatusError) JSON() string {
    return se.Err.Error()
}

type Handler struct {
    Context context.Context
    H func(ctx context.Context, w http.ResponseWriter, r *http.Request) error
}

type JSONErrorType struct {
    Code int `json:"code"`
    Error string `json:"error"`
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    err := h.H(h.Context, w, r)
    if err != nil {
        var jsonErr []byte
        var status int
        var msg string
        switch e := err.(type) {
        case Error:
            status = e.Status()
            msg = e.JSON()
            jsonErr, err = json.Marshal(JSONErrorType{Code: status, Error: msg})
            if err != nil {
                panic(err)
            }
            log.Printf("... ERROR: JSON message %s", jsonErr)
            log.Printf("... ERROR: HTTP %d - %s", status, e)
        default:
            status = http.StatusInternalServerError
            msg = http.StatusText(status)
            jsonErr, err = json.Marshal(JSONErrorType{Code: status, Error: msg})
            if err != nil {
                panic(err)
            }
        }

        JSONErrorResponse(w, jsonErr, &JSONResponseOptionsType{Code: status})
        // http.Error(w, msg, status)
    }
}

type JSONResponseOptionsType struct {
    Code int
}

func JSONErrorResponse(w http.ResponseWriter, resp []byte, options *JSONResponseOptionsType) {
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.Header().Set("X-Content-Type-Options", "nosniff")
    if options != nil {
        if options.Code != 0 {
            w.WriteHeader(options.Code)
        }
    }
    fmt.Fprintln(w, string(resp))
}

func JSONResponse(w http.ResponseWriter, in interface{}, options *JSONResponseOptionsType) {
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.Header().Set("X-Content-Type-Options", "nosniff")
    if options != nil {
        if options.Code != 0 {
            w.WriteHeader(options.Code)
        }
    }

    data, err := json.Marshal(map[string]interface{}{
        "results": in,
        "ok": true,
    })
    if err != nil {
        panic(err)
    }

    fmt.Fprintln(w, string(data))
}
