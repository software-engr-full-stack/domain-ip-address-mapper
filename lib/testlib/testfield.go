package testlib

import (
    "testing"
    "reflect"
    "strings"
    "fmt"
)

type TestFieldType struct {
    Actual reflect.Value
    Expected reflect.Value
    Name string
    Label string
    T *testing.T
    Index int
}

func NewTestField(input TestFieldType) TestFieldType {
    // // TODO: test for nil if possible
    // if input.Actual == nil {
    //     panic("... ERROR: actual is nil.")
    // }

    // // TODO: test for nil if possible
    // if input.Expected == nil {
    //     panic("... ERROR: expected is nil.")
    // }

    return input
}

func (tf TestFieldType) Run(fni... string) {
    fieldNames := fni
    if len(fni) == 0 {
        trimmed := strings.TrimSpace(tf.Name)
        if trimmed == "" {
            panic("... ERROR: must give list of field names or set name field in input struct.")
        }
        fieldNames = append(fieldNames, trimmed)
    }

    fmat := "actual %#v %#v != expected %#v %#v"
    if tf.Index >= 0 {
        fmat = fmt.Sprintf("ix %#v: %s", tf.Index, fmat)
    }

    for _, fName := range fieldNames {
        label := fName
        if tf.Label != "" {
            label = tf.Label
        }

        actual := fieldByName(tf.Actual, fName)
        expected := fieldByName(tf.Expected, fName)

        if actual != expected {
            tf.T.Helper()
            tf.T.Errorf(fmat, label, actual, label, expected)
            // tf.T.Fatalf(fmat, label, actual, label, expected)
        }
    }
}

func fieldByName(value reflect.Value, fName string) interface{} {
    field := reflect.Indirect(value).FieldByName(fName)

    return field.Interface()
}
