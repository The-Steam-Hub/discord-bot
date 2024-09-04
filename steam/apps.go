package steam

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	"github.com/gocolly/colly"
)

type Steam struct {
	Key string
}

type AppPlayerCount struct {
	Current     int
	Peak24Hour  int
	PeakAllTime int
}

type AppGlobalAchievements struct {
	Name    string  `json:"name"`
	Percent float32 `json:"percent"`
}

type AppData struct {
	Name             string   `json:"name"`
	AppID            int      `json:"steam_appid"`
	ShortDescription string   `json:"short_description"`
	Developers       []string `json:"developers"`
	Publishers       []string `json:"publishers"`
	HeaderImage      string   `json:"header_image"`
	IsFree           bool     `json:"is_free"`
	DLC              []string `json:"dlc"`
	PriceOverview    struct {
		FinalFormatted   string `json:"final_formatted"`
		InitialFormatted string `json:"initial_formatted"`
		DiscountPercent  int    `json:"discount_percent"`
	} `json:"price_overview"`
	ReleaseDate struct {
		ComingSoon bool   `json:"coming_soon"`
		Date       string `json:"date"`
	} `json:"release_date"`
	Genres []struct {
		ID          string `json:"id"`
		Description string `json:"description"`
	} `json:"genres"`
}

type AppPlayTime struct {
	AppID                  int    `json:"appid"`
	Name                   string `json:"name"`
	PlayTimeForever        int    `json:"playtime_forever"`
	PlayTimeWindowsForever int    `json:"playtime_windows_forever"`
	PlayTimeMacForever     int    `json:"playtime_mac_forever"`
	PlayTimeLinuxForever   int    `json:"playtime_linux_forever"`
	PlayTimeDeckForever    int    `json:"playtime_deck_forever"`
	PlayTime2Weeks         int    `json:"playtime_2weeks"`
}

const (
	SteamWebAPI                = "http://api.steampowered.com/"
	SteamPoweredAPI            = "https://store.steampowered.com/"
	SteamCommunityAPI          = "https://steamcommunity.com/"
	SteamChartsAPI             = "https://steamcharts.com/"
	SteamWebAPIIPlayerService  = SteamWebAPI + "IPlayerService/"
	SteamWebAPIISteamUser      = SteamWebAPI + "ISteamUser/"
	SteamWebAPIISteamUserStats = SteamWebAPI + "ISteamUserStats/"
)

var (
	ErrNoAppsProvided = errors.New("no apps provided")
	ErrUserNotFound   = errors.New("player not found")
	ErrAppNotFound    = errors.New("app not found")
)

var (
	AppURLRegex = `https:\/\/store.steampowered.com\/app\/(\d+)\/`
)

func (s Steam) AppIDResolve(input string) (int, error) {
	if _, err := strconv.ParseUint(input, 10, 32); err == nil {
		formattedID, _ := strconv.Atoi(input)
		return formattedID, nil
	}

	re := regexp.MustCompile(`\/app\/(\d+)\/`)
	match := re.FindStringSubmatch(input)
	if len(match) > 1 {
		formattedID, _ := strconv.Atoi(match[1])
		return formattedID, nil
	}

	return s.AppSearch(input)

}

func (s Steam) AppsOwned(steamID string) (*[]AppPlayTime, error) {
	baseURL, _ := url.Parse(SteamWebAPIIPlayerService)
	baseURL.Path += "GetOwnedGames/v0001"

	params := url.Values{}
	params.Add("key", s.Key)
	params.Add("steamid", steamID)
	params.Add("format", "json")
	params.Add("include_appinfo", "true")
	params.Add("include_free_games", "true")
	baseURL.RawQuery = params.Encode()

	resp, err := http.Get(baseURL.String())
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP request failed with status code %d", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response struct {
		Games struct {
			PlayTimeStatistics []AppPlayTime `json:"games"`
		} `json:"response"`
	}

	json.Unmarshal(b, &response)
	return &response.Games.PlayTimeStatistics, nil
}

