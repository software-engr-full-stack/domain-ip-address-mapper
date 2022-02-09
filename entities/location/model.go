package location

import (
    "context"
    "reflect"
    "strings"

    "github.com/pkg/errors"
    "gorm.io/gorm"
    "github.com/iancoleman/strcase"

    "fmt"

    "demo/db/query"
    dbutil "demo/db/util"
    "demo/entities/urbanarea"
)

type Model struct {
    CountryCode string `json:"country_code"`
    Country string `json:"country" gorm:"NOT NULL"`
    RegionCode string `json:"region_code"`
    Region string `json:"region"`
    City string `json:"city" gorm:"NOT NULL"`
    ZipCode string `json:"zip_code"`
    TimeZone string `json:"time_zone"`
    Latitude float64 `json:"latitude"`
    Longitude float64 `json:"longitude"`

    // ...
    Error        error  `json:"error" gorm:"-"`
    ErrorMessage string `json:"error_message"`

    FullName string `json:"full_name" gorm:"unique"`
    ViewFullName string `json:"view_full_name"`

    // ...
    Population int64 `json:"population"`
    TeleportGeoNameID int64 `json:"teleport_geo_name_id"`
    TeleportGeoHash string `json:"teleport_geo_hash"`
    TeleportName string `json:"teleport_name"`
    TeleportFullName string `json:"teleport_full_name"`
    TeleportLatitude float64 `json:"teleport_latitude"`
    TeleportLongitude float64 `json:"teleport_longitude"`

    // Belongs to:
    // TODO: on both UA fields, setting tag gorm:"default:NULL" seems to have no effect so it's omitted.
    UrbanAreaID uint `json:"urban_area_id" gorm:"default:NULL"`
    UrbanArea urbanarea.Model `json:"urban_area"`

    gorm.Model
}

func (Model) TableName() string {
    return "locations"
}

func (model *Model) BeforeCreate(tx *gorm.DB) (error) {
    if model.Error != nil {
        model.ErrorMessage = model.Error.Error()
    }
    return nil
}

func (model *Model) AfterFind(tx *gorm.DB) (error) {
    if model.ErrorMessage != "" {
        model.Error = fmt.Errorf("... ERROR: from database %#v", model.ErrorMessage)
    }
    return nil
}

var preloadStrings = []string{dbutil.PreloadString(reflect.TypeOf(Model{}).PkgPath())}

func (model Model) IsEmpty() (bool) {
    fieldNames := []string{
        "City",
        "Region",
        "Country",
    }

    for _, fName := range fieldNames {
        rval := reflect.ValueOf(&model)
        field := reflect.Indirect(rval).FieldByName(fName)

        if strings.TrimSpace(field.String()) != "" {
            return false
        }
    }

    return true
}

var queryFields = []string{
    "CountryCode",
    "Country",
    "RegionCode",
    "Region",
    "City",
    "ZipCode",
    "TimeZone",
    "FullName",
    "ViewFullName",
    "TeleportGeoNameID",
    "TeleportGeoHash",
    "TeleportName",
    "TeleportFullName",
    "UrbanAreaID",
}
func (model Model) ToMap() (map[string]interface{}, error) {
    var empty map[string]interface{}
    mp := map[string]interface{}{}
    for _, qf := range queryFields {
        value := reflect.ValueOf(&model)
        field := reflect.Indirect(value).FieldByName(qf)

        ifv := field.Interface()
        switch typ := ifv.(type) {
        case string:
            trimmed := strings.TrimSpace(ifv.(string))
            if trimmed == "" {
                continue
            }
        case int, uint, int64, float64:
            // TODO: DEBUG, I don't know about this if this is a good practice in general
            //   regarding struct default values
            if ifv == 0 || ifv == uint(0) || ifv == int64(0) || ifv == float64(0) {
                continue
            }
        default:
            return empty, errors.New(fmt.Sprintf(
                "unsupported type %#v for field %#v in %#v",
                typ, qf, model,
            ))
        }
        mp[strcase.ToSnake(qf)] = ifv
    }

    return mp, nil
}

type ModelsType []Model
func (mt *ModelsType) Count() int {
    return len([]Model(*mt))
}

func All(ctx context.Context) ([]Model, error) {
    var empty []Model
    var models ModelsType
    err := query.All(query.InputType{
        Context: ctx,
        PreloadStrings: preloadStrings,
        StoreResultsHere: &models,
    })
    if err != nil {
        return empty, errors.WithStack(err)
    }

    return []Model(models), nil
}

func Where(ctx context.Context, searchUsingThisModel Model) (_found []Model, _didFind bool, _err error) {
    var empty0 []Model
    var empty1 bool
    if searchUsingThisModel.IsEmpty() {
        return empty0, empty1, errors.New("model must not be empty")
    }

    mp, err := searchUsingThisModel.ToMap()
    if err != nil {
        return empty0, empty1, errors.WithStack(err)
    }
    var found ModelsType
    didFind, err := query.Where(mp, query.InputType{
        Context: ctx,
        PreloadStrings: preloadStrings,
        StoreResultsHere: &found,
    })
    if err != nil {
        return empty0, empty1, errors.WithStack(err)
    }

    return []Model(found), didFind, nil
}

// TODO: write test, not covered
func DeleteAll(ctx context.Context) (count int, deferr error) {
    empty := -1
    var models ModelsType
    count, err := query.DeleteAll(query.InputType{
        Context: ctx,
        PreloadStrings: preloadStrings,
        StoreResultsHere: &models,
    })
    if err != nil {
        return empty, errors.WithStack(err)
    }

    return count, nil
}

type UniqueType struct {
    AreRecordsFound bool
    IsFound bool
    Records []Model
    Record Model
}

func (model Model) Unique(ctx context.Context) (UniqueType, error) {
    found, didFind, err := Where(ctx, model)
    var empty UniqueType
    if err != nil {
        return empty, errors.WithStack(err)
    }

    if !didFind {
        return UniqueType{AreRecordsFound: false}, nil
    }

    if len(found) > 1 {
        return UniqueType{AreRecordsFound: true, IsFound: false, Records: found}, nil
    }

    return UniqueType{AreRecordsFound: true, IsFound: true, Records: found, Record: found[0]}, nil
}
