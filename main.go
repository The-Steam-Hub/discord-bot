package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/KevinFagan/steam-stats/steam"
	"github.com/bwmarrin/discordgo"
)

var (
	discordToken = os.Getenv("DISCORD_BOT_TOKEN")
	steamKey     = os.Getenv("STEAM_API_KEY")
	steamClient  = steam.Steam{Key: steamKey}
)

func main() {
	dg, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(messageCreate)
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}

// This function will be called every time a new message is created on
// any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	args := strings.Split(m.Content, " ")
	if len(args) != 3 {
		return // invalid command
	}

	if args[0] == "!stats" {
		fmt.Printf("Received command: %s\n", strings.Join(args, " "))

		// Getting player information
		steamID := steamClient.ResolveID(args[2])
		player, err := steamClient.PlayerWithDetails(steamID)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Unable to retrieve user information")
			fmt.Println(err)
			return
		}

		// Sending message depending on the command
		if args[1] == "profile" {
			message := messageProfile(player)
			s.ChannelMessageSendEmbed(m.ChannelID, &message)
			return
		}
		if args[1] == "friends" {
			message := messageFriends(steamClient, player)
			s.ChannelMessageSendEmbed(m.ChannelID, &message)
			return
		}
		if args[1] == "games" {
			message := messageGames(steamClient, player)
			s.ChannelMessageSendEmbed(m.ChannelID, &message)
			return
		}
		if args[1] == "bans" {
			message := messageBans(player)
			s.ChannelMessageSendEmbed(m.ChannelID, &message)
			return
		}
	}
}

func messageFriends(steam steam.Steam, player steam.Player) discordgo.MessageEmbed {
	newestFriend := "-"
	oldestFriend := "-"

	friends, err := steam.Friends(player.SteamID)
	if err != nil {
		fmt.Println(err)
	} else {
		newest, err := steam.Player(friends.Newest().ID)
		if err != nil {
			fmt.Println(err)
		}
		newestFriend = newest.Name
		oldest, err := steam.Player(friends.Oldest().ID)
		if err != nil {
			fmt.Println(err)
		}
		oldestFriend = oldest.Name
	}

	embed := &discordgo.MessageEmbed{
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Friend information is dependent upon the users privacy settings.",
		},
		Color: 0x66c0f4,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: player.AvatarFull,
		},
		Author: &discordgo.MessageEmbedAuthor{
			Name: fmt.Sprintf("%s %s", player.Status(), player.Name),
			URL:  player.ProfileURL,
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "# of Friends",
				Value:  strconv.Itoa(friends.Count()),
				Inline: true,
			},
			{
				Name:   "Newest Friend",
				Value:  newestFriend,
				Inline: true,
			},
			{
				Name:   "Oldest Friend",
				Value:  oldestFriend,
				Inline: true,
			},
		},
	}
	return *embed
}

func messageGames(steam steam.Steam, player steam.Player) discordgo.MessageEmbed {
	allGames, err := steam.Games(player.SteamID)
	if err != nil {
		fmt.Println(err)
	}
	recentGames, err := steam.RecentGames(player.SteamID)
	if err != nil {
		fmt.Println(err)
	}
	mostPlayed := allGames.MostPlayed().Name
	if mostPlayed == "" {
		mostPlayed = "-"
	}

	embed := &discordgo.MessageEmbed{
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Game information is dependent upon the users privacy settings.",
		},
		Color: 0x66c0f4,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: player.AvatarFull,
		},
		Author: &discordgo.MessageEmbedAuthor{
			Name: fmt.Sprintf("%s %s", player.Status(), player.Name),
			URL:  player.ProfileURL,
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Most Played Game",
				Value:  mostPlayed,
				Inline: true,
			},
			{
				Name:   "Total Playtime",
				Value:  fmt.Sprintf("%dh", allGames.TotalHoursPlayed()),
				Inline: true,
			},
			// Empty field to make the embed look better
			{
				Name:   "",
				Value:  "",
				Inline: true,
			},
			{
				Name:   "Games Owned",
				Value:  strconv.Itoa(len(allGames.Games)),
				Inline: true,
			},
			{
				Name:   "Games Played",
				Value:  strconv.Itoa(allGames.GamesPlayed()),
				Inline: true,
			},
			{
				Name:   "Games Not Played",
				Value:  strconv.Itoa(allGames.GamesNotPlayed()),
				Inline: true,
			},
			{
				Name:   "Last 2 Week Playtime",
				Value:  fmt.Sprintf("%dh", recentGames.HoursPlayed2Weeks()),
				Inline: true,
			},
			{
				Name:   "Last 2 Week Games Played",
				Value:  strconv.Itoa(len(recentGames.Games)),
				Inline: true,
			},
			// Empty field to make the embed look better
			{
				Name:   "",
				Value:  "",
				Inline: true,
			},
		},
	}
	return *embed
}

