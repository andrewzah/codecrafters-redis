package main

import (
	"fmt"
	"os"
)

func LogAtLevel(level, format string, a ...interface{}) {
	msg := fmt.Sprintf("[%s] %s\n", level, format)
	fmt.Fprintf(os.Stdout, msg, a...)
}

func Fatalf(format string, a ...interface{})   { LogAtLevel("FTL", format, a...) }
func Warningf(format string, a ...interface{}) { LogAtLevel("WRN", format, a...) }
func Debugf(format string, a ...interface{})   { LogAtLevel("DBG", format, a...) }
func Errorf(format string, a ...interface{})   { LogAtLevel("ERR", format, a...) }
func Infof(format string, a ...interface{})    { LogAtLevel("INF", format, a...) }
