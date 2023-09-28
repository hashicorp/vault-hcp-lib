package hcpvaultengine

import (
	"github.com/google/uuid"
	hcprmo "github.com/hashicorp/hcp-sdk-go/clients/cloud-resource-manager/stable/2019-12-10/client/organization_service"
	"github.com/hashicorp/hcp-sdk-go/clients/cloud-resource-manager/stable/2019-12-10/models"
	"github.com/hashicorp/hcp-vault-engine-poc/mocks"
	"github.com/mitchellh/cli"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func Test_getOrganization(t *testing.T) {
	mockUi := cli.NewMockUi()
	cmd := HCPConnectCommand{Ui: mockUi}

	mockRmOrgClient := mocks.NewClientService(t)

	expectedOrganizationID := uuid.New().String()

	tests := map[string]struct {
		organizationServiceListResponse *hcprmo.OrganizationServiceListOK
		error                           error
	}{
		// Test single organization
		// No UI interaction required
		"single organization": {
			organizationServiceListResponse: &hcprmo.OrganizationServiceListOK{
				Payload: &models.HashicorpCloudResourcemanagerOrganizationListResponse{
					Organizations: []*models.HashicorpCloudResourcemanagerOrganization{
						{
							ID:    expectedOrganizationID,
							Name:  "mock-organization-1",
							State: models.NewHashicorpCloudResourcemanagerOrganizationOrganizationState(models.HashicorpCloudResourcemanagerOrganizationOrganizationStateACTIVE),
						},
					},
				},
			},
			error: nil,
		},
	}

	for n, tst := range tests {
		t.Run(n, func(t *testing.T) {
			mockRmOrgClient.
				On("OrganizationServiceList", mock.Anything, nil).
				Return(tst.organizationServiceListResponse, tst.error)

			organizationID, err := cmd.getOrganization(mockRmOrgClient)
			assert.NoError(t, err)
			assert.Equal(t, expectedOrganizationID, organizationID)
		})
	}

}
