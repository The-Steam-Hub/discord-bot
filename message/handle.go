package message

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

func HandleEmbeddedMessage(embMsg *discordgo.MessageEmbed, session *discordgo.Session, interaction *discordgo.InteractionCreate, logs logrus.Fields, err error) {
	if err == nil {
		err = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					embMsg,
				},
			},
		})

		if err != nil {
			logs["error"] = err
			logrus.WithFields(logs).Error("unable to send message")
		}
	}

	if err != nil {
		logs["error"] = err
		logrus.WithFields(logs).Error("unable to complete request")
		err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Unable to retrieve player information",
			},
		})

		if err != nil {
			logs["error"] = err
			logrus.WithFields(logs).Error("unable to send message")
		}
	}
}
