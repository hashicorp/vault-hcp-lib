package vaulthcplib

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/hashicorp/hcp-sdk-go/auth"
	hcprmo "github.com/hashicorp/hcp-sdk-go/clients/cloud-resource-manager/stable/2019-12-10/client/organization_service"
	hcprmp "github.com/hashicorp/hcp-sdk-go/clients/cloud-resource-manager/stable/2019-12-10/client/project_service"
	hcprmm "github.com/hashicorp/hcp-sdk-go/clients/cloud-resource-manager/stable/2019-12-10/models"
	hcpvs "github.com/hashicorp/hcp-sdk-go/clients/cloud-vault-service/stable/2020-11-25/client/vault_service"
	hcpvsm "github.com/hashicorp/hcp-sdk-go/clients/cloud-vault-service/stable/2020-11-25/models"
	"github.com/hashicorp/hcp-sdk-go/config"
	"github.com/hashicorp/hcp-sdk-go/httpclient"
	"github.com/mitchellh/cli"
)

var (
	_ cli.Command = (*HCPConnectCommand)(nil)
)

type HCPConnectCommand struct {
	Ui cli.Ui

	flagNonInteractiveMethod bool
	flagClientID             string
	flagSecretID             string
	flagOrganizationID       string
	flagProjectID            string
	flagClusterID            string

	// for testing
	rmOrgClient     hcprmo.ClientService
	vsClient        hcpvs.ClientService
	rmProjClient    hcprmp.ClientService
	testAuthSession auth.Session
}

