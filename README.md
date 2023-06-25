# Goshelf

Goshelf is a combined CLI/REST API server that allows you to save your favorite books and collections of books to a database.

## Usage

To use, run go from the command line:

`go run ./cmd/main [FLAGS] [CLI Function]`

Goshelf has two modes of operation: CLI and API.

### Dockerfile

Goshelf can be ran within docker for extra portability.

`docker run --rm -it $(docker build -q -f dockerfile .) app -a`


## Building

To build the application

`go build ./cmd/main`

## Flags

Goshelf utilizes command flags to configure itself. The following flags are available:

```
-a Run in API mode, default false
-dh Database address, default 0.0.0.0
-dn Database name, default postgres
-dp Database port, default 5432
-dpw Database password, default ""
-ds Database SSL mode, default postgres
-du Database user, default postgres
-p API mode: Host port, default 8080
-s API mode: Host address, default 0.0.0.0
```

### API

To run in API mode, the application must have the `-a` flag set. If done properly Stdout will output the address the API is hosted on.

### CLI

To run in CLI mode, an argument must be passed that matches the API being used. 

`go run ./... BookDelete -dpw "mydatabasepassword"`...

The CLI will walk the user through each step.