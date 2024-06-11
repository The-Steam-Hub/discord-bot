package main

import (
	"fmt"

	"github.com/KevinFagan/steam-stats/steam"
)

func main() {
	steam := steam.Steam{Key: "24CC43B3FDACCE9EE85C3BD16660D255"}
	id := "76561198101884982"

	profile, err := steam.Player(id)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Player Info:")
	fmt.Printf("Name: %s\n", profile.Name)
	fmt.Printf("Account Age: %s\n", profile.ProfileAge())
	fmt.Printf("Country: %s\n", profile.CountryCode)
	fmt.Printf("State: %s\n", profile.StateCode)
	fmt.Printf("Real Name: %s\n", profile.RealName)
	fmt.Printf("Badge Count: %d\n", len(profile.Badges))
	fmt.Printf("Player Level: %d\n", profile.PlayerLevel)
	fmt.Printf("Player Level Percentile: %.2f\n", profile.PlayerLevelPercentile)

	fmt.Println("\nBan Info:")
	fmt.Printf("Community Banned: %t\n", profile.CommunityBanned)
	fmt.Printf("VAC Banned: %t\n", profile.VACBanned)
	fmt.Printf("Number of VAC Bans: %d\n", profile.NumOfVacBans)
	fmt.Printf("Days Since Last Ban: %d\n", profile.DaysSinceLast)
	fmt.Printf("Number of Game Bans: %d\n", profile.NumOfGameBans)
	fmt.Printf("Economy Ban: %s\n", profile.EconomyBan)

	friends, err := steam.Friends(id)
	if err != nil {
		fmt.Println(err)
	}
	oldestFriend, err := steam.Player(friends.Oldest().ID)
	if err != nil {
		fmt.Println(err)
	}
	newestFriend, err := steam.Player(friends.Newest().ID)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("\nFriend Info:")
	fmt.Printf("Friend Count: %d\n", friends.Count())
	fmt.Printf("Oldest Friend: %s\n", oldestFriend.Name)
	fmt.Printf("Newest Friend: %s\n", newestFriend.Name)

	games, err := steam.Games(id)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("\nGame Info:")
	fmt.Printf("Total Games: %d\n", games.Count())
	fmt.Printf("Games Played: %d\n", games.GamesPlayed())
	fmt.Printf("Games Not Played: %d\n", games.GamesNotPlayed())
	fmt.Printf("Total Hours Played: %dh\n", games.TotalHoursPlayed())
	fmt.Printf("Most Played Game: \"%s\"\n", games.MostPlayed().Name)

	fmt.Println("\nOS Info:")
	fmt.Printf("Most Used OS: %s\n", games.MostUsedOS())
	fmt.Printf("Windows Playtime: %dh\n", games.WindowsPlaytime())
	fmt.Printf("Mac Playtime: %dh\n", games.MacPlaytime())
	fmt.Printf("Linux Playtime: %dh\n", games.LinuxPlaytime())
	fmt.Printf("Deck Playtime: %dh\n", games.DeckPlaytime())

}
