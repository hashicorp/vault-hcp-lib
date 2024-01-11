// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package vaulthcplib

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	hcprmm "github.com/hashicorp/hcp-sdk-go/clients/cloud-resource-manager/stable/2019-12-10/models"
	hcpvsm "github.com/hashicorp/hcp-sdk-go/clients/cloud-vault-service/stable/2020-11-25/models"

	"github.com/hashicorp/cli"
	hcprmo "github.com/hashicorp/hcp-sdk-go/clients/cloud-resource-manager/stable/2019-12-10/client/organization_service"
	hcprmp "github.com/hashicorp/hcp-sdk-go/clients/cloud-resource-manager/stable/2019-12-10/client/project_service"
	hcpvs "github.com/hashicorp/hcp-sdk-go/clients/cloud-vault-service/stable/2020-11-25/client/vault_service"
	"github.com/hashicorp/hcp-sdk-go/config"
	"github.com/hashicorp/hcp-sdk-go/httpclient"
)

var (
	_ cli.Command = (*HCPConnectCommand)(nil)

	ErrorProxyDisabled = errors.New("proxy is disabled")
)

type HCPConnectCommand struct {
	Ui cli.Ui

	flagClientID       string
	flagSecretID       string
	flagOrganizationID string
	flagProjectID      string
	flagClusterID      string

	// for testing
	rmOrgClient  hcprmo.ClientService
	vsClient     hcpvs.ClientService
	rmProjClient hcprmp.ClientService
}

func (c *HCPConnectCommand) Help() string {
	helpText := `
Usage: vault hcp connect [options]
  
  Authenticates users or machines to HCP using either provided arguments or retrieved token through
  browser login. A successful authentication results in an HCP token and an HCP Vault address being
  locally cached. 

  The default authentication method is an interactive one, redirecting users to the HCP login browser.
  If a set of service principal credential is supplied, which is generated through the HCP portal with 
  the necessary capabilities to access the organization, project, and HCP Vault cluster chosen.

      $ vault hcp connect -client-id=client-id-value -secret-id=secret-id-value
  
  Additionally, the organization identification, project identification, and cluster name can be passed in to
  directly connect to a specific HCP Vault cluster without interacting with the CLI.
  
      $ vault hcp connect -client-id=client-id-value -secret-id=secret-id-value -organization-id=org-UUID -project-id=proj-UUID -cluster-id=cluster-name
`
	return strings.TrimSpace(helpText)
}

func (c *HCPConnectCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	if err := c.setupClients(); err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	proxyAddr, err := c.getProxyAddr(c.rmOrgClient, c.rmProjClient, c.vsClient)
	if err != nil {
		if errors.Is(err, ErrorProxyDisabled) {
			c.Ui.Error("\nFailed to connect to HCP Vault Cluster: HTTP proxy feature not enabled.")
			return 1
		}
		c.Ui.Error(fmt.Sprintf("\n%s", err.Error()))
		return 1
	}

	err = writeConfig(proxyAddr, c.flagClientID, c.flagSecretID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("\nFailed to connect to HCP Vault Cluster: %s", err))
		return 1
	}

	c.Ui.Info("\nConnected to cluster via HCP proxy. Login with 'vault login' or export a VAULT_TOKEN to access this Vault cluster.")
	return 0
}

func (c *HCPConnectCommand) setupClients() error {
	var opts []config.HCPConfigOption

	if c.rmOrgClient == nil && c.rmProjClient == nil && c.vsClient == nil {
		opts = []config.HCPConfigOption{config.FromEnv()}

		if c.flagClientID != "" && c.flagSecretID == "" {
			return errors.New("secret-id is required when client-id is provided")
		} else if c.flagSecretID != "" && c.flagClientID == "" {
			return errors.New("client-id is required when secret-id is provided")
		} else if c.flagClientID != "" && c.flagSecretID != "" {
			opts = append(opts, config.WithClientCredentials(c.flagClientID, c.flagSecretID))
			opts = append(opts, config.WithoutBrowserLogin())
		}

		cfg, err := config.NewHCPConfig(opts...)
		if err != nil {
			return errors.New(fmt.Sprintf("Failed to connect to HCP: %s", err))
		}

		hcpHttpClient, err := httpclient.New(httpclient.Config{HCPConfig: cfg})
		if err != nil {
			return errors.New(fmt.Sprintf("Failed to connect to HCP: %s", err))
		}

		c.rmOrgClient = hcprmo.New(hcpHttpClient, nil)
		c.rmProjClient = hcprmp.New(hcpHttpClient, nil)
		c.vsClient = hcpvs.New(hcpHttpClient, nil)
	}

	return nil
}

