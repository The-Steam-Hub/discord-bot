package main

import (
	"os"
	"os/signal"

	"github.com/KevinFagan/steam-stats/cmd"
	"github.com/KevinFagan/steam-stats/steam"
	"github.com/bwmarrin/discordgo"
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
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "stats",
			Description: "Fetches statistics related to one of the subcommand",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "player",
					Description: "Fetches player statistics",
					Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
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
			},
		},
	}

	commandHandler = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"stats": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			for _, st := range i.ApplicationCommandData().Options {
				switch st.Name {
				case "player":
					for _, p := range st.Options {
						switch p.Name {
						case "profile":
							id, _ := steamClient.ResolveID(p.Options[0].StringValue())
							cmd.Profile(s, i, steamClient, id)
						case "games":
							id, _ := steamClient.ResolveID(p.Options[0].StringValue())
							cmd.Games(s, i, steamClient, id)
						case "friends":
							id, _ := steamClient.ResolveID(p.Options[0].StringValue())
							cmd.Friends(s, i, steamClient, id)
						case "bans":
							id, _ := steamClient.ResolveID(p.Options[0].StringValue())
							cmd.Bans(s, i, steamClient, id)
						case "id":
							id, _ := steamClient.ResolveID(p.Options[0].StringValue())
							cmd.Profile(s, i, steamClient, id)
						}
					}
				}
			}
		},
	}
)

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
