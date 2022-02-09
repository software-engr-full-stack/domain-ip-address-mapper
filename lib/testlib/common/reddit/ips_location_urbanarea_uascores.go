package reddit

import (
    "testing"
    "reflect"
    "sort"

    "demo/lib/testlib"

    "demo/entities/ip"
)

func testIPsLocationUrbanAreaUAScores(t *testing.T, actualList []ip.Model, expectedList []ip.Model) {
    actualCount := len(actualList)
    expectedCount := len(expectedList)
    if actualCount != expectedCount {
        t.Fatalf("count actual %#v != expected %#v", actualCount, expectedCount)
    }

    sort.SliceStable(actualList, func(i, j int) bool {
        return actualList[i].Address.String() < actualList[j].Address.String()
    })

    sort.SliceStable(expectedList, func(i, j int) bool {
        return expectedList[i].Address.String() < expectedList[j].Address.String()
    })

    for ix := range expectedList {
        actual := actualList[ix]
        expected := expectedList[ix]

        if actual.Address.String() != expected.Address.String() {
            t.Fatalf("actual %#v != expected %#v", actual.Address.String(), expected.Address.String())
        }

        testField := testlib.NewTestField(testlib.TestFieldType{
            Actual: reflect.ValueOf(&actual),
            Expected: reflect.ValueOf(&expected),
            Index: ix,
            T: t,
        })
        testField.Run(
            "Domain",
        )

        // Location
        testField = testlib.NewTestField(testlib.TestFieldType{
            Actual: reflect.ValueOf(&actual.Location),
            Expected: reflect.ValueOf(&expected.Location),
            Index: ix,
            T: t,
        })
        testField.Run(
            "City",
            "RegionCode",
            "CountryCode",
            "Latitude",
            "Longitude",
            "FullName",
            "TeleportFullName",
            "TeleportGeoNameID",
        )

        actualUA := actual.Location.UrbanArea
        expectedUA := expected.Location.UrbanArea
        testField = testlib.NewTestField(testlib.TestFieldType{
            Actual: reflect.ValueOf(&actualUA),
            Expected: reflect.ValueOf(&expectedUA),
            Index: ix,
            T: t,
        })
        testField.Run(
            "Name",
            "Link",
        )

        actualUAS := actualUA.Scores[:]
        expectedUAS := expectedUA.Scores[:]

        actualCount = len(actualUAS)
        expectedCount = len(expectedUAS)
        if actualCount != expectedCount {
            t.Fatalf("count actual %#v != expected %#v", actualCount, expectedCount)
        }

        sort.SliceStable(actualUAS, func(i, j int) bool {
            return actualUAS[i].Name < actualUAS[j].Name
        })
        sort.SliceStable(expectedUAS, func(i, j int) bool {
            return expectedUAS[i].Name < expectedUAS[j].Name
        })

        for ixix := range expectedUAS {
            testField = testlib.NewTestField(testlib.TestFieldType{
                Actual: reflect.ValueOf(&actualUAS[ixix]),
                Expected: reflect.ValueOf(&expectedUAS[ixix]),
                Index: ix,
                T: t,
            })
            testField.Run(
                "Name",
                "ScoreOutOf10",
            )
        }
    }
}
