# Cludo - Cloud Sudo

Distributing AWS credentials is painful and dangerous.  Leaked credentials result in hours to days of operations headaches and developing an automated rotation system is expensive. Enter - `cludo`

> Never distribute or rotate AWS credentials again, with `cludo`

`cludo` is run locally on the developer machine.  It gets temporary credentials from the centralized cludo-server, and provides them locally to arbitrary subprocesses.

`cludo` currently supports the following _environments_:

- AWS

## Installation

TODO: Installation

`go install github.com/superorbital/cludo/cmd/cludo/cludo`

`go install github.com/superorbital/cludo/cmd/cludo-server/cludo-server`

## Setup

`cludo` requires a `cludo.yaml` file in your current working directory (`./.cludo/cludo.yaml`) or your home directory (`~/.cludo/cludo.yaml`) or configuration provided through environment variables prefixed with `CLUDO_`

The following configuration options are supported:

Environment Variable | YAML path | Description
-------------------- | --------- | -----------
`CLUDO_SERVER_URL` | `server_url` | URL of the `cludo-server` instance to connect to.
`CLUDO_PRIVATE_KEY` | `private_key` | Path to a private key for authentication.
`CLUDO_SHELL` | `shell` | Path to the shell to launch when requested. Defaults to user's login shell.

## Usage

```
cludo <command> [options]

Main commands:

    exec    - Runs the provided command with credentials provided through environment variables
    shell   - Opens an interactive shell with credentials provided through environment variables
    version - Prints the cludo client and server versions
```


We recommend adding `.cludo/` to your `.gitignore` files to avoid accidentally committing secrets.

## Environments

`cludo` environments configure different environment variables. `cludo` currently only supports the AWS environment.

### Environments - AWS

When enabled, the AWS environment provides the following environment variables:

Environment Variable | Description
-------------------- | -----------
`AWS_ACCESS_KEY_ID` |
`AWS_SECRET_ACCESS_KEY` |
`AWS_REGION` |

## Running a server

1. Install `cludo-server`:

   ```
   go install github.com/superorbital/cludo/cmd/cludo-server/cludo-server
   ```

2. Configure `cludo-server` by providing a `cludo-server.yaml` file.
3. Run `cludo-server -c /path/to/cludo-server.yaml`

`cludo-server` supports the following configuration options:

```
# cludo-server.yaml

client:
  default:
    server_url: https://www.example.com/
    key_path: ~/.ssh/id_rsa
    shell_path: /usr/local/bin/bash
    roles: default
server:
  port: 443
  plugins:
    aws:
      sessionLength: 10m
      awsAccessKeyID: ...
      awsSecretAccessKey: ...
      awsAssumeRoleARN: ...
  users:
    - publicKey: "ssh-rsa aisudpoifueuyrlkjhflkyhaosiduyflakjsdhflkjashdf7898798765489..."
      roles:
        aws:
          so_org:
            arn: "aws:arn:iam:..."
          so_dev:
            arn: "aws:arn:iam:..."
      default_role: aws_so_org
    - publicKey: "ssh-rsa ..."
      roles:
        aws:
          so_org:
            arn: "aws:arn:iam:..."
          so_dev:
            arn: "aws:arn:iam:..."
      default_role: aws_so_dev
```

We also provide a docker image (`superorbital/cludo-server`) with `cludo-server` pre-installed. Just provide a `/etc/cludo-server/cludo-server.yaml` config file.


# Development

## Release

$ go get github.com/mitchellh/gox
$ gox -osarch='darwin/amd64 darwin/arm64 linux/386 linux/amd64 linux/arm linux/arm64 windows/386 windows/amd64' -output './builds/wordchain_{{.OS}}_{{.Arch}}'

* Create a release in Github with the resulting binaries.

# Acknowledgements

* Cobra & Viper integration code heavily inspired by:
  * [https://github.com/carolynvs/stingoftheviper](https://github.com/carolynvs/stingoftheviper)
  * **License**: MIT
