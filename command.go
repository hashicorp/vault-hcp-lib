// Copyright IBM Corp. 2023, 2025
// SPDX-License-Identifier: MPL-2.0

package vaulthcplib

import "github.com/hashicorp/cli"

// InitOption configures the commands returned by InitHCPCommand.
type InitOption func(*initOptions)

type initOptions struct {
	binaryName string
}

// WithBinaryName sets the executable the commands show in help and output. It
// defaults to "vault".
func WithBinaryName(name string) InitOption {
	return func(o *initOptions) {
		o.binaryName = name
	}
}

func InitHCPCommand(ui cli.Ui, opts ...InitOption) map[string]cli.CommandFactory {
	o := &initOptions{}
	for _, opt := range opts {
		opt(o)
	}

	return map[string]cli.CommandFactory{
		"hcp connect": func() (cli.Command, error) {
			return &HCPConnectCommand{
				Ui:         ui,
				BinaryName: o.binaryName,
			}, nil
		},
		"hcp disconnect": func() (cli.Command, error) {
			return &HCPDisconnectCommand{
				Ui:         ui,
				BinaryName: o.binaryName,
			}, nil
		},
	}
}

func binaryNameOrDefault(name string) string {
	if name == "" {
		return "vault"
	}

	return name
}
