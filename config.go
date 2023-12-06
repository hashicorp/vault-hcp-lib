// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package vaulthcplib

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/hcp-sdk-go/config"
)

const (
	defaultDirectory     = ".config/hcp/hvd"
	testDirectory        = "hcptest"
	fileName             = "hvd_proxy_config.json"
	directoryPermissions = 0755

	envVarCacheTestMode = "HCP_CACHE_TEST_MODE"
)

type HCPToken struct {
	AccessToken       string    `json:"access_token,omitempty"`
	AccessTokenExpiry time.Time `json:"access_token_expiry,omitempty"`
	ProxyAddr         string    `json:"proxy_addr,omitempty"`
}

type HCPTokenHelper interface {
	GetHCPToken() (*HCPToken, error)
}

var _ HCPTokenHelper = (*InternalHCPTokenHelper)(nil)

type InternalHCPTokenHelper struct{}

func (h InternalHCPTokenHelper) GetHCPToken() (*HCPToken, error) {
	configCache, err := readConfig()
	if err != nil {
		return nil, err
	}
	// no valid connection to hcp
	if configCache == nil || configCache.ProxyAddr == "" {
		return nil, nil
	}

	opts := []config.HCPConfigOption{
		config.WithoutLogging(),
		config.WithoutBrowserLogin(),
		config.FromEnv(),
	}
	if configCache.ClientID != "" && configCache.SecretID != "" {
		opts = append(opts, config.WithClientCredentials(configCache.ClientID, configCache.SecretID))
	}
	hcp, err := config.NewHCPConfig(opts...)

	tk, err := hcp.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve the HCP token: %w", err)
	}

	return &HCPToken{
		AccessToken:       tk.AccessToken,
		AccessTokenExpiry: tk.Expiry,
		ProxyAddr:         configCache.ProxyAddr,
	}, nil
}

var _ HCPTokenHelper = (*TestingHCPTokenHelper)(nil)

type TestingHCPTokenHelper struct {
	ValidCache bool
}

func (h TestingHCPTokenHelper) GetHCPToken() (*HCPToken, error) {
	userHome, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	credentialDir := filepath.Join(userHome, testDirectory)
	err = os.RemoveAll(credentialDir)
	if err != nil {
		return nil, err
	}

	if h.ValidCache {
		err = writeConfig("https://hcp-proxy.addr:8200", "", "")
		if err != nil {
			return nil, err
		}

		configCache, err := readConfig()
		if err != nil {
			return nil, err
		}
		if configCache == nil {
			return nil, nil
		}

		tkSrc := &TestTokenSource{}
		tk, _ := tkSrc.Token()

		return &HCPToken{
			AccessToken:       tk.AccessToken,
			AccessTokenExpiry: tk.Expiry,
			ProxyAddr:         configCache.ProxyAddr,
		}, nil
	}

	return nil, nil
}

type HCPConfigCache struct {
	ClientID  string
	SecretID  string
	ProxyAddr string
}

// Write saves HCP auth data in a common location in the home directory.
func writeConfig(addr string, clientID string, secretID string) error {
	credentialPath, credentialDirectory, err := getConfigPaths()
	if err != nil {
		return fmt.Errorf("failed to retrieve credential path and directory: %v", err)
	}

	err = os.MkdirAll(credentialDirectory, directoryPermissions)
	if err != nil {
		return fmt.Errorf("failed to create credential directory: %v", err)
	}

	cache := &HCPConfigCache{
		ClientID:  clientID,
		SecretID:  secretID,
		ProxyAddr: addr,
	}
	cacheJSON, err := json.Marshal(cache)
	if err != nil {
		return fmt.Errorf("failed to marshal the struct to json: %v", err)
	}

	err = os.WriteFile(credentialPath, cacheJSON, directoryPermissions)
	if err != nil {
		return fmt.Errorf("failed to write config to the cache file: %v", err)
	}

	return nil
}

// readConfig opens the saved HCP auth data and returns the token.
func readConfig() (*HCPConfigCache, error) {
	configPath, _, err := getConfigPaths()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve config path and directory: %v", err)
	}

	var cache HCPConfigCache
	if _, err = os.Stat(configPath); os.IsNotExist(err) {
		return nil, nil
	}

	rawJSON, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file from user's config path: %v", err)
	}

	err = json.Unmarshal(rawJSON, &cache)
	if err != nil {
		return nil, err
	}
	return &cache, nil
}

func eraseConfig() error {
	_, credentialDirectory, err := getConfigPaths()
	if err != nil {
		return fmt.Errorf("failed to retrieve credential path and directory: %v", err)
	}

	err = os.RemoveAll(credentialDirectory)
	if err != nil {
		return fmt.Errorf("failed to remove config directory: %v", err)
	}

	return nil
}

// getCredentialPaths returns the complete credential path and directory.
func getConfigPaths() (configPath string, configDirectory string, err error) {
	// Get the user's home directory.
	userHome, err := os.UserHomeDir()
	if err != nil {
		return "", "", fmt.Errorf("failed to retrieve user's home directory path: %v", err)
	}

	directoryName := defaultDirectory
	// If in test mode, use test directory.
	if testMode, ok := os.LookupEnv(envVarCacheTestMode); ok {
		if testMode == "true" {
			directoryName = testDirectory
		}
	}

	// Determine absolute path to config file and directory.
	configDirectory = filepath.Join(userHome, directoryName)
	configPath = filepath.Join(userHome, directoryName, fileName)

	return configPath, configDirectory, nil
}
