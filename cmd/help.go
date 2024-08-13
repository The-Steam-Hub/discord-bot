package cmd

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func Help(s *discordgo.Session, m *discordgo.MessageCreate) {
	logs := logrus.Fields{
		"command": "help",
		"author":  m.Author.Username,
		"uuid":    uuid.New().String(),
	}

	logrus.WithFields(logs).Info("command recieved")

	var sb strings.Builder
	w := tabwriter.NewWriter(&sb, 0, 0, 4, ' ', 0)
	fmt.Fprint(w, "```\n")
	fmt.Fprint(w, "Usage:\n")
	fmt.Fprint(w, "\t!stats [command] [Steam ID, Steam ID3, Steam ID64, Steam URL]\n\n")
	fmt.Fprint(w, "Commands:\n")
	fmt.Fprint(w, "\thelp\tPrints information about the Steam Stats bot\n")
	fmt.Fprint(w, "\tprofile\tPrints information about a user profile\n")
	fmt.Fprint(w, "\tfriends\tPrints information about a users friends list\n")
	fmt.Fprint(w, "\tgames\tPrints information about a users game library\n")
	fmt.Fprint(w, "\tbans\tPrints ban informaton about a user\n")
	fmt.Fprint(w, "\tid\tPrints various Steam IDs for a user\n\n")
	fmt.Fprint(w, "Examples:\n")
	fmt.Fprint(w, "\t!stats profile [U:1:110439373]\n")
	fmt.Fprint(w, "\t!stats profile 76561198070705101\n")
	fmt.Fprint(w, "\t!stats profile STEAM_1:1:55219686\n")
	fmt.Fprint(w, "\t!stats profile https://steamcommunity.com/id/TheLordSquirrel/\n\n")
	fmt.Fprint(w, "```\n")
	w.Flush()

	s.ChannelMessageSend(m.ChannelID, sb.String())

}
