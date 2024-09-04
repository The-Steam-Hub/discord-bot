package player

import (
	"fmt"
	"strconv"

	"github.com/KevinFagan/steam-stats/cmd"
	"github.com/KevinFagan/steam-stats/steam"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func PlayerGames(session *discordgo.Session, interaction *discordgo.InteractionCreate, steamClient steam.Steam, input string) {
	logs := logrus.Fields{
		"input":  input,
		"author": interaction.Member.User.Username,
		"uuid":   uuid.New(),
	}

	id, err := steamClient.ResolveID(input)
	if err != nil {
		logs["error"] = err
		errMsg := "unable to resolve player ID"
		logrus.WithFields(logs).Error(errMsg)
		cmd.HandleErrorMessage(session, interaction, &logs, errMsg)
		return
	}

	player, err := steamClient.PlayerSummaries(id)
	if err != nil {
		logs["error"] = err
		errMsg := "unable to retrieve player summary"
		logrus.WithFields(logs).Error(errMsg)
		cmd.HandleErrorMessage(session, interaction, &logs, errMsg)
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

	cmd.HandleOkMessage(embMsg, session, interaction, &logs)
}

func DefaultAppValue(value *steam.AppPlayTime) string {
	if value == nil {
		return "-"
	}
	return value.Name
}
