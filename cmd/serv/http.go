package serv

import (
    "net/http"
    "log"
    "os"
    "io"

    "fmt"
    "time"

    "github.com/rs/cors"

    "demo/config"
    "demo/cmd/serv/handler"
    "demo/cmd/serv/middleware"
)

func HTTP() {
    logFile, err := os.OpenFile(
        "/var/log/app/domain-ip-address-mapper_backend.log",
        os.O_RDWR | os.O_CREATE | os.O_APPEND, 0644,
    )
    if err != nil {
        log.Fatalf("error opening file: %#v", err)
    }
    defer logFile.Close()
    mw := io.MultiWriter(os.Stdout, logFile)
    log.SetOutput(mw)

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

    mux := http.NewServeMux()

    cobj := cors.New(cors.Options{
        AllowedOrigins: setup.Config.CORSAllowedOrigins,
    })

    mux.HandleFunc("/index", index)

    mux.Handle(
        "/api",
        handler.Handler{Context: ctx, H: handler.Index, Method: http.MethodPost},
    )

    feHandler := http.FileServer(http.Dir("frontend/build"))
    mux.Handle("/", feHandler)

    srv := &http.Server{
        Handler:      cobj.Handler(middleware.WithLogging(mux)),
        Addr:         "0.0.0.0:8000",
        WriteTimeout: 15 * time.Second,
        ReadTimeout:  15 * time.Second,
    }

    fmt.Println("Server started on PORT 8000")
    log.Fatal(srv.ListenAndServe())
}

func index(wri http.ResponseWriter, req *http.Request) {
    http.ServeFile(wri, req, "frontend/build/index.html")
}
