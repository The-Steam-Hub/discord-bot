package main

import (
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"

	"github.com/KevinFagan/steam-stats/steam"
	"github.com/bwmarrin/discordgo"
)

var (
	discordToken = os.Getenv("DISCORD_BOT_TOKEN")
	steamKey     = os.Getenv("STEAM_API_KEY")
)

func main() {
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	args := strings.Split(m.Content, " ")
	if len(args) != 2 {
		return // invalid command
	}

	if args[0] == "!stats" {
		steam := steam.Steam{Key: steamKey}

		vanityRegex := regexp.MustCompile(`https:\/\/steamcommunity\.com\/id\/([^\/]+)`)
		vanityMatch := vanityRegex.FindStringSubmatch(args[1])
		IDRegex := regexp.MustCompile(`https:\/\/steamcommunity\.com\/profiles\/(\d+)`)
		IDMatch := IDRegex.FindStringSubmatch(args[1])
		finalID := ""

		if len(vanityMatch) > 1 {
			vanity, err := steam.ResolveVanityURL(vanityMatch[1])
			if err != nil {
				return // Unable to resolve vanity URL
			}
			finalID = vanity.SteamID
		}

		if len(IDMatch) > 1 {
			finalID = IDMatch[1]
		}

		if finalID == "" {
			return // Invalid Steam URL
		}

		s.ChannelMessageSend(m.ChannelID, finalID)
	}
}
