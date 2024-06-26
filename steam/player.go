package steam

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type PlayerResponse struct {
	Response Response `json:"response"`
}

type Response struct {
	Profiles              []Player `json:"players"`
	Badges                []Badge  `json:"badges"`
	PlayerLevel           int      `json:"player_level"`
	PlayerLevelPercentile float64  `json:"player_level_percentile"`
}

type Player struct {
	SteamID               string `json:"steamid"`
	Name                  string `json:"personaname"`
	TimeCreated           int64  `json:"timecreated"`
	CountryCode           string `json:"loccountrycode"`
	StateCode             string `json:"locstatecode"`
	AvatarFull            string `json:"avatarfull"`
	RealName              string `json:"realname"`
	CommunityBanned       bool   `json:"CommunityBanned"`
	VACBanned             bool   `json:"VACBanned"`
	NumOfVacBans          int    `json:"NumberOfVACBans"`
	DaysSinceLastBan      int    `json:"DaysSinceLastBan"`
	NumOfGameBans         int    `json:"NumberOfGameBans"`
	EconomyBan            string `json:"EconomyBan"`
	ProfileURL            string `json:"profileurl"`
	LastLogOff            int    `json:"lastlogoff"`
	PlayerLevel           int
	PlayerLevelPercentile float64
	PersonaState          int
	Badges                []Badge
}

type Badge struct {
	BadgeID        int `json:"badgeid"`
	Level          int `json:"level"`
	CompletionTime int `json:"completion_time"`
	XP             int `json:"xp"`
	Scarcity       int `json:"scarcity"`
}

func (s Steam) Player(ID string) (Player, error) {
	url := fmt.Sprintf("http://api.steampowered.com/ISteamUser/GetPlayerSummaries/v0002/?key=%s&steamids=%s", s.Key, ID)
	resp, err := http.Get(url)
	if err != nil {
		return Player{}, err
	}

	if resp.StatusCode != 200 {
		return Player{}, fmt.Errorf("HTTP request failed with status code %d", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return Player{}, err
	}

	list := PlayerResponse{}
	json.Unmarshal(b, &list)

	// Steam will still return a 200 if the user is not found
	// so we need to check if the response is empty
	if len(list.Response.Profiles) == 0 {
		return Player{}, fmt.Errorf("no player found with ID %s", ID)
	}

	return list.Response.Profiles[0], nil

}

func (s Steam) PlayerWithDetails(ID string) (Player, error) {
	player, err := s.Player(ID)
	if err != nil {
		return Player{}, err
	}
	err = bans(s, &player)
	if err != nil {
		return Player{}, err
	}
	err = badges(s, &player)
	if err != nil {
		return Player{}, err
	}
	err = profileLevel(s, &player)
	if err != nil {
		return Player{}, err
	}
	err = playerLevelPercentile(s, &player)
	if err != nil {
		return Player{}, err
	}
	return player, nil
}

func (p Player) Status() string {
	var statusEmoji string
	switch p.PersonaState {
	case 0:
		statusEmoji = "⚫️" // Black circle for Offline
	case 1:
		statusEmoji = "🟢" // Green circle for Online
	case 2:
		statusEmoji = "🔴" // Red circle for Busy
	case 3:
		statusEmoji = "🟡" // Yellow circle for Away
	case 4:
		statusEmoji = "💤" // Snooze emoji for Snooze
	case 5:
		statusEmoji = "🔄" // Arrow circle emoji for Looking to trade
	case 6:
		statusEmoji = "🎮" // Video game controller emoji for Looking to play
	}
	return statusEmoji
}

func (p Player) ProfileAge() string {
	return UnixToDate(p.TimeCreated)
}

func (p Player) LastSeen() string {
	return UnixToDate(int64(p.LastLogOff))
}

func bans(s Steam, p *Player) error {
	url := fmt.Sprintf("https://api.steampowered.com/ISteamUser/GetPlayerBans/v1/?key=%s&steamids=%s", s.Key, p.SteamID)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTP request failed with status code %d", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	bans := Response{}
	json.Unmarshal(b, &bans)
	p.CommunityBanned = bans.Profiles[0].CommunityBanned
	p.VACBanned = bans.Profiles[0].VACBanned
	p.NumOfVacBans = bans.Profiles[0].NumOfVacBans
	p.DaysSinceLastBan = bans.Profiles[0].DaysSinceLastBan
	p.NumOfGameBans = bans.Profiles[0].NumOfGameBans
	p.EconomyBan = bans.Profiles[0].EconomyBan
	return nil
}

func badges(s Steam, p *Player) error {
	url := fmt.Sprintf("https://api.steampowered.com/IPlayerService/GetBadges/v1/?key=%s&steamid=%s", s.Key, p.SteamID)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTP request failed with status code %d", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	badges := PlayerResponse{}
	json.Unmarshal(b, &badges)
	p.Badges = badges.Response.Badges
	return nil
}

func profileLevel(s Steam, p *Player) error {
	url := fmt.Sprintf("https://api.steampowered.com/IPlayerService/GetSteamLevel/v1/?key=%s&steamid=%s", s.Key, p.SteamID)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTP request failed with status code %d", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	level := PlayerResponse{}
	json.Unmarshal(b, &level)
	p.PlayerLevel = level.Response.PlayerLevel
	return nil
}

func playerLevelPercentile(s Steam, p *Player) error {
	url := fmt.Sprintf("https://api.steampowered.com/IPlayerService/GetSteamLevelDistribution/v1/?key=%s&player_level=%d", s.Key, p.PlayerLevel)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTP request failed with status code %d", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	level := PlayerResponse{}
	json.Unmarshal(b, &level)
	p.PlayerLevelPercentile = level.Response.PlayerLevelPercentile
	return nil
}

func UnixToDate(ut int64) string {
	if ut == 0 {
		return "0y 0d 0h"
	}

	now := time.Now()
	givenTime := time.Unix(int64(ut), 0)
	duration := now.Sub(givenTime)

	// Calculate years, days, and hours. This doesn't account for leap years.
	years := int(duration.Hours() / (24 * 365))
	days := int(duration.Hours()/24) % 365
	hours := int(duration.Hours()) % 24

	return fmt.Sprintf("%dy %dd %dh", int(years), int(days), int(hours))
}
