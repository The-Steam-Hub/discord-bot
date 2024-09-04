package cmd

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

func HandleErrorMessage(session *discordgo.Session, interaction *discordgo.InteractionCreate, logs *logrus.Fields, errMsg string) {
	err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: errMsg,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	if err != nil {
		(*logs)["error"] = err
		logrus.WithFields(*logs).Error("unable to send message")

	}
}

func HandleOkMessage(embMsg *discordgo.MessageEmbed, session *discordgo.Session, interaction *discordgo.InteractionCreate, logs *logrus.Fields) {
	err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				embMsg,
			},
		},
	})

	if err != nil {
		(*logs)["error"] = err
		logrus.WithFields(*logs).Error("unable to send message")
	}
}

func HandleDefaultString(value string) string {
	if value == "" {
		return "-"
	}
	return value
}
