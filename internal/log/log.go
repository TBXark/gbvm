package log

import (
	"fmt"
	"io"
	"os"
)

var (
	debugWriter io.Writer = os.Stdout
	errorWriter io.Writer = os.Stderr
)

func Debug(args ...any) {
	_, _ = fmt.Fprintln(debugWriter, args...)
}

func Debugf(format string, args ...any) {
	_, _ = fmt.Fprintf(debugWriter, format+"\n", args...)
}

func Error(args ...any) {
	_, _ = fmt.Fprintln(errorWriter, args...)
}

func Errorf(format string, args ...any) {
	_, _ = fmt.Fprintf(errorWriter, format+"\n", args...)
}
