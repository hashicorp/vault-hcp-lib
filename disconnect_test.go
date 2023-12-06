// Copyright (c) HashiCorp, Inc.
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

func Test_HCPDisconnectCommand(t *testing.T) {
	err := os.Setenv(envVarCacheTestMode, "true")
	if err != nil {
		t.Error(err)
	}

	_, cmd := testHCPDisconnectCommand()

	result := cmd.Run([]string{})
	assert.Equal(t, 0, result)
}
