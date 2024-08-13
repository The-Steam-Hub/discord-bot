package main

import (
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/KevinFagan/steam-stats/cmd"
	"github.com/KevinFagan/steam-stats/steam"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

const (
	command = 1
	ID      = 2
)

var (
	prefix       = ""
	steamKey     = ""
	discordToken = ""
	steamClient  = steam.Steam{}
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})

	logrus.Info("loading configurations...")
	loadEnvironment()

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

	if !strings.HasPrefix(m.Content, prefix) {
		return
	}

	args := strings.Split(strings.TrimSpace(m.Content), " ")

	if len(args) == 1 {
		cmd.Help(s, m)
		return
	}
	if args[command] == "profile" && len(args) == 3 {
		id, _ := steamClient.ParseSteamID(args[ID])
		cmd.Profile(s, m, steamClient, id)
		return
	}
	if args[command] == "id" && len(args) == 3 {
		id, _ := steamClient.ParseSteamID(args[ID])
		cmd.ID(s, m, steamClient, id)
		return
	}
	if args[command] == "friends" && len(args) == 3 {
		id, _ := steamClient.ParseSteamID(args[ID])
		cmd.Friends(s, m, steamClient, id)
		return
	}
	if args[command] == "games" && len(args) == 3 {
		id, _ := steamClient.ParseSteamID(args[ID])
		cmd.Games(s, m, steamClient, id)
		return
	}
	if args[command] == "bans" && len(args) == 3 {
		id, _ := steamClient.ParseSteamID(args[ID])
		cmd.Bans(s, m, steamClient, id)
		return
	}
	cmd.Help(s, m)
}

func loadEnvironment() {
	env := os.Getenv("BOT_ENV")
	if env == "" {
		env = "development"
	}

	err := godotenv.Load(".env." + env)
	if err != nil {
		logrus.Fatalf("error loading .env file: %s", err)
	}

	prefix = os.Getenv("PREFIX")
	steamKey = os.Getenv("STEAM_API_KEY")
	discordToken = os.Getenv("DISCORD_BOT_TOKEN")
	steamClient = steam.Steam{Key: steamKey}
}
