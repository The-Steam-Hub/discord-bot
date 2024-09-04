package game

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/KevinFagan/steam-stats/cmd"
	"github.com/KevinFagan/steam-stats/steam"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func AppSearch(session *discordgo.Session, interaction *discordgo.InteractionCreate, steamClient steam.Steam, input string) {
	logs := logrus.Fields{
		"input":  input,
		"author": interaction.Member.User.Username,
		"uuid":   uuid.New(),
	}

	appID, err := steamClient.AppSearch(input)
	if err != nil {
		logs["error"] = err
		errMsg := "game not found"
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

	embMsg := &discordgo.MessageEmbed{
		Title:       appData.Name,
		URL:         steam.SteamPoweredAPI + "app/" + strconv.Itoa(appID),
		Description: appData.ShortDescription,
		Image: &discordgo.MessageEmbedImage{
			URL: appData.HeaderImage,
		},
		Color: 0x66c0f4,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Price",
				Value:  cmd.HandleDefaultString(formatPrice(*appData)),
				Inline: true,
			},
			{
				Name:   "Release Date",
				Value:  cmd.HandleDefaultString(appData.ReleaseDate.Date),
				Inline: true,
			},
			{
				Name:   "# DLC",
				Value:  strconv.Itoa(len(appData.DLC)),
				Inline: true,
			},
			{
				Name:   "Developers",
				Value:  cmd.HandleDefaultString(strings.Join(appData.Developers, ", ")),
				Inline: true,
			},
			{
				Name:   "Publishers",
				Value:  cmd.HandleDefaultString(strings.Join(appData.Publishers, ", ")),
				Inline: true,
			},
			{
				Name:   "Genres",
				Value:  cmd.HandleDefaultString(formatGenres(*appData)),
				Inline: true,
			},
		},
	}
	cmd.HandleOkMessage(embMsg, session, interaction, &logs)
}

func formatPrice(appData steam.AppData) string {
	if appData.IsFree {
		return "Free"
	}

	iFormat := appData.PriceOverview.InitialFormatted
	fFormat := appData.PriceOverview.FinalFormatted

	if iFormat != "" && iFormat != fFormat {
		return fmt.Sprintf("~~%s~~ %s", iFormat, fFormat)
	}

	return appData.PriceOverview.FinalFormatted
}

func formatGenres(appData steam.AppData) string {
	var format string
	for _, v := range appData.Genres {
		format += v.Description + ", "
	}

	return strings.TrimSuffix(format, ", ")
}
