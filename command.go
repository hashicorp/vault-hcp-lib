package hcpvaultengine

import "github.com/mitchellh/cli"

func InitHCPCommand(ui cli.Ui) map[string]cli.CommandFactory {
	return map[string]cli.CommandFactory{
		"connect": func() (cli.Command, error) {
			return &HCPConnectCommand{
				Ui: ui,
			}, nil
		},
	}
}
