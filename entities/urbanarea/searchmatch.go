package urbanarea

import (
    "net/http"
    "net/url"
    "io/ioutil"
    "encoding/json"
    "strings"
    "regexp"

    "github.com/pkg/errors"
)

type citySearchResultType struct {
    GeoNameIDURL string
    MatchingFullName string
}

func search(searchTerm string) ([]citySearchResultType, error) {
    var empty []citySearchResultType
    if strings.TrimSpace(searchTerm) == "" {
        return empty, errors.New("search term must not be empty")
    }
    link := url.URL{
        Scheme: "https",
        Host:   "api.teleport.org",
        Path: "/api/cities/",
    }

    query := link.Query()
    query.Set("search", searchTerm)
    link.RawQuery = query.Encode()

    req, err := http.NewRequest(http.MethodGet, link.String(), nil)
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

    var t0 map[string]interface{}
    if err := json.Unmarshal(body, &t0); err != nil {
        return empty, errors.WithStack(err)
    }

    t1 := t0["_embedded"].(map[string]interface{})["city:search-results"].
                          ([]interface{})

    var results []citySearchResultType
    for _, item := range t1 {
        i0 := item.(map[string]interface{})
        i1 := i0["_links"].(map[string]interface{})["city:item"].
                           (map[string]interface{})["href"]

        results = append(results, citySearchResultType{
            MatchingFullName: i0["matching_full_name"].(string),
            GeoNameIDURL: strings.TrimSpace(i1.(string)),
        })
    }

    return results, nil
}

type matchType struct {
    IsFound bool
    GeoNameIDURL string
    FullName string
    TeleportFullName string
}

func match(
    city string, region string, country string,
    citiesFoundByAPI []citySearchResultType,
) (matchType, error) {
    var empty matchType
    var found string
    var fullName string
    var apiFullName string
    for _, result := range citiesFoundByAPI {
        matchToCityRegionCountry := strings.Join([]string{city, region, country}, ", ")

        noparen, err := removeParen(result.MatchingFullName)
        if err != nil {
            return empty, errors.WithStack(err)
        }
        normSearchResultFromAPI := normalize(noparen)
        if normalize(matchToCityRegionCountry) == normSearchResultFromAPI {
            found = result.GeoNameIDURL
        }

        if found != "" {
            apiFullName = result.MatchingFullName
            fullName = matchToCityRegionCountry
            break
        }

        // TODO: Singapore hack (city, country)
        if region == "" {
            matchToCityCountry := strings.Join([]string{city, country}, ", ")
            noRegion := removeRegion(normSearchResultFromAPI)
            if normalize(matchToCityCountry) == noRegion {
                found = result.GeoNameIDURL
            }
            if found != "" {
                apiFullName = result.MatchingFullName
                fullName = matchToCityCountry
                break
            }
        }
    }

    if found == "" {
        return matchType{IsFound: false}, nil
    }

    return matchType{
        IsFound: true, GeoNameIDURL: found, FullName: fullName, TeleportFullName: apiFullName,
    }, nil
}

// ... WHY? Because some results look like this "Montréal, Quebec, Canada (Montreal)" instead of just
//   "Montréal, Quebec, Canada"
func removeParen(str string) (string, error) {
    fields := strings.Fields(str)
    length := len(fields)

    matched, err := regexp.MatchString(`\(.*?\)`, fields[length - 1])
    var empty string
    if err != nil {
        return empty, errors.WithStack(err)
    }

    if matched {
        return strings.Join(fields[:length - 1], ""), nil
    }

    return strings.Join(fields, ""), nil
}

func normalize(str string) string {
    type replaceType struct {
        Original string
        Replacement string
    }
    replace := []replaceType{
        {"é", "e"},
    }
    for _, rep := range replace {
        str = strings.ReplaceAll(str, rep.Original, rep.Replacement)
    }
    return strings.ToLower(strings.TrimSpace(strings.Join(strings.Fields(str), "")))
}

func removeRegion(location string) string {
    words := strings.Split(location, ",")
    return strings.Join([]string{words[0], words[len(words) - 1]}, ",")
}
