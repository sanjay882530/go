package scraper

import (
	"testing"
	"time"

	auroraclient "github.com/hcnet/go/clients/auroraclient"
	hProtocol "github.com/hcnet/go/protocols/aurora"
	"github.com/hcnet/go/support/errors"
	"github.com/hcnet/go/support/log"
	"github.com/stretchr/testify/assert"
)

func Test_ScraperConfig_FetchAllTrades_doesntCrashWhenReceivesAnError(t *testing.T) {
	auroraClient := &auroraclient.MockClient{}
	auroraClient.
		On("Trades", auroraclient.TradeRequest{Limit: 200, Order: auroraclient.OrderDesc}).
		Return(hProtocol.TradesPage{}, errors.New("something went wrong"))

	sc := ScraperConfig{
		Logger: log.DefaultLogger,
		Client: auroraClient,
	}

	trades, err := sc.FetchAllTrades(time.Now(), 0)
	assert.EqualError(t, err, "something went wrong")
	assert.Empty(t, trades)
}

func Test_ScraperConfig_FetchAllTrades_doesntCrashWhenReceivesEmptyList(t *testing.T) {
	auroraClient := &auroraclient.MockClient{}
	auroraClient.
		On("Trades", auroraclient.TradeRequest{Limit: 200, Order: auroraclient.OrderDesc}).
		Return(hProtocol.TradesPage{}, nil)

	sc := ScraperConfig{
		Logger: log.DefaultLogger,
		Client: auroraClient,
	}

	trades, err := sc.FetchAllTrades(time.Now(), 0)
	assert.NoError(t, err)
	assert.Empty(t, trades)
}
