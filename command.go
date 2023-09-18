package vaulthcplib

import "github.com/mitchellh/cli"

func InitHCPCommand(ui cli.Ui) map[string]cli.CommandFactory {
	return map[string]cli.CommandFactory{
		"hcp connect": func() (cli.Command, error) {
			return &HCPConnectCommand{
				Ui: ui,
			}, nil
		},
		"hcp disconnect": func() (cli.Command, error) {
			return &HCPDisconnectCommand{
				Ui: ui,
			}, nil
		},
	}
}
