package models

type Project struct {
	ID             int64  `db:"id" json:"id"`
	Handle         string `db:"handle" json:"handle"`
	SourceUrl      string `db:"source_url" json:"username"`
	SourceProvider string `db:"source_provider" json:"source_provider"`
}

func (db *DB) SaveProject(p *Project, u *User) error {
	var id int64
	q := `INSERT INTO project(handle,source_url,source_provider) VALUES($1,$2,$3) RETURNING id`
	tx := db.MustBegin()
	r := tx.QueryRowx(q, p.Handle, p.SourceUrl, p.SourceProvider)
	err := r.Scan(&id)
	if err != nil {
		return err
	}
	p.ID = id
	tx.MustExec(`INSERT INTO user_project (user_id,project_id) VALUES($1,$2)`, u.ID, p.ID)
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

/*
func (db *DB) GetProjects() ([]Project, error) {
	projects := []Project{}
	err := db.Select(&projects, `SELECT * FROM project`)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (db *DB) DeleteProject(id int64) error {
	q := `DELETE FROM project WHERE id=$1`
	_, err := db.Exec(q, id)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) GetProjectById(id int64) (*Project, error) {
	project := &Project{}
	err := db.Get(project, `SELECT * FROM project WHERE id=$1`, id)
	if err != nil {
		return nil, err
	}
	return project, nil
}
*/
