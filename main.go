package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/chrismrivera/cmd"
	"github.com/mcfly-svc/mcfly/api"
	"github.com/mcfly-svc/mcfly/config"
	"github.com/mcfly-svc/mcfly/db"
	"github.com/mcfly-svc/mcfly/logging"
	"github.com/mcfly-svc/mcfly/mq"
	"github.com/mcfly-svc/mcfly/provider"

	_ "github.com/mattes/migrate/driver/postgres"
)

var cmdr *cmd.App = cmd.NewApp()
var cfg *config.Config

func main() {

	_cfg, err := config.NewConfigFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}
	cfg = _cfg

	cmdr.AddCommand(runServerCmd)
	cmdr.AddCommand(dbMigrateCmd)
	cmdr.AddCommand(dbCreateCmd)
	cmdr.AddCommand(dbInitCmd)
	cmdr.AddCommand(dbCleanCmd)
	cmdr.AddCommand(dbSeedCmd)

	cmdr.Description = "McFly CLI"
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
	"run", "Server", "Runs the mcfly server",
	func(cmd *cmd.Command) {},
	func(cmd *cmd.Command) error {
		RunServer()
		return nil
	},
)

var dbMigrateCmd = cmd.NewCommand(
	"migrate", "Database", "Run migration scripts",
	func(cmd *cmd.Command) {
		cmd.AppendArg("name", "The name of the database")
		cmd.AppendArg("direction", "(up|down)")
	},
	func(cmd *cmd.Command) error {
		name := cmd.Arg("name")
		dir := cmd.Arg("direction")
		log.Printf("Running database migration %s", dir)
		newDbConn(name).RunMigrate(dir)
		return nil
	},
)

var dbInitCmd = cmd.NewCommand(
	"init-db", "Database", "Initializes the database schema",
	func(cmd *cmd.Command) {
		cmd.AppendArg("name", "The name of the database")
	},
	func(cmd *cmd.Command) error {
		name := cmd.Arg("name")
		log.Printf("Initializing the %s database\n", name)
		newDbConn(name).Init()
		return nil
	},
)

var dbCreateCmd = cmd.NewCommand(
	"create-db", "Database", "Creates a new database",
	func(cmd *cmd.Command) {
		cmd.AppendArg("name", "The name of the database")
	},
	func(cmd *cmd.Command) error {
		name := cmd.Arg("name")
		log.Printf("Creating the %s database\n", name)
		newDbConn("").Create(name)
		return nil
	},
)

var dbCleanCmd = cmd.NewCommand(
	"clean-db", "Database", "Removes all data from the database",
	func(cmd *cmd.Command) {
		cmd.AppendArg("name", "The name of the database")
	},
	func(cmd *cmd.Command) error {
		name := cmd.Arg("name")
		log.Println("Cleaning database")
		newDbConn(name).Clean()
		return nil
	},
)

var dbSeedCmd = cmd.NewCommand(
	"seed-db", "Database", "Adds seed data to the database",
	func(cmd *cmd.Command) {
		cmd.AppendArg("name", "The name of the database")
	},
	func(cmd *cmd.Command) error {
		name := cmd.Arg("name")
		log.Println("Seeding database")
		newDbConn(name).Seed()
		return nil
	},
)

// RunServer runs the HTTP server
func RunServer() {
	msgChannel := mq.CreateChannel(cfg.RabbitMQUrl)
	defer msgChannel.CloseConnection()
	defer msgChannel.CloseChannel()

	srcProviderCfg := provider.SourceProviderConfig{
		ProjectUpdateHookUrlFmt: fmt.Sprintf("%s/api/0/webhooks/{provider}/project-update", cfg.ApiUrl),
		WebhookSecret:           cfg.WebhookSecret,
	}

	authProviders := provider.GetAuthProviders()
	sourceProviders := provider.GetSourceProviders(&srcProviderCfg)

	router := api.NewRouter(
		cfg,
		logging.HttpRequestLogger{},
		msgChannel,
		generateAccessToken,
		authProviders,
		sourceProviders,
	)

	log.Fatal(http.ListenAndServe(":8081", router))
}

func newDbConn(databaseName string) *db.McflyDB {
	mdb := db.NewMcflyDB(cfg.DatabaseUrl, databaseName, cfg.DatabaseUseSSL)
	return mdb
}

func generateAccessToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
