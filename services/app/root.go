package app

import (
    "context"

    "github.com/pkg/errors"
    "gorm.io/gorm"

    "demo/entities/user"
)

func createRoot(ctx context.Context, tx *gorm.DB) (user.Model, error) {
    root, isFound, err := user.Root(ctx)
    var empty user.Model
    if err != nil {
        return empty, errors.WithStack(err)
    }
    if isFound {
        return root, nil
    }

    root = user.Model{Name: "root"}
    if err := tx.Create(&root).Error; err != nil {
        return empty, errors.WithStack(err)
    }

    return root, nil
}
