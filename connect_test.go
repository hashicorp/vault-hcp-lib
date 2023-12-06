// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package vaulthcplib

import (
	"errors"
	"io"
	"testing"

	"github.com/hashicorp/cli"
	clustermocks "github.com/hashicorp/vault-hcp-lib/mocks/cluster"

	orgmocks "github.com/hashicorp/vault-hcp-lib/mocks/organization"
	projmocks "github.com/hashicorp/vault-hcp-lib/mocks/project"

	"github.com/google/uuid"
	hcprmo "github.com/hashicorp/hcp-sdk-go/clients/cloud-resource-manager/stable/2019-12-10/client/organization_service"
	hcprmp "github.com/hashicorp/hcp-sdk-go/clients/cloud-resource-manager/stable/2019-12-10/client/project_service"
	"github.com/hashicorp/hcp-sdk-go/clients/cloud-resource-manager/stable/2019-12-10/models"
	hcpvs "github.com/hashicorp/hcp-sdk-go/clients/cloud-vault-service/stable/2020-11-25/client/vault_service"
	hcpvsm "github.com/hashicorp/hcp-sdk-go/clients/cloud-vault-service/stable/2020-11-25/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func testHCPConnectCommand() (*cli.MockUi, *HCPConnectCommand) {
	ui := cli.NewMockUi()
	return ui, &HCPConnectCommand{Ui: ui}
}

func Test_HCPConnectCommand(t *testing.T) {
	_, cmd := testHCPConnectCommand()

	mockRmOrgClient := orgmocks.NewClientService(t)
	mockRmOrgClient.
		On("OrganizationServiceList", mock.Anything, nil).
		Return(&hcprmo.OrganizationServiceListOK{
			Payload: &models.HashicorpCloudResourcemanagerOrganizationListResponse{
				Organizations: []*models.HashicorpCloudResourcemanagerOrganization{
					{
						ID:   uuid.New().String(),
						Name: "mock-organization-1",
						State: models.NewHashicorpCloudResourcemanagerOrganizationOrganizationState(
							models.HashicorpCloudResourcemanagerOrganizationOrganizationStateACTIVE,
						),
					},
				},
			},
		}, nil)

	mockRmProjClient := projmocks.NewClientService(t)
	mockRmProjClient.
		On("ProjectServiceList", mock.Anything, nil).
		Return(&hcprmp.ProjectServiceListOK{
			Payload: &models.HashicorpCloudResourcemanagerProjectListResponse{
				Projects: []*models.HashicorpCloudResourcemanagerProject{
					{
						ID:   uuid.New().String(),
						Name: "mock-project-1",
						State: models.NewHashicorpCloudResourcemanagerProjectProjectState(
							models.HashicorpCloudResourcemanagerProjectProjectStateACTIVE,
						),
					},
				},
			},
		}, nil)

	mockVsClient := clustermocks.NewClientService(t)
	mockVsClient.
		On("Get", mock.Anything, nil).
		Return(&hcpvs.GetOK{
			Payload: &hcpvsm.HashicorpCloudVault20201125GetResponse{
				Cluster: &hcpvsm.HashicorpCloudVault20201125Cluster{
					ID:       "cluster-1",
					DNSNames: &hcpvsm.HashicorpCloudVault20201125ClusterDNSNames{Proxy: "hcp-proxy-cluster-1.addr:8200"},
					State: hcpvsm.NewHashicorpCloudVault20201125ClusterState(
						hcpvsm.HashicorpCloudVault20201125ClusterStateRUNNING,
					),
				},
			},
		}, nil)

	cmd.rmOrgClient = mockRmOrgClient
	cmd.rmProjClient = mockRmProjClient
	cmd.vsClient = mockVsClient

	result := cmd.Run([]string{"-cluster-id", "cluster-1"})
	assert.Equal(t, 0, result)
}

