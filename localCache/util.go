package main

import (
    "fmt"
)

func Printf(format string, a ...any) (n int, err error) {
    return fmt.Printf(fmt.Sprintf("%s\n", format), a...)
}
