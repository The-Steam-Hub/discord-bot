package cmd

import (
	"fmt"
	"strconv"

	"github.com/KevinFagan/steam-stats/steam"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func ID(s *discordgo.Session, i *discordgo.InteractionCreate, steamClient steam.Steam, steamID string) {
	logs := logrus.Fields{
		"command": "id",
		"player":  steamID,
		"uuid":    uuid.New().String(),
	}

	logrus.WithFields(logs).Info("command recieved")

	player, err := steamClient.GetPlayerSummaries(steamID)
	if err != nil {
		logs["error"] = err
		logrus.WithFields(logs).Error("unable to retrieve player information")
		s.ChannelMessageSend(i.ChannelID, "unable to retrieve player information")
		return
	}

	steamIDInt, _ := strconv.ParseUint(player[0].SteamID, 10, 64)
	steamID64 := steam.SteamID64ToSteamID(steamIDInt)
	steamID3 := steam.SteamID64ToSteamID3(steamIDInt)

	embedMessage := &discordgo.MessageEmbed{
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
				Value:  DefaultStringValue(steamID64),
				Inline: true,
			},
			{
				Name:   "Steam ID3",
				Value:  DefaultStringValue(steamID3),
				Inline: true,
			},
			{
				Name:   "Steam ID64",
				Value:  DefaultStringValue(player[0].SteamID),
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
