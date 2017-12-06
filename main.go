package main

import (
	"github.com/scott-wilson/dosbot"
	discord "github.com/scott-wilson/dosbot-connector-discord"
)

func main() {
	dosbot.RegisterAction(dosbot.EventDirectedMessage, baselineRecalibration, "Baseline Recalibration", "Get the baseline recalibration", "Within cells interlinked")
	dosbot.RegisterAction(dosbot.EventDirectedMessage, addRole, "Add Role to Self", "Give yourself a role with @dos give me role <rolename>", "Dos can only give out non-coloured roles.")
	dosbot.RegisterAction(dosbot.EventDirectedMessage, removeRole, "Remove Role from Self", "Remove a role you have with @dos remove my role <rolename>", "Dos can remove any role assocated to you.")
	dosbot.RegisterAction(dosbot.EventDirectedMessage, listRoles, "List Roles", "List available roles with @dos list roles", "All of the roles listed are valid roles that Dos can give you.")
	dosbot.RegisterConnector(discord.DiscordConnector)
	dosbot.Run()
}
