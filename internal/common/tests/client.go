package tests

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"toto-config-api/internal/common/client/skuconfig"
)

type SKUConfigHTTPClient struct {
	client *skuconfig.ClientWithResponses
}

func NewSKUConfigHTTPClient(t *testing.T) SKUConfigHTTPClient {
	addr := os.Getenv("SKUCONFIG_HTTP_ADDR")
	ok := WaitForPort(addr)

	require.True(t, ok, "SKUConfig HTTP timed out")

	url := fmt.Sprintf("http://%v/api", addr)

	client, err := skuconfig.NewClientWithResponses(url)
	require.NoError(t, err)

	return SKUConfigHTTPClient{client: client}
}

func (c SKUConfigHTTPClient) GetSKU(t *testing.T, pkg string, country string) string {
	return ""
}
