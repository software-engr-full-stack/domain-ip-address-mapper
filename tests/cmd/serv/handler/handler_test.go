package handler

import (
    "testing"
    "net/http"
    "net/http/httptest"
    "context"
    "encoding/json"
    "io/ioutil"

    "fmt"

    "demo/lib/testlib"

    "demo/cmd/serv/handler"
)

func TestHandler(t *testing.T) {
    t.Parallel()
    setup := testlib.Setup()
    setup.TestDB.Setup()
    defer setup.TestDB.CleanUp()

    instance, err := setup.DBObj.Open()
    if err != nil {
        panic(err)
    }
    defer func() {
        if err = setup.DBObj.Close(); err != nil {
            panic(err)
        }
    }()

    ctx := setup.ContextWithDBInstance(instance)

    type errorJSONType struct {
        Code int `json:"code"`
        Error string `json:"error"`
    }

    type testType struct {
        code int
        msg string
        title string
        isNoError bool
    }

    tests := []testType{
        testType{title: "StatusError type", code: 501, msg: "TEST ERROR fbe0ebd0-360b-4664-a05f-5c5568d22811"},
        testType{title: "default error type", code: http.StatusInternalServerError, msg: http.StatusText(http.StatusInternalServerError)},
        testType{title: "no error", code: http.StatusOK, isNoError: true},
    }

    for _, test := range tests {
        req := httptest.NewRequest(http.MethodGet, "https://example.com", nil)
        wri := httptest.NewRecorder()

        hfunc := handlerFactory(handler.StatusError{test.code, fmt.Errorf(test.msg)}, ctx, wri, req)
        if test.isNoError {
            hfunc = handlerFactory(nil, ctx, wri, req)
        }

        hobj := handler.Handler{Context: ctx, H: hfunc, Method: http.MethodGet}

        hobj.ServeHTTP(wri, req)

        res := wri.Result()
        if actual, expected := res.StatusCode, test.code; actual != expected {
            t.Fatalf("%#v, response status code: actual %#v != expected %#v", test.title, actual, expected)
        }

        defer res.Body.Close()
        body, err := ioutil.ReadAll(res.Body)
        if err != nil {
            panic(err)
        }

        if test.isNoError {
            if actual, expected := len(body), 0; actual != expected {
                t.Fatalf("%#v, response body length: actual %#v != expected %#v", test.title, actual, expected)
            }
            continue
        }

        var errJSON errorJSONType
        if err := json.Unmarshal(body, &errJSON); err != nil {
            panic(err)
        }

        if actual, expected := errJSON.Code, test.code; actual != expected {
            t.Fatalf("%#v, JSON status code: actual %#v != expected %#v", test.title, actual, expected)
        }
        if actual, expected := errJSON.Error, test.msg; actual != expected {
            t.Fatalf("%#v, JSON error message: actual %#v != expected %#v", test.title, actual, expected)
        }
    }
}

func handlerFactory(
    err error,
    ctx context.Context, wri http.ResponseWriter, req *http.Request,
) (func (ctx context.Context, wri http.ResponseWriter, req *http.Request) error) {
    return func(ctx context.Context, wri http.ResponseWriter, req *http.Request) error {
        return err
    }
}
