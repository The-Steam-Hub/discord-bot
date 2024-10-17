package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/the-steam-bot/discord-bot/cmd/game"
	"github.com/the-steam-bot/discord-bot/cmd/player"
	"github.com/the-steam-bot/discord-bot/steam"
)

var discordSession *discordgo.Session

var (
	steamToken   string
	discordToken string
	steamClient  steam.Steam
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "player",
			Description: "Fetches player statistics",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "profile",
					Description: "Fetches statistics about a players profile",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "value",
							Description: "Steam Identifier",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
					},
				},
				{
					Name:        "games",
					Description: "Fetches statistics about a players game library",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "value",
							Description: "Steam Identifier",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
					},
				},
				{
					Name:        "bans",
					Description: "Fetches statistics about a players bans histroy",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "value",
							Description: "Steam Identifier",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
					},
				},
				{
					Name:        "friends",
					Description: "Fetches statistics about a players friends list",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "value",
							Description: "Steam Identifier",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
					},
				},
				{
					Name:        "id",
					Description: "Fetches multiple formats of the players Steam ID",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "value",
							Description: "Steam Identifier",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
					},
				},
			},
		},
		{
			Name:        "game",
			Description: "Fetches game information",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "search",
					Description: "Fetches information about a game",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "value",
							Description: "Game Identifier",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
					},
				},
				{
					Name:        "player-count",
					Description: "Fetches player count",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "value",
							Description: "Game Identifier",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
					},
				},
				{
					Name:        "news",
					Description: "Fetches latest news about a game",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "value",
							Description: "Game Identifier",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
					},
				},
			},
		},
	}

	commandHandler = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"player": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			for _, o := range i.ApplicationCommandData().Options {
				v := o.Options[0].StringValue()
				switch o.Name {
				case "profile":
					player.PlayerProfile(s, i, steamClient, v)
				case "games":
					player.PlayerGames(s, i, steamClient, v)
				case "friends":
					player.PlayerFriends(s, i, steamClient, v)
				case "bans":
					player.PlayerBans(s, i, steamClient, v)
				case "id":
					player.PlayerID(s, i, steamClient, v)
				}
			}
		},
		"game": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			for _, o := range i.ApplicationCommandData().Options {
				v := o.Options[0].StringValue()
				switch o.Name {
				case "search":
					game.AppSearch(s, i, steamClient, v)
				case "player-count":
					game.AppPlayerCount(s, i, steamClient, v)
				case "news":
					game.AppNews(s, i, steamClient, v)
				}
			}
		},
	}
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
}

func init() {
	env := os.Getenv("BOT_ENV")
	if env == "" {
		env = "development"
	}

	logrus.Info("loading configurations...")
	err := godotenv.Load(".env." + env)
	if err != nil {
		logrus.Fatalf("error loading .env file: %s", err)
	}
	logrus.Infof("launching in %s mode...", env)

	steamToken = os.Getenv("STEAM_API_KEY")
	discordToken = os.Getenv("DISCORD_BOT_TOKEN")

	fmt.Println(discordToken)

	steamClient = steam.Steam{Key: steamToken}
}

func init() {
	logrus.Info("creating Discord session...")
	var err error
	discordSession, err = discordgo.New("Bot " + discordToken)
	if err != nil {
		logrus.Fatalf("error creating Discord session: %s", err)
	}

	discordSession.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandler[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	discordSession.AddHandler(func(s *discordgo.Session, event *discordgo.Ready) {
		logrus.Infof("logging in as %s#%s", s.State.User.Username, s.State.User.Discriminator)
	})
}

func main() {
	logrus.Info("opening websocket connection to Discord...")
	err := discordSession.Open()
	if err != nil {
		logrus.Fatalf("error opening connection: %s", err)
	}

	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := discordSession.ApplicationCommandCreate(discordSession.State.User.ID, "", v)
		logrus.Infof("creating command: %s", v.Name)
		if err != nil {
			logrus.Fatalf("cannot create %s command: %s", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	err = discordSession.UpdateStatusComplex(discordgo.UpdateStatusData{
		Activities: []*discordgo.Activity{
			{
				Name: "Slash Commands",
				Type: discordgo.ActivityTypeListening,
			},
		},
	})

	if err != nil {
		log.Fatalf("cannot set status: %v", err)
	}

	logrus.Info("Steam Stats is now running. Press CTRL+C to exit.")
	defer discordSession.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}
