# cludo - Development

## Download Source Code

```sh
git clone https://github.com/superorbital/cludo
```

* [Code Documentation](https://pkg.go.dev/github.com/superorbital/cludo)

## Build Locally

Development is driven by make (see the [Makefile](https://github.com/superorbital/cludo/blob/main/Makefile)). Running `make` will run the tests, build the `cludo` and `cludod` binaries, and create the docker images.

```sh
make
```

There's also:

```sh
make build
make test
make -j $(nproc) docker-local-arch-build
make swagger
# and more...
```

* Binaries are cross-compiled for various OS's and architectures and are available in the `builds/` directory.

## Running the server locally

If you have a `cludod.yaml` file in `~/.cludod` then a local copy of the server can be built and spun up with:

```sh
docker compose up --build -d
```

## Running the CLI locally

You can run the CLI in a container using something like this:

```sh
docker run -ti \
  -v ${PWD}/data/cludo.yaml:/root/.cludo/cludo.yaml \
  -v ${PWD}/data/id_ed25519_TEST_KEY_PP:/app/id_ed25519_TEST_KEY_PP \
  ${IMAGE_ID} \
  "cludo shell"
```

## Release

* Create a new branch from main
* Make your changes
* Update [VERSION](./VERSION), [CHANGELOG.md](./CHANGELOG.md), [the kustomization.yaml files under k8s/kustomize/overlays](./k8s/kustomize/overlays), and the helm [Chart.yaml in superorbital/helm-charts](.[/k8s/helm/cludo/Chart.yaml](https://github.com/superorbital/helm-charts/blob/main/charts/cludod/Chart.yaml) to reflect/install the new version.
* Create a PR
* Get approval
* Merge the PR
* Sit back and let the CI pipeline do it's job

## Kubernetes & ContainerD

If you are running Kubernetes locally with containerd via something like [colima](https://github.com/abiosoft/colima) or [Rancher Desktop](https://rancher.com/downloads/rancher-desktop/) then you can use something like the following to automate building and running the client and server.

* *NOTE*: You will need:
  * A valid `cludod.yaml` config file stored at `k8s/kustomize/base/files-secrets/secret-cludod.yaml` for this to work as expected.
  * `nerdctl` aliased to the correct version of `nerdctl` for your setup.

```sh
make nc && { kukectl delete -k ./k8s/kustomize/overlays/local ; kubectl apply -k ./k8s/kustomize/overlays/local && sleep 10; echo -e "\n\n"; }  && ./builds/darwin_amd64_cludo exec aws sts get-caller-identity
```

## Github Actions

### Dependency Management (`dependabot.yaml`)

Dependabot creates pull requests for any dependencies that are in need of an update.

* [Documentation](https://docs.github.com/en/code-security/supply-chain-security/keeping-your-dependencies-updated-automatically/about-dependabot-version-updates)

### Code Analysis (`codeql-analysis.yaml`)

CodeQL is a  semantic code analysis engine that is used to discover vulnerabilities across a codebase. CodeQL lets you query code as though it were data. Write a query to find all variants of a vulnerability and eradicating it forever. Then share your query to help others do the same.

* [Documentation](https://codeql.github.com/)

### CI Pipeline (`ci.yaml`)

Our CI pipeline uses [GitHub Actions](https://github.com/features/actions) and [GitHub Secrets](https://docs.github.com/en/actions/security-guides/encrypted-secrets) to test, build, and release our code.

*Note: With the exception of GITHUB_TOKEN, secrets are not passed to the runner when a workflow is triggered from a forked repository.*

The workflow looks something like this:

* On any push:
  * Set up `qemu` and `binfmt` to support multi-architecture container image builds.
  * Checkout source code.
  * Setup `go` environment.
  * Install `snyk` CLI
  * Setup `docker buildx` and run `docker login`
  * Install and Run Kustomize
  * List Kustomize Generated Manifests
  * Install `go` dependencies.
  * Detect if this is a Pull Request (PR).
  * Run `go` tests.
  * Run `snyk` to check for Go vulnerabilities.
  * Upload Snyk's `sarif` report to Github.
  * Determine the most recent version tag in `git`
  * Parse the Change Log.
  * Build and push `cludo` and `cludod` container images
    * We only do this step if:
      * We **ARE NOT** on the `main` branch **OR**
      * We **ARE** on the `main` branch and a new version has been identified in `CHANGELOG.md`.
  * Run `snyk` to check for `cludo` and `cludod` image container vulnerabilities.
  * Upload Snyk's `sarif` reports to Github.
  * Build  `cludo` and `cludod` binaries for Github release
    * We only do this step if:
      * We **ARE** on the `main` branch and a new version has been identified in `CHANGELOG.md`.
  * Compress Binaries via UPX
    * We only do this step if:
      * We **ARE** on the `main` branch and a new version has been identified in `CHANGELOG.md`.
  * Create a **non-production release** on Github
    * We only do this step if:
      * We **ARE** on the `main` branch, a new version has been identified in `CHANGELOG.md`, and the release version has a suffix (*e.g. `0.0.1-alpha`*)
  * Create a **production release** on Github
    * We only do this step if:
      * We **ARE** on the `main` branch, a new version has been identified in `CHANGELOG.md`, and the release version does not have a suffix (*e.g. `1.0.0`*)
  * Install Cosign
    * We only do this step if:
      * We **ARE** on the `main` branch, a new version has been identified in `CHANGELOG.md`, and the release version does not have a suffix (*e.g. `1.0.0`*)
  * Add cosign signature to cludo latest
    * We only do this step if:
      * We **ARE** on the `main` branch, a new version has been identified in `CHANGELOG.md`, and the release version does not have a suffix (*e.g. `1.0.0`*)
  * Add cosign signature to cludod latest
    * We only do this step if:
      * We **ARE** on the `main` branch, a new version has been identified in `CHANGELOG.md`, and the release version does not have a suffix (*e.g. `1.0.0`*)
  * Trigger a Slack message via a workflow.
    * We only do this step if:
      * A Github release was created.

### Pull Request Closure (`pr_closed.yaml`)

The workflow looks something like this:

* On any PR closure:
  * Install a fork of Docker 's `hub-tool`
  * Remove all PR-related container images via `hub-tool`
