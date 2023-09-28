package hcpvaultengine

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"

	hcprmo "github.com/hashicorp/hcp-sdk-go/clients/cloud-resource-manager/stable/2019-12-10/client/organization_service"
	hcprmp "github.com/hashicorp/hcp-sdk-go/clients/cloud-resource-manager/stable/2019-12-10/client/project_service"
	hcprmm "github.com/hashicorp/hcp-sdk-go/clients/cloud-resource-manager/stable/2019-12-10/models"
	hcpvs "github.com/hashicorp/hcp-sdk-go/clients/cloud-vault-service/stable/2020-11-25/client/vault_service"
	hcpvsm "github.com/hashicorp/hcp-sdk-go/clients/cloud-vault-service/stable/2020-11-25/models"
	"github.com/mitchellh/cli"
)

var (
	_ cli.Command = (*HCPConnectCommand)(nil)
)

type HCPConnectCommand struct {
	Ui cli.Ui

	flagNonInteractiveMethod bool
}

func (c *HCPConnectCommand) Help() string {
	//TODO implement me
	panic("implement me")
}

func (c *HCPConnectCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	connectHandler := c.connectHandlerFactory()

	hcpHttpClient, err := connectHandler.Connect(args)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Failed to connect to HCP: %s", err))
		return 1
	}

	// List orgs
	// If list is greater than 1, ask for user input c.Ui.Ask()
	//
	// should we add pagination?
	organizationID, err := c.getOrganization(hcprmo.New(hcpHttpClient, nil))
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Failed to get HCP organization information: %s", err))
		return 1
	}

	// List projects for chosen org
	// If list is greater than 1, ask for user input c.Ui.Ask()
	//
	// should we add pagination?
	projectID, err := c.getProject(hcprmp.New(hcpHttpClient, nil))
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Failed to get HCP project information: %s", err))
		return 1
	}

	// List clusters for org+project
	// If list is greater than 1, ask for user input c.Ui.Ask()
	//
	// should we add pagination?
	err = c.getCluster(organizationID, projectID, hcpvs.New(hcpHttpClient, nil))
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Failed to get HCP Vault Cluster information: %s", err))
		return 1
	}

	// Cache details -- in memory? in disk?
	// In memory for POC

	return 0
}

func (c *HCPConnectCommand) Synopsis() string {
	return "Authenticate to HCP"
}

func (c *HCPConnectCommand) Flags() *flag.FlagSet {
	mainSet := flag.NewFlagSet("", flag.ContinueOnError)
	mainSet.Var(&boolValue{target: &c.flagNonInteractiveMethod}, "non-interactive", "")
	return mainSet
}

func (c *HCPConnectCommand) getOrganization(rmOrgClient hcprmo.ClientService) (organizationID string, err error) {
	organizationsResp, err := rmOrgClient.OrganizationServiceList(&hcprmo.OrganizationServiceListParams{}, nil)
	switch {
	case err != nil:
		return "", err
	case organizationsResp.GetPayload() == nil:
		return "", errors.New("payload is nil")
	case len(organizationsResp.GetPayload().Organizations) < 1:
		return "", errors.New("no organizations available")
	case len(organizationsResp.GetPayload().Organizations) > 1:
		var orgs []string
		for i, org := range organizationsResp.GetPayload().Organizations {
			if org.State == hcprmm.HashicorpCloudResourcemanagerOrganizationOrganizationStateACTIVE.Pointer() {
				orgs = append(orgs, fmt.Sprintf("%d: Organization name: %s Organization ID: %s", i, org.Name, org.ID))
			}
		}
		userInput, err := c.Ui.Ask(fmt.Sprintf("Choose one of the following organizations: %s", strings.Join(orgs, "\n")))
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Failed to get HCP organization information: %s", err))
			return "", err
		}
		// convert userInput to int
		var index int
		index, err = strconv.Atoi(userInput)
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Failed to get HCP organization information: %s", err))
			return "", err
		}
		// if conversion fails, return an error
		// else validate that the index is within boundaries of the organization slice
		if index > len(orgs) || index < len(orgs) {
			return "", errors.New("invalid organization chosen")
		}

		// set the org ID
		organizationID = organizationsResp.GetPayload().Organizations[index].ID
		organizationName := organizationsResp.GetPayload().Organizations[index].Name
		c.Ui.Info(fmt.Sprintf("HCP Organization: %s", organizationName))

		break
	case len(organizationsResp.GetPayload().Organizations) == 1:
		organization := organizationsResp.GetPayload().Organizations[0]
		if *organization.State != hcprmm.HashicorpCloudResourcemanagerOrganizationOrganizationStateACTIVE {
			return "", errors.New("organization is not active")
		}
		organizationID = organization.ID
		c.Ui.Info(fmt.Sprintf("HCP Organization: %s", organization.Name))
	}
	return organizationID, nil
}

