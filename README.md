# Bugsnag Build Resource

Implements a [Concourse CI](https://concourse-ci.org/) resource type that reports a new build to [Bugsnag](https://bugsnag.com/) using the [Bugsnag Build API](https://docs.bugsnag.com/api/build/). 

## Source Configuration

* `api_key`: *Required.* The Bugsnag Project API KEY

## Behavior

### `check`: None.

### `in`: None.

### `out`: Notifies Bugsnag of a new Build.

#### Parameters

* `app_version`: *Required.* The app version. It can be a string or a path to a file (ex. `master-code/.git/ref`)
* `release_stage`: *Optional.* The stage (ex. `staging`, `production`)
* `source_control`: *Optional.* Provide information about the source code
  * `provider`: *Required.* One of: `github`, `github-enterprise`, `bitbucket`, `bitbucket-server`, `gitlab`, `gitlab-onpremise`
  * `repository`: *Required.* The url of the repository (ex. `https://github.com/dcsg/bugsnag-build-resource.git`)
  * `revision`: *Required.* The commit reference. It can be a string or a filepath (ex. `master-code/.git/ref`)

## Example

```yaml
---
resource_types:
- name: bugsnag-build
  type: registry-image
  source:
    repository: dcsg/bugsnag-build-resource

resources:
- name: bugsnag-build
  type: bugsnag-build
  source:
    api_key: '<your bugsnag project api key>'

- name: master-code
  type: git
  icon: git
  source:
    uri: https://github.com/dcsg/bugsnag-build-resource.git
    branch: master

jobs:
  - name: notify-bugsnag
    plan:
      - get: master-code
      - put: bugsnag-build
        params:
          app_version: master-code/.git/ref
          release_stage: staging
          source_control:
            provider: gitlab
            repository: https://github.com/dcsg/bugsnag-build-resource.git
            revision: master-code/.git/ref
```

## Not implemented (yet)

The Bugsnag Build API have the following properties that are not yet implemented:

* appBundleVersion
* builderName
* metadata
* autoAssignRelease

## Development

### Prerequisites

* golang is *required* - version 1.14.x is tested; earlier versions may also
  work.
* docker is *required* - version 19.03.x is tested; earlier versions may also
  work.

### Running the tests

The tests have been embedded with the `Dockerfile`; ensuring that the testing
environment is consistent across any `docker` enabled platform. When the docker
image builds, the test are run inside the docker container, on failure they
will stop the build.

Run the tests with the following command:

```sh
docker build -t bugsnag-build-resource -f dockerfiles/Dockerfile .
```

### Contributing

Please make all pull requests to the `master` branch.
