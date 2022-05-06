package secret

import (
    "os"

    "github.com/pkg/errors"
    "gopkg.in/yaml.v2"

    "demo/config"
)

type Type struct {
    FreeGeoIPAppKey string `yaml:"free_geo_ip_app_key"`
}

func New() (Type, error) {
    var empty Type

    env, err := config.New(config.Type{})
    if err != nil {
        return empty, errors.WithStack(err)
    }

    fileData, err := os.ReadFile(env.SecretsFilePath)
    if err != nil {
        return empty, errors.WithStack(err)
    }

    sec := Type{}
    err = yaml.Unmarshal(fileData, &sec)
    if err != nil {
        return empty, errors.WithStack(err)
    }

    return sec, nil
}
