package tools

import (
	"fmt"
	"os"
)

func init() {
	afterInit()
}

func afterInit() {
	loadConfig()
}

// Danger puts a error message
func Danger(step string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "[%s] ", step)
	fmt.Fprintln(os.Stderr, args...)
}

// Log puts a log message
func Log(step string, args ...interface{}) {
	fmt.Printf("[%s] ", step)
	fmt.Println(args...)
}
