package hcpvaultengine

import (
	"testing"

	"github.com/mitchellh/cli"
	"github.com/stretchr/testify/assert"
)

func Test_InitHCPCommand(t *testing.T) {
	cmdMap := InitHCPCommand(&cli.MockUi{})

	assert.Contains(t, cmdMap, "hcp connect")
	assert.Contains(t, cmdMap, "hcp disconnect")
}
