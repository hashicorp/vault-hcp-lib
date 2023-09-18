package hcpvaultengine

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetHCPConfiguration(t *testing.T) {
	cases := map[string]struct {
		Valid bool
	}{
		"valid hcp configuration": {
			Valid: true,
		},
		"empty hcp configuration": {
			Valid: false,
		},
	}

	for n, tst := range cases {
		t.Run(n, func(t *testing.T) {
			CacheSetup(t, tst.Valid)

			tk, err := GetHCPToken()

			assert.NoError(t, err)

			if tst.Valid {
				assert.Equal(t, "https://hcp-proxy.addr:8200", tk.ProxyAddr)
				assert.Contains(t, tk.AccessToken, "Test.Access.Token")
				assert.NotEmpty(t, tk.AccessTokenExpiry)
			} else {
				assert.Nil(t, tk)
				assert.Nil(t, err)
			}

		})
	}
}

func Test_GetHCPConfiguration_EraseConfig(t *testing.T) {
	err := os.Setenv(envVarCacheTestMode, "true")
	assert.NoError(t, err)

	err = eraseConfig()
	assert.NoError(t, err)

	_, err = GetHCPToken()
	assert.NoError(t, err)
}
