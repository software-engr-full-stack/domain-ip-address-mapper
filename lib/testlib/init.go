package testlib

import (
    "os"
)

func init() {
    os.Setenv("APP_ENV", "APP_ENV_TEST")
}
