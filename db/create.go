package db

// Create destroys and creates the database by running all migrations down and then up
func Create(dbUrl string) {
	RunMigrate(dbUrl, "down")
	RunMigrate(dbUrl, "up")
}
