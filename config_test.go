package vaulthcplib

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
			tkHelper := &TestingHCPTokenHelper{validCache: tst.Valid}
			tk, err := tkHelper.GetHCPToken()

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

	tkHelper := &TestingHCPTokenHelper{}
	_, err = tkHelper.GetHCPToken()
	assert.NoError(t, err)
}
