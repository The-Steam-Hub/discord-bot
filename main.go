package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"

	"github.com/KevinFagan/steam-stats/cmd"
	"github.com/KevinFagan/steam-stats/cmd/player"
	"github.com/KevinFagan/steam-stats/steam"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var s *discordgo.Session

var (
	steamToken   string
	discordToken string
	steamClient  steam.Steam
)

var (
	logs = logrus.Fields{
		"uuid": uuid.New().String(),
	}

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
					Name:        "global-achievements",
					Description: "Fetches global achievement statistics",
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
			},
		},
	}

	commandHandler = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"player": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			for _, o := range i.ApplicationCommandData().Options {
				value := o.Options[0].StringValue()
				switch o.Name {
				case "profile":
					embMsg, err := player.PlayerProfile(steamClient, value)
					cmd.HandleEmbeddedMessage(embMsg, s, i, logs, err)
				case "games":
					embMsg, err := player.PlayerGames(steamClient, value)
					cmd.HandleEmbeddedMessage(embMsg, s, i, logs, err)
				case "friends":
					embMsg, err := player.PlayerFriends(steamClient, value)
					cmd.HandleEmbeddedMessage(embMsg, s, i, logs, err)
				case "bans":
					embMsg, err := player.PlayerBans(steamClient, value)
					cmd.HandleEmbeddedMessage(embMsg, s, i, logs, err)
				case "id":
					embMsg, err := player.PlayerID(steamClient, value)
					cmd.HandleEmbeddedMessage(embMsg, s, i, logs, err)
				}
			}
		},
		"game": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			for _, o := range i.ApplicationCommandData().Options {
				value := o.Options[0].StringValue()
				switch o.Name {
				case "search":
					ad, err := steamClient.AppSearch(value)
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println(ad)
				case "global-achievements":
					intValue, _ := strconv.Atoi(value)
					achievements, err := steamClient.AppGlobalAchievements(intValue)
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println(achievements)
				case "player-count":
					intValue, _ := strconv.Atoi(value)
					pc, err := steamClient.AppPlayerCount(intValue)
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println(pc)
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
	steamClient = steam.Steam{Key: steamToken}
}

func init() {
	logrus.Info("creating Discord session...")
	var err error
	s, err = discordgo.New("Bot " + discordToken)
	if err != nil {
		logrus.Fatalf("error creating Discord session: %s", err)
		return
	}

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandler[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	s.AddHandler(func(s *discordgo.Session, event *discordgo.Ready) {
		logrus.Infof("logging in as %s#%s", s.State.User.Username, s.State.User.Discriminator)
	})
}

func main() {
	logrus.Info("opening websocket connection to Discord...")
	err := s.Open()
	if err != nil {
		logrus.Fatalf("error opening connection: %s", err)
		return
	}

	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, "", v)
		logrus.Infof("creating command: %s", v.Name)
		if err != nil {
			logrus.Fatalf("cannot create %s command: %s", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	err = s.UpdateStatusComplex(discordgo.UpdateStatusData{
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
	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	// for _, v := range registeredCommands {
	// 	logrus.Infof("deleting command: %s", v.Name)
	// 	err := s.ApplicationCommandDelete(s.State.User.ID, "", v.ID)
	// 	if err != nil {
	// 		logrus.Errorf("cannot delete command %s: %s", v.Name, err)
	// 	}
	// }
}
