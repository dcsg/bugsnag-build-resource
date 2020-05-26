package resource

type Source struct {
	ApiKey string `json:"api_key"`
}

type SourceControl struct {
	Provider      string `json:"provider"`
	RepositoryUrl string `json:"repository"`
	Revision      string `json:"revision"`
}

type Metadata struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Version struct {
	Build string `json:"build"`
}

type CheckRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

type InRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

type InResponse struct {
	Version  Version    `json:"version"`
	Metadata []Metadata `json:"metadata,omitempty"`
}

type OutParams struct {
	AppVersion    string        `json:"app_version"`
	ReleaseStage  string        `json:"release_stage,omitempty"`
	SourceControl SourceControl `json:"source_control,omitempty"`
}

type OutRequest struct {
	Source    Source    `json:"source"`
	OutParams OutParams `json:"params"`
}

type OutResponse struct {
	Version  Version    `json:"version"`
	Metadata []Metadata `json:"metadata,omitempty"`
}

type BugsnagBuildParams struct {
	ApiKey        string        `json:"apiKey"`
	AppVersion    string        `json:"appVersion"`
	ReleaseStage  string        `json:"releaseStage,omitempty"`
	SourceControl SourceControl `json:"sourceControl"`
}

type BugsnagResponse struct {
	Status   string   `json:"string"`
	Warnings []string `json:"warnings,omitempty"`
	Errors   []string `json:"errors,omitempty"`
}
