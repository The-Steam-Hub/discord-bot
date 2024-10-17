package game

import (
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/the-steam-bot/discord-bot/cmd"
	"github.com/the-steam-bot/discord-bot/steam"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func AppPlayerCount(session *discordgo.Session, interaction *discordgo.InteractionCreate, steamClient steam.Steam, input string) {
	logs := logrus.Fields{
		"input":  input,
		"author": interaction.Member.User.Username,
		"uuid":   uuid.New(),
	}

	appID, err := steamClient.AppSearch(input)
	if err != nil {
		logs["error"] = err
		errMsg := "unable to find game"
		logrus.WithFields(logs).Error(errMsg)
		cmd.HandleMessageError(session, interaction, &logs, errMsg)
		return
	}

	appPlayerCount, err := steamClient.AppPlayerCount(appID)
	if err != nil {
		logs["error"] = err
		errMsg := "unable to retrieve player count"
		logrus.WithFields(logs).Error(errMsg)
		cmd.HandleMessageError(session, interaction, &logs, errMsg)
		return
	}

	appData, err := steamClient.AppDetailedData(appID)
	if err != nil {
		logs["error"] = err
		errMsg := "unable to retrieve game data"
		logrus.WithFields(logs).Error(errMsg)
		cmd.HandleMessageError(session, interaction, &logs, errMsg)
		return
	}

	p := message.NewPrinter(language.English)

	embMsg := &discordgo.MessageEmbed{
		Title: trimTitle(appData.Name),
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
	cmd.HandleMessageOk(embMsg, session, interaction, &logs)
}