func (c *HCPConnectCommand) getProxyAddr(organizationClient hcprmo.ClientService, projectClient hcprmp.ClientService, clusterClient hcpvs.ClientService) (string, error) {
	var err error

	var organizationID string
	if c.flagOrganizationID != "" {
		organizationID = c.flagOrganizationID
	} else {
		organizationID, err = c.getOrganization(organizationClient)
		if err != nil {
			return "", errors.New(fmt.Sprintf("Failed to get HCP organization information: %s", err))
		}
	}

	var projectID string
	if c.flagProjectID != "" {
		projectID = c.flagProjectID
	} else {
		projectID, err = c.getProject(organizationID, projectClient)
		if err != nil {
			return "", errors.New(fmt.Sprintf("Failed to get HCP project information: %s", err))
		}
	}

	proxyAddr, err := c.getCluster(organizationID, projectID, c.flagClusterID, clusterClient)
	if err != nil {
		return "", err
	}
	return proxyAddr, nil
}

func (c *HCPConnectCommand) Synopsis() string {
	return "Connect to an HCP Vault Cluster"
}

func (c *HCPConnectCommand) Flags() *flag.FlagSet {
	mainSet := flag.NewFlagSet("", flag.ContinueOnError)

	mainSet.StringVar(&c.flagClientID, "client-id", "", "")
	mainSet.StringVar(&c.flagSecretID, "secret-id", "", "")
	mainSet.StringVar(&c.flagOrganizationID, "organization-id", "", "")
	mainSet.StringVar(&c.flagProjectID, "project-id", "", "")
	mainSet.StringVar(&c.flagClusterID, "cluster-id", "", "")

	return mainSet
}

func (c *HCPConnectCommand) getOrganization(rmOrgClient hcprmo.ClientService) (organizationID string, err error) {
	organizationsResp, err := rmOrgClient.OrganizationServiceList(hcprmo.NewOrganizationServiceListParams().WithDefaults(), nil)
	switch {
	case err != nil:
		return "", err
	case organizationsResp.GetPayload() == nil:
		return "", errors.New("payload is nil")
	case len(organizationsResp.GetPayload().Organizations) < 1:
		return "", errors.New("no organizations available")
	case len(organizationsResp.GetPayload().Organizations) > 1:
		title := "Available organizations:"
		u := strings.Repeat("-", len(title))
		c.Ui.Info(fmt.Sprintf("%s\n%s\n", u, title))

		orgs := make(map[string]*hcprmm.HashicorpCloudResourcemanagerOrganization, len(organizationsResp.GetPayload().Organizations))
		for _, org := range organizationsResp.GetPayload().Organizations {
			if *org.State == hcprmm.HashicorpCloudResourcemanagerOrganizationOrganizationStateACTIVE {
				c.Ui.Info(fmt.Sprintf("Organization name: %s", org.Name))
				name := strings.ToLower(org.Name)
				orgs[name] = org
			}
		}
		userInput, err := c.Ui.Ask(fmt.Sprintf("\nChoose a organization: "))
		if err != nil {
			return "", err
		}
		chosenOrg, ok := orgs[userInput]
		if !ok {
			return "", errors.New(fmt.Sprintf("invalid HCP organization: %s", userInput))
		}
		return chosenOrg.ID, nil
	default:
		organization := organizationsResp.GetPayload().Organizations[0]
		if *organization.State != hcprmm.HashicorpCloudResourcemanagerOrganizationOrganizationStateACTIVE {
			return "", errors.New("organization is not active")
		}
		return organization.ID, nil
	}
}

func (c *HCPConnectCommand) getProject(organizationID string, rmProjClient hcprmp.ClientService) (projectID string, err error) {
	scopeType := "ORGANIZATION"
	projectListReq := hcprmp.
		NewProjectServiceListParams().
		WithDefaults().
		WithScopeType(&scopeType).
		WithScopeID(&organizationID)
	projectResp, err := rmProjClient.ProjectServiceList(projectListReq, nil)
	switch {
	case err != nil:
		return "", err
	case projectResp.GetPayload() == nil:
		return "", errors.New("payload is nil")
	case len(projectResp.GetPayload().Projects) < 1:
		return "", errors.New("no projects available")
	case len(projectResp.GetPayload().Projects) > 1:
		title := "Available projects:"
		u := strings.Repeat("-", len(title))
		c.Ui.Info(fmt.Sprintf("%s\n%s\n", u, title))

		projs := make(map[string]*hcprmm.HashicorpCloudResourcemanagerProject, len(projectResp.GetPayload().Projects))
		for _, proj := range projectResp.GetPayload().Projects {
			if *proj.State == hcprmm.HashicorpCloudResourcemanagerProjectProjectStateACTIVE {
				c.Ui.Info(fmt.Sprintf("Project name: %s", proj.Name))
				name := strings.ToLower(proj.Name)
				projs[name] = proj
			}
		}
		userInput, err := c.Ui.Ask(fmt.Sprintf("\nChoose a project: "))
		if err != nil {
			return "", err
		}
		chosenProj, ok := projs[userInput]
		if !ok {
			return "", errors.New(fmt.Sprintf("invalid HCP project: %s", userInput))
		}
		return chosenProj.ID, nil
	default:
		project := projectResp.GetPayload().Projects[0]
		if *project.State != hcprmm.HashicorpCloudResourcemanagerProjectProjectStateACTIVE {
			return "", errors.New("project is not active")
		}
		return project.ID, nil
	}
}

