package cmd

import (
	"fmt"
	"strconv"

	"github.com/KevinFagan/steam-stats/steam"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func Profile(s *discordgo.Session, m *discordgo.MessageCreate, steamClient steam.Steam, steamID string) {
	logs := logrus.Fields{
		"command": "profile",
		"player":  steamID,
		"author":  m.Author.Username,
		"uuid":    uuid.New().String(),
	}

	logrus.WithFields(logs).Info("command recieved")

	player, err := steamClient.GetPlayerSummariesWithExtra(steamID)
	if err != nil {
		logs["error"] = err
		logrus.WithFields(logs).Error("unable to retrieve player information")
		s.ChannelMessageSend(m.ChannelID, "unable to retrieve player information")
		return
	}

	embedInfo := EmbededInfo(player)
	embedMessage := &discordgo.MessageEmbed{
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Profile information is dependent upon the user's privacy settings.",
		},
		Color:     embedInfo.Color,
		Thumbnail: embedInfo.Thumbnail,
		Author:    embedInfo.Author,
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
	s.ChannelMessageSendEmbed(m.ChannelID, embedMessage)
}

func EmbededInfo(p steam.Player) discordgo.MessageEmbed {
	return discordgo.MessageEmbed{
		Color: 0x66c0f4,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: p.AvatarFull,
		},
		Author: &discordgo.MessageEmbedAuthor{
			Name: fmt.Sprintf("%s %s", p.Status(), p.Name),
			URL:  p.ProfileURL,
		},
	}
}

func DefaultStringValue(value string) string {
	if value == "" {
		return "-"
	}
	return value
}
