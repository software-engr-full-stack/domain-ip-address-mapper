package urbanarea

import (
    "context"
    "reflect"
    "strings"
    "fmt"

    "github.com/pkg/errors"
    "gorm.io/gorm"
    "github.com/iancoleman/strcase"

    "demo/db/query"
    dbutil "demo/db/util"
)

type Model struct {
    Name string `json:"urban_area_name" gorm:"unique"`
    Link string `json:"urban_area_link" gorm:"unique"`

    Scores []ScoreModel `json:"urban_area_scores" gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE;"`
    Summary string `json:"urban_area_summary"`

    gorm.Model
}

type ScoreModel struct {
    Name string `json:"name"`
    ScoreOutOf10 float64 `json:"score_out_of_10"`
    Color string `json:"color"`

    ParentID uint `json:"parent_id"`

    gorm.Model
}

func (Model) TableName() string {
    return "urban_areas"
}

func (ScoreModel) TableName() string {
    return "urban_area_scores"
}

var preloadStrings = []string{dbutil.PreloadString(reflect.TypeOf(Model{}).PkgPath())}

func (model Model) IsEmpty() (bool) {
    fieldNames := []string{
        "Name",
        "Link",
        "Summary",
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
    "Name",
    "Link",
    "Summary",
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
