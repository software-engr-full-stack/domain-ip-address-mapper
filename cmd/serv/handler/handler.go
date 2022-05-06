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
    Method string
}

func (h Handler) ServeHTTP(wri http.ResponseWriter, req *http.Request) {
    if h.Method == "" {
        panic("HTTP method must be defined")
    }

    if !h.isMethodValid(wri, req) {
        return
    }

    err := h.H(h.Context, wri, req)
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

        JSONErrorResponse(wri, jsonErr, &JSONResponseOptionsType{Code: status})
    }
}

func (h Handler) isMethodValid(wri http.ResponseWriter, req *http.Request) bool {
    if req.Method == h.Method {
        return true
    }

    msg := fmt.Sprintf(
        "request method %#v != to expected method %#v for this path",
        req.Method, h.Method,
    )
    code := http.StatusOK
    jsonErr, err := json.Marshal(JSONErrorType{
        Code: code,
        Error: msg,
    })
    if err != nil {
        panic(err)
    }

    JSONErrorResponse(
        wri,
        jsonErr,
        &JSONResponseOptionsType{Code: code},
    )
    log.Printf("... ERROR: %s", msg)
    return false
}

type JSONErrorType struct {
    Code int `json:"code"`
    Error string `json:"error"`
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
