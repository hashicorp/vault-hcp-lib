package hcpvaultengine

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"
)

var (
	_ cli.Command = (*HCPConnectCommand)(nil)
)

type HCPConnectCommand struct {
	Ui cli.Ui

	flagNonInteractiveMethod bool

	connectHandler ConnectHandler
}

func (c *HCPConnectCommand) Help() string {
	//TODO implement me
	panic("implement me")
}

func (c *HCPConnectCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	if c.flagNonInteractiveMethod {
		c.connectHandler = &nonInteractiveConnectHandler{}
	} else {
		c.connectHandler = &interactiveConnectHandler{}
	}

	authStatus, err := c.connectHandler.Connect(args)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Failed to connect to HCP: %s", err))
		return 1
	}

	return 0
}

func (c *HCPConnectCommand) Synopsis() string {
	return "Authenticate to HCP"
}

func (c *HCPConnectCommand) Flags() *flag.FlagSet {
	mainSet := flag.NewFlagSet("", flag.ContinueOnError)
	mainSet.Var(&boolValue{target: &c.flagNonInteractiveMethod}, "non-interactive", "")
	return mainSet
}
