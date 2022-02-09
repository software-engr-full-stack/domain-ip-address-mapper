package config

import (
    "os"
    "strings"
    "fmt"

    "github.com/pkg/errors"
)

const (
    appEnvKey        = "APP_ENV"
    appEnvDev        = "APP_ENV_DEVELOPMENT"
    appEnvTest       = "APP_ENV_TEST"
    appEnvProduction = "APP_ENV_PRODUCTION"
)

type EnvType struct {
    FullName          string
    BaseName          string
    LowerCaseBaseName string
}

func NewEnv() (EnvType, error) {
    envStr := ""
    var empty EnvType
    switch true {
    case IsDevEnv():
        envStr = appEnvDev
    case IsTestEnv():
        envStr = appEnvTest
    case IsProductionEnv():
        envStr = appEnvProduction
    default:
        return empty, errors.New("should be unreachable")
    }

    bn := strings.Split(envStr, "_")[2]
    return EnvType{
        FullName: envStr,
        BaseName: bn,
        LowerCaseBaseName: strings.ToLower(bn),
    }, nil
}

func TestEnvLowerCaseBaseName() string {
    bn := strings.Split(appEnvTest, "_")[2]
    return strings.ToLower(bn)
}

func CurrentEnv() string {
    return strings.TrimSpace(os.Getenv(appEnvKey))
}

func SetEnv(env string) error {
    trimmed := strings.TrimSpace(env)
    switch trimmed {
    case "", appEnvDev, appEnvTest, appEnvProduction:
        return os.Setenv(appEnvKey, trimmed)
    }

    return errors.New(fmt.Sprintf("unsupported env %#v", trimmed))
}

func SetEnvToTest() error {
    return SetEnv(appEnvTest)
}

func IsDevEnv() bool {
    env := CurrentEnv()

    if env == appEnvDev || env == "" {
        return true
    }

    return false
}

func IsTestEnv() bool {
    env := CurrentEnv()

    if env == appEnvTest {
        return true
    }

    return false
}

func IsProductionEnv() bool {
    env := CurrentEnv()

    if env == appEnvProduction {
        return true
    }

    return false
}
