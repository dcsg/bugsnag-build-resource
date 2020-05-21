package resource

type Source struct {
	ApiKey string `json:"api_key"`
}

type SourceControl struct {
	Provider   string `json:"provider"`
	Repository string `json:"repository_url"`
	Revision   string `json:"revision"`
}

type OutParams struct {
	AppVersion    string        `json:"app_version"`
	ReleaseStage  string        `json:"release_stage,omitempty"`
	SourceControl SourceControl `json:"source_control"`
}

type OutPayload struct {
	Source    Source    `json:"source"`
	OutParams OutParams `json:"params"`
}

type Version struct {
	Id string `json:"id"`
}

type InResponse struct {
	Version Version `json:"version"`
}

type OutResponse struct {
	Version Version `json:"version"`
}

type BugsnagBuild struct {
	ApiKey        string        `json:"apiKey"`
	AppVersion    string        `json:"appVersion"`
	ReleaseStage  string        `json:"releaseStage,omitempty"`
	SourceControl SourceControl `json:"sourceControl"`
}
