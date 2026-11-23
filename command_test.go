// Copyright IBM Corp. 2023, 2025
// SPDX-License-Identifier: MPL-2.0

package vaulthcplib

import (
	"testing"

	"github.com/hashicorp/cli"
	"github.com/stretchr/testify/assert"
)

func Test_InitHCPCommand(t *testing.T) {
	cmdMap := InitHCPCommand(&cli.MockUi{})

	assert.Contains(t, cmdMap, "hcp connect")
	assert.Contains(t, cmdMap, "hcp disconnect")
}

// Test_InitHCPCommand_WithBinaryName verifies that WithBinaryName reaches the
// help of every command InitHCPCommand returns.
func Test_InitHCPCommand_WithBinaryName(t *testing.T) {
	t.Parallel()

	for name, factory := range InitHCPCommand(&cli.MockUi{}, WithBinaryName("vault-client")) {
		cmd, err := factory()
		assert.NoError(t, err)
		assert.Contains(t, cmd.Help(), "vault-client "+name)
	}
}
