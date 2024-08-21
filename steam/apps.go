package steam

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Steam struct {
	Key string
}

type AppNews struct{}
type AppData struct{}

type AppPlayTimeStatistics struct {
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
	SteamAPI               = "http://api.steampowered.com/"
	SteamCommunity         = "https://steamcommunity.com"
	SteamAPIIPlayerService = SteamAPI + "IPlayerService/"
	SteamAPIISteamUser     = SteamAPI + "ISteamUser/"
)

var (
	ErrNoAppsProvided = errors.New("no apps provided")
	ErrUserNotFound   = errors.New("player not found")
)

func (s *Steam) AppsOwned(ID string) (*[]AppPlayTimeStatistics, error) {
	baseURL, _ := url.Parse(SteamAPIIPlayerService)
	baseURL.Path += "GetOwnedGames/v0001"

	params := url.Values{}
	params.Add("key", s.Key)
	params.Add("steamid", ID)
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
			PlayTimeStatistics []AppPlayTimeStatistics `json:"games"`
		} `json:"response"`
	}

	json.Unmarshal(b, &response)
	return &response.Games.PlayTimeStatistics, nil
}

func (s *Steam) AppsRecentlyPlayed(ID string) (*[]AppPlayTimeStatistics, error) {
	baseURL, _ := url.Parse(SteamAPIIPlayerService)
	baseURL.Path += "GetRecentlyPlayedGames/v0001"

	params := url.Values{}
	params.Add("key", s.Key)
	params.Add("steamid", ID)
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
			PlayTimeStatistics []AppPlayTimeStatistics `json:"games"`
		} `json:"response"`
	}

	json.Unmarshal(b, &response)
	return &response.Games.PlayTimeStatistics, nil
}

func AppsMostPlayed(appStats []AppPlayTimeStatistics) (*AppPlayTimeStatistics, error) {
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

func AppsLeastPlayed(appStats []AppPlayTimeStatistics) (*AppPlayTimeStatistics, error) {
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

func AppsPlayed(appStats []AppPlayTimeStatistics) []AppPlayTimeStatistics {
	apps := []AppPlayTimeStatistics{}
	for _, game := range appStats {
		if game.PlayTimeForever > 0 {
			apps = append(apps, game)
		}
	}
	return apps
}

func AppsNotPlayed(appStats []AppPlayTimeStatistics) []AppPlayTimeStatistics {
	apps := []AppPlayTimeStatistics{}
	for _, game := range appStats {
		if game.PlayTimeForever == 0 {
			apps = append(apps, game)
		}
	}
	return apps
}

func AppsTotalHoursPlayed(appStats []AppPlayTimeStatistics) int {
	total := 0
	for _, game := range appStats {
		total += game.PlayTimeForever
	}
	return total / 60
}

func AppsRecentHoursPlayed(appStats []AppPlayTimeStatistics) int {
	total := 0
	for _, game := range appStats {
		total += game.PlayTime2Weeks
	}
	return total / 60
}
