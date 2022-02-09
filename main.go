package main

import (
    "os"
    "path/filepath"
    "strings"

    "fmt"

    "demo/cmd/serv"
    "demo/cmd/db"
    "demo/cmd/debug"
)

func main() {
    op, ignore := parseop()
    if ignore {
        return
    }

    switch true {
    case strings.Contains(op, "serv/http"):
        serv.HTTP()

    case strings.Contains(op, "db/setup"):
        err := db.Setup(op)
        if err != nil {
            frkout(err)
        }
    case strings.Contains(op, "db/testseed"):
        err := db.TestSeed(op)
        if err != nil {
            frkout(err)
        }

    case strings.Contains(op, "debug"):
        err := debug.Debug(op)
        if err != nil {
            frkout(err)
        }

    default:
        frkout(fmt.Errorf("invalid op %#v", op))
    }
}

func frkout(err error) {
    fmt.Printf("... ERROR: %+v\n", err)
    panic("...")
}

func parseop() (op string, ignore bool) {
    arg0 := os.Args[0]

    var empty string
    // Live reloading binary, ignore
    if strings.Contains(arg0, "tmp/main") {
        return empty, true
    }

    dir, err := filepath.Abs(filepath.Dir(arg0))
    if err != nil {
        panic(err)
    }
    symlinkFullPath := filepath.Join(dir, filepath.Base(arg0))

    afterBin := strings.Split(symlinkFullPath, "bin")

    components := strings.Split(afterBin[1], string(os.PathSeparator))
    noEmptyStrings := []string{}
    for _, str := range components {
        trimmed := strings.TrimSpace(str)
        if trimmed != "" {
            noEmptyStrings = append(noEmptyStrings, trimmed)
        }
    }

    return strings.Join(noEmptyStrings, string(os.PathSeparator)), false
}
