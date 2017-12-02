package main

import (
	"github.com/scott-wilson/dosbot"
	discord "github.com/scott-wilson/dosbot-connector-discord"
)

func main() {
	dosbot.RegisterAction(dosbot.EventDirectedMessage, baselineRecalibration)
	dosbot.RegisterConnector(discord.DiscordConnector)
	dosbot.Run()
}
