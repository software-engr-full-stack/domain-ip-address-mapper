package serv

import (
    "net/http"
    "log"

    "github.com/rs/cors"

    "demo/config"
    "demo/cmd/serv/handler"
)

const testURLPath = "/test/"

func HTTP() {
    setup, err := config.Setup()
    if err != nil {
        panic(err)
    }

    instance, err := setup.DBObj.Open()
    if err != nil {
        panic(err)
    }
    defer func() {
        if err := setup.DBObj.Close(); err != nil {
            panic(err)
        }
    }()

    ctx, err := setup.ContextWithDBInstance(instance)
    if err != nil {
        panic(err)
    }

    http.HandleFunc(testURLPath, serveFile)

    cobj := cors.New(cors.Options{
        AllowedOrigins: []string{"http://localhost:3000"},
    })
    // TODO: require POST request then create a handler for GET
    http.Handle(
        "/",
        cobj.Handler(handler.Handler{Context: ctx, H: handler.Index}),
    )

    log.Fatal(http.ListenAndServe(":8000", nil))
}

func serveFile(w http.ResponseWriter, r *http.Request) {
    p := "." + r.URL.Path
    if r.URL.Path == testURLPath {
        p = "./frontend/static/leaflet.html"
    }
    http.ServeFile(w, r, p)
}
