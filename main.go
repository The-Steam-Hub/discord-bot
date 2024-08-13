package main

import (
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/KevinFagan/steam-stats/cmd"
	"github.com/KevinFagan/steam-stats/steam"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

const (
	command = 1
	ID      = 2
)

var (
	discordToken = os.Getenv("DISCORD_BOT_TOKEN_DEV")
	steamKey     = os.Getenv("STEAM_API_KEY")
	steamClient  = steam.Steam{Key: steamKey}
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})

	logrus.Info("creating Discord session...")
	dg, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		logrus.Fatalf("error creating Discord session: %s\n", err)
		return
	}

	dg.AddHandler(messageCreate)
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	logrus.Info("opening websocket connection to Discord...")
	err = dg.Open()
	if err != nil {
		logrus.Fatalf("error opening Discord websocket connection: %s\n", err)
		return
	}

	logrus.Info("Steam Stats is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}
	if !strings.HasPrefix(m.Content, "!dev-stats") {
		return
	}

	args := strings.Split(strings.TrimSpace(m.Content), " ")

	if len(args) == 1 {
		cmd.Help(s, m)
		return
	}
	if args[command] == "profile" && len(args) == 3 {
		resolvedID := steamClient.ResolveIDFromURL(args[ID])
		cmd.Profile(s, m, steamClient, resolvedID)
		return
	}
	if args[command] == "id" && len(args) == 3 {
		resolvedID := steamClient.ResolveIDFromURL(args[ID])
		cmd.ID(s, m, steamClient, resolvedID)
		return
	}
	if args[command] == "friends" && len(args) == 3 {
		resolvedID := steamClient.ResolveIDFromURL(args[ID])
		cmd.Friends(s, m, steamClient, resolvedID)
		return
	}
	if args[command] == "games" && len(args) == 3 {
		resolvedID := steamClient.ResolveIDFromURL(args[ID])
		cmd.Games(s, m, steamClient, resolvedID)
		return
	}
	if args[command] == "bans" && len(args) == 3 {
		resolvedID := steamClient.ResolveIDFromURL(args[ID])
		cmd.Bans(s, m, steamClient, resolvedID)
		return
	}
	cmd.Help(s, m)
}
