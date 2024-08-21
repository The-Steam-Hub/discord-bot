package player

import (
	"fmt"
	"strconv"

	"github.com/KevinFagan/steam-stats/steam"
	"github.com/bwmarrin/discordgo"
)

func PlayerProfile(steamClient steam.Steam, steamID string) (*discordgo.MessageEmbed, error) {
	id, err := steamClient.ResolveID(steamID)
	if err != nil {
		// log error
		return nil, err
	}

	player, err := steamClient.PlayerSummaries(id)
	if err != nil {
		// log error
		return nil, err
	}

	err = steamClient.PlayerBadges(&player[0])
	if err != nil {
		// log error
	}

	err = steamClient.PlayerLevelDistribution(&player[0])
	if err != nil {
		// log error
	}

	embMsg := &discordgo.MessageEmbed{
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Profile information is dependent upon the user's privacy settings.",
		},
		Color: 0x66c0f4,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: player[0].AvatarFull,
		},
		Author: &discordgo.MessageEmbedAuthor{
			Name: fmt.Sprintf("%s %s", player[0].Status(), player[0].Name),
			URL:  player[0].ProfileURL,
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Real Name",
				Value:  DefaultStringValue(player[0].RealName),
				Inline: true,
			},
			{
				Name:   "Country Code",
				Value:  DefaultStringValue(player[0].CountryCode),
				Inline: true,
			},
			{
				Name:   "State Code",
				Value:  DefaultStringValue(player[0].StateCode),
				Inline: true,
			},
			{
				Name:   "Profile Age",
				Value:  DefaultStringValue(player[0].ProfileAge()),
				Inline: true,
			},
			{
				Name:   "Last Seen",
				Value:  DefaultStringValue(player[0].LastSeen()),
				Inline: true,
			},
			{
				Name:   "Level",
				Value:  strconv.Itoa(player[0].PlayerLevel),
				Inline: true,
			},
			{
				Name:   "Level Percentile",
				Value:  strconv.FormatFloat(player[0].PlayerLevelPercentile, 'f', 2, 64),
				Inline: true,
			},
			{
				Name:   "Total XP",
				Value:  strconv.Itoa(player[0].PlayerXP),
				Inline: true,
			},
			{
				Name:   "XP To Next Level",
				Value:  strconv.Itoa(player[0].PlayerXPNeededToLevelUp),
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
