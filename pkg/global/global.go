package global

import "fmt"

var IsVerbose bool

const CLI_Border = "--------------------------------------------------\n"

func PrintVerboseInfo(format string, args ...interface{}) {
	if IsVerbose {
		fmt.Printf(format, args...)
	}
}
