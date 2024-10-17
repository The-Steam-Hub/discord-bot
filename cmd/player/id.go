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

func PlayerID(session *discordgo.Session, interaction *discordgo.InteractionCreate, steamClient steam.Steam, input string) {
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

	steamIDInt, _ := strconv.ParseUint(player[0].SteamID, 10, 64)

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
				Name:   "Steam ID",
				Value:  cmd.HandleStringDefault(steam.SteamID64ToSteamID(steamIDInt)),
				Inline: true,
			},
			{
				Name:   "Steam ID3",
				Value:  cmd.HandleStringDefault(steam.SteamID64ToSteamID3(steamIDInt)),
				Inline: true,
			},
			{
				Name:   "Steam ID64",
				Value:  cmd.HandleStringDefault(player[0].SteamID),
				Inline: true,
			},
		},
	}
	cmd.HandleMessageOk(embMsg, session, interaction, &logs)
}
