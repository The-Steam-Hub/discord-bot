package cmd

import (
	"fmt"
	"strconv"

	"github.com/KevinFagan/steam-stats/steam"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func Games(s *discordgo.Session, i *discordgo.InteractionCreate, steamClient steam.Steam, steamID string) {
	logs := logrus.Fields{
		"command": "games",
		"player":  steamID,
		"uuid":    uuid.New().String(),
	}

	logrus.WithFields(logs).Info("command recieved")

	player, err := steamClient.GetPlayerSummaries(steamID)
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

	allGames, err := steamClient.GetOwnedGames(player[0].SteamID)
	if err != nil {
		logs["error"] = err
		logrus.WithFields(logs).Error("unable to retrieve game information")
	}

	recentGames, err := steamClient.GetRecentlyPlayedGames(player[0].SteamID)
	if err != nil {
		logs["error"] = err
		logrus.WithFields(logs).Error("unable to retrieve recent game information")
	}

	embedMessage := &discordgo.MessageEmbed{
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Game information is dependent upon the user's privacy settings.",
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
				Name:   "Total Playtime",
				Value:  fmt.Sprintf("%dh", allGames.TotalHoursPlayed()),
				Inline: true,
			},
			{
				Name:   "Most Played Game",
				Value:  DefaultStringValue(allGames.MostPlayed().Name),
				Inline: true,
			},
			{
				Name:   "Least Played Game",
				Value:  DefaultStringValue(allGames.LeastPlayed().Name),
				Inline: true,
			},
			{
				Name:   "Games Owned",
				Value:  strconv.Itoa(len(allGames.Games)),
				Inline: true,
			},
			{
				Name:   "Games Played",
				Value:  strconv.Itoa(allGames.GamesPlayed()),
				Inline: true,
			},
			{
				Name:   "Games Not Played",
				Value:  strconv.Itoa(allGames.GamesNotPlayed()),
				Inline: true,
			},
			{
				Name:   "Recent Playtime",
				Value:  fmt.Sprintf("%dh", recentGames.RecentHoursPlayed()),
				Inline: true,
			},
			{
				Name:   "Recent Games Played",
				Value:  strconv.Itoa(len(recentGames.Games)),
				Inline: true,
			},
			{
				Name:   "",
				Value:  "",
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