func (c *HCPConnectCommand) getProject(rmProjClient hcprmp.ClientService) (projectID string, err error) {
	projectResp, err := rmProjClient.ProjectServiceList(&hcprmp.ProjectServiceListParams{}, nil)
	switch {
	case err != nil:
		return "", err
	case projectResp.GetPayload() == nil:
		return "", errors.New("payload is nil")
	case len(projectResp.GetPayload().Projects) < 1:
		return "", errors.New("no projects available")
	case len(projectResp.GetPayload().Projects) > 1:
		var projs []string
		for i, proj := range projectResp.GetPayload().Projects {
			if *proj.State == hcprmm.HashicorpCloudResourcemanagerProjectProjectStateACTIVE {
				projs = append(projs, fmt.Sprintf("%d: Project name: %s Project ID: %s", i, proj.Name, proj.ID))
			}
		}
		userInput, err := c.Ui.Ask(fmt.Sprintf("Choose one of the following projects: %s", strings.Join(projs, "\n")))
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Failed to get HCP project information: %s", err))
			return "", err
		}
		// convert userInput to int
		var index int
		index, err = strconv.Atoi(userInput)
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Failed to get HCP project information: %s", err))
			return "", err
		}
		// else validate that the index is within boundaries of the organization slice
		if index > len(projs) || index < len(projs) {
			return "", errors.New("invalid project chosen")
		}

		// set the org ID
		projectID = projectResp.GetPayload().Projects[index].ID
		projectName := projectResp.GetPayload().Projects[index].Name
		c.Ui.Info(fmt.Sprintf("HCP Project: %s", projectName))

		break
	case len(projectResp.GetPayload().Projects) == 1:
		project := projectResp.GetPayload().Projects[0]
		if *project.State != hcprmm.HashicorpCloudResourcemanagerProjectProjectStateACTIVE {
			return "", errors.New("organization is not active")
		}
		projectID = project.ID
		c.Ui.Info(fmt.Sprintf("HCP Project: %s", project.Name))
	}
	return projectID, nil
}

func (c *HCPConnectCommand) getCluster(organizationID string, projectID string, vsClient hcpvs.ClientService) error {
	clustersResp, err := vsClient.List(&hcpvs.ListParams{LocationOrganizationID: organizationID, LocationProjectID: projectID}, nil)
	switch {
	case err != nil:
		return err
	case clustersResp.GetPayload() == nil:
		return errors.New("payload is nil")
	case len(clustersResp.GetPayload().Clusters) < 1:
		return errors.New("no projects available")
	case len(clustersResp.GetPayload().Clusters) > 1:
		var clusters []string
		for i, cluster := range clustersResp.GetPayload().Clusters {
			if *cluster.State == hcpvsm.HashicorpCloudVault20201125ClusterStateRUNNING {
				clusters = append(clusters, fmt.Sprintf("%d: HCP Vault Cluster name: %s HCP Vault Cluster ID: %s", i, cluster.ID, cluster.ResourceID))
			}
		}
		userInput, err := c.Ui.Ask(fmt.Sprintf("Choose one of the following HCP Vault Clusters: %s", strings.Join(clusters, "\n")))
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Failed to get HCP Vault Cluster information: %s", err))
			return err
		}
		// convert userInput to int
		var index int
		index, err = strconv.Atoi(userInput)
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Failed to get HCP Vault Cluster information: %s", err))
			return err
		}
		// else validate that the index is within boundaries of the organization slice
		if index > len(clusters) || index < len(clusters) {
			return errors.New("invalid cluster chosen")
		}

		// set the org ID
		clusterName := clustersResp.GetPayload().Clusters[index].ID
		c.Ui.Info(fmt.Sprintf("HCP Vault Cluster: %s", clusterName))

		cache.Address = clustersResp.GetPayload().Clusters[index].DNSNames.Proxy

		break
	case len(clustersResp.GetPayload().Clusters) == 1:
		cluster := clustersResp.GetPayload().Clusters[0]
		if *cluster.State != hcpvsm.HashicorpCloudVault20201125ClusterStateRUNNING {
			return errors.New("cluster is not running")
		}
		projectID = cluster.ResourceID
		c.Ui.Info(fmt.Sprintf("HCP Vault Cluster: %s", cluster.ID))

		cache.Address = clustersResp.GetPayload().Clusters[0].DNSNames.Proxy

		break
	}

	return nil
}

func (c *HCPConnectCommand) connectHandlerFactory() ConnectHandler {
	if c.flagNonInteractiveMethod {
		return &nonInteractiveConnectHandler{}
	}
	return &interactiveConnectHandler{}
}
