package cmd

import (
	"github.com/KevinFagan/steam-stats/steam"
	"github.com/bwmarrin/discordgo"
)

func Games(s *discordgo.Session, m *discordgo.MessageCreate, steam steam.Steam) {
	return
}

// func messageGames(m *discordgo.MessageCreate, s steam.Steam, p steam.Player) discordgo.MessageEmbed {
// 	allGames, err := s.GetOwnedGames(p.SteamID)
// 	if err != nil {
// 		logrus.WithFields(logrus.Fields{
// 			"author":  m.Message.Author.Username,
// 			"channel": m.ChannelID,
// 			"command": m.Content,
// 			"error":   err,
// 		}).Error("unable to retrieve game information")
// 	}

// 	recentGames, err := s.GetRecentlyPlayedGames(p.SteamID)
// 	if err != nil {
// 		logrus.WithFields(logrus.Fields{
// 			"author":  m.Message.Author.Username,
// 			"channel": m.ChannelID,
// 			"command": m.Content,
// 			"error":   err,
// 		}).Error("unable to retrieve game information")
// 	}

// 	mostPlayed := allGames.MostPlayed().Name
// 	if mostPlayed == "" {
// 		mostPlayed = "-"
// 	}
// 	leastPlayed := allGames.LeastPlayed().Name
// 	if leastPlayed == "" {
// 		leastPlayed = "-"
// 	}

// 	embedInfo := embededInfo(p)
// 	embed := &discordgo.MessageEmbed{
// 		Footer: &discordgo.MessageEmbedFooter{
// 			Text: "Game information is dependent upon the user's privacy settings.",
// 		},
// 		Color:     embedInfo.Color,
// 		Thumbnail: embedInfo.Thumbnail,
// 		Author:    embedInfo.Author,
// 		Fields: []*discordgo.MessageEmbedField{
// 			{
// 				Name:   "Total Playtime",
// 				Value:  fmt.Sprintf("%dh", allGames.TotalHoursPlayed()),
// 				Inline: true,
// 			},
// 			{
// 				Name:   "Most Played Game",
// 				Value:  mostPlayed,
// 				Inline: true,
// 			},
// 			{
// 				Name:   "Least Played Game",
// 				Value:  leastPlayed,
// 				Inline: true,
// 			},
// 			{
// 				Name:   "Games Owned",
// 				Value:  strconv.Itoa(len(allGames.Games)),
// 				Inline: true,
// 			},
// 			{
// 				Name:   "Games Played",
// 				Value:  strconv.Itoa(allGames.GamesPlayed()),
// 				Inline: true,
// 			},
// 			{
// 				Name:   "Games Not Played",
// 				Value:  strconv.Itoa(allGames.GamesNotPlayed()),
// 				Inline: true,
// 			},
// 			{
// 				Name:   "2 Week Playtime",
// 				Value:  fmt.Sprintf("%dh", recentGames.RecentHoursPlayed()),
// 				Inline: true,
// 			},
// 			{
// 				Name:   "2 Week Games Played",
// 				Value:  strconv.Itoa(len(recentGames.Games)),
// 				Inline: true,
// 			},
// 			{
// 				Name:   "",
// 				Value:  "",
// 				Inline: true,
// 			},
// 		},
// 	}
// 	return *embed
// }
