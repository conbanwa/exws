package zero

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	. "github.com/rs/zerolog"
)

var Writer Logger

func init() {
	Writer = Colored(os.Stderr)
}

func Json(writer io.Writer, domain ...string) Logger {
	ctx := New(writer).With()
	if len(domain) > 0 {
		return ctx.Str("domain", domain[0]).Logger()
	}
	return ctx.Logger()
}

func Plain(writer io.Writer, formatT bool, domain ...string) Logger {
	return PlainWithT(writer, false, domain...)
}

func PlainWithT(writer io.Writer, formatT bool, domain ...string) Logger {
	cwd, err := os.Getwd()
	w := ConsoleWriter{
		Out:             writer,
		NoColor:         true,
		FormatTimestamp: func(i interface{}) string { return "" },
		FormatCaller: func(i interface{}) string {
			if c, _ := i.(string); len(c) > 0 && err == nil {
				if c, err = filepath.Rel(cwd, c); err == nil {
					return strings.Split(c, ":")[0]
				}
			}
			return ""
		},
	}
	if formatT {
		w.FormatPrepare = handler
	}
	ctx := New(w).With().Caller()
	if len(domain) > 0 {
		return ctx.Str("domain", domain[0]).Logger()
	}
	return ctx.Logger()
}

func Colored(writer io.Writer, domain ...string) Logger {
	return ColoredWithT(writer, false, domain...)
}

func ColoredWithT(writer io.Writer, formatT bool, domain ...string) Logger {
	w := ConsoleWriter{
		Out:        writer,
		TimeFormat: time.TimeOnly + ".0",
	}
	if formatT {
		w.FormatPrepare = handler
	}
	ctx := New(w).With().Timestamp().Caller()
	if len(domain) > 0 {
		return ctx.Str("domain", domain[0]).Logger()
	}
	return ctx.Logger()
}

func handler(evt map[string]interface{}) error {
	if ts, ok := evt["t"].([]interface{}); ok {
		evt["t"] = fmt.Sprint(ts) + "ttt"
	}
	return nil
}
