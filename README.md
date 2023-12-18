# Snake

## Description

Snake is a Go library designed to help build tools faster. It combines the binding logic of the Kong CLI library with the Cobra CLI library, providing a powerful and efficient way to create command-line interfaces.

## Installation

To install the Snake library, use the `go get` command:

```bash
go get github.com/yourusername/snake
```

## Usage

Define a new command by creating a new struct that implements the `snake.Command` interface:

```go
type MyCommand struct {
    // ...
}
```