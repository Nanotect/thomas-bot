package game

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

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
