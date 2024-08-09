package cmd

import (
	"fmt"

	"github.com/KevinFagan/steam-stats/steam"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func Friends(s *discordgo.Session, m *discordgo.MessageCreate, steamClient steam.Steam, steamID string) {
	logs := logrus.Fields{
		"command": "friends",
		"player":  steamID,
		"author":  m.Author.Username,
		"uuid":    uuid.New().String(),
	}

	logrus.WithFields(logs).Info("command recieved")

	player, err := steamClient.GetPlayerSummaries(steamID)
	if err != nil {
		logs["error"] = err
		logrus.WithFields(logs).Error("unable to retrieve player information")
		s.ChannelMessageSend(m.ChannelID, "unable to retrieve player information")
		return
	}

	friendsList, err := steamClient.GetFriendsList(player[0].SteamID)
	if err != nil {
		logs["error"] = err
		logrus.WithFields(logs).Error("unable to retrieve friend information")
	}

	friendIDs := make([]string, len(friendsList.Friends))
	for i, friend := range friendsList.Friends {
		friendIDs[i] = friend.ID
	}

	friendInfo, err := steamClient.GetPlayerSummaries(friendIDs...)
	if err != nil {
		logs["error"] = err
		logrus.WithFields(logs).Error("unable to retrieve player information")
	}

	// TODO: Before creating the list below, we need to sort out the top 50 frinds
	names := ""
	statuses := ""
	friendsFor := ""
	newest := friendsList.Newest().ID
	oldest := friendsList.Oldest().ID

	for _, fi := range friendInfo {
		names += fmt.Sprintf("%s\n", fi.Name)
		statuses += fmt.Sprintf("%s\n", fi.Status())
		for _, f := range friendsList.Friends {
			if fi.SteamID == f.ID {
				friendsFor += fmt.Sprintf("%s\n", steam.UnixToDate(f.FriendsSince))
			}
			if fi.SteamID == newest {
				newest = fi.Name
			}
			if fi.SteamID == oldest {
				oldest = fi.Name
			}
		}
	}

	embedInfo := EmbededInfo(player[0])
	embedMessage := &discordgo.MessageEmbed{
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Friend information is dependent upon the user's privacy settings.",
		},
		Color:     embedInfo.Color,
		Thumbnail: embedInfo.Thumbnail,
		Author:    embedInfo.Author,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Newest",
				Value:  DefaultStringValue(newest),
				Inline: true,
			},
			{
				Name:   "Oldest",
				Value:  DefaultStringValue(oldest),
				Inline: true,
			},
			{
				Name:   "Count",
				Value:  fmt.Sprintf("%d", len(friendIDs)),
				Inline: true,
			},
			{
				Name:   "Friends",
				Value:  DefaultStringValue(names),
				Inline: true,
			},
			{
				Name:   "Friends For",
				Value:  DefaultStringValue(friendsFor),
				Inline: true,
			},
			{
				Name:   "Status",
				Value:  DefaultStringValue(statuses),
				Inline: true,
			},
		},
	}
	s.ChannelMessageSendEmbed(m.ChannelID, embedMessage)
}
