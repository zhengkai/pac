// Package log 把 fmt 包装简化
package log

import (
	"fmt"
	"os"
)

func J(a ...any) {
	fmt.Fprintln(os.Stdout, a...)
}

func F(format string, a ...any) {
	fmt.Fprintf(os.Stdout, format, a...)
	fmt.Fprint(os.Stdout, "\n")
}

func W(a ...any) {
	fmt.Fprintln(os.Stderr, a...)
}

func WF(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format, a...)
	fmt.Fprint(os.Stderr, "\n")
}
