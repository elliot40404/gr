# gr

A simple go tool to run programs in the `cmd/` directory. All it does is run `go run cmd/program/main.go` for you.

- [x] Supports specifying a binary to run with the `--bin` flag if you have multiple directories in `cmd/`.
- [x] Also can pass arguments. 

## Why?
Because I'm lazy and I don't want to type `go run cmd/program/main.go` every time I want to run a program.

## Installation
```console
$ go install github.com/elliot40404/gr@latest
```

## Usage


```console
$ gr
```

```console
$ gr --bin <program>
```

# LICENSE
MIT

