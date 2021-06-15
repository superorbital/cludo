# Cludo - Cloud Sudo

Distributing AWS credentials is painful and dangerous.  Leaked credentials result in hours to days of operations headaches and developing an automated rotation system is expensive. Enter - `cludo`

> Never distribute or rotate AWS credentials again, with `cludo`

`cludo` is run locally on the developer machine.  It gets temporary credentials from the centralized cludo-server, and provides them locally to arbitrary subprocesses.

`cludo` currently supports the following _environments_:

- AWS

## Installation

To install or update `cludo`/`cludod`:

```shell
go get -u github.com/superorbital/cludo/cmd/cludo/cludo
go get -u github.com/superorbital/cludo/cmd/cludod/cludod
```

## Setup

`cludo` requires a `cludo.yaml` file in your current working directory (`./.cludo/cludo.yaml`) or your home directory (`~/.cludo/cludo.yaml`) or configuration provided through environment variables prefixed with `CLUDO_`

The following configuration options are supported:

```yaml
# cludo.yaml

client:
  default:
    server_url: "https://www.example.com/"
    key_path: "~/.ssh/id_rsa"
    shell_path: "/usr/local/bin/bash"
    roles: ["default"]
```

Environment Variable | YAML path | Description
-------------------- | --------- | -----------
`CLUDO_CLIENT_DEFAULT_SERVER_URL` | `client.default.server_url` | URL of the `cludo-server` instance to connect to.
`CLUDO_CLIENT_DEFAULT_KEY_PATH` | `client.default.key_path` | Path to a private key for authentication.
`CLUDO_CLIENT_DEFAULT_SHELL_PATH` | `client.default.shell_path` | Path to the shell to launch when requested. Defaults to user's login shell.
`CLUDO_CLIENT_DEFAULT_ROLES` | `client.default.roles` | List of roles to apply to cludo environment when generated. Currently only supports one role at a time. Role ids should correspond to role ids assigned to user in cludod.

## Usage

```shell
cludo <command> [options]

Main commands:

    exec    - Runs the provided command with credentials provided through environment variables
    shell   - Opens an interactive shell with credentials provided through environment variables
    version - Prints the cludo client and server versions
```

We recommend adding `.cludo/` to your `.gitignore` files to avoid accidentally committing secrets.

### Examples


#### Server

- `cludod --scheme=http --host=0.0.0.0 --port=8080`

#### Client

List AWS EC2 instances using a cludo environment:

```sh
cludo exec aws ec2 describe-instances
```

You can add `--debug` to get some extra debugging output.

We also provide a docker image (`superorbital/cludo`). Just provide a `/etc/cludo/cludo.yaml` config file!

## Environments

`cludo` environments configure different environment variables. `cludo` currently only supports the AWS environment.

### Environments - AWS

When enabled, the AWS environment provides the following environment variables:

Environment Variable | Description
-------------------- | -----------
`AWS_ACCESS_KEY_ID` | `cludo` environment AWS access key.
`AWS_SECRET_ACCESS_KEY` | `cludo` environment AWS secret access key.
`AWS_SESSION_TOKEN` | AWS session token generated for `cludo` environment.
`AWS_SESSION_EXPIRATION` | Time when generated AWS session token will expire.

Each time a `cludo` command that uses an environment is run, a new AWS session token is generated.

## Running cludod

1. Install `cludod`:

   ```shell
   go install github.com/superorbital/cludo/cmd/cludod/cludod
   ```

2. Configure `cludod` by providing a `cludod.yaml` file.
3. Run `cludod --scheme=http --port=8080 --host=0.0.0.0 -c /path/to/cludod.yaml` or just `cludod --scheme=http --port=8080 --host=0.0.0.0` if your `cludod.yaml` file is placed in: `/etc/cludod/cludod.yaml`, `~/.cludod/cludod.yaml`, `./.cludod/cludod.yaml`, or `./cludod.yaml`

`cludod` supports the following configuration options:

```yaml
# cludod.yaml
server:
    targets:
        prod:
          aws:
            arn: "aws:arn:iam:..."
            session_duration: "20m"
            access_key_id: "456DEF..."
            secret_access_key: "UVW789..."
          ssh:
            ...
        dev:
          aws:
            arn: "aws:arn:iam:..."
            session_duration: "8h"
            access_key_id: "123ABC..."
            secret_access_key: "ZXY098..."
        sean:
        robert:
        qa:
        prod_frontend:
        prod_backend:
        prod_db:
  users:
    - public_key: "ssh-rsa aisudpoifueuyrlkjhflkyhaosiduyflakjsdhflkjashdf7898798765489..."
      targets: ["prod", "dev"]
```

We also provide a docker image (`superorbital/cludod`) with `cludod` pre-installed. Just provide a `/etc/cludod/cludod.yaml` config file.

## Running `cludo` client

In order to access a running cludod server, create a `cludo.yaml` file in the root of your application repository. This allows the repository to specify targets users should use when developing a particular application. An example `cludo.yaml` file contains a single key:

```yaml
target: https://my-cludod-server.myorg.com/staging
```

The final fragment of the url path is used as the profile name set in the cludod server config. The above example would send a request to the URL `https://my-cludod-server.myorg.com` using the target profile `profile`.

Users may want to configure the SSH keys they use for authentication. This can be done globally for a single user in the files `~/.cludo/config.yaml` or `~/.config/cludod/config.yaml` This file contains user specific metadata about how to authenticate them to the cludod server:

```yaml
client:
  key_path: "~/.ssh/my_ssh_key"
  shell_path: "/usr/local/bin/bash"
```

These values can also be set on the command line via options on the command.

## Development

To build/test `cludo`/`cludod`:

```shell
make
```

Binaries cross compiled for various OS's and architectures are available in the `builds/` directory.

### Release

- Checkout the `main` branch
- `make all docker-push`
- Create a Github release
- Attach the binaries for all platforms to the release
- List the fully qualified image tags in the release description.

**TODO**: Automate the creation of a release in Github with the resulting binaries.
