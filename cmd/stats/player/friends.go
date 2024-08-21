package player

import (
	"fmt"
	"math"

	"github.com/KevinFagan/steam-stats/steam"
	"github.com/bwmarrin/discordgo"
)

const (
	cap = 50
)

type FriendData struct {
	Friend steam.Friend
	Player steam.Player
}

func PlayerFriends(steamClient steam.Steam, steamID string) (*discordgo.MessageEmbed, error) {
	id, err := steamClient.ResolveID(steamID)
	if err != nil {
		// log error
		return nil, err
	}

	player, err := steamClient.PlayerSummaries(id)
	if err != nil {
		// log error
		return nil, err
	}

	friendsList, _ := steamClient.FriendsList(player[0].SteamID)
	// Sorting the friends list so we display the oldest friends first
	sortedFriendsList := steam.FriendsSort(*friendsList)
	// Capping the friends list to avoid message overflow issues with Discord
	sortedCappedFriendsList := sortedFriendsList[:int(math.Min(float64(len(sortedFriendsList)), cap))]
	// Length may be zero if the players account is private
	if len(sortedCappedFriendsList) > 0 {
		// Assigning the newest friend to the last index. This allows us to grab the name of the newest friend in the same API call as the other 49 friends
		sortedCappedFriendsList[len(sortedCappedFriendsList)-1] = sortedFriendsList[len(sortedFriendsList)-1]
	}

	// Getting player information for all friends within the cap range
	players, _ := steamClient.PlayerSummaries(steam.FriendIDs(sortedFriendsList)[:len(sortedCappedFriendsList)]...)

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

	embMsg := &discordgo.MessageEmbed{
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Friend information is dependent upon the user's privacy settings.",
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
				Name:   "Top 50 Friends",
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

	return embMsg, nil
}
