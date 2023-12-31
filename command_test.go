// Copyright (c) HashiCorp, Inc.
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