func Test_getOrganization(t *testing.T) {
	organizationID := uuid.New().String()
	organizationIDTwo := uuid.New().String()
	organizationIDThree := uuid.New().String()

	tests := map[string]struct {
		userInputOrganizationName string
		expectedOrganizationID    string

		organizationServiceListResponse *hcprmo.OrganizationServiceListOK

		expectedError error
	}{
		// Test single organization
		// No UI interaction required
		"single organization": {
			expectedOrganizationID: organizationID,
			organizationServiceListResponse: &hcprmo.OrganizationServiceListOK{
				Payload: &models.HashicorpCloudResourcemanagerOrganizationListResponse{
					Organizations: []*models.HashicorpCloudResourcemanagerOrganization{
						{
							ID:    organizationID,
							Name:  "mock-organization-1",
							State: models.NewHashicorpCloudResourcemanagerOrganizationOrganizationState(models.HashicorpCloudResourcemanagerOrganizationOrganizationStateACTIVE),
						},
					},
				},
			},
		},

		// Test multiple organizations
		// UI interaction required
		"multiple organizations": {
			userInputOrganizationName: "mock-organization-2\n",
			expectedOrganizationID:    organizationIDTwo,
			organizationServiceListResponse: &hcprmo.OrganizationServiceListOK{
				Payload: &models.HashicorpCloudResourcemanagerOrganizationListResponse{
					Organizations: []*models.HashicorpCloudResourcemanagerOrganization{
						{
							ID:    organizationID,
							Name:  "mock-organization-1",
							State: models.NewHashicorpCloudResourcemanagerOrganizationOrganizationState(models.HashicorpCloudResourcemanagerOrganizationOrganizationStateACTIVE),
						},
						{
							ID:    organizationIDTwo,
							Name:  "mock-organization-2",
							State: models.NewHashicorpCloudResourcemanagerOrganizationOrganizationState(models.HashicorpCloudResourcemanagerOrganizationOrganizationStateACTIVE),
						},
						{
							ID:    organizationIDThree,
							Name:  "mock-organization-3",
							State: models.NewHashicorpCloudResourcemanagerOrganizationOrganizationState(models.HashicorpCloudResourcemanagerOrganizationOrganizationStateACTIVE),
						},
					},
				},
			},
		},

		// Test invalid organization
		// UI interaction required
		"invalid organization": {
			userInputOrganizationName: "mock-organization-4",
			organizationServiceListResponse: &hcprmo.OrganizationServiceListOK{
				Payload: &models.HashicorpCloudResourcemanagerOrganizationListResponse{
					Organizations: []*models.HashicorpCloudResourcemanagerOrganization{
						{
							ID:    organizationID,
							Name:  "mock-organization-1",
							State: models.NewHashicorpCloudResourcemanagerOrganizationOrganizationState(models.HashicorpCloudResourcemanagerOrganizationOrganizationStateACTIVE),
						},
						{
							ID:    organizationIDTwo,
							Name:  "mock-organization-2",
							State: models.NewHashicorpCloudResourcemanagerOrganizationOrganizationState(models.HashicorpCloudResourcemanagerOrganizationOrganizationStateACTIVE),
						},
						{
							ID:    organizationIDThree,
							Name:  "mock-organization-3",
							State: models.NewHashicorpCloudResourcemanagerOrganizationOrganizationState(models.HashicorpCloudResourcemanagerOrganizationOrganizationStateACTIVE),
						},
					},
				},
			},
			expectedError: errors.New("invalid HCP organization: mock-organization-4"),
		},

		// Test generic expectedError returned
		"expectedError": {
			expectedError: errors.New("error getting organization"),
		},
	}

	for n, tst := range tests {
		t.Run(n, func(t *testing.T) {
			ui, cmd := testHCPConnectCommand()
			stdinR, stdinW := io.Pipe()
			go func() {
				stdinW.Write([]byte(tst.userInputOrganizationName))
				stdinW.Close()
			}()
			ui.InputReader = stdinR

			mockRmOrgClient := orgmocks.NewClientService(t)

			mockRmOrgClient.
				On("OrganizationServiceList", mock.Anything, nil).
				Return(tst.organizationServiceListResponse, tst.expectedError)

			orgID, err := cmd.getOrganization(mockRmOrgClient)
			if tst.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tst.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tst.expectedOrganizationID, orgID)
			}
		})
	}

}

