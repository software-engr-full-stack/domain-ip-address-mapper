package testlib

import (
    "os"
    "time"
    "database/sql"
    "github.com/lib/pq"

    "fmt"

    "gorm.io/gorm"
    "gorm.io/driver/postgres"

    "demo/entities/user"
    "demo/entities/data"
    "demo/entities/subdomain"
    "demo/entities/urbanarea"
    "demo/entities/location"
    "demo/entities/ip"

    "demo/db"
)

type TestDbType struct {
    Name string

    adapter string
    host string
    encoding string
}

func NewTestDbType() TestDbType {
    randName, err := db.RandomDBName("test")
    if err != nil {
        panic(err)
    }

    return TestDbType{
        Name: randName,

        adapter: "postgres",
        host: "/var/run/postgresql/",
        encoding: "utf-8",
    }
}

func (tdt *TestDbType) Setup() {
    appEnv := os.Getenv("APP_ENV")
    expectedAppEnv := "APP_ENV_TEST"
    if appEnv != expectedAppEnv {
        panic(fmt.Errorf("... ERROR: app env != %#v", expectedAppEnv))
    }

    tdt.dropDatabasePostgres()
    tdt.createDatabasePostgres()
    tdt.migrate()
}

func (tdt *TestDbType) CleanUp() {
    tdt.dropDatabasePostgres()
}

func (tdt *TestDbType) createDatabasePostgres() {
    tdt.manage("CREATE", "")
}

func (tdt *TestDbType) dropDatabasePostgres() {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()

    timeout := time.Duration(10)
    timeoutExceeded := time.After(timeout * time.Second)
    ix := 0
    for {
        select {
        case <-timeoutExceeded:
            panic(fmt.Errorf("... ERROR: failed to drop database after %#v timeout", timeout))

        case <-ticker.C:
            err := tdt._manage("DROP", "IF EXISTS")
            ix += 1
            if err != nil {
                // fmt.Printf("... DEBUG: trying to drop database try number %#v, error %#v\n", ix, err)
                if pqerr, ok := err.(*pq.Error); ok {
                    if pqerr.Code == "55006" {
                        tdt.closeAllConnections()
                    }
                }

                continue
            }
            return
        }
    }
}

func (tdt *TestDbType) manage(operation string, cond string) {
    err := tdt._manage(operation, cond)
    if err != nil {
        panic(err)
    }
}

func (tdt *TestDbType) _manage(operation string, cond string) error {
    connectionSpec := fmt.Sprintf(
        "host=%s dbname=%s sslmode=disable",
        tdt.host,
        "postgres",
    )
    instance, err := sql.Open(tdt.adapter, connectionSpec)
    if err != nil {
        return err
    }

    _, err = instance.Exec(fmt.Sprintf(
        "%s DATABASE %s %s",
        operation,
        cond,
        tdt.Name,
    ))
    if err != nil {
        return err
    }

    return nil
}

func (tdt *TestDbType) migrate() {
    connectionSpec := fmt.Sprintf(
        "host=%s dbname=%s sslmode=disable",
        tdt.host,
        tdt.Name,
    )
    dbobj, err := gorm.Open(postgres.Open(connectionSpec), &gorm.Config{})
    if err != nil {
        panic(err)
    }
    defer func() {
        sqldb, err := dbobj.DB()
        if err != nil {
            panic(err)
        }
        err = sqldb.Close()
        if err != nil {
            panic(err)
        }
    }()

    models := []interface{}{
        &user.Model{},
        &data.Model{},
        &subdomain.Model{},
        &urbanarea.Model{},
        &urbanarea.ScoreModel{},
        &location.Model{},
        &ip.Model{},
    }

    for _, mdl := range models {
        err := dbobj.Migrator().CreateTable(mdl)
        if err != nil {
            panic(err)
        }
    }
}

func (tdt *TestDbType) closeAllConnections() {
    connectionSpec := fmt.Sprintf(
        "host=%s dbname=%s sslmode=disable",
        tdt.host,
        "postgres",
    )
    instance, err := sql.Open(tdt.adapter, connectionSpec)
    if err != nil {
        panic(err)
    }

    _, err = instance.Exec(fmt.Sprintf(
        "select pg_terminate_backend(pid) from pg_stat_activity where datname='%s';",
        tdt.Name,
    ))
    if err != nil {
        panic(err)
    }
}
