# go-log

Opinionated Go package for doing minimal structured logging and prefixing of log messages with Emoji for easier filtering. It's possible this package will become irrelevant if and when Go [slog](https://github.com/golang/go/issues/56345) package because part of "core". Until then it does what I need.

## Documentation

[![Go Reference](https://pkg.go.dev/badge/github.com/aaronland/go-log.svg)](https://pkg.go.dev/github.com/aaronland/go-log)

## Example

```
import (
	"log"

	aa_log "github.com/aaronland/go-log/v2"
)	

func main(){

	logger := log.Default()

	aa_log.SetMinLevelWithPrefix(aa_log.WARNING_PREFIX)

	// No output
	aa_log.Debug(logger, "This is a test")

	aa_log.UnsetMinLevel()

	// prints "💬 Hello, world"
	aa_log.Info(logger, "Hello, %w", "world")
	
	// prints "🪵 This is a second test"
	aa_log.Debug(logger, "This is a second test")

	// prints "🧯 This is an error"
	aa_log.Warning(logger, fmt.Errorf("This is an error"))

	// Emits errors using the default Go *log.Logger instance
	// prints "{YYYY}/{MM}/{DD} {HH}:{MM}:{SS} 🧯 This is a second error"
	aa_log.Warning(fmt.Errorf("This is a second error"))
}
```

## Prefixes

| Log level | Prefix |
| --- | --- |
| debug | 🪵 |
| info | 💬 |
| warning | 🧯 |
| error | 🔥 |
| fatal | 💥 |
