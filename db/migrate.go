package db

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/mattes/migrate/migrate"
)

// RunMigrate runs the migration scripts from db/migrations/. direction can be `up` or `down`.
// Migration scripts are copied as binary data in migrations.go using github.com/jteeuwen/go-bindata
func RunMigrate(
	dbUrl string,
	direction string,
) {
	var doMigrate func(string, string) ([]error, bool)
	switch direction {
	case "up":
		doMigrate = migrate.UpSync
	case "down":
		doMigrate = migrate.DownSync
	default:
		log.Fatal(fmt.Errorf("No migrate direction `%s`. Use `up` or `down`", direction))
	}

	fmt.Println("Running migrate", direction)

	tmpDir := createTmpDir()

	defer removeDir(tmpDir)

	assetPath := "db/migrations"
	assets, err := AssetDir(assetPath)
	check(err)

	for _, assetFile := range assets {
		writeAssetToTmpDir(assetPath, assetFile, tmpDir)
	}

	errs, ok := doMigrate(dbUrl, fmt.Sprintf("./%s", tmpDir))
	if !ok {
		for _, err := range errs {
			fmt.Println("Migration Error: ", err)
		}
		return
	}
}

func writeAssetToTmpDir(assetPath string, assetFile string, tmpDir string) {
	d, err := Asset(fmt.Sprintf("%s/%s", assetPath, assetFile))
	check(err)

	newPath := fmt.Sprintf("./%s/%s", tmpDir, assetFile)
	err = ioutil.WriteFile(newPath, d, 0777)
	check(err)
}

func removeDir(dir string) {
	err := os.RemoveAll(dir)
	check(err)
}

func createTmpDir() string {
	tmpDir := fmt.Sprintf("tmp_%s", randString())
	removeDir(tmpDir)
	createDir(tmpDir)
	return tmpDir
}

func createDir(dir string) {
	err := os.Mkdir(dir, 0777)
	check(err)
}

func randString() string {
	b := make([]byte, 32)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
