package models

type Project struct {
	ID     int64  `db:"id" json:"id"`
	Handle string `db:"handle" json:"handle"`
	// TODO : fix this and test
	SourceUrl      string `db:"source_url" json:"username"`
	SourceProvider string `db:"source_provider" json:"source_provider"`
}

func (db *DB) SaveProject(p *Project, u *User) error {
	var id int64
	q := `INSERT INTO project(handle,source_url,source_provider) VALUES($1,$2,$3) RETURNING id`
	tx := db.MustBegin()

	// TODO: wrap this to handle QueryError like the *DB methods
	// ... or figure out a better way to do this?
	r := tx.QueryRowx(q, p.Handle, p.SourceUrl, p.SourceProvider)
	err := r.Scan(&id)
	if err != nil {
		qErr := NewQueryError(q, err, []interface{}{p.Handle, p.SourceUrl, p.SourceProvider})
		if qErr.Name == "unique_violation" {
			return ErrDuplicate
		} else {
			return qErr
		}
	}

	p.ID = id
	tx.MustExec(`INSERT INTO user_project (user_id,project_id) VALUES($1,$2)`, u.ID, p.ID)
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

// GetAllProjects returns all projects for all users
func (db *DB) GetAllProjects() ([]Project, error) {
	projects := []Project{}
	err := db.Select(&projects, `SELECT * FROM project`)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

// GetUserProjects returns projects owned by a given user
func (db *DB) GetUserProjects(user *User) ([]Project, error) {
	projects := []Project{}
	err := db.Select(
		&projects,
		`SELECT project.* FROM project
		 INNER JOIN user_project 
		 ON user_project.project_id=project.id 
		 WHERE user_project.user_id=$1`,
		user.ID,
	)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

// GetProject returns projects by handle and source provider
func (db *DB) GetProject(projectHandle string, sourceProvider string) (*Project, error) {
	p := Project{}
	q := `SELECT * FROM project
			  WHERE handle=$1
			  AND source_provider=$2`
	err := db.Get(&p, q, projectHandle, sourceProvider)
	if err != nil {
		if err.NoRows {
			return nil, ErrNotFound
		} else {
			return nil, err
		}
	}
	return &p, nil
}

func (db *DB) DeleteUserProject(user *User, provider string, handle string) error {
	result, qErr := db.Exec(
		`DELETE FROM project
		 USING user_project
		 WHERE id=user_project.project_id
		 AND source_provider=$1
		 AND handle=$2
		 AND user_project.user_id=$3`,
		provider,
		handle,
		user.ID,
	)
	if qErr != nil {
		return qErr
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
