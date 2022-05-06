package db

import (
    "strings"
    _ "embed"
    "os"
    "fmt"

    "github.com/pkg/errors"
    "gopkg.in/yaml.v2"

    "demo/lib"
)

const (
    cfgFileEnv = "APP_ENV_DB_CONFIG_FILE"
)

//go:embed postgres.yml
var defaultCfgFileData []byte

type ConfigType struct {
    Host         string `yaml:"host"`
    Name         string `yaml:"name"`
    Adapter      string `yaml:"adapter"`
    Encoding     string `yaml:"encoding"`
    MaxOpenConns uint   `yaml:"max_open_conns"`
    SearchPath   string `yaml:"search_path"`
    Pool         uint   `yaml:"pool"`
    User         string `yaml:"user"`
    Password     string `yaml:"password"`
    Port         int    `yaml:"port"`
}

func cfgData() ([]byte, error) {
    var empty []byte
    dcf := strings.TrimSpace(os.Getenv(cfgFileEnv))
    if dcf == "" {
        return defaultCfgFileData, nil
    }

    fc, err := os.ReadFile(dcf)
    if err != nil {
        return empty, errors.WithStack(err)
    }

    return fc, nil
}

func NewConfig(in InputType) (ConfigType, error) {
    empty := ConfigType{}

    temp := map[string]ConfigType{}
    cfgFileData, err := cfgData()
    if err != nil {
        return empty, errors.WithStack(err)
    }
    err = yaml.Unmarshal(cfgFileData, &temp)
    if err != nil {
        return empty, errors.WithStack(err)
    }

    dbcfg, ok := temp[in.Env]
    if !ok {
        return empty, fmt.Errorf("env %#v not supported", in.Env)
    }

    trimmed := strings.TrimSpace(in.DBName)
    if trimmed != "" {
        dbcfg.Name = trimmed
    }

    switch true {
    case strings.TrimSpace(dbcfg.Host) == "":
        return empty, fmt.Errorf("database host config value should not be empty")

    case strings.TrimSpace(dbcfg.Name) == "":
        return empty, fmt.Errorf("database name config value should not be empty")

    case strings.TrimSpace(dbcfg.Adapter) == "":
        return empty, fmt.Errorf("database adapter config value should not be empty")
    }

    return dbcfg, nil
}

func RandomDBName(env string) (string, error) {
    var empty string

    temp := map[string]ConfigType{}
    cfgFileData, err := cfgData()
    if err != nil {
        return empty, errors.WithStack(err)
    }
    err = yaml.Unmarshal(cfgFileData, &temp)
    if err != nil {
        return empty, errors.WithStack(err)
    }

    dbcfg, ok := temp[env]
    if !ok {
        return empty, errors.WithStack(fmt.Errorf("env %#v not supported", env))
    }

    rhex, err := lib.RandomHex(8)
    if err != nil {
        return empty, errors.WithStack(err)
    }
    return strings.Join([]string{dbcfg.Name, rhex}, "_"), nil
}
