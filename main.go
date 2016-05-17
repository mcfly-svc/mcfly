package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/chrismrivera/cmd"
	"github.com/mikec/msplapi/api"
	"github.com/mikec/msplapi/config"
	"github.com/mikec/msplapi/db"
	"github.com/mikec/msplapi/logging"
	"github.com/mikec/msplapi/mq"
	"github.com/mikec/msplapi/provider"

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
		dir := cmd.Arg("direction")
		log.Printf("Running database migration %s", dir)
		db.RunMigrate(cfg.DatabaseUrl, dir)
		return nil
	},
)

var dbCreateCmd = cmd.NewCommand(
	"create-db", "Database", "Creates the database",
	func(cmd *cmd.Command) {},
	func(cmd *cmd.Command) error {
		log.Println("Creating database")
		db.Create(cfg.DatabaseUrl)
		return nil
	},
)

var dbCleanCmd = cmd.NewCommand(
	"clean-db", "Database", "Removes all data from the database",
	func(cmd *cmd.Command) {},
	func(cmd *cmd.Command) error {
		log.Println("Cleaning database")
		db.Clean(db.Connect(cfg.DatabaseUrl))
		return nil
	},
)

var dbSeedCmd = cmd.NewCommand(
	"seed-db", "Database", "Adds seed data to the database",
	func(cmd *cmd.Command) {},
	func(cmd *cmd.Command) error {
		log.Println("Seeding database")
		db.Seed(db.Connect(cfg.DatabaseUrl))
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

	github := provider.GitHub{
		GitHubClient:         &provider.GoGitHubClient{},
		SourceProviderConfig: &srcProviderCfg,
	}
	dropbox := provider.Dropbox{
		SourceProviderConfig: &srcProviderCfg,
	}

	authProviders := make(map[string]provider.AuthProvider)
	authProviders[github.Key()] = &github

	sourceProviders := make(map[string]provider.SourceProvider)
	sourceProviders[github.Key()] = &github
	sourceProviders[dropbox.Key()] = &dropbox

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

func generateAccessToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
