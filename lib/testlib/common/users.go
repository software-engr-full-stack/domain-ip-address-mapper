package common

import (
    "testing"
    "context"

    "demo/entities/user"
)

const gExpectedRootName = "root"

func TestUsers(t *testing.T, ctx context.Context, expectedAllUsers []user.Model) {
    if actual, expected := len(expectedAllUsers), 1; actual != expected {
        t.Fatalf("actual %#v != expected %#v", actual, expected)
    }
    usr := expectedAllUsers[0]
    if actual, expected := usr.Name, gExpectedRootName; actual != expected {
        t.Fatalf("actual %#v != expected %#v", actual, expected)
    }
}

func TestRootUser(t *testing.T, ctx context.Context, rootUser user.Model, isFound bool) {
    if !isFound {
        t.Fatalf("%#v user must be found", gExpectedRootName)
    }

    if actual, expected := rootUser.Name, gExpectedRootName; actual != expected {
        t.Fatalf("actual %#v != expected %#v", actual, expected)
    }
}
