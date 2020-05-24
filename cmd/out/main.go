package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	resource "github.com/dcsg/bugsnag-build-resource"
)

func main() {
	var (
		payload resource.OutPayload
	)

	artifactDirectory := os.Args[1]

	err := json.NewDecoder(os.Stdin).Decode(&payload)
	if err != nil {
		panic(err)
	}

	if true == isFileExists(payload.OutParams.AppVersion, artifactDirectory) {
		payload.OutParams.AppVersion = retrieveFileText(payload.OutParams.AppVersion, artifactDirectory)
	}

	if payload.OutParams.SourceControl != (resource.SourceControl{}) {
		if false == isSourceControlProviderValid(payload.OutParams.SourceControl.Provider) {
			panic("invalid provider. It must be one of the following: github, github-enterprise, bitbucket, bitbucket-server, gitlab, and gitlab-onpremise")
		}

		if true == isFileExists(payload.OutParams.SourceControl.Revision, artifactDirectory) {
			payload.OutParams.SourceControl.Revision = retrieveFileText(payload.OutParams.SourceControl.Revision, artifactDirectory)
		}
	}

	bugsnagHost := os.Getenv("BUGSNAG_HOST")
	if bugsnagHost == "" {
		bugsnagHost = "https://build.bugsnag.com/"
	}

	requestPayload, err := json.Marshal(ConvertParamsToBugsnagBuild(payload))
	resp, err := http.Post(bugsnagHost, "application/json", bytes.NewBuffer(requestPayload))
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	bodyString := string(bodyBytes)

	if resp.StatusCode != http.StatusOK {
		panic(errors.New(bodyString))
	}

	response := resource.OutResponse{
		Version: resource.Version{Id: payload.OutParams.AppVersion},
	}

	err = json.NewEncoder(os.Stdout).Encode(response)
	if err != nil {
		panic(err)
	}
}

func ConvertParamsToBugsnagBuild(p resource.OutPayload) resource.BugsnagBuild {
	var (
		b resource.BugsnagBuild
	)

	b.ApiKey = p.Source.ApiKey
	b.AppVersion = p.OutParams.AppVersion
	b.ReleaseStage = p.OutParams.ReleaseStage
	b.SourceControl = p.OutParams.SourceControl

	return b
}

func isSourceControlProviderValid(p string) bool {
	switch p {
	case "github", "github-enterprise", "bitbucket", "bitbucket-server", "gitlab", "gitlab-onpremise":
		return true
	}

	return false
}

func isFileExists(fp string, artifactDirectory string) bool {
	path := filepath.Join(artifactDirectory, fp)
	// Check if file already exists
	if _, err := os.Stat(path); err == nil {
		return true
	}

	return false
}

func retrieveFileText(fp string, artifactDirectory string) string {
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
