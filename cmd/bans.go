package cmd

import (
	"fmt"
	"strconv"

	"github.com/KevinFagan/steam-stats/steam"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func Bans(s *discordgo.Session, m *discordgo.MessageCreate, steamClient steam.Steam, steamID string) {
	logs := logrus.Fields{
		"command": "bans",
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
	}

	embedInfo := EmbededInfo(player)
	embedMessage := &discordgo.MessageEmbed{
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Ban information is dependent upon the user's privacy settings.",
		},
		Color:     embedInfo.Color,
		Thumbnail: embedInfo.Thumbnail,
		Author:    embedInfo.Author,
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
	s.ChannelMessageSendEmbed(m.ChannelID, embedMessage)
}