func messageBans(player steam.Player) discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Ban information is dependent upon the users privacy settings.",
		},
		Color: 0x66c0f4,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: player.AvatarFull,
		},
		Author: &discordgo.MessageEmbedAuthor{
			Name: fmt.Sprintf("%s %s", player.Status(), player.Name),
			URL:  player.ProfileURL,
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "VAC Banned",
				Value:  strconv.FormatBool(player.VACBanned),
				Inline: true,
			},
			{
				Name:   "# Of VAC Bans",
				Value:  strconv.Itoa(player.NumOfVacBans),
				Inline: true,
			},
			{
				Name:   "# Of Game Bans",
				Value:  strconv.Itoa(player.NumOfGameBans),
				Inline: true,
			},
			{
				Name:   "Days Since Last Ban",
				Value:  fmt.Sprintf("%dd", player.DaysSinceLastBan),
				Inline: true,
			},
			{
				Name:   "Community Banned",
				Value:  strconv.FormatBool(player.CommunityBanned),
				Inline: true,
			},
			{
				Name:   "Economy Banned",
				Value:  player.EconomyBan,
				Inline: true,
			},
		},
	}
	return *embed
}

func messageProfile(player steam.Player) discordgo.MessageEmbed {
	realName := player.RealName
	if realName == "" {
		realName = "-"
	}
	countryCode := player.CountryCode
	if countryCode == "" {
		countryCode = "-"
	}
	stateCode := player.StateCode
	if stateCode == "" {
		stateCode = "-"
	}

	embed := &discordgo.MessageEmbed{
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Profile information is dependent upon the users privacy settings.",
		},
		Color: 0x66c0f4,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: player.AvatarFull,
		},
		Author: &discordgo.MessageEmbedAuthor{
			Name: fmt.Sprintf("%s %s", player.Status(), player.Name),
			URL:  player.ProfileURL,
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Real Name",
				Value:  realName,
				Inline: true,
			},
			{
				Name:   "Country Code",
				Value:  countryCode,
				Inline: true,
			},
			{
				Name:   "State Code",
				Value:  stateCode,
				Inline: true,
			},
			{
				Name:   "Profile Age",
				Value:  player.ProfileAge(),
				Inline: true,
			},
			{
				Name:   "Last Seen",
				Value:  player.LastSeen(),
				Inline: true,
			},
			// Empty field to make the embed look better
			{
				Name:   "",
				Value:  "",
				Inline: true,
			},
			{
				Name:   "Level Percentile",
				Value:  strconv.FormatFloat(player.PlayerLevelPercentile, 'f', 2, 64),
				Inline: true,
			},
			{
				Name:   "Level",
				Value:  strconv.Itoa(player.PlayerLevel),
				Inline: true,
			},
			{
				Name:   "Badges",
				Value:  strconv.Itoa(int(len(player.Badges))),
				Inline: true,
			},
		},
	}
	return *embed
}
