package cmd

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	auroraclient "github.com/hcnet/go/clients/auroraclient"
	hlog "github.com/hcnet/go/support/log"
)

var DatabaseURL string
var Client *auroraclient.Client
var UseTestNet bool
var Logger = hlog.New()

var defaultDatabaseURL = getEnv("DB_URL", "postgres://localhost:5432/hcnetticker01?sslmode=disable")

var rootCmd = &cobra.Command{
	Use:   "ticker",
	Short: "Hcnet Development Foundation Ticker.",
	Long:  `A tool to provide Hcnet Asset and Market data.`,
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(
		&DatabaseURL,
		"db-url",
		"d",
		defaultDatabaseURL,
		"database URL, such as: postgres://user:pass@localhost:5432/ticker",
	)
	rootCmd.PersistentFlags().BoolVar(
		&UseTestNet,
		"testnet",
		false,
		"use the Hcnet Test Network, instead of the Hcnet Public Network",
	)

	Logger.SetLevel(logrus.DebugLevel)
}

func initConfig() {
	if UseTestNet {
		Logger.Debug("Using Hcnet Default Test Network")
		Client = auroraclient.DefaultTestNetClient
	} else {
		Logger.Debug("Using Hcnet Default Public Network")
		Client = auroraclient.DefaultPublicNetClient
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
