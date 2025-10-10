package tools

import (
	"fmt"
	"os"
)

func init() {
	loadConfig()
	setCsrf()
}

// Danger puts a error message
func Danger(step string, args ...any) {
	fmt.Fprintf(os.Stderr, "[%s] ", step)
	fmt.Fprintln(os.Stderr, args...)
}

// Log puts a log message
func Log(step string, args ...any) {
	fmt.Printf("[%s] ", step)
	fmt.Println(args...)
}
