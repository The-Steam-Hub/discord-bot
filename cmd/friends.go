package cmd

import (
	"github.com/KevinFagan/steam-stats/steam"
	"github.com/bwmarrin/discordgo"
)

func Friends(s *discordgo.Session, m *discordgo.MessageCreate, steam steam.Steam) {
	return
}

// func messageFriends(m *discordgo.MessageCreate, s steam.Steam, p steam.Player) discordgo.MessageEmbed {
// 	// Retrieving friend IDs for a player
// 	f, err := s.GetFriendsList(p.SteamID)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	// Retrieving friend information for each friend
// 	fIDs := make([]string, len(f.Friends))
// 	for i, friend := range f.Friends {
// 		fIDs[i] = friend.ID
// 	}
// 	fInfo, err := s.GetPlayerSummaries(fIDs...)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	// Building strings for the embed fields
// 	fNameString, fStatusString, fForString := "", "", ""
// 	fCount := f.Count()
// 	fNewest := f.Newest().ID
// 	fOldest := f.Oldest().ID

// 	for _, fi := range fInfo {
// 		fNameString += fmt.Sprintf("%s\n", fi.Name)
// 		fStatusString += fmt.Sprintf("%s\n", fi.Status())
// 		for _, f := range f.Friends {
// 			if fi.SteamID == f.ID {
// 				fForString += fmt.Sprintf("%s\n", steam.UnixToDate(f.FriendsSince))
// 			}
// 			if fi.SteamID == fNewest {
// 				fNewest = fi.Name
// 			}
// 			if fi.SteamID == fOldest {
// 				fOldest = fi.Name
// 			}
// 		}
// 	}
