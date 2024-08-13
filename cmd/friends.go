package cmd

import (
	"fmt"
	"math"

	"github.com/KevinFagan/steam-stats/steam"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const (
	cap = 50
)

type FriendData struct {
	Friend steam.Friend
	Player steam.Player
}

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

	// Retrieving list of friends for a given player, there is no known limit to how many friends
	// will be returned within a single request
	friendsList, err := steamClient.GetFriendsList(player[0].SteamID)
	if err != nil {
		logs["error"] = err
		logrus.WithFields(logs).Error("unable to retrieve friend information")
	}

	// Sorting the friends list so we display the oldest friends first
	sortedFriendsList := steam.SortFriends(friendsList)
	// Capping the friends list to avoid message overflow issues with Discord
	sortedCappedFriendsList := sortedFriendsList[:int(math.Min(float64(len(sortedFriendsList)), cap))]
	// Length may be zero if the players account is private
	if len(sortedCappedFriendsList) > 0 {
		// Assigning the newest friend to the last index. This allows us to grab the name of the newest friend in the same API call as the other 49 friends
		sortedCappedFriendsList[len(sortedCappedFriendsList)-1] = sortedFriendsList[len(sortedFriendsList)-1]
	}

	// Getting player information for all friends within the cap range
	players, err := steamClient.GetPlayerSummaries(steam.GetFriendIDs(sortedFriendsList)[:len(sortedCappedFriendsList)]...)
	if err != nil {
		logs["error"] = err
		logrus.WithFields(logs).Error("unable to retrieve player information")
	}

	// Friend data and Player data exists in two seperate API calls, and so, we need to tie the data together
	// The data is already sorted and is persisted in the friendData slice
	friendData := make([]FriendData, len(players))
	for _, v := range players {
		for k, j := range sortedCappedFriendsList {
			if v.SteamID == j.ID {
				friendData[k] = FriendData{
					Friend: j,
					Player: v,
				}
			}
		}
	}

	names, statuses, friendsSince, oldest, newest := "", "", "", "", ""

	// Length may be zero if the players account is private
	if len(sortedCappedFriendsList) > 0 {
		oldest = sortedCappedFriendsList[0].ID
		newest = sortedCappedFriendsList[len(sortedCappedFriendsList)-1].ID
	}

	for k, i := range friendData {
		// Avoid adding the last entry (newest friend) if we are at the cap
		if k < cap-1 {
			names += fmt.Sprintf("%s\n", i.Player.Name)
			statuses += fmt.Sprintf("%s\n", i.Player.Status())
			friendsSince += fmt.Sprintf("%s\n", steam.UnixToDate(i.Friend.FriendsSince))
		}
		// Finding the name that belongs to the newest ID
		if i.Player.SteamID == newest {
			newest = i.Player.Name
		}
		// Finding the name that belongs to the oldest ID
		if i.Player.SteamID == oldest {
			oldest = i.Player.Name
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
				Value:  fmt.Sprintf("%d", len(sortedFriendsList)),
				Inline: true,
			},
			{
				Name:   "Friends (Top 50)",
				Value:  DefaultStringValue(names),
				Inline: true,
			},
			{
				Name:   "Friends For",
				Value:  DefaultStringValue(friendsSince),
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
