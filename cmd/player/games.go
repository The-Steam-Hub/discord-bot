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

func PlayerGames(session *discordgo.Session, interaction *discordgo.InteractionCreate, steamClient steam.Steam, input string) {
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

	ownedApps, err := steamClient.AppsOwned(player[0].SteamID)
	if err != nil {
		logs["error"] = err
		logrus.WithFields(logs).Error("unable to retrieve owned games")
	}

	recentApps, err := steamClient.AppsRecentlyPlayed(player[0].SteamID)
	if err != nil {
		logs["error"] = err
		logrus.WithFields(logs).Error("unable to retireve recently played games")
	}

	mostPlayed, _ := steam.AppsMostPlayed(*ownedApps)
	leastPlayed, _ := steam.AppsLeastPlayed(*ownedApps)

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
				Value:  fmt.Sprintf("%dh", steam.AppsTotalHoursPlayed(*ownedApps)),
				Inline: true,
			},
			{
				Name:   "Most Played Game",
				Value:  DefaultAppValue(mostPlayed),
				Inline: true,
			},
			{
				Name:   "Least Played Game",
				Value:  DefaultAppValue(leastPlayed),
				Inline: true,
			},
			{
				Name:   "Games Owned",
				Value:  strconv.Itoa(len(*ownedApps)),
				Inline: true,
			},
			{
				Name:   "Games Played",
				Value:  strconv.Itoa(len(steam.AppsPlayed(*ownedApps))),
				Inline: true,
			},
			{
				Name:   "Games Not Played",
				Value:  strconv.Itoa(len(steam.AppsNotPlayed(*ownedApps))),
				Inline: true,
			},
			{
				Name:   "Recent Playtime",
				Value:  fmt.Sprintf("%dh", steam.AppsRecentHoursPlayed(*ownedApps)),
				Inline: true,
			},
			{
				Name:   "Recent Games Played",
				Value:  strconv.Itoa(len(*recentApps)),
				Inline: true,
			},
			{
				Name:   "",
				Value:  "",
				Inline: true,
			},
		},
	}

	cmd.HandleMessageOk(embMsg, session, interaction, &logs)
}

func DefaultAppValue(value *steam.AppPlayTime) string {
	if value == nil {
		return "-"
	}
	return value.Name
}
