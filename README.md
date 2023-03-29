# DevCycle Terraform Provider

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.17

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.


## Testing  the Provider
Tests use the Hashicorp [Terraform Acceptance Tests](https://developer.hashicorp.com/terraform/plugin/sdkv2/testing/acceptance-tests).
The test suite also requires the correct DevCycle ids and secrets to run.

In order to run the full suite of Acceptance tests locally run:
```shell
make testacc DEVCYCLE_CLIENT_ID=<id> DEVCYCLE_CLIENT_SECRET=<secret> DEVCYCLE_SERVER_TOKEN=<token>
```

*Note:* Acceptance tests create real resources, and often cost money to run.
