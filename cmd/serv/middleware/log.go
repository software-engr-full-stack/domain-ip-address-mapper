package middleware

import (
    "net/http"
    "log"
    "strings"

    "fmt"
)

func WithLogging(h http.Handler) http.Handler {
    loggingFn := func(wri http.ResponseWriter, req *http.Request) {
        possibleForwarded := []string{
            req.Header.Get("X-Forwarded-For"),
            req.Header.Get("x-forwarded-for"),
            req.Header.Get("X-FORWARDED-FOR"),
        }
        noEmpty := []string{}
        for _, fw := range possibleForwarded {
            trimmed := strings.TrimSpace(fw)
            if trimmed == "" {
                continue
            }
            noEmpty = append(noEmpty, trimmed)
        }

        log.Printf(
            "... %#v %#v %#v %#v",
            req.RemoteAddr,
            req.RequestURI,
            req.Method,
            fmt.Sprintf("forwarded=%v", noEmpty),
        )
        h.ServeHTTP(wri, req)
    }
    return http.HandlerFunc(loggingFn)
}
