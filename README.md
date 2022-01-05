# Cludo - Cloud Sudo

* [Github Repo](https://github.com/superorbital/cludo)
* [cludo](https://hub.docker.com/r/superorbital/cludo) and  [cludod](https://hub.docker.com/r/superorbital/cludod) container images

Distributing AWS credentials is painful and dangerous. Leaked credentials result in hours to days of operations headaches and developing an automated rotation system is expensive. Enter `cludo`!

> Never distribute or force developers to rotate AWS credentials again.

The `cludo` command is run locally on the developer machine. It gets temporary AWS credentials from the centralized `cludod` server, and then provides the credentials locally via environment variables to arbitrary commands. Those credentials expire after a short amount of time, so any leaked credentials are already outdated.

`cludo` currently only supports AWS, but we plan to expand to many other backends in the future.

This README documents the client.  See also [SERVER.md](SERVER.md) and [DEVELOPMENT.md](DEVELOPMENT.md).

## Installation

``` bash
go get -u github.com/superorbital/cludo/cmd/cludo/cludo
```

## Configuration

The `cludo` client will read _both_ your user's `~/.cludo/cludo.yaml` file and the `cludo.yaml` file in your current working directory.  This allows you to configure per-repo and per-user aspects separately.

For example, it's common to have the following in your `~/.cludo/cludo.yaml` file to configure your user's SSH keys for authenticating with `cludod`:

``` yaml
ssh_key_paths: 
- ~/.ssh/superorbial_cludo
- ~/.ssh/bigco_cludo
```

Then your team would include this `cludo.yaml` file in a directory in the project's git
repository to configure the target `cludod` server and environment:

``` yaml
target: https://cludo.bigco.com/staging
```

Alternatively, you can provide the values as environment variables: 

``` console
$ export CLUDO_TARGET=https://cludo.bigco.com/staging
```

Currently, only the following two values are configurable for the client:

Key             |  Description                                        | Environment Variable 
---------       |  -----------                                        | -------------------- 
`target`        |  URL of the `cludo-server` instance to connect to.  | `CLUDO_TARGET`
`ssh_key_paths` |  Paths to the private keys used for authentication. | `CLUDO_SSH_KEY_PATHS`

## Authentication with the `cludod` server

Cludo uses SSH keys for authentication.  The client will try all of the keys listed in the `ssh_key_paths` setting when authenticating with the server until one succeeds (or they all fail).

## Usage

```
cludo <command> [options]

Main commands:

    exec    - Runs the provided command with credentials provided through environment variables
    shell   - Opens an interactive shell with credentials provided through environment variables
    version - Prints the cludo client and server versions
```

For example, to list AWS EC2 instances using the currently configured target:

``` console
$ cludo exec aws ec2 describe-instances
```

You can add `--debug` to get some extra debugging output.

#### Docker

We also provide a docker image (`superorbital/cludo`). Just provide a `/etc/cludo/cludo.yaml` config file!

## AWS

The AWS backend provides the following environment variables:

Environment Variable | Description
-------------------- | -----------
`AWS_ACCESS_KEY_ID` | `cludo` environment AWS access key.
`AWS_SECRET_ACCESS_KEY` | `cludo` environment AWS secret access key.
`AWS_SESSION_TOKEN` | AWS session token generated for `cludo` environment.
`AWS_SESSION_EXPIRATION` | Time when generated AWS session token will expire.

Each time a `cludo` command that uses an environment is run, a new AWS session token is generated.

## Comparisons to other tools

Cludo is heavily inspired by [the venerable `aws-vault` tool](https://github.com/99designs/aws-vault).  `aws-vault` is entirely client-side, meaning you don't need a centralized authentication server.  But this also means that each developer is responsible for configuring the tool correctly and consistently.  This also requires that the master credentials be stored on each workstation (via one of many encrypted backends).  If you're a solo developer, then Cludo is overkill, and `aws-vault` is the right tool for you.

## Development

### Download Source Code

```sh
git clone https://github.com/superorbital/cludo
```

### Build Locally

```sh
make
```

### Github Actions

#### Dependency Management (`dependabot.yaml`)

Dependabot creates pull requests for any dependencies that are in need of an update.

* [Documentation](https://docs.github.com/en/code-security/supply-chain-security/keeping-your-dependencies-updated-automatically/about-dependabot-version-updates)

#### Code Analysis (`codeql-analysis.yaml`)

CodeQL is a  semantic code analysis engine that is used to discover vulnerabilities across a codebase. CodeQL lets you query code as though it were data. Write a query to find all variants of a vulnerability and eradicating it forever. Then share your query to help others do the same.

* [Documentation](https://codeql.github.com/)

#### CI Pipeline (`ci.yaml`)

Our CI pipeline uses [GitHub Actions](https://github.com/features/actions) and [GitHub Secrets](https://docs.github.com/en/actions/security-guides/encrypted-secrets) to test, build, and release our code.

*Note: With the exception of GITHUB_TOKEN, secrets are not passed to the runner when a workflow is triggered from a forked repository.*

The workflow looks something like this:

* On any push:
  * Set up `qemu` and `binfmt` to support multi-architecture container image builds.
  * Checkout source code.
  * Setup `go` environment.
  * Setup `docker buildx` and run `docker login`
  * Install `go` dependencies.
  * Detect if this is a Pull Request (PR).
  * Run `go` tests.
  * Determine the most recent version tag in `git`
  * Parse the Change Log.
  * Build and push `cludo` and `cludod` container images
    * We only do this step if:
      * We **ARE NOT** on the `main` branch **OR**
      * We **ARE** on the `main` branch and a new version has been identified in `CHANGELOG.md`.
  * Build  `cludo` and `cludod` binaries for Github release
    * We only do this step if:
      * We **ARE** on the `main` branch and a new version has been identified in `CHANGELOG.md`.
  * Create a **non-production release** on Github
    * We only do this step if:
      * We **ARE** on the `main` branch, a new version has been identified in `CHANGELOG.md`, and the release version has a suffix (*e.g. `v0.0.1-alpha`*)
  * Create a **production release** on Github
    * We only do this step if:
      * We **ARE** on the `main` branch, a new version has been identified in `CHANGELOG.md`, and the release version does not have a suffix (*e.g. `v1.0.0`*)
  * Trigger a Slack message via [a workflow created from this file](https://github.com/superorbital/cludo/blob/main/.slack/workflows/cludo_github_actions.slackworkflow).
    * We only do this step if:
      * A Github release was created.