func (s Steam) AppSearch(appName string) (int, error) {
	baseURL, _ := url.Parse(SteamPoweredAPI)
	baseURL.Path += "api/storesearch"

	params := url.Values{}
	params.Add("term", appName)
	params.Add("l", "english")
	params.Add("cc", "US")
	baseURL.RawQuery = params.Encode()

	resp, err := http.Get(baseURL.String())
	if err != nil {
		return 0, err
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var response struct {
		Items []struct {
			ID int `json:"id"`
		} `json:"items"`
	}

	json.Unmarshal(b, &response)
	if len(response.Items) == 0 {
		return -1, ErrAppNotFound
	}

	return response.Items[0].ID, nil
}

func (s Steam) AppGlobalAchievements(appID int) (*[]AppGlobalAchievements, error) {
	baseURL, _ := url.Parse(SteamWebAPIISteamUserStats)
	baseURL.Path += "GetGlobalAchievementPercentagesForApp/v0002"

	params := url.Values{}
	params.Add("format", "json")
	params.Add("gameid", strconv.Itoa(appID))
	baseURL.RawQuery = params.Encode()

	resp, err := http.Get(baseURL.String())
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response struct {
		AchievementPercentages struct {
			AppGlobalAchievements []AppGlobalAchievements `json:"achievements"`
		} `json:"achievementpercentages"`
	}

	json.Unmarshal(b, &response)
	return &response.AchievementPercentages.AppGlobalAchievements, nil
}

func (s Steam) AppPlayerCount(appID int) (*AppPlayerCount, error) {
	c := colly.NewCollector()
	playerCount := AppPlayerCount{}

	c.OnHTML(".app-stat span", func(e *colly.HTMLElement) {
		switch e.Index {
		case 0:
			pc, _ := strconv.Atoi(e.Text)
			playerCount.Current = pc
		case 1:
			pc, _ := strconv.Atoi(e.Text)
			playerCount.Peak24Hour = pc
		case 2:
			pc, _ := strconv.Atoi(e.Text)
			playerCount.PeakAllTime = pc
		}
	})

	var scrapeError error
	c.OnError(func(_ *colly.Response, err error) {
		scrapeError = err
	})

	err := c.Visit(SteamChartsAPI + "app/" + strconv.Itoa(appID))
	if err != nil {
		return nil, err
	}

	c.Wait()

	if scrapeError != nil {
		return nil, scrapeError
	}

	return &playerCount, nil
}

func (s Steam) AppData(appID int) (*AppData, error) {
	baseURL, _ := url.Parse(SteamPoweredAPI)
	baseURL.Path += "api/appdetails"

	params := url.Values{}
	params.Add("appids", strconv.Itoa(appID))
	params.Add("l", "english")
	params.Add("cc", "US")
	baseURL.RawQuery = params.Encode()

	resp, err := http.Get(baseURL.String())
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response map[string]struct {
		AppData AppData `json:"data"`
	}

	json.Unmarshal(b, &response)
	appData := response[strconv.Itoa(appID)].AppData
	return &appData, nil
}

func (s Steam) AppsRecentlyPlayed(steamID string) (*[]AppPlayTime, error) {
	baseURL, _ := url.Parse(SteamWebAPIIPlayerService)
	baseURL.Path += "GetRecentlyPlayedGames/v0001"

	params := url.Values{}
	params.Add("key", s.Key)
	params.Add("steamid", steamID)
	params.Add("format", "json")
	baseURL.RawQuery = params.Encode()

	resp, err := http.Get(baseURL.String())
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP request failed with status code %d", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response struct {
		Games struct {
			PlayTimeStatistics []AppPlayTime `json:"games"`
		} `json:"response"`
	}

	json.Unmarshal(b, &response)
	return &response.Games.PlayTimeStatistics, nil
}

func AppsMostPlayed(appStats []AppPlayTime) (*AppPlayTime, error) {
	if len(appStats) == 0 {
		return nil, ErrNoAppsProvided
	}

	highest := appStats[0]
	for _, game := range appStats {
		if game.PlayTimeForever > highest.PlayTimeForever {
			highest = game
		}
	}

	return &highest, nil
}

func AppsLeastPlayed(appStats []AppPlayTime) (*AppPlayTime, error) {
	if len(appStats) == 0 {
		return nil, ErrNoAppsProvided
	}

	lowest := appStats[0]
	for _, game := range appStats {
		if (game.PlayTimeForever < lowest.PlayTimeForever) && game.PlayTimeForever > 0 {
			lowest = game
		}
	}

	return &lowest, nil
}

func AppsPlayed(appStats []AppPlayTime) []AppPlayTime {
	apps := []AppPlayTime{}
	for _, game := range appStats {
		if game.PlayTimeForever > 0 {
			apps = append(apps, game)
		}
	}
	return apps
}

func AppsNotPlayed(appStats []AppPlayTime) []AppPlayTime {
	apps := []AppPlayTime{}
	for _, game := range appStats {
		if game.PlayTimeForever == 0 {
			apps = append(apps, game)
		}
	}
	return apps
}

func AppsTotalHoursPlayed(appStats []AppPlayTime) int {
	total := 0
	for _, game := range appStats {
		total += game.PlayTimeForever
	}
	return total / 60
}

func AppsRecentHoursPlayed(appStats []AppPlayTime) int {
	total := 0
	for _, game := range appStats {
		total += game.PlayTime2Weeks
	}
	return total / 60
}
