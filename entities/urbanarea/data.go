package urbanarea

import (
    "net/http"
    "io/ioutil"
    "encoding/json"
    "strings"

    "github.com/pkg/errors"
)

type geoDataType struct {
    Name string `json:"name"`
    FullName string `json:"full_name"`
    GeoNameID int64 `json:"geoname_id"`
    Population int64 `json:"population"`

    UrbanAreaLink string
    UrbanAreaName string
    GeoHash string
    Latitude float64
    Longitude float64
}

func geoData(geoNameIDURL string) (geoDataType, error) {
    var empty geoDataType

    req, err := http.NewRequest(http.MethodGet, geoNameIDURL, nil)
    if err != nil {
        return empty, errors.WithStack(err)
    }

    req.Header.Set("User-Agent", "Mozilla/4.0 (compatible; MSIE 6.0; Windows 98)")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return empty, errors.WithStack(err)
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return empty, errors.WithStack(err)
    }

    var blob map[string]interface{}
    if err := json.Unmarshal(body, &blob); err != nil {
        return empty, errors.WithStack(err)
    }

    var result geoDataType
    if err := json.Unmarshal(body, &result); err != nil {
        return empty, errors.WithStack(err)
    }

    location := blob["location"].(map[string]interface{})
    latlon := location["latlon"].(map[string]interface{})

    cua := blob["_links"].(map[string]interface{})["city:urban_area"]

    result.GeoHash = location["geohash"].(string)
    result.Latitude = latlon["latitude"].(float64)
    result.Longitude = latlon["longitude"].(float64)

    if cua != nil {
        cuaif := cua.(map[string]interface{})
        result.UrbanAreaLink = cuaif["href"].(string)
        result.UrbanAreaName = cuaif["name"].(string)
    }

    return result, nil
}

type uaType struct {
    UrbanAreaScores []ScoreModel
    UrbanAreaSummary string
}

type scoreType struct {
    Name string
    ScoreOutOf10 float64
    Color string
}

func urbanArea(cuaLink string) (uaType, error) {
    var empty uaType

    scrLink := strings.Join([]string{cuaLink, "scores"}, "")
    req, err := http.NewRequest(http.MethodGet, scrLink, nil)
    if err != nil {
        return empty, errors.WithStack(err)
    }

    req.Header.Set("User-Agent", "Mozilla/4.0 (compatible; MSIE 6.0; Windows 98)")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return empty, errors.WithStack(err)
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return empty, errors.WithStack(err)
    }

    var blob map[string]interface{}
    if err := json.Unmarshal(body, &blob); err != nil {
        return empty, errors.WithStack(err)
    }

    var scores []ScoreModel
    for _, catif := range blob["categories"].([]interface{}) {
        cat := catif.(map[string]interface{})
        scores = append(scores, ScoreModel{
            Name: cat["name"].(string),
            ScoreOutOf10: cat["score_out_of_10"].(float64),
            Color: cat["color"].(string),
        })
    }

    return uaType{
        UrbanAreaSummary: blob["summary"].(string),
        UrbanAreaScores: scores,
    }, nil
}
