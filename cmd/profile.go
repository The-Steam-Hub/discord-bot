package cmd

import (
	"fmt"
	"strconv"

	"github.com/KevinFagan/steam-stats/steam"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func Profile(s *discordgo.Session, i *discordgo.InteractionCreate, steamClient steam.Steam, steamID string) {
	logs := logrus.Fields{
		"command": "profile",
		"player":  steamID,
		"uuid":    uuid.New().String(),
	}
	logrus.WithFields(logs).Info("command recieved")

	player, err := steamClient.GetPlayerSummariesWithExtra(steamID)

	if err != nil {
		logs["error"] = err
		logrus.WithFields(logs).Error("unable to retrieve player information")

		// attempt to send the error back to the user
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Unable to retrieve player information",
			},
		})

		// unable to send the message. This could be due to discord permission settings
		logs["error"] = err
		if err != nil {
			logrus.WithFields(logs).Error("unable to send message")
		}

		return
	}

	embedMessage := &discordgo.MessageEmbed{
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

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embedMessage},
		},
	})
}

func DefaultStringValue(value string) string {
	if value == "" {
		return "-"
	}
	return value
}
