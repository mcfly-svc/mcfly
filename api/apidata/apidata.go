package apidata

type ApiResponse struct {
	Status string `json:"status"`
}

type ApiError struct {
	Error string `json:"error"`
}

type LoginReq struct {
	Token    string `json:"token" validate:"nonzero"`
	Provider string `json:"provider" validate:"nonzero"`
}

func (lr *LoginReq) AuthProvider() string {
	return lr.Provider
}

type LoginResp struct {
	Name        *string `json:"name,omitempty"`
	AccessToken string  `json:"access_token"`
}

type BuildReq struct {
	Handle        string  `json:"handle"`
	ProjectHandle string  `json:"project_handle"`
	Provider      string  `json:"provider"`
	ProviderUrl   *string `json:"provider_url,omitempty"`
}

func (br *BuildReq) SourceProvider() string {
	return br.Provider
}

type DeployReq struct {
	BuildHandle         string `json:"build_handle" validate:"nonzero"`
	SourceProjectHandle string `json:"project_handle" validate:"nonzero"`
	Provider            string `json:"provider" validate:"nonzero"`
}

func (r *DeployReq) SourceProvider() string {
	return r.Provider
}

func (r *DeployReq) ProjectHandle() string {
	return r.SourceProjectHandle
}

type ProjectReq struct {
	Handle   string `json:"handle" validate:"nonzero"`
	Provider string `json:"provider" validate:"nonzero"`
}

func (pr *ProjectReq) SourceProvider() string {
	return pr.Provider
}

type ProjectResp struct {
	Handle   string `json:"handle"`
	Url      string `json:"url"`
	Provider string `json:"provider"`
}
