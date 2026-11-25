// Copyright IBM Corp. 2023, 2025
// SPDX-License-Identifier: MPL-2.0

package vaulthcplib

import (
	"os"
	"testing"

	"github.com/hashicorp/cli"
	"github.com/stretchr/testify/assert"
)

func testHCPDisconnectCommand() (*cli.MockUi, *HCPDisconnectCommand) {
	ui := cli.NewMockUi()
	return ui, &HCPDisconnectCommand{Ui: ui}
}

// Test_HCPDisconnectCommand_Help verifies that Help uses BinaryName in the
// commands it shows and falls back to "vault" when it's unset.
func Test_HCPDisconnectCommand_Help(t *testing.T) {
	t.Parallel()

	assert.Contains(t, (&HCPDisconnectCommand{}).Help(), "vault hcp disconnect")

	help := (&HCPDisconnectCommand{BinaryName: "vault-client"}).Help()
	assert.Contains(t, help, "vault-client hcp disconnect")
	assert.NotContains(t, help, "vault hcp")
}

func Test_HCPDisconnectCommand(t *testing.T) {
	err := os.Setenv(envVarCacheTestMode, "true")
	if err != nil {
		t.Error(err)
	}

	_, cmd := testHCPDisconnectCommand()

	result := cmd.Run([]string{})
	assert.Equal(t, 0, result)
}
