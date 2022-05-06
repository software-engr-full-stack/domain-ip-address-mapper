package reddit

import (
    "testing"
    "context"
    "reflect"
    "net"

    "demo/lib/testlib"

    "demo/entities/data"
    "demo/entities/ip"
    "demo/entities/urbanarea"
    "demo/entities/location"
)

func TestData(t *testing.T, ctx context.Context, dataList []data.Model) {
    actualCount := len(dataList)
    expectedCount := 1
    if actualCount != expectedCount {
        t.Fatalf("count actual %#v != expected %#v", actualCount, expectedCount)
    }

    testField := testlib.NewTestField(testlib.TestFieldType{
        Actual: reflect.ValueOf(&dataList[0]),
        Expected: reflect.ValueOf(&data.Model{
            Domain: "reddit.com",
        }),
        T: t,
    })
    testField.Run(
        "Domain",
    )

    rDomain := "reddit.com"
    rCity := "Montreal"
    rRegionCode := "QC"
    rCountryCode := "CA"
    rLat := 45.5017
    rLon := -73.5673

    fullName := "Montreal, Quebec, Canada, H4X"

    locDataFullName := "Montréal, Quebec, Canada"
    locDataGeoNameID := int64(6077243)

    var expectedList []ip.Model
    ipstrs := []string{"151.101.1.140", "151.101.129.140", "151.101.193.140", "151.101.65.140"}
    for _, ipstr := range ipstrs {
        expectedList = append(expectedList, ip.Model{
            Address: net.ParseIP(ipstr), Domain: rDomain,
            Location: location.Model{
                City: rCity, RegionCode: rRegionCode, CountryCode: rCountryCode,
                Latitude: rLat, Longitude: rLon,
                FullName: fullName,
                TeleportFullName: locDataFullName,
                TeleportGeoNameID: locDataGeoNameID,
                UrbanArea: urbanarea.Model{
                    Name: "Montreal",
                    Link: "https://api.teleport.org/api/urban_areas/slug:montreal/",
                    Scores: []urbanarea.ScoreModel{
                        urbanarea.ScoreModel{Name: "Housing", ScoreOutOf10: 7.392000000000001},
                        urbanarea.ScoreModel{Name: "Cost of Living", ScoreOutOf10: 5.948},
                        urbanarea.ScoreModel{Name: "Startups", ScoreOutOf10: 8.102500000000001},
                        urbanarea.ScoreModel{Name: "Venture Capital", ScoreOutOf10: 5.773000000000001},
                        urbanarea.ScoreModel{Name: "Travel Connectivity", ScoreOutOf10: 3.443},
                        urbanarea.ScoreModel{Name: "Commute", ScoreOutOf10: 5.107},
                        urbanarea.ScoreModel{Name: "Business Freedom", ScoreOutOf10: 8.966},
                        urbanarea.ScoreModel{Name: "Safety", ScoreOutOf10: 7.822},
                        urbanarea.ScoreModel{Name: "Healthcare", ScoreOutOf10: 8.325666666666667},
                        urbanarea.ScoreModel{Name: "Education", ScoreOutOf10: 7.299000000000001},
                        urbanarea.ScoreModel{Name: "Environmental Quality", ScoreOutOf10: 7.7215},
                        urbanarea.ScoreModel{Name: "Economy", ScoreOutOf10: 5.8405000000000005},
                        urbanarea.ScoreModel{Name: "Taxation", ScoreOutOf10: 7.2745000000000015},
                        urbanarea.ScoreModel{Name: "Internet Access", ScoreOutOf10: 4.478000000000001},
                        urbanarea.ScoreModel{Name: "Leisure & Culture", ScoreOutOf10: 6.9174999999999995},
                        urbanarea.ScoreModel{Name: "Tolerance", ScoreOutOf10: 8.193},
                        urbanarea.ScoreModel{Name: "Outdoors", ScoreOutOf10: 5.303},
                    },
                },
            },
        })
    }
    testIPsLocationUrbanAreaUAScores(t, dataList[0].IPs, expectedList)
}
