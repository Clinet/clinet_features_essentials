package essentials

//NOTE: It is very important to defer registering
// the essentials feature until *after* everything
// else has been registered. This ensures a proper
// command list, as well as anything else prepped
// ahead of time.

import (
	"fmt"
	"math"

	"github.com/Clinet/clinet_cmds"
	"github.com/Clinet/clinet_features"
	"github.com/Clinet/clinet_storage"
)

var cmdHelp = cmds.NewCmd("help", "Returns a help page for a feature", nil)
var cmdCmds = cmds.NewCmd("cmds", "Returns a page of available commands", nil)

var Feature = features.Feature{
	Help: `Using the essentials feature is easy!

We provide 2 commands:

/help (feature)
- Returns the specified feature's help text
 - Ex: /help essentials
 - Ex: /help moderation
 - Ex: /help voice
 - Ex: /help music
/cmds (page/cmd)
- Returns the specified page of commands from the command list
- Also returns the specified command's description and usage requirements
 - Ex: /cmds help
 - Ex: /cmds cmds
 - Ex: /cmds ban`,
	Name: "essentials",
	Cmds: []*cmds.Cmd{cmdHelp, cmdCmds},
	Init: func() error {
		cfg := &storage.Storage{}
		if err := cfg.LoadFrom("essentials"); err != nil {
			return err
		}

		cmdsPerPage := float64(5)
		rawCmdsPerPage, err := cfg.ConfigGet("cmds", "cmdsPerPage")
		if err != nil {
			cfg.ConfigSet("cmds", "cmdsPerPage", cmdsPerPage)
		} else {
			cmdsPerPage = rawCmdsPerPage.(float64)
		}

		for i := 0; i < len(features.FM.Features); i++ {
			feature := features.FM.Features[i]
			if feature.Help == "" {
				continue //Skip features that don't provide end-user help texts
			}
			cmdHelp.AddSubCmds(cmds.NewCmd(feature.Name, "Explains how to use " + feature.Name, handleHelp))
		}

		//Init is called before UpdateFeature is called on us, so add our help topic
		cmdHelp.AddSubCmds(cmds.NewCmd("essentials", "Explains how to use essentials", handleHelp))

		//Init is called after loading the commands in, so we don't need to add our own commands here
		pageCount := int(math.Ceil(float64(len(cmds.Commands)) / float64(cmdsPerPage)))
		for i := 0; i < pageCount; i++ {
			pageNum := fmt.Sprintf("%d", i+1)
			cmdCmds.AddSubCmds(cmds.NewCmd(pageNum, fmt.Sprintf("Returns page %s/%d of %d total commands", pageNum, pageCount, len(cmds.Commands)), handleCmds))
		}

		return nil
	},
}

func handleHelp(ctx *cmds.CmdCtx) *cmds.CmdResp {
	return nil
}
func handleCmds(ctx *cmds.CmdCtx) *cmds.CmdResp {
	return nil
}