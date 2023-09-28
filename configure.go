package hcpvaultengine

import (
	"github.com/hashicorp/vault/api"
	"golang.org/x/oauth2"
)

// TODO: reevaluate this cache strategy
// Disk? Memory? Memdb?
var cache = &HCPVClusterCache{}

type HCPVClusterCache struct {
	// Memory cache of the token source
	Source oauth2.TokenSource

	// Memory cache of the cluster address
	Address string

	// Memory cache of the cluster ID
	ID string
}

// ConfigureHCPProxy adds a client-side middleware, an implementation of http.RoundTripper on top of the base transport,
// that will add a cookie to every request made from the CLI client. Additionally, it overrides the configuration's address
// The address will be that of the proxy by default and the cookie will have the HCP access token data necessary to make requests to
// the cluster through HCP.
//
// TODO: is there a better way to change the configuration without parametizing the Vault Config?
func ConfigureHCPProxy(client *api.Client) error {
	// TODO: reevaluate this. Which scheme? https?
	addr := "https://" + cache.Address
	err := client.SetAddress(addr)
	if err != nil {
		return err
	}

	// TODO: understand and reevaluate exactly what it means to get the token from the source or to get the TokenSource.
	token, err := cache.Source.Token()
	if err != nil {
		return err
	}

	client.SetHCPToken(token)

	return nil
}
