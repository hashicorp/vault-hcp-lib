package vaulthcplib

import (
	"github.com/hashicorp/hcp-sdk-go/auth"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func CacheSetup(t *testing.T, validCache bool) {
	err := os.Setenv(envVarCacheTestMode, "true")
	if err != nil {
		t.Error(err)
		return
	}

	userHome, err := os.UserHomeDir()
	if err != nil {
		t.Error(err)
		return
	}
	credentialDir := filepath.Join(userHome, testDirectory)
	err = os.RemoveAll(credentialDir)
	if err != nil {
		t.Error(err)
		return
	}

	if validCache {
		now := time.Now()
		cache := auth.Cache{
			AccessToken:       "Test.Access.Token",
			RefreshToken:      "TestRefreshToken",
			AccessTokenExpiry: now.Add(time.Hour * 2),
			SessionExpiry:     time.Now().Add(time.Hour * 24),
		}
		err = auth.Write(cache)
		if err != nil {
			t.Error(err)
		}

		err = writeConfig("https://hcp-proxy.addr:8200", "", "")
		if err != nil {
			t.Error(err)
		}
	}

	return
}
