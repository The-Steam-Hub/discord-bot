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
}

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

func (g *GamesList) Count() int {
	return len(g.Games)
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

func (g *GamesList) MostUsedOS() string {
	osMap := make(map[string]int)
	for _, game := range g.Games {
		osMap["-"] = 0
		osMap["windows"] += game.PlayTimeWindowsForever
		osMap["mac"] += game.PlayTimeMacForever
		osMap["linux"] += game.PlayTimeLinuxForever
		osMap["deck"] += game.PlayTimeDeckForever
	}

	highest := "-"
	for os, playtime := range osMap {
		if playtime > osMap[highest] {
			highest = os
		}
	}

	return highest
}

func (g *GamesList) WindowsPlaytime() int {
	total := 0
	for _, game := range g.Games {
		total += game.PlayTimeWindowsForever
	}
	return total / 60
}

func (g *GamesList) MacPlaytime() int {
	total := 0
	for _, game := range g.Games {
		total += game.PlayTimeMacForever
	}
	return total / 60
}

func (g *GamesList) LinuxPlaytime() int {
	total := 0
	for _, game := range g.Games {
		total += game.PlayTimeLinuxForever
	}
	return total / 60
}

func (g *GamesList) DeckPlaytime() int {
	total := 0
	for _, game := range g.Games {
		total += game.PlayTimeDeckForever
	}
	return total / 60
}

func (g *GamesList) TotalHoursPlayed() int {
	total := 0
	for _, game := range g.Games {
		total += game.PlayTimeForever
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
