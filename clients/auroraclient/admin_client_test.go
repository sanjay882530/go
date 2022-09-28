package auroraclient

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultAdminHostPort(t *testing.T) {
	auroraAdminClient, err := NewAdminClient(0, "", 0)

	fullAdminURL := auroraAdminClient.getIngestionFiltersURL("test")
	require.NoError(t, err)
	assert.Equal(t, "http://localhost:4200/ingestion/filters/test", fullAdminURL)
}

func TestOverrideAdminHostPort(t *testing.T) {
	auroraAdminClient, err := NewAdminClient(1234, "127.0.0.1", 0)

	fullAdminURL := auroraAdminClient.getIngestionFiltersURL("test")
	require.NoError(t, err)
	assert.Equal(t, "http://127.0.0.1:1234/ingestion/filters/test", fullAdminURL)
}
