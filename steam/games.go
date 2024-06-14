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
	url := fmt.Sprintf("http://api.steampowered.com/IPlayerService/GetOwnedGames/v0001/?key=%s&steamid=%s&format=json&include_appinfo=true", s.Key, ID)
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

func (g *GamesResponse) Count() int {
	return len(g.GamesList.Games)
}

func (g *GamesResponse) MostPlayed() GameStats {
	highest := g.GamesList.Games[0]
	for _, game := range g.GamesList.Games {
		if game.PlayTimeForever > highest.PlayTimeForever {
			highest = game
		}
	}
	return highest
}

func (g *GamesResponse) MostUsedOS() string {
	osMap := make(map[string]int)
	for _, game := range g.GamesList.Games {
		osMap["windows"] += game.PlayTimeWindowsForever
		osMap["mac"] += game.PlayTimeMacForever
		osMap["linux"] += game.PlayTimeLinuxForever
		osMap["deck"] += game.PlayTimeDeckForever
	}

	highest := "windows"
	for os, playtime := range osMap {
		if playtime > osMap[highest] {
			highest = os
		}
	}

	return highest
}

func (g *GamesResponse) WindowsPlaytime() int {
	total := 0
	for _, game := range g.GamesList.Games {
		total += game.PlayTimeWindowsForever
	}
	return total / 60
}

func (g *GamesResponse) MacPlaytime() int {
	total := 0
	for _, game := range g.GamesList.Games {
		total += game.PlayTimeMacForever
	}
	return total / 60
}

func (g *GamesResponse) LinuxPlaytime() int {
	total := 0
	for _, game := range g.GamesList.Games {
		total += game.PlayTimeLinuxForever
	}
	return total / 60
}

func (g *GamesResponse) DeckPlaytime() int {
	total := 0
	for _, game := range g.GamesList.Games {
		total += game.PlayTimeDeckForever
	}
	return total / 60
}

func (g *GamesResponse) TotalHoursPlayed() int {
	total := 0
	for _, game := range g.GamesList.Games {
		total += game.PlayTimeForever
	}
	return total / 60
}

func (g *GamesResponse) GamesPlayed() int {
	played := 0
	for _, game := range g.GamesList.Games {
		if game.PlayTimeForever > 0 {
			played++
		}
	}
	return played
}

func (g *GamesResponse) GamesNotPlayed() int {
	notPlayed := 0
	for _, game := range g.GamesList.Games {
		if game.PlayTimeForever == 0 {
			notPlayed++
		}
	}
	return notPlayed
}
