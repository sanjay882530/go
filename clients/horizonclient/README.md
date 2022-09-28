# auroraclient


`auroraclient` is a [Hcnet Go SDK](https://developers.hcnet.org/api/) package that provides client access to a aurora server. It supports all endpoints exposed by the [aurora API](https://developers.hcnet.org/api/introduction/).

This project is maintained by the Hcnet Development Foundation.

## Getting Started
This library is aimed at developers building Go applications that interact with the [Hcnet network](https://www.hcnet.org/). It allows users to query the network and submit transactions to the network. The recommended transaction builder for Go programmers is [txnbuild](https://github.com/hcnet/go/tree/master/txnbuild). Together, these two libraries provide a complete Hcnet SDK.

* The [auroraclient API reference](https://godoc.org/github.com/hcnet/go/clients/auroraclient).
* The [txnbuild API reference](https://godoc.org/github.com/hcnet/go/txnbuild).

### Prerequisites
* Go (this repository is officially supported on the last two releases of Go)
* [Modules](https://github.com/golang/go/wiki/Modules) to manage dependencies

### Installing
* `go get github.com/hcnet/go/clients/auroraclient`

### Usage

``` golang
    ...
    import hClient "github.com/hcnet/go/clients/auroraclient"
    ...

    // Use the default pubnet client
    client := hClient.DefaultPublicNetClient

    // Create an account request
    accountRequest := hClient.AccountRequest{AccountID: "GCLWGQPMKXQSPF776IU33AH4PZNOOWNAWGGKVTBQMIC5IMKUNP3E6NVU"}

    // Load the account detail from the network
    account, err := client.AccountDetail(accountRequest)
    if err != nil {
        fmt.Println(err)
        return
    }
    // Account contains information about the hcnet account
    fmt.Print(account)
```
For more examples, refer to the [documentation](https://godoc.org/github.com/hcnet/go/clients/auroraclient).

## Running the tests
Run the unit tests from the package directory: `go test`

## Contributing
Please read [Code of Conduct](https://github.com/hcnet/.github/blob/master/CODE_OF_CONDUCT.md) to understand this project's communication rules.

To submit improvements and fixes to this library, please see [CONTRIBUTING](../CONTRIBUTING.md).

## License
This project is licensed under the Apache License - see the [LICENSE](../../LICENSE) file for details.
