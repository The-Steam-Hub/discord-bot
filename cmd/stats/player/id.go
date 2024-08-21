package player

import (
	"fmt"
	"strconv"

	"github.com/KevinFagan/steam-stats/steam"
	"github.com/bwmarrin/discordgo"
)

func IDEmbeddedMessage(steamClient steam.Steam, steamID string) (*discordgo.MessageEmbed, error) {
	id, err := steamClient.ResolveID(steamID)
	if err != nil {
		return nil, err
	}

	player, err := steamClient.GetPlayerSummaries(id)
	if err != nil {
		return nil, err
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
				Value:  DefaultStringValue(steam.SteamID64ToSteamID(steamIDInt)),
				Inline: true,
			},
			{
				Name:   "Steam ID3",
				Value:  DefaultStringValue(steam.SteamID64ToSteamID3(steamIDInt)),
				Inline: true,
			},
			{
				Name:   "Steam ID64",
				Value:  DefaultStringValue(player[0].SteamID),
				Inline: true,
			},
		},
	}
	return embMsg, nil
}
