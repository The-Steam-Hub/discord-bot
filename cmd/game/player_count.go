package game

import (
	"strconv"

	"github.com/KevinFagan/steam-stats/cmd"
	"github.com/KevinFagan/steam-stats/steam"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func AppPlayerCount(session *discordgo.Session, interaction *discordgo.InteractionCreate, steamClient steam.Steam, input string) {
	logs := logrus.Fields{
		"input":  input,
		"author": interaction.Member.User.Username,
		"uuid":   uuid.New(),
	}

	appID, err := steamClient.AppIDResolve(input)
	if err != nil {
		logs["error"] = err
		errMsg := "unable to resolve game ID"
		logrus.WithFields(logs).Error(errMsg)
		cmd.HandleErrorMessage(session, interaction, &logs, errMsg)
		return
	}

	appPlayerCount, err := steamClient.AppPlayerCount(appID)
	if err != nil {
		logs["error"] = err
		errMsg := "unable to retrieve player count"
		logrus.WithFields(logs).Error(errMsg)
		cmd.HandleErrorMessage(session, interaction, &logs, errMsg)
		return
	}

	appData, err := steamClient.AppData(appID)
	if err != nil {
		logs["error"] = err
		errMsg := "unable to retrieve game data"
		logrus.WithFields(logs).Error(errMsg)
		cmd.HandleErrorMessage(session, interaction, &logs, errMsg)
		return
	}

	p := message.NewPrinter(language.English)

	embMsg := &discordgo.MessageEmbed{
		Title: appData.Name,
		URL:   steam.SteamPoweredAPI + "app/" + strconv.Itoa(appID),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: appData.HeaderImage,
		},
		Color: 0x66c0f4,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Players",
				Value:  p.Sprintf("%d", appPlayerCount.Current),
				Inline: true,
			},
			{
				Name:   "24h Peak",
				Value:  p.Sprintf("%d", appPlayerCount.Peak24Hour),
				Inline: true,
			},
			{
				Name:   "All-Time Peak",
				Value:  p.Sprintf("%d", appPlayerCount.PeakAllTime),
				Inline: true,
			},
		},
	}
	cmd.HandleOkMessage(embMsg, session, interaction, &logs)
}
