package game

import (
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Check for alphanumeric characters
var reg = regexp.MustCompile(`^[A-Za-z0-9 ]+$`)

// Boolean return types
func isValidGamingRole(selectedId string, role *discordgo.Role) bool {
	return selectedId == role.ID && role.Color != 0x9c9c9c
}

func playercountAndTimeReached(message *discordgo.Message, activePlayers []string, playersRequired int) bool {
	if message.Embeds[0].Fields[2] != nil {
		return message.Embeds[0].Fields[2].Value == "Now!" && isPlayerCountReached(activePlayers, playersRequired)
	}
	return false
}

func playerNotInList(index int) bool {
	return index == -1
}

func playerInList(index int) bool {
	return index != -1
}

func isBackupPlayer(player string) bool {
	return strings.HasSuffix(player, "\u200b") && player != "\u200b"
}

func isPlayerCountReached(activePlayers []string, playersRequired int) bool {
	return len(activePlayers) >= playersRequired
}

func areBackupsPresent(players []string) bool {
	return len(players) != 0
}

func intWithinLimits(toCheck int, lowerBound int, upperBound int) bool {
	return lowerBound <= toCheck && toCheck <= upperBound
}

func floatWithinLimits(toCheck float64, lowerBound float64, upperBound float64) bool {
	return lowerBound <= toCheck && toCheck <= upperBound
}

// switch contents
func getGameErrorResponse(option *discordgo.ApplicationCommandInteractionDataOption) string {
	var name string
	var ok bool
	name, ok = option.Value.(string)

	if !ok {
		return "Please enter a valid name."
	}
	if intWithinLimits(len(name), 2, 25) {
		return "Your game needs to be between 2-25 characters long"
	}
	if !reg.MatchString(name) {
		return "Your game cannot contain any special characters"
	}

	return ""
}

func getAmountErrorResponse(option *discordgo.ApplicationCommandInteractionDataOption) string {
	var amount float64
	var ok bool

	amount, ok = option.Value.(float64)

	if !ok {
		return "Please enter a valid amount."
	}
	if floatWithinLimits(amount, 2, 40) {
		return "Your game needs to contain between 2-40 players."
	}
	return ""
}

func getNotifyroleErrorResponse(option *discordgo.ApplicationCommandInteractionDataOption, roles []*discordgo.Role) string {
	var selectedRoleId string
	var ok bool

	selectedRoleId, ok = option.Value.(string)

	if !ok {
		return "Please enter a valid role."
	}

	for _, role := range roles {
		if isValidGamingRole(selectedRoleId, role) {
			return "Please enter a valid gaming role."
		}
	}

	return ""
}

func getTimeErrorResponse(option *discordgo.ApplicationCommandInteractionDataOption) string {
	var timeString string
	var ok bool

	timeString, ok = option.Value.(string)

	if !ok {
		return "Please enter a valid time in format hh:mm."
	}
	if _, err := time.Parse("15:04", timeString); err != nil {
		return "Please enter your time in format hh:mm (For example 15:50)"
	}

	return ""
}
