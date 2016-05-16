package models

type Build struct {
	// TODO: get rid of json: for all models?
	ID        int64  `db:"id" json:"id"`
	Handle    string `db:"handle" json:"handle"`
	ProjectID int64  `db:"project_id"`
}

func (db *DB) SaveBuild(b *Build, p *Project) error {
	q := `INSERT INTO build(project_id, handle) VALUES($1,$2)`
	_, err := db.Exec(q, p.ID, b.Handle)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) GetProjectBuilds(p *Project) ([]Build, error) {
	builds := []Build{}
	err := db.Select(
		&builds,
		`SELECT build.* FROM build
		 INNER JOIN project 
		 ON project.id=build.project_id 
		 WHERE project.handle=$1
		 AND project.source_provider=$2`,
		p.Handle,
		p.SourceProvider,
	)
	if err != nil {
		return nil, err
	}
	return builds, nil
}
