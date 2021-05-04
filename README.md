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

## Usage

```
cludo <command> [options]

Main commands:

    exec    - Runs the provided command with credentials provided through environment variables
    login   - Authenticate with a `cludo-server` instance
    logout  - Forget current authentication state (if any)
    shell   - Opens an interactive shell with credentials provided through environment variables
    status  - Prints the current authentication state for cludo

Administration:

    users-list    - Print all users
    users-new     - Create a new user and assign its role
    users-delete  - Delete a user
    users-assign  - Assign a role to a user

    roles-list    - Print all roles
    roles-new     - Creates a new role using $EDITOR
    roles-delete  - Removes a role. Also unassigns role from all users.
    roles-edit    - Opens an $EDITOR instance to edit a role definition.
```


We recommend adding `.cludo/` to your `.gitignore` files to avoid accidentally committing secrets.

## Environments

`cludo` environments configure different environment variables. `cludo` currently only supports the AWS environment.

### Environments - AWS

When enabled, the AWS environment provides the following environment variables:

Environment Variable | Description
-------------------- | -----------
`CLUDO_SERVER_URL` | URL of the `cludo-server` instance to connect to.

## Running a server

1. Install `cludo-server`:

   ```
   go install github.com/superorbital/cludo/cmd/cludo-server/cludo-server
   ```

2. Configure `cludo-server` by providing a `cludo-server.yaml` file.
3. Run `cludo-server -c /path/to/cludo-server.yaml`

We also provide a docker image (`superorbital/cludo-server`) with `cludo-server` pre-installed. Just provide a `/etc/cludo-server/cludo-server.yaml` config file.
