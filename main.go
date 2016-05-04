package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/chrismrivera/cmd"
	"github.com/mikec/msplapi/api"
	"github.com/mikec/msplapi/db"
	"github.com/mikec/msplapi/logging"
	"github.com/mikec/msplapi/provider"

	_ "github.com/mattes/migrate/driver/postgres"
)

var cmdr *cmd.App = cmd.NewApp()
var dbUrl = "postgres://localhost:5432/marsupi_test?sslmode=disable"

func main() {

	cmdr.AddCommand(runServerCmd)
	cmdr.AddCommand(dbMigrateCmd)
	cmdr.AddCommand(dbCreateCmd)
	cmdr.AddCommand(dbCleanCmd)
	cmdr.AddCommand(dbSeedCmd)

	cmdr.Description = "Marsupi CLI"
	if err := cmdr.Run(os.Args); err != nil {
		if ue, ok := err.(*cmd.UsageErr); ok {
			ue.ShowUsage()
		} else {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err.Error())
		}
		os.Exit(1)
	}

}

var runServerCmd = cmd.NewCommand(
	"run", "Server", "Runs the msplapi server",
	func(cmd *cmd.Command) {},
	func(cmd *cmd.Command) error {
		RunServer()
		return nil
	},
)

var dbMigrateCmd = cmd.NewCommand(
	"migrate", "Database", "Run migration scripts",
	func(cmd *cmd.Command) {
		cmd.AppendArg("direction", "(up|down)")
	},
	func(cmd *cmd.Command) error {
		db.RunMigrate(dbUrl, cmd.Arg("direction"))
		return nil
	},
)

var dbCreateCmd = cmd.NewCommand(
	"create-db", "Database", "Creates the database",
	func(cmd *cmd.Command) {},
	func(cmd *cmd.Command) error {
		db.Create(dbUrl)
		return nil
	},
)

var dbCleanCmd = cmd.NewCommand(
	"clean-db", "Database", "Removes all data from the database",
	func(cmd *cmd.Command) {},
	func(cmd *cmd.Command) error {
		db.Clean(dbUrl)
		return nil
	},
)

var dbSeedCmd = cmd.NewCommand(
	"seed-db", "Database", "Adds seed data to the database",
	func(cmd *cmd.Command) {},
	func(cmd *cmd.Command) error {
		db.Seed(dbUrl)
		return nil
	},
)

// RunServer runs the HTTP server
func RunServer() {
	github := provider.GitHub{}
	dropbox := provider.Dropbox{}

	authProviders := make(map[string]provider.AuthProvider)
	authProviders[github.Key()] = &github
	authProviders[dropbox.Key()] = &dropbox

	router := api.NewRouter(
		dbUrl,
		logging.HttpRequestLogger{},
		authProviders,
	)

	log.Fatal(http.ListenAndServe(":8081", router))
}
