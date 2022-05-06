package query

import (
    "context"

    "github.com/pkg/errors"

    dbutil "demo/db/util"
)

type ResultsStorageType interface {
    Count() int
}

type InputType struct {
    Context context.Context
    PreloadStrings []string
    StoreResultsHere ResultsStorageType
}

func Where(searchUsingThisMap map[string]interface{}, in InputType) (_ bool, deferr error) {
    var empty bool
    db, dbobj, err := dbutil.FindOpenConnection(in.Context)
    if err != nil {
        return empty, errors.WithStack(err)
    }
    if db == nil {
        db, err = dbobj.Open()
        if err != nil {
            return empty, errors.WithStack(err)
        }
        defer func() {
            if tempErr := dbobj.Close(); tempErr != nil {
                deferr = errors.WithStack(tempErr)
            }
        }()
    }

    for _, pstr := range in.PreloadStrings {
        db = db.Preload(pstr)
    }

    if err = db.Where(searchUsingThisMap).Find(in.StoreResultsHere).Error; err != nil {
        return empty, errors.WithStack(err)
    }

    if in.StoreResultsHere.Count() > 0 {
        return true, deferr
    }

    return false, deferr
}

func All(in InputType) (deferr error) {
    db, dbobj, err := dbutil.FindOpenConnection(in.Context)
    if err != nil {
        return errors.WithStack(err)
    }
    if db == nil {
        db, err = dbobj.Open()
        if err != nil {
            return errors.WithStack(err)
        }
        defer func() {
            if tempErr := dbobj.Close(); tempErr != nil {
                deferr = errors.WithStack(tempErr)
            }
        }()
    }

    for _, pstr := range in.PreloadStrings {
        db = db.Preload(pstr)
    }

    if err = db.Find(in.StoreResultsHere).Error; err != nil {
        return errors.WithStack(err)
    }

    return deferr
}

func DeleteAll(in InputType) (count int, deferr error) {
    empty := -1
    db, dbobj, err := dbutil.FindOpenConnection(in.Context)
    if err != nil {
        return empty, errors.WithStack(err)
    }
    if db == nil {
        db, err = dbobj.Open()
        if err != nil {
            return empty, errors.WithStack(err)
        }
        defer func() {
            if tempErr := dbobj.Close(); tempErr != nil {
                deferr = errors.WithStack(tempErr)
            }
        }()
    }

    err = All(in)
    if err != nil {
        return empty, errors.WithStack(err)
    }

    toDeleteCount := in.StoreResultsHere.Count()
    if toDeleteCount == 0 {
        return 0, deferr
    }

    // // Soft delete or...
    // result := db.Delete(models)
    // Permanent delete
    if err = db.Unscoped().Delete(in.StoreResultsHere).Error; err != nil {
        return empty, errors.WithStack(err)
    }

    return toDeleteCount, deferr
}
