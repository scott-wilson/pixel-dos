package main

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/scott-wilson/dosbot"
	discord "github.com/scott-wilson/dosbot-connector-discord"
)

var (
	baselineRecalibrationMessageResponse = map[*regexp.Regexp]string{
		regexp.MustCompile(`(?i)recite your baseline`):                                                    "And blood-black nothingness began to spin... A system of cells interlinked within cells interlinked within cells interlinked within one stem... And dreadfully distinct against the dark, a tall white fountain played.",
		regexp.MustCompile(`(?i)cells`):                                                                   "Cells.",
		regexp.MustCompile(`(?i)have you been in an institution`):                                         "Cells.",
		regexp.MustCompile(`(?i)do they keep you in a cell`):                                              "Cells.",
		regexp.MustCompile(`(?i)when you're not performing your duties do they keep you in a little box`): "Cells.",
		regexp.MustCompile(`(?i)interlinked`):                                                             "Interlinked.",
		regexp.MustCompile(`(?i)what's it like to hold the hand of someone you love`):                     "Interlinked.",
		regexp.MustCompile(`(?i)did they teach you how to feel finger to finger`):                         "Interlinked.",
		regexp.MustCompile(`(?i)do you long for having your heart interlinked`):                           "Interlinked.",
		regexp.MustCompile(`(?i)do you dream about being interlinked`):                                    "Interlinked.",
		regexp.MustCompile(`(?i)what's it like to hold your child in your arms`):                          "Interlinked.",
		regexp.MustCompile(`(?i)do you feel that there's a part of you that's missing`):                   "Interlinked.",
		regexp.MustCompile(`(?i)within cells interlinked`):                                                "Within cells interlinked.",
		regexp.MustCompile(`(?i)why don't you say within cells interlinked three times`):                  "Within cells interlinked. Within cells interlinked. Within cells interlinked.",
	}
	addRoleMessageRegex    = regexp.MustCompile(`(?i)give me role (.+)`)
	removeRoleMessageRegex = regexp.MustCompile(`(?i)remove my role (.+)`)
	listRolesMessageRegex  = regexp.MustCompile(`(?i)list roles`)
)

func baselineRecalibration(event dosbot.Event) error {
	bot := event.Bot()
	room := event.Room()
	sender := event.Sender()
	message := event.Message()

	for key, value := range baselineRecalibrationMessageResponse {
		if key.FindString(message) != "" {
			return bot.SendDirectMessage(room, sender, value)
		}
	}

	return nil
}

func addRole(event dosbot.Event) error {
	room := event.Room()
	message := event.Message()
	result := addRoleMessageRegex.FindAllStringSubmatch(message, -1)

	if len(result) == 0 {
		return nil
	}

	roleName := result[0][1]

	bot := event.Bot().(discord.Bot)
	sender := event.Sender()
	session := bot.Session()
	channel, err := session.Channel(room.ID().(string))

	if err != nil {
		return err
	}

	guild, err := session.Guild(channel.GuildID)

	if err != nil {
		return err
	}

	for _, role := range guild.Roles {
		if role.Name == roleName {
			if !isValidRole(role) {
				bot.SendDirectMessage(room, sender, "I'm sorry, but I don't have permissions to give you that role. :frowning:")
				return nil
			}

			if err := session.GuildMemberRoleAdd(channel.GuildID, sender.ID().(string), role.ID); err != nil {
				return err
			}

			bot.SendDirectMessage(room, sender, "Done! :grinning:")

			return nil
		}
	}

	bot.SendDirectMessage(room, sender, "I'm sorry, but I didn't find that role. :frowning:")

	return nil
}

func removeRole(event dosbot.Event) error {
	room := event.Room()
	message := event.Message()
	result := removeRoleMessageRegex.FindAllStringSubmatch(message, -1)

	if len(result) == 0 {
		return nil
	}

	roleName := result[0][1]

	bot := event.Bot().(discord.Bot)
	sender := event.Sender()
	session := bot.Session()
	channel, err := session.Channel(room.ID().(string))

	if err != nil {
		return err
	}

	guild, err := session.Guild(channel.GuildID)

	if err != nil {
		return err
	}

	for _, role := range guild.Roles {
		if role.Name == roleName {
			if err := session.GuildMemberRoleRemove(channel.GuildID, sender.ID().(string), role.ID); err != nil {
				return err
			}

			bot.SendDirectMessage(room, sender, "Done! :grinning:")

			return nil
		}
	}

	bot.SendDirectMessage(room, sender, "I'm sorry, but I didn't find that role. :frowning:")

	return nil
}

func listRoles(event dosbot.Event) error {
	room := event.Room()
	message := event.Message()
	if result := listRolesMessageRegex.FindString(message); result == "" {
		return nil
	}

	bot := event.Bot().(discord.Bot)
	sender := event.Sender()
	session := bot.Session()
	channel, err := session.Channel(room.ID().(string))

	if err != nil {
		return err
	}

	guild, err := session.Guild(channel.GuildID)

	if err != nil {
		return err
	}

	roles := make([]string, 0)

	for _, role := range guild.Roles {
		if !isValidRole(role) {
			continue
		}

		roles = append(roles, fmt.Sprintf(" - %s", role.Name))
	}

	sort.Strings(roles)
	message = strings.Join([]string{
		"Here are the available roles:",
		strings.Join(roles, "\n"),
	}, "\n")

	bot.SendDirectMessage(room, sender, message)

	return nil
}

func isValidRole(role *discordgo.Role) bool {
	if role.Managed || role.Color != 0 || role.Name == "@everyone" {
		return false
	}
	return true
}
