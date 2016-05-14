package models

type Build struct {
	ID        int64  `db:"id" json:"id"`
	Handle    string `db:"handle" json:"handle"`
	ProjectID int64  `db:"project_id" json:"project_id"`
}

func (db *DB) SaveBuild(b *Build) error {
	q := `INSERT INTO build(project_id, handle) VALUES($1,$2)`
	_, err := db.Exec(q, b.ProjectID, b.Handle)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) SaveBuilds(builds []*Build) []error {
	errs := make([]error, 0)
	for _, b := range builds {
		err := db.SaveBuild(b)
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errs
	} else {
		return nil
	}
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
