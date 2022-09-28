package main

import (
	"fmt"

	"github.com/hcnet/go/ingest/ledgerbackend"
)

var (
	config = captiveCoreConfig()
)

func captiveCoreConfig() ledgerbackend.CaptiveCoreConfig {
	archiveURLs := []string{
		"https://history.hcnet.org/prd/core-testnet/core_testnet_001",
		"https://history.hcnet.org/prd/core-testnet/core_testnet_002",
		"https://history.hcnet.org/prd/core-testnet/core_testnet_003",
	}
	networkPassphrase := "Test SDF Network ; September 2015"
	captiveCoreToml, err := ledgerbackend.NewCaptiveCoreToml(ledgerbackend.CaptiveCoreTomlParams{
		NetworkPassphrase:  networkPassphrase,
		HistoryArchiveURLs: archiveURLs,
	})
	panicIf(err)

	return ledgerbackend.CaptiveCoreConfig{
		// Change these based on your environment:
		BinaryPath:         "/usr/local/bin/hcnet-core",
		NetworkPassphrase:  networkPassphrase,
		HistoryArchiveURLs: archiveURLs,
		Toml:               captiveCoreToml,
	}
}

func panicIf(err error) {
	if err != nil {
		panic(fmt.Errorf("An error occurred, panicking: %s\n", err))
	}
}
