package player

import (
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/the-steam-bot/discord-bot/cmd"
	"github.com/the-steam-bot/discord-bot/steam"
)

func PlayerProfile(session *discordgo.Session, interaction *discordgo.InteractionCreate, steamClient steam.Steam, input string) {
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

	err = steamClient.PlayerBadges(&player[0])
	if err != nil {
		logs["error"] = err
		logrus.WithFields(logs).Error("unable to retireve player badges")
	}

	err = steamClient.PlayerLevelDistribution(&player[0])
	if err != nil {
		logs["error"] = err
		logrus.WithFields(logs).Error("unable to retireve player level distribution")
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
				Value:  cmd.HandleStringDefault(player[0].RealName),
				Inline: true,
			},
			{
				Name:   "Country Code",
				Value:  cmd.HandleStringDefault(player[0].CountryCode),
				Inline: true,
			},
			{
				Name:   "State Code",
				Value:  cmd.HandleStringDefault(player[0].StateCode),
				Inline: true,
			},
			{
				Name:   "Profile Age",
				Value:  cmd.HandleStringDefault(player[0].ProfileAge()),
				Inline: true,
			},
			{
				Name:   "Last Seen",
				Value:  cmd.HandleStringDefault(player[0].LastSeen()),
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
	cmd.HandleMessageOk(embMsg, session, interaction, &logs)
}
