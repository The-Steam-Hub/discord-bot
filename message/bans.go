package message

import (
	"fmt"
	"strconv"

	"github.com/KevinFagan/steam-stats/steam"
	"github.com/bwmarrin/discordgo"
)

func BanEmbeddedMessage(steamClient steam.Steam, steamID string) (*discordgo.MessageEmbed, error) {
	id, err := steamClient.ResolveID(steamID)
	if err != nil {
		return nil, err
	}

	player, err := steamClient.GetPlayerSummariesWithExtra(id)
	if err != nil {
		return nil, err
	}

	embMsg := &discordgo.MessageEmbed{
		Color: 0x66c0f4,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: player.AvatarFull,
		},
		Author: &discordgo.MessageEmbedAuthor{
			Name: fmt.Sprintf("%s %s", player.Status(), player.Name),
			URL:  player.ProfileURL,
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "VAC Banned",
				Value:  strconv.FormatBool(player.VACBanned),
				Inline: true,
			},
			{
				Name:   "# Of VAC Bans",
				Value:  strconv.Itoa(player.NumOfVacBans),
				Inline: true,
			},
			{
				Name:   "# Of Game Bans",
				Value:  strconv.Itoa(player.NumOfGameBans),
				Inline: true,
			},
			{
				Name:   "Days Since Last Ban",
				Value:  fmt.Sprintf("%dd", player.DaysSinceLastBan),
				Inline: true,
			},
			{
				Name:   "Community Banned",
				Value:  strconv.FormatBool(player.CommunityBanned),
				Inline: true,
			},
			{
				Name:   "Economy Banned",
				Value:  player.EconomyBan,
				Inline: true,
			},
		},
	}

	return embMsg, nil
}
