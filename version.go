package main

import (
	"fmt"
	"runtime"
)

var (
	builddatetime string
	revision      string
	modified      string
)

func version() string {
	ret := ""
	ret += "// by iloy today\n"
	ret += fmt.Sprintln("Built:", builddatetime)
	ret += fmt.Sprintln("Go:", runtime.Version())
	ret += fmt.Sprintln("GOARCH:", runtime.GOARCH)
	ret += fmt.Sprintln("GOOS:", runtime.GOOS)
	ret += fmt.Sprintf("Revision: %s %s", revision, modified)

	return ret
}
