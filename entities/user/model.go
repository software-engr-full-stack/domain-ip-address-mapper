package user

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
    "demo/entities/data"
)

type Model struct {
    Name string       `json:"name" gorm:"unique"`
    Pw   string       `json:"pw"`

    Data []data.Model `json:"data" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`

    gorm.Model
}

// // DOC: Probably a bad idea using init to automigrate: it won't allow you to drop the database.
// // It seems like it creates an open database connection that never closes.
// func init() {
//     postgresConfig, err := pgcfg.New()
//     if err != nil {
//         panic(err)
//     }

//     connectionSpec := fmt.Sprintf(
//         "host=%s dbname=%s sslmode=disable",
//         postgresConfig.Host,
//         postgresConfig.Name,
//     )
//     db, err := gorm.Open(postgres.Open(connectionSpec), &gorm.Config{})
//     if err != nil {
//         panic(err)
//     }

//     db.AutoMigrate(&Model{})
// }

// How to set singular name for a table in gorm
// https://stackoverflow.com/questions/44589060/how-to-set-singular-name-for-a-table-in-gorm
func (Model) TableName() string {
    return "users"
}

var preloadStrings = []string{
    dbutil.PreloadStringForData(reflect.TypeOf(Model{}).PkgPath()),
    dbutil.PreloadStringForSubDomains(reflect.TypeOf(Model{}).PkgPath()),
}

var queryFields = []string{
    "Name",
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

const RootName = "root"

func Root(ctx context.Context) (Model, bool, error) {
    empty0 := Model{}
    empty1 := false

    searchUsingThisModel := Model{Name: RootName}
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

    users := []Model(found)

    if !didFind {
        return empty0, false, nil
    }

    if len(users) == 1 {
        return users[0], true, nil
    }

    return empty0, empty1, errors.New(fmt.Sprintf(
        "there should only be one %#v user, found %#v",
        RootName,
        len(users),
    ))
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