func (c *HCPConnectCommand) getCluster(organizationID string, projectID string, clusterID string, vsClient hcpvs.ClientService) (proxyAddr string, err error) {
	if clusterID == "" {
		return c.listClusters(organizationID, projectID, vsClient)
	}

	clusterGetReq := hcpvs.NewGetParams().
		WithDefaults().
		WithLocationOrganizationID(organizationID).
		WithLocationProjectID(projectID).
		WithClusterID(clusterID)
	clusterResp, err := vsClient.Get(clusterGetReq, nil)
	switch {
	case err != nil:
		return "", err
	case clusterResp.GetPayload() == nil:
		return "", errors.New("payload is nil")
	default:
		cluster := clusterResp.GetPayload().Cluster

		if *cluster.Config.NetworkConfig.HTTPProxyOption == hcpvsm.HashicorpCloudVault20201125HTTPProxyOptionDISABLED {
			return "", ErrorProxyDisabled
		}

		title := "HCP Vault Cluster:"
		u := strings.Repeat("-", len(title))
		c.Ui.Info(fmt.Sprintf("%s\n%s: %s\n", u, title, cluster.ID))

		proxyAddr = "https://" + cluster.DNSNames.Proxy
		return proxyAddr, nil
	}
}

func (c *HCPConnectCommand) listClusters(organizationID string, projectID string, vsClient hcpvs.ClientService) (proxyAddr string, err error) {
	clusterListReq := hcpvs.NewListParams().
		WithDefaults().
		WithLocationOrganizationID(organizationID).
		WithLocationProjectID(projectID)

	// Purposely calling List instead of ListAll because we are only interested in HVD clusters.
	clustersResp, err := vsClient.List(clusterListReq, nil)
	switch {
	case err != nil:
		return "", err
	case clustersResp.GetPayload() == nil:
		return "", errors.New("payload is nil")
	case len(clustersResp.GetPayload().Clusters) < 1:
		return "", errors.New("no clusters available")
	case len(clustersResp.GetPayload().Clusters) > 1:
		title := "Available clusters:"
		u := strings.Repeat("-", len(title))
		c.Ui.Info(fmt.Sprintf("%s\n%s\n", u, title))

		clusters := make(map[string]*hcpvsm.HashicorpCloudVault20201125Cluster, len(clustersResp.GetPayload().Clusters))
		for _, cluster := range clustersResp.GetPayload().Clusters {
			if *cluster.State == hcpvsm.HashicorpCloudVault20201125ClusterStateRUNNING {
				c.Ui.Info(fmt.Sprintf("Cluster identification: %s", cluster.ID))
				id := strings.ToLower(cluster.ID)
				clusters[id] = cluster
			}
		}
		userInput, err := c.Ui.Ask("\nChoose a cluster:")
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Failed to get cluster information: %s", err))
			return "", err
		}

		// set the cluster
		cluster, ok := clusters[userInput]
		if !ok {
			return "", errors.New(fmt.Sprintf("invalid cluster: %s", userInput))
		}
		if *cluster.Config.NetworkConfig.HTTPProxyOption == hcpvsm.HashicorpCloudVault20201125HTTPProxyOptionDISABLED {
			return "", ErrorProxyDisabled
		}

		proxyAddr = "https://" + cluster.DNSNames.Proxy
		return proxyAddr, nil

	default:
		cluster := clustersResp.GetPayload().Clusters[0]
		if *cluster.State != hcpvsm.HashicorpCloudVault20201125ClusterStateRUNNING {
			return "", errors.New("cluster is not running")
		}
		if *cluster.Config.NetworkConfig.HTTPProxyOption == hcpvsm.HashicorpCloudVault20201125HTTPProxyOptionDISABLED {
			return "", ErrorProxyDisabled
		}

		c.Ui.Info(fmt.Sprintf("HCP Vault Cluster: %s", cluster.ID))

		proxyAddr = "https://" + cluster.DNSNames.Proxy
		return proxyAddr, nil
	}
}