func Test_getProject(t *testing.T) {
	projectID := uuid.New().String()
	projectIDTwo := uuid.New().String()
	projectIDThree := uuid.New().String()

	tests := map[string]struct {
		userInputProjectName string
		expectedProjectID    string

		projectServiceListResponse *hcprmp.ProjectServiceListOK

		expectedError error
	}{
		// Test single project
		// No UI interaction required
		"single project": {
			expectedProjectID: projectID,
			projectServiceListResponse: &hcprmp.ProjectServiceListOK{
				Payload: &models.HashicorpCloudResourcemanagerProjectListResponse{
					Projects: []*models.HashicorpCloudResourcemanagerProject{
						{
							ID:    projectID,
							Name:  "mock-project-1",
							State: models.NewHashicorpCloudResourcemanagerProjectProjectState(models.HashicorpCloudResourcemanagerProjectProjectStateACTIVE),
						},
					},
				},
			},
		},

		// Test multiple projects
		// UI interaction required
		"multiple projects": {
			userInputProjectName: "mock-project-2\n",
			expectedProjectID:    projectIDTwo,
			projectServiceListResponse: &hcprmp.ProjectServiceListOK{
				Payload: &models.HashicorpCloudResourcemanagerProjectListResponse{
					Projects: []*models.HashicorpCloudResourcemanagerProject{
						{
							ID:    projectID,
							Name:  "mock-project-1",
							State: models.NewHashicorpCloudResourcemanagerProjectProjectState(models.HashicorpCloudResourcemanagerProjectProjectStateACTIVE),
						},
						{
							ID:    projectIDTwo,
							Name:  "mock-project-2",
							State: models.NewHashicorpCloudResourcemanagerProjectProjectState(models.HashicorpCloudResourcemanagerProjectProjectStateACTIVE),
						},
						{
							ID:    projectIDThree,
							Name:  "mock-project-3",
							State: models.NewHashicorpCloudResourcemanagerProjectProjectState(models.HashicorpCloudResourcemanagerProjectProjectStateACTIVE),
						},
					},
				},
			},
		},

		// Test invalid project
		// UI interaction required
		"invalid project": {
			userInputProjectName: "mock-project-4",
			expectedProjectID:    projectID,
			projectServiceListResponse: &hcprmp.ProjectServiceListOK{
				Payload: &models.HashicorpCloudResourcemanagerProjectListResponse{
					Projects: []*models.HashicorpCloudResourcemanagerProject{
						{
							ID:    projectID,
							Name:  "mock-project-1",
							State: models.NewHashicorpCloudResourcemanagerProjectProjectState(models.HashicorpCloudResourcemanagerProjectProjectStateACTIVE),
						},
						{
							ID:    projectIDTwo,
							Name:  "mock-project-2",
							State: models.NewHashicorpCloudResourcemanagerProjectProjectState(models.HashicorpCloudResourcemanagerProjectProjectStateACTIVE),
						},
						{
							ID:    projectIDThree,
							Name:  "mock-project-3",
							State: models.NewHashicorpCloudResourcemanagerProjectProjectState(models.HashicorpCloudResourcemanagerProjectProjectStateACTIVE),
						},
					},
				},
			},
			expectedError: errors.New("invalid HCP project: mock-project-4"),
		},

		// Test generic expectedError returned
		"expectedError": {
			expectedError: errors.New("error getting project"),
		},
	}

	for n, tst := range tests {
		t.Run(n, func(t *testing.T) {
			ui, cmd := testHCPConnectCommand()
			stdinR, stdinW := io.Pipe()
			go func() {
				stdinW.Write([]byte(tst.userInputProjectName))
				stdinW.Close()
			}()
			ui.InputReader = stdinR

			mockRmProjClient := projmocks.NewClientService(t)

			mockRmProjClient.
				On("ProjectServiceList", mock.Anything, nil).
				Return(tst.projectServiceListResponse, tst.expectedError)

			projID, err := cmd.getProject("", mockRmProjClient)
			if tst.expectedError != nil {
				assert.Error(t, tst.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tst.expectedProjectID, projID)
			}
		})
	}

}

