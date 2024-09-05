package player

import (
	"fmt"
	"strconv"

	"github.com/KevinFagan/steam-stats/cmd"
	"github.com/KevinFagan/steam-stats/steam"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func PlayerBans(session *discordgo.Session, interaction *discordgo.InteractionCreate, steamClient steam.Steam, input string) {
	logs := logrus.Fields{
		"input":  input,
		"author": interaction.Member.User.Username,
		"uuid":   uuid.New(),
	}

	id, err := steamClient.ResolveSteamID(input)
	if err != nil {
		logs["error"] = err
		errMsg := "unable to resolve player ID"
		logrus.WithFields(logs).Error(errMsg)
		cmd.HandleMessageError(session, interaction, &logs, errMsg)
		return
	}

	player, err := steamClient.PlayerSummaries(id)
	if err != nil {
		logs["error"] = err
		errMsg := "unable to retrieve player summary"
		logrus.WithFields(logs).Error(errMsg)
		cmd.HandleMessageError(session, interaction, &logs, errMsg)
		return
	}

	err = steamClient.PlayerBans(&player[0])
	if err != nil {
		logs["error"] = err
		logrus.WithFields(logs).Error("unable to retieve player ban information")
	}

	embMsg := &discordgo.MessageEmbed{
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
				Name:   "VAC Banned",
				Value:  strconv.FormatBool(player[0].VACBanned),
				Inline: true,
			},
			{
				Name:   "# Of VAC Bans",
				Value:  strconv.Itoa(player[0].NumOfVacBans),
				Inline: true,
			},
			{
				Name:   "# Of Game Bans",
				Value:  strconv.Itoa(player[0].NumOfGameBans),
				Inline: true,
			},
			{
				Name:   "Days Since Last Ban",
				Value:  fmt.Sprintf("%dd", player[0].DaysSinceLastBan),
				Inline: true,
			},
			{
				Name:   "Community Banned",
				Value:  strconv.FormatBool(player[0].CommunityBanned),
				Inline: true,
			},
			{
				Name:   "Economy Banned",
				Value:  player[0].EconomyBan,
				Inline: true,
			},
		},
	}
	cmd.HandleMessageOk(embMsg, session, interaction, &logs)
}
