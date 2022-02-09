package util

import (
    "strings"
    "path/filepath"
    "fmt"
)

func deleteME() {
    fmt.Printf("....")
}

const (
    joiner     = "."
    common     = "IPs.Location.UrbanArea.Scores"
    data       = "Data" + joiner + "IPs.Location.UrbanArea.Scores"
    subdomains = "Data" + joiner + "SubDomains" + joiner + "IPs.Location.UrbanArea.Scores"
)

var pkgBaseToRelationNameMap = map[string]string{
    "data": "Data",
    "ip": "IPs",
    "location": "Location",
    "urbanarea": "UrbanArea",
}

func PreloadStringForData(pkgPath string) string {
    base := filepath.Base(pkgPath)

    if relationName, ok := pkgBaseToRelationNameMap[base]; ok {
        rnames := strings.Split(data, strings.Join([]string{relationName, joiner}, ""))[1:]

        return strings.Join(rnames, joiner)
    }

    return data
}

func PreloadStringForSubDomains(pkgPath string) string {
    base := filepath.Base(pkgPath)

    if relationName, ok := pkgBaseToRelationNameMap[base]; ok {
        rnames := strings.Split(subdomains, strings.Join([]string{relationName, joiner}, ""))[1:]

        return strings.Join(rnames, joiner)
    }

    return subdomains
}

func PreloadString(pkgPath string) string {
    base := filepath.Base(pkgPath)

    if relationName, ok := pkgBaseToRelationNameMap[base]; ok {
        rnames := strings.Split(subdomains, strings.Join([]string{relationName, joiner}, ""))[1:]

        return strings.Join(rnames, joiner)
    }

    return common
}
