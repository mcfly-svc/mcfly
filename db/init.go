package db

// Create destroys and creates the database by running all migrations down and then up
func (mdb *McflyDB) Init() {
	mdb.RunMigrate("down")
	mdb.RunMigrate("up")
}
