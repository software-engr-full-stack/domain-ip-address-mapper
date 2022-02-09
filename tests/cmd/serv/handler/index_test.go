package handler

import (
    "testing"
    "net/http"
    "net/http/httptest"
    "context"
    "strings"
    "encoding/json"
    "bytes"
    "io/ioutil"

    "demo/lib/testlib"
    "demo/lib/testlib/common"
    "demo/lib/testlib/common/reddit"

    "demo/cmd/serv/handler"
    "demo/entities/user"
    "demo/entities/data"
)

func TestIndex(t *testing.T) {
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

    testBadRequestResponse(t, ctx)

    testResponse(t, ctx)

    // Make the same request. Make sure it pulls the data from the database.
    //   Make sure no additional data is entered into the database.
    testResponse(t, ctx)
}

func testBadRequestResponse(t *testing.T, ctx context.Context) {
    req := httptest.NewRequest(http.MethodGet, "https://example.com", nil)
    wri := httptest.NewRecorder()
    err := handler.Index(ctx, wri, req)
    if actual, expected := err.Error(), "empty request body"; !strings.Contains(err.Error(), expected) {
        t.Fatalf("actual %#v should contain %#v", actual, expected)
    }

    if actual, expected := wri.Result().StatusCode, http.StatusOK; actual != expected {
        t.Fatalf("actual %#v != expected %#v", actual, expected)
    }

    payload := map[string]interface{}{
        "not-a-domain": "...",
    }
    body, err := json.Marshal(payload)
    if err != nil {
        panic(err)
    }
    req = httptest.NewRequest(http.MethodGet, "https://example.com", bytes.NewReader(body))
    wri = httptest.NewRecorder()
    err = handler.Index(ctx, wri, req)
    if actual, expected := err.Error(), `"domain" not found in request body`; !strings.Contains(err.Error(), expected) {
        t.Fatalf("actual %#v should contain %#v", actual, expected)
    }
}

func testResponse(t *testing.T, ctx context.Context) {
    payload := map[string]interface{}{
        "domain": "reddit.com",
    }
    body, err := json.Marshal(payload)
    if err != nil {
        panic(err)
    }
    req := httptest.NewRequest(http.MethodGet, "https://example.com", bytes.NewReader(body))
    wri := httptest.NewRecorder()
    err = handler.Index(ctx, wri, req)
    if err != nil {
        panic(err)
    }
    if actual, expected := wri.Result().StatusCode, http.StatusOK; actual != expected {
        t.Fatalf("actual %#v != %#v", actual, expected)
    }

    res := wri.Result()
    defer res.Body.Close()
    body, err = ioutil.ReadAll(res.Body)
    if err != nil {
        panic(err)
    }

    var temp map[string]interface{}
    if err := json.Unmarshal(body, &temp); err != nil {
        panic(err)
    }

    type wrapperType struct {
        OK bool `json:"ok"`
        Results data.Model `json:"results"`
    }
    var wrapper wrapperType
    if err := json.Unmarshal(body, &wrapper); err != nil {
        panic(err)
    }
    if actual, expected := wrapper.OK, true; actual != expected {
        t.Fatalf("actual %#v != expected %#v", actual, expected)
    }

    // Test if the data model returned by the service is correct.
    reddit.TestData(t, ctx, []data.Model{wrapper.Results})
    reddit.TestSubDomains(t, ctx, []data.Model{wrapper.Results})

    // Test if the things are stored in the database properly
    dataList, err := data.All(ctx)
    if err != nil {
        panic(err)
    }
    if actualCount, expectedCount := len(dataList), 1; actualCount != expectedCount {
        t.Fatalf("actual %#v != expected %#v", actualCount, expectedCount)
    }

    users, err := user.All(ctx)
    if err != nil {
        panic(err)
    }
    common.TestUsers(t, ctx, users)

    rootUser, isFound, err := user.Root(ctx)
    if err != nil {
        panic(err)
    }
    common.TestRootUser(t, ctx, rootUser, isFound)

    reddit.TestData(t, ctx, rootUser.Data)
    reddit.TestSubDomains(t, ctx, rootUser.Data)
}
