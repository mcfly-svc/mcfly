package logging

import (
    "log"
)

func LogInternalError(method string, err error) {
	log.Printf("%s: %v", method, err)
}

func LogFatal(err error) {
	log.Fatal(err)
}