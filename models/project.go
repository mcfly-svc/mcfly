package models

type Project struct {
	ID						int64			`db:"id" 				json:"id"`
	Name					string		`db:"name" 			json:"name"`
	Username			string		`db:"username" 	json:"username"`
	Service				string		`db:"service" 	json:"service"`
}

func (db *DB) SaveProject(p *Project) error {
	q := `INSERT INTO project(name,username,service) VALUES($1,$2,$3)`
	_, err := db.Exec(q, p.Name, p.Username, p.Service)
	if err != nil {
		return &QueryExecError{"SaveProject", q, err}
	}
	return nil
}

func (db *DB) GetProjects() (*[]Project, error) {
	projects := &[]Project{}
	err := db.Select(projects, `SELECT * FROM project`)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (db *DB) DeleteProject(id int64) error {
	q := `DELETE FROM project WHERE id=$1`
	_, err := db.Exec(q, id)
	if err != nil {
		return &QueryExecError{"DeleteProject", q, err}
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