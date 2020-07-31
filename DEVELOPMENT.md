# Development

Want to contribute code to this project?
Here's what you need to know.

## Tools

These are the tools used in this project:

* [Go 1.14](https://golang.org/) as the platform
* [Docker](https://www.docker.com/) (or [Podman](https://podman.io/) if you prefer it) for packaging the tool and running various tasks
* [golangci-lint](https://golangci-lint.run/) for finding smelly code
* [GoReleaser](https://goreleaser.com/) for releasing software

## Repository

The code for this project is hosted in [Gitlab](https://gitlab.com/polarsquad/eks-auth-sync).
You can contribute code (using merge requests) and tickets there.

## Continuous Integration

[Gitlab CI](https://gitlab.com/polarsquad/eks-auth-sync/-/pipelines) is used for building and publishing this project.
All merge requests must pass the CI pipeline before they can be merged.
The pipeline is configured in [`.gitlab-ci.yml`](.gitlab-ci.yml) file.

## Helper scripts

The repository contains a few helper scripts to aid in development.

* `bin/build.sh`:
  Build the project using your computer's architecture.
* `bin/build-in-docker.sh`:
  Build the project using a standard Go Docker image, and build a Docker image for it.
  This is handy when you don't have the Go tools installed.
* `bin/lint.sh`:
  Run lint checks.
* `bin/test.sh`:
  Run tests and generate a code coverage report.

## Releasing

To perform a release, create a new Git tag in format `vX.X.X` and push it to the repository.

```bash
git tag vX.X.X -m "Version X.X.X"
git push --tags
```

Gitlab CI will take care of publishing the release artifacts.
