# credly-go

[![Go Report Card](https://goreportcard.com/badge/github.com/isovalent/credly-go)](https://goreportcard.com/report/github.com/isovalent/credly-go)
[![Go Reference](https://pkg.go.dev/badge/github.com/isovalent/credly-go.svg)](https://pkg.go.dev/github.com/isovalent/credly-go)
[![License: Apache 2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

`credly-go` is a Go client library for interacting with the Credly platform. It provides a simple and convenient way to programmatically access Credly's APIs and handle badges and templates.

## Features

- **Badge Management**: Issue, retrieve, and manage badges using the Credly API.

## Installation

To install the `credly-go` library, run:

```shell
go get github.com/isovalent/credly-go
```


## Example Usage


```go
package main

import (
    "github.com/isovalent/credly-go/credly"
)

func main() {
    // Initialize the Credly client
    client := credly.NewClient("your-api-token", "your-credly-org")

    // Get all badges for user joe@example.com
    badges, err := client.GetBadges("joe@example.com")
}
```

## Contributing

We welcome contributions! Please follow these steps to contribute:

1. Fork the repository.
2. Create a new branch with your feature or bug fix.
3. Make your changes and add tests.
4. Submit a pull request with a detailed description of your changes.

## Running Tests

To run the tests, use:

```shell
go test ./...
```


Make sure to write tests for any new functionality and ensure that all existing tests pass.

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Support

If you have any questions or need help, feel free to open an issue in the [GitHub repository](https://github.com/isovalent/credly-go/issues).
