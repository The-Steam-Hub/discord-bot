package cmd

import (
	"strconv"

	"github.com/KevinFagan/steam-stats/steam"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func ID(s *discordgo.Session, m *discordgo.MessageCreate, steamClient steam.Steam, steamID string) {
	logs := logrus.Fields{
		"command": "id",
		"player":  steamID,
		"author":  m.Author.Username,
		"uuid":    uuid.New().String(),
	}

	logrus.WithFields(logs).Info("command recieved")

	player, err := steamClient.GetPlayerSummaries(steamID)
	if err != nil {
		logs["error"] = err
		logrus.WithFields(logs).Error("unable to retrieve player information")
		s.ChannelMessageSend(m.ChannelID, "unable to retrieve player information")
	}

	steamIDInt, _ := strconv.ParseUint(player[0].SteamID, 10, 64)
	steamID64 := steam.SteamID64ToSteamID(steamIDInt)
	steamID3 := steam.SteamID64ToSteamID3(steamIDInt)

	embedInfo := EmbededInfo(player[0])
	embedMessage := &discordgo.MessageEmbed{
		Color:     embedInfo.Color,
		Thumbnail: embedInfo.Thumbnail,
		Author:    embedInfo.Author,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Steam ID",
				Value:  DefaultValue(steamID64),
				Inline: true,
			},
			{
				Name:   "Steam ID3",
				Value:  DefaultValue(steamID3),
				Inline: true,
			},
			{
				Name:   "Steam ID64",
				Value:  DefaultValue(player[0].SteamID),
				Inline: true,
			},
		},
	}
	s.ChannelMessageSendEmbed(m.ChannelID, embedMessage)
}
