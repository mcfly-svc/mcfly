package main

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/chrismrivera/cmd"
	"github.com/mikec/marsupi-api/api"
	"github.com/mikec/marsupi-api/logging"
	"github.com/mikec/marsupi-api/provider"

	_ "github.com/mattes/migrate/driver/postgres"
	"github.com/mattes/migrate/migrate"
)

var cmdr *cmd.App = cmd.NewApp()
var pgDriver = "postgres://localhost:5432/marsupi_test?sslmode=disable"

func main() {

	cmdr.AddCommand(runCommand)
	cmdr.AddCommand(migrateCommand)

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

var runCommand = cmd.NewCommand(
	"run", "Server", "Runs the marsupi-api server",
	func(cmd *cmd.Command) {},
	func(cmd *cmd.Command) error {
		runServer()
		return nil
	},
)

var migrateCommand = cmd.NewCommand(
	"migrate", "Database", "Run migration scripts",
	func(cmd *cmd.Command) {
		cmd.AppendArg("direction", "up || down")
	},
	func(cmd *cmd.Command) error {
		runMigrate(cmd.Arg("direction"))
		return nil
	},
)

func runServer() {
	github := provider.GitHub{}
	dropbox := provider.Dropbox{}

	authProviders := make(map[string]provider.AuthProvider)
	authProviders[github.Key()] = &github
	authProviders[dropbox.Key()] = &dropbox

	router := api.NewRouter(
		pgDriver,
		logging.HttpRequestLogger{},
		authProviders,
	)

	log.Fatal(http.ListenAndServe(":8081", router))
}

func runMigrate(direction string) {
	var doMigrate func(string, string) ([]error, bool)
	switch direction {
	case "up":
		doMigrate = migrate.UpSync
	case "down":
		doMigrate = migrate.DownSync
	default:
		return
	}

	tmpDir := fmt.Sprintf("tmp_%s", randString())

	removeDir(tmpDir)
	createDir(tmpDir)

	assets, err := AssetDir("db/migrations")
	check(err)

	for _, asset := range assets {
		d, err := Asset(fmt.Sprintf("db/migrations/%s", asset))
		check(err)

		newPath := fmt.Sprintf("./%s/%s", tmpDir, asset)
		err = ioutil.WriteFile(newPath, d, 0777)
		check(err)
	}

	errs, ok := doMigrate(pgDriver, fmt.Sprintf("./%s", tmpDir))
	if !ok {
		for _, err := range errs {
			fmt.Println("runMigrate Error: ", err)
		}
	}

	removeDir(tmpDir)
}

func removeDir(dir string) {
	err := os.RemoveAll(dir)
	check(err)
}

func createDir(dir string) {
	err := os.Mkdir(dir, 0777)
	check(err)
}

/*func checkPathError(err error, allowedMessage string) {
	if err != nil {
		pErr, ok := err.(*os.PathError)
		if !ok {
			panic(err)
		}
		if fmt.Sprint(pErr.Err) != allowedMessage {
			panic(err)
		}
	}
}*/

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func randString() string {
	b := make([]byte, 32)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
