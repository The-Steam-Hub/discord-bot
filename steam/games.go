package steam

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type GamesResponse struct {
	GamesList GamesList `json:"response"`
}

type GamesList struct {
	Games []GameStats `json:"games"`
}

type GameStats struct {
	AppID                  int    `json:"appid"`
	Name                   string `json:"name"`
	PlayTimeForever        int    `json:"playtime_forever"`
	PlayTimeWindowsForever int    `json:"playtime_windows_forever"`
	PlayTimeMacForever     int    `json:"playtime_mac_forever"`
	PlayTimeLinuxForever   int    `json:"playtime_linux_forever"`
	PlayTimeDeckForever    int    `json:"playtime_deck_forever"`
	PlayTime2Weeks         int    `json:"playtime_2weeks"`
}

// Games retrives a list of all games owned by the user, incuding free to play games
func (s *Steam) Games(ID string) (GamesList, error) {
	url := fmt.Sprintf("http://api.steampowered.com/IPlayerService/GetOwnedGames/v0001/?key=%s&steamid=%s&format=json&include_appinfo=true&include_played_free_games=true", s.Key, ID)
	resp, err := http.Get(url)
	if err != nil {
		return GamesList{}, err
	}

	if resp.StatusCode != 200 {
		return GamesList{}, fmt.Errorf("HTTP request failed with status code %d", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return GamesList{}, err
	}

	gamesResponse := GamesResponse{}
	json.Unmarshal(b, &gamesResponse)
	return gamesResponse.GamesList, nil
}

// RecentGames retrives a list of games played by the user in the last 2 weeks
func (s *Steam) RecentGames(ID string) (GamesList, error) {
	url := fmt.Sprintf("http://api.steampowered.com/IPlayerService/GetRecentlyPlayedGames/v0001/?key=%s&steamid=%s&format=json", s.Key, ID)
	resp, err := http.Get(url)
	if err != nil {
		return GamesList{}, err
	}

	if resp.StatusCode != 200 {
		return GamesList{}, fmt.Errorf("HTTP request failed with status code %d", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return GamesList{}, err
	}

	gamesResponse := GamesResponse{}
	json.Unmarshal(b, &gamesResponse)
	return gamesResponse.GamesList, nil
}

func (g *GamesList) MostPlayed() GameStats {
	if len(g.Games) != 0 {
		highest := g.Games[0]
		for _, game := range g.Games {
			if game.PlayTimeForever > highest.PlayTimeForever {
				highest = game
			}
		}
		return highest
	}
	return GameStats{}
}

func (g *GamesList) TotalHoursPlayed() int {
	total := 0
	for _, game := range g.Games {
		total += game.PlayTimeForever
	}
	return total / 60
}

func (g *GamesList) HoursPlayed2Weeks() int {
	total := 0
	for _, game := range g.Games {
		total += game.PlayTime2Weeks
	}
	return total / 60
}

func (g *GamesList) GamesPlayed() int {
	played := 0
	for _, game := range g.Games {
		if game.PlayTimeForever > 0 {
			played++
		}
	}
	return played
}

func (g *GamesList) GamesNotPlayed() int {
	notPlayed := 0
	for _, game := range g.Games {
		if game.PlayTimeForever == 0 {
			notPlayed++
		}
	}
	return notPlayed
}
