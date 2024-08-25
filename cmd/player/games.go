package player

import (
	"fmt"
	"strconv"

	"github.com/KevinFagan/steam-stats/steam"
	"github.com/bwmarrin/discordgo"
)

func PlayerGames(steamClient steam.Steam, steamID string) (*discordgo.MessageEmbed, error) {
	id, err := steamClient.ResolveID(steamID)
	if err != nil {
		// log error
		return nil, err
	}

	player, err := steamClient.PlayerSummaries(id)
	if err != nil {
		// log error
		return nil, err
	}

	ownedApps, _ := steamClient.AppsOwned(player[0].SteamID)
	recentApps, _ := steamClient.AppsRecentlyPlayed(player[0].SteamID)

	fields := []*discordgo.MessageEmbedField{}
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "Total Playtime",
		Value:  fmt.Sprintf("%dh", steam.AppsTotalHoursPlayed(*ownedApps)),
		Inline: true,
	})

	mostPlayed, _ := steam.AppsMostPlayed(*ownedApps)
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "Most Played Game",
		Value:  DefaultAppValue(mostPlayed),
		Inline: true,
	})

	leastPlayed, _ := steam.AppsLeastPlayed(*ownedApps)
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "Least Played Game",
		Value:  DefaultAppValue(leastPlayed),
		Inline: true,
	})

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "Games Owned",
		Value:  strconv.Itoa(len(*ownedApps)),
		Inline: true,
	})

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "Games Played",
		Value:  strconv.Itoa(len(steam.AppsPlayed(*ownedApps))),
		Inline: true,
	})

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "Games Not Played",
		Value:  strconv.Itoa(len(steam.AppsNotPlayed(*ownedApps))),
		Inline: true,
	})

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "Recent Playtime",
		Value:  fmt.Sprintf("%dh", steam.AppsRecentHoursPlayed(*ownedApps)),
		Inline: true,
	})

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "Recent Games Played",
		Value:  strconv.Itoa(len(*recentApps)),
		Inline: true,
	})

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "",
		Value:  "",
		Inline: true,
	})

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
		Fields: fields,
	}

	return embMsg, nil
}

func DefaultAppValue(value *steam.AppPlayTime) string {
	if value == nil {
		return "-"
	}
	return value.Name
}
