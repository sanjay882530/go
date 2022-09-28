---
title: Overview
---

The Go SDK is a set of packages for interacting with most aspects of the Hcnet ecosystem. The primary component is the Aurora SDK, which provides convenient access to Aurora services. There are also packages for other Hcnet services such as [TOML support](https://github.com/hcnet/hcnet-protocol/blob/master/ecosystem/sep-0001.md) and [federation](https://github.com/hcnet/hcnet-protocol/blob/master/ecosystem/sep-0002.md).

## Aurora SDK

The Aurora SDK is composed of two complementary libraries: `txnbuild` + `auroraclient`.
The `txnbuild` ([source](https://github.com/hcnet/go/tree/master/txnbuild), [docs](https://godoc.org/github.com/hcnet/go/txnbuild)) package enables the construction, signing and encoding of Hcnet [transactions](https://developers.hcnet.org/docs/glossary/transactions/) and [operations](https://developers.hcnet.org/docs/start/list-of-operations/) in Go. The `auroraclient` ([source](https://github.com/hcnet/go/tree/master/clients/auroraclient), [docs](https://godoc.org/github.com/hcnet/go/clients/auroraclient)) package provides a web client for interfacing with [Aurora](https://developers.hcnet.org/docs/start/introduction/) server REST endpoints to retrieve ledger information, and to submit transactions built with `txnbuild`.

## List of major SDK packages

- `auroraclient` ([source](https://github.com/hcnet/go/tree/master/clients/auroraclient), [docs](https://godoc.org/github.com/hcnet/go/clients/auroraclient)) - programmatic client access to Aurora
- `txnbuild` ([source](https://github.com/hcnet/go/tree/master/txnbuild), [docs](https://godoc.org/github.com/hcnet/go/txnbuild)) - construction, signing and encoding of Hcnet transactions and operations
- `hcnettoml` ([source](https://github.com/hcnet/go/tree/master/clients/hcnettoml), [docs](https://godoc.org/github.com/hcnet/go/clients/hcnettoml)) - parse [Hcnet.toml](../../guides/concepts/hcnet-toml.md) files from the internet
- `federation` ([source](https://godoc.org/github.com/hcnet/go/clients/federation)) - resolve federation addresses  into hcnet account IDs, suitable for use within a transaction

