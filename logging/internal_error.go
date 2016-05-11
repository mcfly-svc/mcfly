package logging

import (
	"fmt"
	"runtime/debug"
)

func InternalError(err error) {
	//TODO: log to Sentry or similar service. if we end up here, that's bad!
	fmt.Println("[MSPLAPI ERROR]:", err)
	debug.PrintStack()
}