func Test_getCluster(t *testing.T) {
	tests := map[string]struct {
		expectedProxyAddr string
		userParamCluster  string
		userInputCluster  string

		getClusterServiceListResponse   *hcpvs.GetOK
		listClustersServiceListResponse *hcpvs.ListOK

		expectedError error
	}{
		// Test using cluster id received as a parameter
		// No UI interaction required
		"parameter cluster": {
			userParamCluster:  "cluster-1",
			expectedProxyAddr: "https://hcp-proxy-cluster-1.addr:8200",
			getClusterServiceListResponse: &hcpvs.GetOK{
				Payload: &hcpvsm.HashicorpCloudVault20201125GetResponse{
					Cluster: &hcpvsm.HashicorpCloudVault20201125Cluster{
						ID:       "cluster-1",
						DNSNames: &hcpvsm.HashicorpCloudVault20201125ClusterDNSNames{Proxy: "hcp-proxy-cluster-1.addr:8200"},
						State:    hcpvsm.NewHashicorpCloudVault20201125ClusterState(hcpvsm.HashicorpCloudVault20201125ClusterStateRUNNING),
					},
				},
			},
		},

		// Test single project
		// No UI interaction required
		"single cluster": {
			expectedProxyAddr: "https://hcp-proxy-cluster-1.addr:8200",
			listClustersServiceListResponse: &hcpvs.ListOK{
				Payload: &hcpvsm.HashicorpCloudVault20201125ListResponse{
					Clusters: []*hcpvsm.HashicorpCloudVault20201125Cluster{
						{
							ID:       "cluster-1",
							DNSNames: &hcpvsm.HashicorpCloudVault20201125ClusterDNSNames{Proxy: "hcp-proxy-cluster-1.addr:8200"},
							State:    hcpvsm.NewHashicorpCloudVault20201125ClusterState(hcpvsm.HashicorpCloudVault20201125ClusterStateRUNNING),
						},
					},
				},
			},
		},

		// Test multiple projects
		// UI interaction required
		"multiple clusters": {
			expectedProxyAddr: "https://hcp-proxy-cluster-2.addr:8200",
			userInputCluster:  "cluster-2\n",
			listClustersServiceListResponse: &hcpvs.ListOK{
				Payload: &hcpvsm.HashicorpCloudVault20201125ListResponse{
					Clusters: []*hcpvsm.HashicorpCloudVault20201125Cluster{
						{
							ID:       "cluster-1",
							DNSNames: &hcpvsm.HashicorpCloudVault20201125ClusterDNSNames{Proxy: "hcp-proxy-cluster-1.addr:8200"},
							State:    hcpvsm.NewHashicorpCloudVault20201125ClusterState(hcpvsm.HashicorpCloudVault20201125ClusterStateRUNNING),
						},
						{
							ID:       "cluster-2",
							DNSNames: &hcpvsm.HashicorpCloudVault20201125ClusterDNSNames{Proxy: "hcp-proxy-cluster-2.addr:8200"},
							State:    hcpvsm.NewHashicorpCloudVault20201125ClusterState(hcpvsm.HashicorpCloudVault20201125ClusterStateRUNNING),
						},
						{
							ID:       "cluster-3",
							DNSNames: &hcpvsm.HashicorpCloudVault20201125ClusterDNSNames{Proxy: "hcp-proxy-cluster-3.addr:8200"},
							State:    hcpvsm.NewHashicorpCloudVault20201125ClusterState(hcpvsm.HashicorpCloudVault20201125ClusterStateRUNNING),
						},
					},
				},
			},
		},

		// Test invalid project
		// UI interaction required
		"invalid cluster": {
			userInputCluster: "cluster-4",
			listClustersServiceListResponse: &hcpvs.ListOK{
				Payload: &hcpvsm.HashicorpCloudVault20201125ListResponse{
					Clusters: []*hcpvsm.HashicorpCloudVault20201125Cluster{
						{
							ID:       "cluster-1",
							DNSNames: &hcpvsm.HashicorpCloudVault20201125ClusterDNSNames{Proxy: "hcp-proxy-cluster-1.addr:8200"},
							State:    hcpvsm.NewHashicorpCloudVault20201125ClusterState(hcpvsm.HashicorpCloudVault20201125ClusterStateRUNNING),
						},
						{
							ID:       "cluster-2",
							DNSNames: &hcpvsm.HashicorpCloudVault20201125ClusterDNSNames{Proxy: "hcp-proxy-cluster-2.addr:8200"},
							State:    hcpvsm.NewHashicorpCloudVault20201125ClusterState(hcpvsm.HashicorpCloudVault20201125ClusterStateRUNNING),
						},
						{
							ID:       "cluster-3",
							DNSNames: &hcpvsm.HashicorpCloudVault20201125ClusterDNSNames{Proxy: "hcp-proxy-cluster-3.addr:8200"},
							State:    hcpvsm.NewHashicorpCloudVault20201125ClusterState(hcpvsm.HashicorpCloudVault20201125ClusterStateRUNNING),
						},
					},
				},
			},
			expectedError: errors.New("invalid cluster: cluster-4"),
		},

		// Test generic expectedError returned
		"expectedError": {
			expectedError: errors.New("error getting cluster"),
		},
	}

	for n, tst := range tests {
		t.Run(n, func(t *testing.T) {
			ui, cmd := testHCPConnectCommand()
			stdinR, stdinW := io.Pipe()
			go func() {
				stdinW.Write([]byte(tst.userInputCluster))
				stdinW.Close()
			}()
			ui.InputReader = stdinR

			mockVsClient := clustermocks.NewClientService(t)

			// in case user pass in the cluster id, we'll request the cluster details from the Get RPC
			// else we'll request a list of clusters and ask to choose one
			if tst.userParamCluster != "" {
				mockVsClient.
					On("Get", mock.Anything, nil).
					Return(tst.getClusterServiceListResponse, tst.expectedError)
			} else {
				mockVsClient.
					On("List", mock.Anything, nil).
					Return(tst.listClustersServiceListResponse, tst.expectedError)
			}

			proxyAddr, err := cmd.getCluster("", "", tst.userParamCluster, mockVsClient)
			if tst.expectedError != nil {
				assert.Error(t, tst.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tst.expectedProxyAddr, proxyAddr)
			}
		})
	}
}

