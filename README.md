# WasaText

WasaText is a university project developed for the Web and Software Architecture course at Sapienza University of Rome.

This project aims to explore and implement concepts of software architecture and web development by creating a text messaging application.

**Note:** This project is developed for educational purposes and may not be suitable for production environments.

## Project Structure

* `cmd/`: Contains executables. Go programs here handle execution-related tasks, such as reading options from the CLI/environment.
    * `cmd/webapi`: Contains the web API server.
* `doc/`: Contains documentation, including OpenAPI files for APIs.
* `service/`: Contains packages for implementing project-specific functionalities.
    * `service/api`: Contains the API server code.
* `vendor/`: Managed by Go, contains a copy of all dependencies.
* `webui/`: Contains the web frontend written in Vue.js.

Other project files include:

* `open-node.sh`: Starts a temporary container using the `node:20` image for safe web frontend development.

## Go Vendoring

This project uses [Go Vendoring](https://go.dev/ref/mod#vendoring). After modifying dependencies (`go get` or `go mod tidy`), you must run `go mod vendor` and include all files in the `vendor/` directory in your commit.

## Node/YARN Vendoring

This project uses `yarn` and a vendoring technique that exploits the ["Offline mirror"](https://yarnpkg.com/features/caching). Dependencies are located within the repository, in the `.yarn` directory.

## How to build

If you're not using the WebUI, or if you don't want to embed the WebUI into the final executable, then:

```shell
go build ./cmd/webapi/
```

If you're using the WebUI and you want to embed it into the final executable:

```shell
./open-node.sh
# (here you're inside the container)
yarn run build-embed
exit
# (outside the container)
go build -tags webui ./cmd/webapi/
```

## How to run (in development mode)

You can launch the backend only using:

```shell
go run ./cmd/webapi/
```

If you want to launch the WebUI, open a new tab and launch:

```shell
./open-node.sh
# (here you're inside the container)
yarn run dev
```

## License

See [LICENSE](LICENSE).
