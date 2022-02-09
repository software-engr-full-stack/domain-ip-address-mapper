package ip

import (
    "context"
    "net"
    "reflect"

    "github.com/pkg/errors"
    "gorm.io/gorm"

    "demo/db/query"
    dbutil "demo/db/util"
    "demo/entities/location"
)

type Model struct {
    Address      net.IP `json:"address" gorm:"-"`
    AddressBytes []byte `json:"address_bytes"`

    ISP string `json:"isp"`
    Org string `json:"org"`
    AS string `json:"as"` // Autonomous System Number https://afrinic.net/asn

    Domain   string `json:"domain"`

    // Belongs to:
    LocationID  uint `json:"location_id" gorm:"default:NULL"`
    Location location.Model `json:"location"`

    ParentDataID      uint `json:"parent_data_id" gorm:"default:NULL"`
    ParentSubDomainID uint `json:"parent_domain_id" gorm:"default:NULL"`

    gorm.Model
}

func (model *Model) BeforeCreate(tx *gorm.DB) (error) {
    if len(model.Address) == 0 {
        return errors.New("can't save invalid data")
    }
    model.AddressBytes = []byte(model.Address)
    return nil
}

func (model *Model) AfterFind(tx *gorm.DB) (error) {
    model.Address = net.IP(model.AddressBytes)
    return nil
}

// How to set singular name for a table in gorm
// https://stackoverflow.com/questions/44589060/how-to-set-singular-name-for-a-table-in-gorm
func (Model) TableName() string {
    return "ips"
}

var preloadStrings = []string{dbutil.PreloadString(reflect.TypeOf(Model{}).PkgPath())}

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
