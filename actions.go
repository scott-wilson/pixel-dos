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

// Regexes dos will pay attention to
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
	addRoleMessageRegex    = regexp.MustCompile(`(?i)addrole (.+)`)
	removeRoleMessageRegex = regexp.MustCompile(`(?i)removerole (.+)`)
	listRolesMessageRegex  = regexp.MustCompile(`(?i)listroles`)
)

// Helper regexes
var (
	splitRoleRegex = regexp.MustCompile(` *, *`)
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
	result := addRoleMessageRegex.FindStringSubmatch(message)

	if len(result) == 0 {
		return nil
	}

	roleNames := splitRoleRegex.Split(result[1], -1)

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

	roleNameMap := make(map[string]*discordgo.Role)

	for _, role := range guild.Roles {
		roleNameMap[role.Name] = role
	}

	validRoles := make([]string, 0)

	for _, roleName := range roleNames {
		role, ok := roleNameMap[roleName]

		if !ok || !isValidRole(role) {
			continue
		}

		if err := session.GuildMemberRoleAdd(channel.GuildID, sender.ID().(string), role.ID); err != nil {
			return err
		}

		validRoles = append(validRoles, roleName)
	}

	if len(validRoles) > 0 {
		sort.Strings(validRoles)
		validRolesMessage := strings.Join(validRoles, ", ")
		bot.SendDirectMessage(room, sender, fmt.Sprintf("I've added the following roles to your account: %s :grinning:", validRolesMessage))
	} else {
		bot.SendDirectMessage(room, sender, "I'm sorry, but I could not do what you asked. :frowning:")
	}

	return nil
}

func removeRole(event dosbot.Event) error {
	room := event.Room()
	message := event.Message()
	result := removeRoleMessageRegex.FindStringSubmatch(message)

	if len(result) == 0 {
		return nil
	}

	roleNames := splitRoleRegex.Split(result[1], -1)

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

	roleNameMap := make(map[string]*discordgo.Role)

	for _, role := range guild.Roles {
		roleNameMap[role.Name] = role
	}

	validRoles := make([]string, 0)

	for _, roleName := range roleNames {
		role, ok := roleNameMap[roleName]

		if !ok {
			continue
		}

		if err := session.GuildMemberRoleRemove(channel.GuildID, sender.ID().(string), role.ID); err != nil {
			return err
		}

		validRoles = append(validRoles, roleName)
	}

	if len(validRoles) > 0 {
		sort.Strings(validRoles)
		validRolesMessage := strings.Join(validRoles, ", ")
		bot.SendDirectMessage(room, sender, fmt.Sprintf("I've removed the following roles from your account: %s :grinning:", validRolesMessage))
	} else {
		bot.SendDirectMessage(room, sender, "I'm sorry, but I could not do what you asked. :frowning:")
	}

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
