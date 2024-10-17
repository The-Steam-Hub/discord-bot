package game

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/the-steam-bot/discord-bot/cmd"
	"github.com/the-steam-bot/discord-bot/steam"
)

func AppNews(session *discordgo.Session, interaction *discordgo.InteractionCreate, steamClient steam.Steam, input string) {
	logs := logrus.Fields{
		"input":  input,
		"author": interaction.Member.User.Username,
		"uuid":   uuid.New(),
	}

	appID, err := steamClient.AppSearch(input)
	if err != nil {
		logs["error"] = err
		errMsg := "game not found"
		logrus.WithFields(logs).Error(errMsg)
		cmd.HandleMessageError(session, interaction, &logs, errMsg)
		return
	}

	appData, err := steamClient.AppDetailedData(appID)
	if err != nil {
		logs["error"] = err
		errMsg := "unable to retrieve game data"
		logrus.WithFields(logs).Error(errMsg)
		cmd.HandleMessageError(session, interaction, &logs, errMsg)
		return
	}

	appNews, err := steamClient.AppNews(appID)
	if err != nil {
		logs["error"] = err
		errMsg := "unable to retrieve game news"
		logrus.WithFields(logs).Error(errMsg)
		cmd.HandleMessageError(session, interaction, &logs, errMsg)
		return
	}

	embMsg := &discordgo.MessageEmbed{
		Title: trimTitle(fmt.Sprintf("%s - %s", appData.Name, appNews.Title)),
		URL:   appNews.URL,
		Image: &discordgo.MessageEmbedImage{
			URL: renderNewsImage(*appData, *appNews),
		},
		Description: trimNewsContents(appNews.Contents),
		Color:       0x66c0f4,
	}

	cmd.HandleMessageOk(embMsg, session, interaction, &logs)
}

func trimTitle(input string) string {
	if len(input) >= 253 {
		return input[:253] + "..."
	}
	return input
}

func renderNewsImage(appData steam.AppDetailedData, appNews steam.AppNews) string {
	regex := regexp.MustCompile(`\[img\](.*)\[\/img\]`)
	matches := regex.FindStringSubmatch(appNews.Contents)

	if len(matches) > 1 {
		newURL := strings.Replace(matches[1], "{STEAM_CLAN_IMAGE}", "https://clan.akamai.steamstatic.com/images/", 1)
		// Some images dont use the {STEAM_CLAN_IMAGE}. When they do not, they do not always start with https:
		if !strings.HasPrefix(newURL, "https:") {
			newURL = "http:" + newURL
		}
		return newURL
	}

	return appData.HeaderImage
}

// This is hacky.... replace with proper tokenization in the future
func trimNewsContents(input string) string {
	replacements := map[*regexp.Regexp]string{
		regexp.MustCompile(`\[b\](.*?)\[/b\]`):           "**$1**",
		regexp.MustCompile(`\[i\](.*?)\[/i\]`):           "*$1*",
		regexp.MustCompile(`\[u\](.*?)\[/u\]`):           "__$1__",
		regexp.MustCompile(`\[strike\](.*?)\[/strike\]`): "~~$1~~",
		regexp.MustCompile(`\[url=(.*?)\](.*?)\[/url\]`): "[$2]($1)",

		// Converting headers to bold
		regexp.MustCompile(`\[h1\](.*?)\[/h1\]`): "**$1**",
		regexp.MustCompile(`\[h2\](.*?)\[/h2\]`): "**$1**",
		regexp.MustCompile(`\[h3\](.*?)\[/h3\]`): "**$1**",

		// Converting lists
		regexp.MustCompile(`\[list\]`):    "",
		regexp.MustCompile(`\[/list\]`):   "",
		regexp.MustCompile(`\[olist\]`):   "",
		regexp.MustCompile(`\[/olist\]`):  "",
		regexp.MustCompile(`\[\*\](.*?)`): "- $1",

		// Delete tags
		regexp.MustCompile(`\[img\](.*?)\[/img\]`):                       "",
		regexp.MustCompile(`\[hr\]\[/hr\]`):                              "",
		regexp.MustCompile(`\[previewyoutube=(.*)\]\[/previewyoutube\]`): "",
	}

	for re, replacement := range replacements {
		input = re.ReplaceAllString(input, replacement)
	}

	// Removing any newlines that might have been added during the tag
	// deletion process
	input = regexp.MustCompile(`(\n{3,})`).ReplaceAllString(input, "\n\n")

	if len(input) > 512 {
		return input[:512] + "..."
	}
	return input
}
