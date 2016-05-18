package models

type Build struct {
	ID           int64   `db:"id"`
	Handle       string  `db:"handle"`
	ProjectID    int64   `db:"project_id"`
	DeployStatus string  `db:"deploy_status"`
	ProviderUrl  *string `db:"provider_url"`
}

func (db *DB) SaveBuild(b *Build, p *Project) error {
	q := `INSERT INTO build(project_id, handle, deploy_status, provider_url) VALUES($1,$2,$3,$4)`
	_, err := db.Exec(q, p.ID, b.Handle, "pending", b.ProviderUrl)
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
