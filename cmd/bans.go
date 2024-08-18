package cmd

import (
	"fmt"
	"strconv"

	"github.com/KevinFagan/steam-stats/steam"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func Bans(s *discordgo.Session, i *discordgo.InteractionCreate, steamClient steam.Steam, steamID string) {
	logs := logrus.Fields{
		"command": "bans",
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
				Name:   "VAC Banned",
				Value:  strconv.FormatBool(player.VACBanned),
				Inline: true,
			},
			{
				Name:   "# Of VAC Bans",
				Value:  strconv.Itoa(player.NumOfVacBans),
				Inline: true,
			},
			{
				Name:   "# Of Game Bans",
				Value:  strconv.Itoa(player.NumOfGameBans),
				Inline: true,
			},
			{
				Name:   "Days Since Last Ban",
				Value:  fmt.Sprintf("%dd", player.DaysSinceLastBan),
				Inline: true,
			},
			{
				Name:   "Community Banned",
				Value:  strconv.FormatBool(player.CommunityBanned),
				Inline: true,
			},
			{
				Name:   "Economy Banned",
				Value:  player.EconomyBan,
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
