package main

/*import (
	"fmt"
	"log"

	_ "github.com/mattes/migrate/driver/postgres"
	"github.com/mattes/migrate/migrate"
)

var pgDriver = "postgres://localhost:5432/marsupi_test?sslmode=disable"

func RunCommands(args []string) {
	cmd := args[0]
	switch cmd {
	case "db":
	case "db.clean":
		db.RunHelperScript("./db/helpers/clean.sh")
	case "db.recreate":
		db.RunHelperScript("./db/helpers/recreate.sh")
	case "db.seed":
		db.RunHelperScript("./db/helpers/seed.sh")
	default:
		log.Fatal(fmt.Errorf("Unknown command %s", cmd))
	}
}

func MigrateUp() {

}

func MigrateDown() {
	errs, ok := migrate.DownSync(pgDriver, "./db/migrations")
	if !ok {
		for _, err := range errs {
			fmt.Println("MigrateDown Error: ", err)
		}
	}
}
*/
