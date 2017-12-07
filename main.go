package main

import (
	"github.com/scott-wilson/dosbot"
	discord "github.com/scott-wilson/dosbot-connector-discord"
)

func main() {
	dosbot.RegisterAction(dosbot.EventDirectedMessage, baselineRecalibration, "recite your baseline", "Get the baseline recalibration")
	dosbot.RegisterAction(dosbot.EventDirectedMessage, addRole, "addrole <role>", "Give yourself a role on the server. Dos can only give out non-coloured roles.")
	dosbot.RegisterAction(dosbot.EventDirectedMessage, removeRole, "removerole <role>", "Remove a role you have. Dos can remove any role assocated to you.")
	dosbot.RegisterAction(dosbot.EventDirectedMessage, listRoles, "listroles", "List available roles. All of the roles listed are valid roles can be given to you.")
	dosbot.RegisterConnector(discord.DiscordConnector)
	dosbot.Run()
}
