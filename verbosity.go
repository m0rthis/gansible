package gansible

import (
	"fmt"
	"os"
)

var Verbosity int

func init() {
	if v := os.Getenv("ANSIBLE_VERBOSITY"); v != "" {
		fmt.Sscanf(v, "%d", &Verbosity)
	}
}

func V(msg string, args ...any) {
	if Verbosity >= 1 {
		fmt.Fprintf(os.Stderr, "[v] "+msg+"\n", args...)
	}
}

func VV(msg string, args ...any) {
	if Verbosity >= 2 {
		fmt.Fprintf(os.Stderr, "[vv] "+msg+"\n", args...)
	}
}

func VVV(msg string, args ...any) {
	if Verbosity >= 3 {
		fmt.Fprintf(os.Stderr, "[vvv] "+msg+"\n", args...)
	}
}

func VVVV(msg string, args ...any) {
	if Verbosity >= 4 {
		fmt.Fprintf(os.Stderr, "[vvvv] "+msg+"\n", args...)
	}
}

func VVVVV(msg string, args ...any) {
	if Verbosity >= 5 {
		fmt.Fprintf(os.Stderr, "[vvvvv] "+msg+"\n", args...)
	}
}
