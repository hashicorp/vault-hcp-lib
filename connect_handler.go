package hcpvaultengine

import (
	"errors"
	"flag"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/hashicorp/hcp-sdk-go/config"
	"github.com/hashicorp/hcp-sdk-go/httpclient"
	"golang.org/x/oauth2"
)

type AuthStatus struct {
	Token oauth2.TokenSource
}

var (
	_ ConnectHandler = (*interactiveConnectHandler)(nil)
	_ ConnectHandler = (*nonInteractiveConnectHandler)(nil)
)

type ConnectHandler interface {
	Connect(args []string) (*httptransport.Runtime, error)
}

type interactiveConnectHandler struct{}

func (h *interactiveConnectHandler) Connect(_ []string) (*httptransport.Runtime, error) {
	// Start a callback listener
	// Define timeout: 2 minutes?
	return nil, nil
}

type nonInteractiveConnectHandler struct {
	flagClientID string
	flagSecretID string
}

func (h *nonInteractiveConnectHandler) Connect(args []string) (*httptransport.Runtime, error) {
	f := h.Flags()

	if err := f.Parse(args); err != nil {
		return nil, err
	}

	if h.flagClientID == "" || h.flagSecretID == "" {
		return nil, errors.New("client ID and Secret ID need to be set in non-interactive mode")
	}

	opts := []config.HCPConfigOption{config.FromEnv()}
	opts = append(opts, config.WithClientCredentials(h.flagClientID, h.flagSecretID))
	opts = append(opts, config.WithoutBrowserLogin())

	cfg, err := config.NewHCPConfig(opts...)
	if err != nil {
		return nil, err
	}

	// Cache token source
	cache.Source = cfg

	hcpClient, err := httpclient.New(httpclient.Config{HCPConfig: cfg})
	if err != nil {
		return nil, err
	}

	return hcpClient, nil
}

func (h *nonInteractiveConnectHandler) Flags() *flag.FlagSet {
	mainSet := flag.NewFlagSet("", flag.ContinueOnError)

	mainSet.Var(&stringValue{target: &h.flagClientID}, "client-id", "")
	mainSet.Var(&stringValue{target: &h.flagSecretID}, "secret-id", "")

	return mainSet
}
