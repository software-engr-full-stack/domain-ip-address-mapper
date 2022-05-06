package util

import (
    "testing"
    "reflect"

    "demo/entities/user"
    "demo/entities/data"
    "demo/entities/subdomain"
    "demo/entities/ip"
    "demo/entities/location"
    "demo/entities/urbanarea"

    dbutil "demo/db/util"
)

func TestPreloadStringForDataGenerator(t *testing.T) {
    actual := dbutil.PreloadStringForData(reflect.TypeOf(user.Model{}).PkgPath())
    expected := "Data.IPs.Location.UrbanArea.Scores"

    if actual != expected {
        t.Fatalf("actual %#v != expected %#v", actual, expected)
    }

    actual = dbutil.PreloadStringForData(reflect.TypeOf(data.Model{}).PkgPath())
    expected = "IPs.Location.UrbanArea.Scores"

    if actual != expected {
        t.Fatalf("actual %#v != expected %#v", actual, expected)
    }
}

func TestPreloadStringForSubDomainsGenerator(t *testing.T) {
    actual := dbutil.PreloadStringForSubDomains(reflect.TypeOf(user.Model{}).PkgPath())
    expected := "Data.SubDomains.IPs.Location.UrbanArea.Scores"

    if actual != expected {
        t.Fatalf("actual %#v != expected %#v", actual, expected)
    }

    actual = dbutil.PreloadStringForSubDomains(reflect.TypeOf(data.Model{}).PkgPath())
    expected = "SubDomains.IPs.Location.UrbanArea.Scores"

    if actual != expected {
        t.Fatalf("actual %#v != expected %#v", actual, expected)
    }
}

func TestPreloadString(t *testing.T) {
    actual := dbutil.PreloadString(reflect.TypeOf(subdomain.Model{}).PkgPath())
    expected := "IPs.Location.UrbanArea.Scores"

    if actual != expected {
        t.Fatalf("actual %#v != expected %#v", actual, expected)
    }

    actual = dbutil.PreloadString(reflect.TypeOf(ip.Model{}).PkgPath())
    expected = "Location.UrbanArea.Scores"

    if actual != expected {
        t.Fatalf("actual %#v != expected %#v", actual, expected)
    }

    actual = dbutil.PreloadString(reflect.TypeOf(location.Model{}).PkgPath())
    expected = "UrbanArea.Scores"

    if actual != expected {
        t.Fatalf("actual %#v != expected %#v", actual, expected)
    }

    actual = dbutil.PreloadString(reflect.TypeOf(urbanarea.Model{}).PkgPath())
    expected = "Scores"

    if actual != expected {
        t.Fatalf("actual %#v != expected %#v", actual, expected)
    }
}