func (c *HCPConnectCommand) Help() string {
	helpText := `
Usage: vault hcp connect [options]
  
  Authenticates users or machines to HCP using either provided arguments or retrieved token through
  browser login. A successful authentication results in an HCP token and an HCP Vault address being
  locally cached. 

  The default authentication method is an interactive one, redirecting users to the HCP login browser.
  If a non-interactive option is supplied, it can be used if provided with a service principal credential
  generated through the HCP portal with the necessary capabilities to access the organization, project, and
  HCP Vault cluster chosen.

      $ vault hcp connect -non-interactive=true -client-id=client-id-value -secret-id=secret-id-value
  
  Additionally, the organization identification, project identification, and cluster name can be passed in to
  directly connect to a specific HCP Vault cluster without interacting with the CLI.
  
      $ vault hcp connect -non-interactive=true -client-id=client-id-value -secret-id=secret-id-value -organization-id=org-UUID -project-id=proj-UUID -cluster-id=cluster-name
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
		c.Ui.Error(err.Error())
		return 1
	}

	err = writeConfig(proxyAddr, c.flagClientID, c.flagSecretID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Failed to connect to HCP Vault Cluster: %s", err))
		return 1
	}

	return 0
}

func (c *HCPConnectCommand) setupClients() error {
	var opts []config.HCPConfigOption
	if c.testAuthSession != nil {
		opts = []config.HCPConfigOption{config.WithSession(c.testAuthSession)}
	} else {
		opts = []config.HCPConfigOption{config.FromEnv()}
		if c.flagClientID != "" && c.flagSecretID != "" {
			opts = append(opts, config.WithClientCredentials(c.flagClientID, c.flagSecretID))
			opts = append(opts, config.WithoutBrowserLogin())
		}
	}

	cfg, err := config.NewHCPConfig(opts...)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to connect to HCP: %s", err))
	}

	hcpHttpClient, err := httpclient.New(httpclient.Config{HCPConfig: cfg})
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to connect to HCP: %s", err))
	}

	if c.rmOrgClient == nil {
		c.rmOrgClient = hcprmo.New(hcpHttpClient, nil)
	}
	if c.rmProjClient == nil {
		c.rmProjClient = hcprmp.New(hcpHttpClient, nil)
	}
	if c.vsClient == nil {
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
		return "", errors.New(fmt.Sprintf("Failed to get HCP Vault Cluster information: %s", err))
	}
	return proxyAddr, nil
}

func (c *HCPConnectCommand) Synopsis() string {
	return "Connect to an HCP Vault Cluster"
}

func (c *HCPConnectCommand) Flags() *flag.FlagSet {
	mainSet := flag.NewFlagSet("", flag.ContinueOnError)

	mainSet.BoolVar(&c.flagNonInteractiveMethod, "non-interactive", false, "")
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
		orgs := make(map[string]*hcprmm.HashicorpCloudResourcemanagerOrganization, len(organizationsResp.GetPayload().Organizations))
		for _, org := range organizationsResp.GetPayload().Organizations {
			if *org.State == hcprmm.HashicorpCloudResourcemanagerOrganizationOrganizationStateACTIVE {
				c.Ui.Info(fmt.Sprintf("\nHCP Organization Name: %s", org.Name))
				name := strings.ToLower(org.Name)
				orgs[name] = org
			}
		}
		userInput, err := c.Ui.Ask(fmt.Sprintf("Choose one organization: "))
		if err != nil {
			return "", err
		}
		chosenOrg, ok := orgs[userInput]
		if !ok {
			return "", errors.New(fmt.Sprintf("invalid HCP organization: %s", userInput))
		}
		// set the org ID
		organizationID = chosenOrg.ID
		organizationName := chosenOrg.ID

		c.Ui.Info(fmt.Sprintf("HCP Organization: %s", organizationName))
		return organizationID, nil
	default:
		organization := organizationsResp.GetPayload().Organizations[0]
		if *organization.State != hcprmm.HashicorpCloudResourcemanagerOrganizationOrganizationStateACTIVE {
			return "", errors.New("organization is not active")
		}
		organizationID = organization.ID
		c.Ui.Info(fmt.Sprintf("HCP Organization: %s", organization.Name))
		return organizationID, nil
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
		projs := make(map[string]*hcprmm.HashicorpCloudResourcemanagerProject, len(projectResp.GetPayload().Projects))
		for _, proj := range projectResp.GetPayload().Projects {
			if *proj.State == hcprmm.HashicorpCloudResourcemanagerProjectProjectStateACTIVE {
				c.Ui.Info(fmt.Sprintf("\nHCP Project Name: %s", proj.Name))
				name := strings.ToLower(proj.Name)
				projs[name] = proj
			}
		}
		userInput, err := c.Ui.Ask(fmt.Sprintf("Choose one project: "))
		if err != nil {
			return "", err
		}
		chosenProj, ok := projs[userInput]
		if !ok {
			return "", errors.New(fmt.Sprintf("invalid HCP project: %s", userInput))
		}

		// set the project ID
		projectID = chosenProj.ID
		projectName := chosenProj.Name

		c.Ui.Info(fmt.Sprintf("HCP Project: %s", projectName))
		return projectID, nil
	default:
		project := projectResp.GetPayload().Projects[0]
		if *project.State != hcprmm.HashicorpCloudResourcemanagerProjectProjectStateACTIVE {
			return "", errors.New("project is not active")
		}
		projectID = project.ID
		c.Ui.Info(fmt.Sprintf("HCP Project: %s", project.Name))

		return projectID, nil
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
		c.Ui.Info(fmt.Sprintf("HCP Vault Cluster: %s", cluster.ID))

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
		clusters := make(map[string]*hcpvsm.HashicorpCloudVault20201125Cluster, len(clustersResp.GetPayload().Clusters))
		for _, cluster := range clustersResp.GetPayload().Clusters {
			if *cluster.State == hcpvsm.HashicorpCloudVault20201125ClusterStateRUNNING {
				c.Ui.Info(fmt.Sprintf("\nHCP Vault Cluster ID: %s", cluster.ID))
				id := strings.ToLower(cluster.ID)
				clusters[id] = cluster
			}
		}
		userInput, err := c.Ui.Ask("\nChoose a Vault cluster:")
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Failed to get HCP Vault Cluster information: %s", err))
			return "", err
		}

		// set the cluster ID
		cluster, ok := clusters[userInput]
		if !ok {
			return "", errors.New(fmt.Sprintf("invalid cluster: %s", userInput))
		}
		c.Ui.Info(fmt.Sprintf("\nHCP Vault Cluster: %s", cluster.ID))

		proxyAddr = "https://" + cluster.DNSNames.Proxy
		return proxyAddr, nil
	default:
		cluster := clustersResp.GetPayload().Clusters[0]
		if *cluster.State != hcpvsm.HashicorpCloudVault20201125ClusterStateRUNNING {
			return "", errors.New("cluster is not running")
		}
		projectID = cluster.ResourceID
		c.Ui.Info(fmt.Sprintf("HCP Vault Cluster: %s", cluster.ID))

		proxyAddr = "https://" + cluster.DNSNames.Proxy
		return proxyAddr, nil
	}
}
