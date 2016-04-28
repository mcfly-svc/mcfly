package main

import (
	"fmt"
	"github.com/mikec/marsupi-api/db"
	"log"
	"os/exec"
)

func RunCommands(args []string) {
	cmd := args[0]
	switch cmd {
	case "db.clean":
		db.RunHelperScript("./db/helpers/clean.sh")
	case "db.recreate":
		db.RunHelperScript("./db/helpers/recreate.sh")
	default:
		log.Fatal(fmt.Errorf("Unknown command %s", cmd))
	}
}