func Test_getProxyAddr(t *testing.T) {
	tests := map[string]struct {
		expectedProxyAddr             string
		userParamCluster              string
		getClusterServiceListResponse *hcpvs.GetOK

		userParamOrgID                  string
		organizationServiceListResponse *hcprmo.OrganizationServiceListOK

		userParamProjID            string
		projectServiceListResponse *hcprmp.ProjectServiceListOK

		expectedError error
	}{
		"success: not parameterized org, not parameterized project": {
			userParamCluster:  "cluster-1",
			expectedProxyAddr: "https://hcp-proxy-cluster-1.addr:8200",
			getClusterServiceListResponse: &hcpvs.GetOK{
				Payload: &hcpvsm.HashicorpCloudVault20201125GetResponse{
					Cluster: &hcpvsm.HashicorpCloudVault20201125Cluster{
						ID:       "cluster-1",
						DNSNames: &hcpvsm.HashicorpCloudVault20201125ClusterDNSNames{Proxy: "hcp-proxy-cluster-1.addr:8200"},
						State:    hcpvsm.NewHashicorpCloudVault20201125ClusterState(hcpvsm.HashicorpCloudVault20201125ClusterStateRUNNING),
					},
				},
			},
			organizationServiceListResponse: &hcprmo.OrganizationServiceListOK{
				Payload: &models.HashicorpCloudResourcemanagerOrganizationListResponse{
					Organizations: []*models.HashicorpCloudResourcemanagerOrganization{
						{
							ID:    uuid.New().String(),
							Name:  "mock-organization-1",
							State: models.NewHashicorpCloudResourcemanagerOrganizationOrganizationState(models.HashicorpCloudResourcemanagerOrganizationOrganizationStateACTIVE),
						},
					},
				},
			},
			projectServiceListResponse: &hcprmp.ProjectServiceListOK{
				Payload: &models.HashicorpCloudResourcemanagerProjectListResponse{
					Projects: []*models.HashicorpCloudResourcemanagerProject{
						{
							ID:    uuid.New().String(),
							Name:  "mock-project-1",
							State: models.NewHashicorpCloudResourcemanagerProjectProjectState(models.HashicorpCloudResourcemanagerProjectProjectStateACTIVE),
						},
					},
				},
			},
		},

		"success: parameterized org": {
			userParamCluster:  "cluster-1",
			expectedProxyAddr: "https://hcp-proxy-cluster-1.addr:8200",
			getClusterServiceListResponse: &hcpvs.GetOK{
				Payload: &hcpvsm.HashicorpCloudVault20201125GetResponse{
					Cluster: &hcpvsm.HashicorpCloudVault20201125Cluster{
						ID:       "cluster-1",
						DNSNames: &hcpvsm.HashicorpCloudVault20201125ClusterDNSNames{Proxy: "hcp-proxy-cluster-1.addr:8200"},
						State:    hcpvsm.NewHashicorpCloudVault20201125ClusterState(hcpvsm.HashicorpCloudVault20201125ClusterStateRUNNING),
					},
				},
			},
			userParamOrgID: uuid.New().String(),
			projectServiceListResponse: &hcprmp.ProjectServiceListOK{
				Payload: &models.HashicorpCloudResourcemanagerProjectListResponse{
					Projects: []*models.HashicorpCloudResourcemanagerProject{
						{
							ID:    uuid.New().String(),
							Name:  "mock-project-1",
							State: models.NewHashicorpCloudResourcemanagerProjectProjectState(models.HashicorpCloudResourcemanagerProjectProjectStateACTIVE),
						},
					},
				},
			},
		},

		"success: parameterized project": {
			userParamCluster:  "cluster-1",
			expectedProxyAddr: "https://hcp-proxy-cluster-1.addr:8200",
			getClusterServiceListResponse: &hcpvs.GetOK{
				Payload: &hcpvsm.HashicorpCloudVault20201125GetResponse{
					Cluster: &hcpvsm.HashicorpCloudVault20201125Cluster{
						ID:       "cluster-1",
						DNSNames: &hcpvsm.HashicorpCloudVault20201125ClusterDNSNames{Proxy: "hcp-proxy-cluster-1.addr:8200"},
						State:    hcpvsm.NewHashicorpCloudVault20201125ClusterState(hcpvsm.HashicorpCloudVault20201125ClusterStateRUNNING),
					},
				},
			},
			organizationServiceListResponse: &hcprmo.OrganizationServiceListOK{
				Payload: &models.HashicorpCloudResourcemanagerOrganizationListResponse{
					Organizations: []*models.HashicorpCloudResourcemanagerOrganization{
						{
							ID:    uuid.New().String(),
							Name:  "mock-organization-1",
							State: models.NewHashicorpCloudResourcemanagerOrganizationOrganizationState(models.HashicorpCloudResourcemanagerOrganizationOrganizationStateACTIVE),
						},
					},
				},
			},
			userParamProjID: uuid.New().String(),
		},

		"success: parameterized org and project": {
			userParamCluster:  "cluster-1",
			expectedProxyAddr: "https://hcp-proxy-cluster-1.addr:8200",
			getClusterServiceListResponse: &hcpvs.GetOK{
				Payload: &hcpvsm.HashicorpCloudVault20201125GetResponse{
					Cluster: &hcpvsm.HashicorpCloudVault20201125Cluster{
						ID:       "cluster-1",
						DNSNames: &hcpvsm.HashicorpCloudVault20201125ClusterDNSNames{Proxy: "hcp-proxy-cluster-1.addr:8200"},
						State:    hcpvsm.NewHashicorpCloudVault20201125ClusterState(hcpvsm.HashicorpCloudVault20201125ClusterStateRUNNING),
					},
				},
			},
			userParamOrgID:  uuid.New().String(),
			userParamProjID: uuid.New().String(),
		},

		"error: parameterized org and project": {
			userParamCluster:  "cluster-1",
			expectedProxyAddr: "https://hcp-proxy-cluster-1.addr:8200",
			userParamOrgID:    "invalid-org",
			userParamProjID:   "invalid-proj",
			expectedError:     errors.New("error getting cluster"),
		},
	}

	for n, tst := range tests {
		t.Run(n, func(t *testing.T) {
			_, cmd := testHCPConnectCommand()
			cmd.flagClusterID = tst.userParamCluster
			cmd.flagOrganizationID = tst.userParamOrgID
			cmd.flagProjectID = tst.userParamProjID

			mockRmOrgClient := orgmocks.NewClientService(t)
			mockRmProjClient := projmocks.NewClientService(t)
			mockVsClient := clustermocks.NewClientService(t)

			// mock vault service response
			if tst.getClusterServiceListResponse != nil {
				mockVsClient.
					On("Get", mock.Anything, nil).
					Return(tst.getClusterServiceListResponse, nil)
			} else {
				mockVsClient.
					On("Get", mock.Anything, nil).
					Return(nil, tst.expectedError)
			}

			// mock resource manager service response
			if tst.organizationServiceListResponse != nil {
				mockRmOrgClient.
					On("OrganizationServiceList", mock.Anything, nil).
					Return(tst.organizationServiceListResponse, nil)
			}
			if tst.projectServiceListResponse != nil {
				mockRmProjClient.
					On("ProjectServiceList", mock.Anything, nil).
					Return(tst.projectServiceListResponse, nil)
			}

			proxyAddr, err := cmd.getProxyAddr(mockRmOrgClient, mockRmProjClient, mockVsClient)
			if tst.expectedError != nil {
				assert.Error(t, tst.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tst.expectedProxyAddr, proxyAddr)
			}
		})
	}
}
