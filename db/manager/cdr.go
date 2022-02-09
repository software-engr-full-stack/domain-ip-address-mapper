package manager

import (
    "context"
    "database/sql"
    _ "github.com/lib/pq"
    "strings"
    "fmt"

    "github.com/pkg/errors"

    "demo/config"
    "demo/db"
)

type managerType struct {
    targetDatabaseConfig db.ConfigType
    adapter string
    connectionSpec string
}

func newManager(targetDatabaseConfig db.ConfigType) (*managerType) {
    connectionSpec := []string{
        fmt.Sprintf("host=%s", targetDatabaseConfig.Host),
        fmt.Sprintf("dbname=%s", "postgres"),
        "sslmode=disable",
    }
    if targetDatabaseConfig.User != "" {
        connectionSpec = append(connectionSpec, fmt.Sprintf("user=%s", targetDatabaseConfig.User))
    }
    if targetDatabaseConfig.Password != "" {
        connectionSpec = append(connectionSpec, fmt.Sprintf("password=%s", targetDatabaseConfig.Password))
    }
    if targetDatabaseConfig.Port != 0 {
        connectionSpec = append(connectionSpec, fmt.Sprintf("port=%d", targetDatabaseConfig.Port))
    }

    return &managerType{
        targetDatabaseConfig: targetDatabaseConfig,
        adapter: targetDatabaseConfig.Adapter,
        connectionSpec: strings.Join(connectionSpec, " "),
    }
}

func CreateDatabase(ctx context.Context) error {
    dbobj := ctx.Value("config").(config.Type).DBObj
    mngr := newManager(dbobj.Config)
    dbOp := "CREATE"
    dbName := mngr.targetDatabaseConfig.Name
    dbExists, err := mngr.targetDatabaseExists()
    if err != nil {
        return errors.WithStack(err)
    }
    if dbExists {
        fmt.Printf("... database %#v already exists, skipping %#v\n", dbName, dbOp)
        return nil
    }
    fmt.Printf("... %#v %#v...\n", dbOp, dbName)
    return mngr.manage(dbOp)
}

func DropDatabase(ctx context.Context) error {
    dbobj := ctx.Value("config").(config.Type).DBObj
    mngr := newManager(dbobj.Config)
    dbOp := "DROP"
    dbName := mngr.targetDatabaseConfig.Name
    dbExists, err := mngr.targetDatabaseExists()
    if err != nil {
        return errors.WithStack(err)
    }
    if !dbExists {
        fmt.Printf("... database %#v does not exist, skipping %#v\n", dbName, dbOp)
        return nil
    }

    fmt.Printf("... %#v %#v...\n", dbOp, dbName)
    return mngr.manage(dbOp)
}

func ResetDatabase(ctx context.Context) error {
    err := DropDatabase(ctx)
    if err != nil {
        return errors.WithStack(err)
    }
    return CreateDatabase(ctx)
}

func (mngr *managerType) manage(operation string) error {
    instance, err := sql.Open(mngr.adapter, mngr.connectionSpec)
    if err != nil {
        return errors.WithStack(err)
    }

    _, err = instance.Exec(fmt.Sprintf("%s DATABASE %s", operation, mngr.targetDatabaseConfig.Name))
    if err != nil {
        return errors.WithStack(err)
    }

    return nil
}

func (mngr *managerType) targetDatabaseExists() (bool, error) {
    empty := false
    instance, err := sql.Open(mngr.adapter, mngr.connectionSpec)
    if err != nil {
        return empty, errors.WithStack(err)
    }

    var exists bool
    query := fmt.Sprintf(
        "SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname='%s')",
        mngr.targetDatabaseConfig.Name,
    )
    err = instance.QueryRow(query).Scan(&exists)
    if err != nil && err != sql.ErrNoRows {
        return empty, errors.WithStack(err)
    }

    return exists, nil
}
