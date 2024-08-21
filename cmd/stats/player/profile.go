package player

import (
	"fmt"
	"strconv"

	"github.com/KevinFagan/steam-stats/steam"
	"github.com/bwmarrin/discordgo"
)

func ProfileEmbeddedMessage(steamClient steam.Steam, steamID string) (*discordgo.MessageEmbed, error) {
	id, err := steamClient.ResolveID(steamID)
	if err != nil {
		return nil, err
	}

	player, err := steamClient.GetPlayerSummariesWithExtra(id)
	if err != nil {
		return nil, err
	}

	embMsg := &discordgo.MessageEmbed{
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Profile information is dependent upon the user's privacy settings.",
		},
		Color: 0x66c0f4,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: player.AvatarFull,
		},
		Author: &discordgo.MessageEmbedAuthor{
			Name: fmt.Sprintf("%s %s", player.Status(), player.Name),
			URL:  player.ProfileURL,
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Real Name",
				Value:  DefaultStringValue(player.RealName),
				Inline: true,
			},
			{
				Name:   "Country Code",
				Value:  DefaultStringValue(player.CountryCode),
				Inline: true,
			},
			{
				Name:   "State Code",
				Value:  DefaultStringValue(player.StateCode),
				Inline: true,
			},
			{
				Name:   "Profile Age",
				Value:  DefaultStringValue(player.ProfileAge()),
				Inline: true,
			},
			{
				Name:   "Last Seen",
				Value:  DefaultStringValue(player.LastSeen()),
				Inline: true,
			},
			{
				Name:   "Level",
				Value:  strconv.Itoa(player.PlayerLevel),
				Inline: true,
			},
			{
				Name:   "Level Percentile",
				Value:  strconv.FormatFloat(player.PlayerLevelPercentile, 'f', 2, 64),
				Inline: true,
			},
			{
				Name:   "Total XP",
				Value:  strconv.Itoa(player.PlayerXP),
				Inline: true,
			},
			{
				Name:   "XP To Next Level",
				Value:  strconv.Itoa(player.PlayerXPNeededToLevelUp),
				Inline: true,
			},
		},
	}
	return embMsg, nil
}

func DefaultStringValue(value string) string {
	if value == "" {
		return "-"
	}
	return value
}
