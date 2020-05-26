package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	resource "github.com/dcsg/bugsnag-build-resource"
)

func main() {
	var (
		payload     resource.OutRequest
		bugsnagResp resource.BugsnagResponse
	)

	artifactDirectory := os.Args[1]

	err := json.NewDecoder(os.Stdin).Decode(&payload)
	if err != nil {
		panic(err)
	}

	if true == IsFileExists(payload.OutParams.AppVersion, artifactDirectory) {
		payload.OutParams.AppVersion = RetrieveFileText(payload.OutParams.AppVersion, artifactDirectory)
	}

	if payload.OutParams.SourceControl != (resource.SourceControl{}) {
		if false == IsSourceControlProviderValid(payload.OutParams.SourceControl.Provider) {
			panic("invalid provider. It must be one of the following: github, github-enterprise, bitbucket, bitbucket-server, gitlab, and gitlab-onpremise")
		}

		if payload.OutParams.SourceControl.RepositoryUrl == "" {
			panic("'repository' is required when source_control is defined")
		}

		if payload.OutParams.SourceControl.Revision == "" {
			panic("'revision' param is required when source_control is defined")
		}

		if true == IsFileExists(payload.OutParams.SourceControl.Revision, artifactDirectory) {
			payload.OutParams.SourceControl.Revision = RetrieveFileText(payload.OutParams.SourceControl.Revision, artifactDirectory)
		}
	}

	bugsnagHost := os.Getenv("BUGSNAG_HOST")
	if bugsnagHost == "" {
		bugsnagHost = "https://build.bugsnag.com/"
	}

	bugsnagPayload, err := json.Marshal(ConvertParamsToBugsnagBuild(payload))
	resp, err := http.Post(bugsnagHost, "application/json", bytes.NewBuffer(bugsnagPayload))
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&bugsnagResp)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != http.StatusOK {
		panic(bugsnagResp.Errors)
	}

	response := resource.OutResponse{
		Version:  resource.Version{Build: payload.OutParams.AppVersion},
		Metadata: GenerateMetadata(payload.OutParams, bugsnagResp),
	}

	err = json.NewEncoder(os.Stdout).Encode(response)
	if err != nil {
		panic(err)
	}
}

func ConvertParamsToBugsnagBuild(p resource.OutRequest) resource.BugsnagBuildParams {
	var (
		b resource.BugsnagBuildParams
	)

	b.ApiKey = p.Source.ApiKey
	b.AppVersion = p.OutParams.AppVersion
	b.ReleaseStage = p.OutParams.ReleaseStage
	b.SourceControl = p.OutParams.SourceControl

	return b
}

func GenerateMetadata(p resource.OutParams, b resource.BugsnagResponse) []resource.Metadata {
	metadata := []resource.Metadata{
		{Name: "app_version", Value: p.AppVersion},
	}

	if p.ReleaseStage != "" {
		metadata = append(metadata, resource.Metadata{Name: "release_stage", Value: p.ReleaseStage})
	}

	if p.SourceControl != (resource.SourceControl{}) {
		metadata = append(
			metadata,
			resource.Metadata{Name: "source_control.provider", Value: p.SourceControl.Provider},
			resource.Metadata{Name: "source_control.repository", Value: p.SourceControl.RepositoryUrl},
			resource.Metadata{Name: "source_control.revision", Value: p.SourceControl.Revision},
		)
	}

	if len(b.Warnings) > 0 {
		warnings, err := json.Marshal(b.Warnings)
		if err != nil {
			panic(err)
		}
		metadata = append(metadata, resource.Metadata{Name: "bugsnag.warnings", Value: string(warnings)})
	}

	return metadata
}

func IsSourceControlProviderValid(p string) bool {
	switch p {
	case "github", "github-enterprise", "bitbucket", "bitbucket-server", "gitlab", "gitlab-onpremise":
		return true
	}

	return false
}

func IsFileExists(fp string, artifactDirectory string) bool {
	path := filepath.Join(artifactDirectory, fp)
	// Check if file already exists
	if _, err := os.Stat(path); err == nil {
		return true
	}

	return false
}

func RetrieveFileText(fp string, artifactDirectory string) string {
	path := filepath.Join(artifactDirectory, fp)
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		return scanner.Text()
	}

	return ""
}

// appVersion, err := ioutil.ReadFile(filepath.Join(artifactDirectory, payload.OutParams.AppVersion))
// string(appVersion)
