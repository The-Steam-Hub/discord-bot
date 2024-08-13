package cmd

import (
	"fmt"
	"strconv"

	"github.com/KevinFagan/steam-stats/steam"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func Games(s *discordgo.Session, m *discordgo.MessageCreate, steamClient steam.Steam, steamID string) {
	logs := logrus.Fields{
		"command": "games",
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

	embedInfo := EmbededInfo(player[0])
	embedMessage := &discordgo.MessageEmbed{
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Game information is dependent upon the user's privacy settings.",
		},
		Color:     embedInfo.Color,
		Thumbnail: embedInfo.Thumbnail,
		Author:    embedInfo.Author,
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
	s.ChannelMessageSendEmbed(m.ChannelID, embedMessage)
}
