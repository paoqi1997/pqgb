package main

import (
    "github.com/paoqi1997/pqgb/util"
)

func Printf(format string, a ...any) (n int, err error) {
    return util.Printf(format, a...)
}
