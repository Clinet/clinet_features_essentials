package essentials

//NOTE: It is very important to defer registering
// the essentials feature until *after* everything
// else has been registered. This ensures a proper
// command list, as well as anything else prepped
// ahead of time.

///help
//**/help topic:voice**
//This wouldn't show up because no Feature.Desc
//**/help topic:moderation**
//Allows easy-going and simple control over a community.
//*/help page:2*

///help page:2
//**/help topic:hellodolly**
//Inspired by the WordPress sample plugin! Responds with a random lyric from Louis Armstrong's Hello, Dolly.
//*/help page:1 | /help page:3*

///help topic:moderation
//**/ban**
//Bans a given user
//**/hackban**
//Hackbans a given user ID
//**/kick**
//Kicks a given user
//*/help cmd:ban | /help topic:moderation page:2*

///help topic:moderation page:2
//**/warn**
//Warns a given user
//*/help topic:moderation page:1 | /help topic:moderation page:3*

///help cmd:ban
//**/ban @user <reason> <rule>**
//*user* - Who to actually ban (required)
//*reason* - Reason for the ban (default "No reason provided.")
//*rule* - Rule broken that led to ban (default -1)
//*/help cmd:ban page:2*

//TODO: Examples for commands will come later!
///help cmd:warn page:2
//**/warn user:@user**
//Warns the user @user, (kicking/banning) them after too many warnings
//**/warn user:@user reason:"Not funny" rule:4**
//Warns the user @user for breaking rule 4, with a custom reason attached
//*/help cmd:ban page:1 | /help cmd:ban page:3*

import (
	"fmt"
	//"math"
	"strings"

	"github.com/Clinet/clinet_cmds"
	"github.com/Clinet/clinet_features"
	"github.com/Clinet/clinet_services"
	//"github.com/Clinet/clinet_storage"
	"github.com/JoshuaDoes/logger"
)

var cmdHelp = cmds.NewCmd("help", "Returns a help page for a topic or a command", handleHelp).AddArgs(
	cmds.NewCmdArg("topic", "The feature topic to browse help pages for", ""),
	cmds.NewCmdArg("cmd", "The feature command to browse help pages for", ""),
	cmds.NewCmdArg("page", "The page to read of the selected topic/command", 1),
)

var Feature = features.Feature{
	Name: "essentials",
	Desc: "Provides essential commands every bot should have!",
	Cmds: []*cmds.Cmd{cmdHelp},
	/*Init: func() error {
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
			if feature.Desc == "" {
				continue //Skip features that don't provide descriptions
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
	},*/
}

var Log *logger.Logger
func init() {
	Log = logger.NewLogger("essentials", 2)
}

var msgErrInvalidCmd = &services.Message{Title: "Help Error", Content: "Invalid command!", Color: 0xB40000}
var msgErrInvalidSubCmd = &services.Message{Title: "Help Error", Content: "Invalid subcommand!", Color: 0xB40000}

func handleHelp(ctx *cmds.CmdCtx) *cmds.CmdResp {
	Log.Trace("help -")
	cmdPrefix := ctx.Service.CmdPrefix()
	Log.Trace("help -")

	argTopic := ctx.GetArg("topic").GetString()
	argCmd := ctx.GetArg("cmd").GetString()
	//argPage := ctx.GetArg("page").GetInt()
	Log.Trace("help -")

	if argCmd != "" {
		Log.Trace("help --")
		cmdSplit := strings.Split(argCmd, " ")
		cmd := cmds.GetCmd(cmdSplit[0])
		if cmd == nil {
			return cmds.CmdRespFromMsg(msgErrInvalidCmd)
		}
		cmdUsage := cmdPrefix + cmd.Name
		if len(cmdSplit) > 1 {
			cmd = cmd.GetSubCmd(cmdSplit[1])
			if cmd == nil {
				return cmds.CmdRespFromMsg(msgErrInvalidSubCmd)
			}
			cmdUsage += " " + cmd.Name
		}
		Log.Trace("help --")

		cmdArgs := "No arguments necessary."

		if len(cmd.Subcommands) > 0 {
			cmdUsage = "Get help with a subcommand."
			cmdArgs = ""
			for i := 0; i < len(cmd.Subcommands); i++ {
				if cmdArgs != "" {
					cmdArgs += "\n"
				}
				cmdArgs += cmdPrefix + ctx.Alias + " cmd:`" + cmd.Name + " " + cmd.Subcommands[i].Name + "`"
			}
		} else if len(cmd.Args) > 0 {
			Log.Trace("help ---")
			cmdArgs = ""
			for i := 0; i < len(cmd.Args); i++ {
				Log.Trace("help ----")
				cmdUsage += " "
				if !cmd.Args[i].Required {
					cmdUsage += "<"
				}
				cmdUsage += cmd.Args[i].Name
				switch cmd.Args[i].Value.(type) {
				case *services.User:
					cmdUsage += ":`@user`"
				case *services.Role:
					cmdUsage += ":`@role`"
				case *services.Channel:
					cmdUsage += ":`#channel`"
				default:
					cmdUsage += fmt.Sprintf(":`%v`", cmd.Args[i].Value)
				}
				if !cmd.Args[i].Required {
					cmdUsage += ">"
				}

				if cmdArgs != "" {
					cmdArgs += "\n"
				}
				cmdArgs += "*" + cmd.Args[i].Name + "* - " + cmd.Args[i].Description
				if cmd.Args[i].Required {
					cmdArgs += " (required)"
				}
			}
		}

		Log.Trace("help --")
		msgCmd := services.NewMessage().
			SetTitle(cmdUsage).
			SetContent(cmdArgs).
			SetColor(0x1C1C1C)

		return cmds.CmdRespFromMsg(msgCmd)
	} else if argTopic != "" {

	}

	listTopic := ""
	enabledFeatures := features.GetEnabledFeatures()
	for i := 0; i < len(enabledFeatures); i++ {
		if listTopic != "" {
			listTopic += ", "
		}
		listTopic += enabledFeatures[i]
	}
	if listTopic == "" {
		listTopic = "No feature topics are available!"
	}

	listCmd := ""
	for i := 0; i < len(cmds.Commands); i++ {
		if listCmd != "" {
			listCmd += ", "
		}
		listCmd += cmds.Commands[i].Name
	}
	if listCmd == "" {
		listCmd = "No commands are available!"
	}

	msgHelp := services.NewMessage().
		SetTitle("Help!").
		SetContent("Get help with a feature topic or a particular command.").
		AddField("/help topic:`name`", "*name* - " + listTopic).
		AddField("/help cmd:`name`", "*name* - " + listCmd).
		SetColor(0x1C1C1C)

	return cmds.CmdRespFromMsg(msgHelp)
}