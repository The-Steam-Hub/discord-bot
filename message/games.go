package message

import (
	"fmt"
	"strconv"

	"github.com/KevinFagan/steam-stats/steam"
	"github.com/bwmarrin/discordgo"
)

func GamesEmbeddedMessage(steamClient steam.Steam, steamID string) (*discordgo.MessageEmbed, error) {
	id, err := steamClient.ResolveID(steamID)
	if err != nil {
		return nil, err
	}

	player, err := steamClient.GetPlayerSummaries(id)
	if err != nil {
		return nil, err
	}

	allGames, _ := steamClient.GetOwnedGames(player[0].SteamID)
	recentGames, _ := steamClient.GetRecentlyPlayedGames(player[0].SteamID)

	embMsg := &discordgo.MessageEmbed{
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

	return embMsg, nil
}
