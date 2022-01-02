# Development

Development is driven by `make` (see the [Makefile](./Makefile)). Running
`make` will run the tests, build the `cludo` and `cludod` binaries, and create
the docker images.

There's also:

- `make build`
- `make test`
- `make docker-push`
- `make swagger`

Binaries cross compiled for various OS's and architectures are available in the `builds/` directory.

### Running the server locally

If you have a `cludod.yaml` file in `~/.cludod` then a local copy of the server can be built and spun up with:

``` bash
docker compose up --build -d
```

### Running the CLI locally

You can run the CLI in a container using something like this:

``` bash
docker run -ti \
  -v ${PWD}/data/cludo.yaml:/root/.cludo/cludo.yaml \
  -v ${PWD}/data/id_rsa_TEST_KEY_PP:/app/id_rsa_TEST_KEY_PP \
  ${IMAGE_ID} \
  "cludo shell"
```

### Release

- Checkout the `main` branch
- `make all docker-push`
- Create a Github release
- Attach the binaries for all platforms to the release
- List the fully qualified image tags in the release description.

**TODO**: Automate the creation of a release in Github with the resulting binaries.
