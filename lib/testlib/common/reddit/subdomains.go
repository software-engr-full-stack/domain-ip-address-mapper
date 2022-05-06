package reddit

import (
    "testing"
    "context"
    "reflect"
    "sort"
    "net"

    "demo/lib/testlib"

    "demo/entities/data"
    "demo/entities/subdomain"
    "demo/entities/ip"
    "demo/entities/urbanarea"
    "demo/entities/location"
)

func TestSubDomains(t *testing.T, ctx context.Context, dataList []data.Model) {
    expectedListUA := urbanarea.Model{
        Name: "Washington, D.C.",
        Link: "https://api.teleport.org/api/urban_areas/slug:washington-dc/",
        Scores: []urbanarea.ScoreModel{
            urbanarea.ScoreModel{Name: "Housing", ScoreOutOf10: 1.2105000000000001},
            urbanarea.ScoreModel{Name: "Cost of Living", ScoreOutOf10: 3.5949999999999993},
            urbanarea.ScoreModel{Name: "Startups", ScoreOutOf10: 8.774000000000001},
            urbanarea.ScoreModel{Name: "Venture Capital", ScoreOutOf10: 8.056},
            urbanarea.ScoreModel{Name: "Travel Connectivity", ScoreOutOf10: 4.505000000000001},
            urbanarea.ScoreModel{Name: "Commute", ScoreOutOf10: 4.457000000000001},
            urbanarea.ScoreModel{Name: "Business Freedom", ScoreOutOf10: 8.671},
            urbanarea.ScoreModel{Name: "Safety", ScoreOutOf10: 2.1915},
            urbanarea.ScoreModel{Name: "Healthcare", ScoreOutOf10: 8.490666666666666},
            urbanarea.ScoreModel{Name: "Education", ScoreOutOf10: 5.968500000000001},
            urbanarea.ScoreModel{Name: "Environmental Quality", ScoreOutOf10: 6.9937499999999995},
            urbanarea.ScoreModel{Name: "Economy", ScoreOutOf10: 6.5145},
            urbanarea.ScoreModel{Name: "Taxation", ScoreOutOf10: 4.062},
            urbanarea.ScoreModel{Name: "Internet Access", ScoreOutOf10: 3.825500000000001},
            urbanarea.ScoreModel{Name: "Leisure & Culture", ScoreOutOf10: 10},
            urbanarea.ScoreModel{Name: "Tolerance", ScoreOutOf10: 6.5495},
            urbanarea.ScoreModel{Name: "Outdoors", ScoreOutOf10: 5.023499999999999},
        },
    }

    expectedList := []subdomain.Model{
        {
            Domain: "mail-p236.reddit.com",
            IPs: []ip.Model{ip.Model{
                Domain: "mail-p236.reddit.com", Address: net.ParseIP("184.173.153.236"),
                Location: location.Model{
                    City: "Chantilly",
                    RegionCode: "VA",
                    CountryCode: "US",
                    Latitude: 38.8879,
                    Longitude: -77.4448,
                    FullName: "Chantilly, Virginia, United States, 20151",
                    TeleportFullName: "Chantilly, Virginia, United States",
                    TeleportGeoNameID: int64(4751935),
                    UrbanArea: expectedListUA,
                },
            }},
        },

        {
            Domain: "mx-02.reddit.com",
            IPs: []ip.Model{ip.Model{
                Domain: "mx-02.reddit.com", Address: net.ParseIP("52.205.61.79"),
                Location: location.Model{
                    City: "Ashburn",
                    RegionCode: "VA",
                    CountryCode: "US",

                    // DEBUG: API data changed in 2022 01 30
                    // Latitude: 39.0438,
                    // Longitude: -77.4874,
                    Latitude: 39.0469,
                    Longitude: -77.4903,

                    FullName: "Ashburn, Virginia, United States, 20149",
                    TeleportFullName: "Ashburn, Virginia, United States",
                    TeleportGeoNameID: int64(4744870),
                    UrbanArea: expectedListUA,
                },
            }},
        },
    }
    actualList := dataList[0].SubDomains[:]

    actualCount := len(actualList)
    expectedCount := len(expectedList)
    if actualCount != expectedCount {
        t.Fatalf("count actual %#v != expected %#v", actualCount, expectedCount)
    }

    sort.SliceStable(expectedList, func(i, j int) bool {
        return expectedList[i].Domain < expectedList[j].Domain
    })
    sort.SliceStable(actualList, func(i, j int) bool {
        return actualList[i].Domain < actualList[j].Domain
    })

    for ix := range expectedList {
        actual := actualList[ix]
        expected := expectedList[ix]

        testIPsLocationUrbanAreaUAScores(t, actual.IPs, expected.IPs)

        testField := testlib.NewTestField(testlib.TestFieldType{
            Actual: reflect.ValueOf(&actual),
            Expected: reflect.ValueOf(&expected),
            T: t,
        })
        testField.Run(
            "Domain",
        )
    }
}
