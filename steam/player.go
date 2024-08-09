package steam

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type PlayerResponse struct {
	Response Response `json:"response"`
}

type Response struct {
	Profiles                   []Player `json:"players"`
	Badges                     []Badge  `json:"badges"`
	PlayerXP                   int      `json:"player_xp"`
	PlayerLevel                int      `json:"player_level"`
	PlayerLevelPercentile      float64  `json:"player_level_percentile"`
	PlayerXPNeededToLevelUp    int      `json:"player_xp_needed_to_level_up"`
	PlayerXPNeededCurrentLevel int      `json:"player_xp_needed_current_level"`
}

type Player struct {
	SteamID                    string `json:"steamid"`
	Name                       string `json:"personaname"`
	TimeCreated                int64  `json:"timecreated"`
	CountryCode                string `json:"loccountrycode"`
	StateCode                  string `json:"locstatecode"`
	AvatarFull                 string `json:"avatarfull"`
	RealName                   string `json:"realname"`
	CommunityBanned            bool   `json:"CommunityBanned"`
	VACBanned                  bool   `json:"VACBanned"`
	NumOfVacBans               int    `json:"NumberOfVACBans"`
	DaysSinceLastBan           int    `json:"DaysSinceLastBan"`
	NumOfGameBans              int    `json:"NumberOfGameBans"`
	EconomyBan                 string `json:"EconomyBan"`
	ProfileURL                 string `json:"profileurl"`
	LastLogOff                 int    `json:"lastlogoff"`
	PlayerXP                   int
	PlayerLevel                int
	PlayerLevelPercentile      float64
	PlayerXPNeededToLevelUp    int
	PlayerXPNeededCurrentLevel int
	PersonaState               int
	Badges                     []Badge
}

type Badge struct {
	BadgeID        int `json:"badgeid"`
	Level          int `json:"level"`
	CompletionTime int `json:"completion_time"`
	XP             int `json:"xp"`
	Scarcity       int `json:"scarcity"`
}

func (s Steam) GetPlayerSummaries(ID ...string) ([]Player, error) {
	cappedID := ID[:min(len(ID), 100)]
	url := fmt.Sprintf("http://api.steampowered.com/ISteamUser/GetPlayerSummaries/v0002/?key=%s&steamids=%s", s.Key, strings.Join(cappedID, ","))
	resp, err := http.Get(url)
	if err != nil {
		return []Player{}, err
	}

	if resp.StatusCode != 200 {
		return []Player{}, fmt.Errorf("HTTP request failed with status code %d", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return []Player{}, err
	}

	list := PlayerResponse{}
	json.Unmarshal(b, &list)

	// Steam will still return a 200 if the user is not found
	// so we need to check if the response is empty
	if len(list.Response.Profiles) == 0 {
		return []Player{}, fmt.Errorf("user not found")
	}

	return list.Response.Profiles, nil
}

func (s Steam) GetPlayerSummariesWithExtra(ID string) (Player, error) {
	player, err := s.GetPlayerSummaries(ID)
	if err != nil {
		return Player{}, err
	}
	err = GetPlayerBans(s, &player[0])
	if err != nil {
		return Player{}, err
	}
	err = GetBadges(s, &player[0])
	if err != nil {
		return Player{}, err
	}
	err = GetSteamLevelDistribution(s, &player[0])
	if err != nil {
		return Player{}, err
	}
	return player[0], nil
}

func GetPlayerBans(s Steam, p *Player) error {
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

func GetBadges(s Steam, p *Player) error {
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
	p.PlayerXP = badges.Response.PlayerXP
	p.PlayerLevel = badges.Response.PlayerLevel
	p.PlayerXPNeededToLevelUp = badges.Response.PlayerXPNeededToLevelUp
	p.PlayerXPNeededCurrentLevel = badges.Response.PlayerXPNeededCurrentLevel
	return nil
}

func GetSteamLevelDistribution(s Steam, p *Player) error {
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

// Status returns the player's status as an emoji
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

// ProfileAge returns the age of the player's profile
//
// Format Example: 18y 0d 0h
func (p Player) ProfileAge() string {
	return UnixToDate(p.TimeCreated)
}

// LastSeen returns the last time the player was seen online
//
// Format Example: 18y 0d 0h
func (p Player) LastSeen() string {
	if p.PersonaState == 0 {
		return UnixToDate(int64(p.LastLogOff))
	}
	return UnixToDate(0)
}

// UnixToDate converts a Unix timestamp to a human-readable string.
// This doesn't account for leap years so it's not 100% accurate
// but it's good enough for this use case.
//
// Example: 1610000000 -> 18y 0d 0h
func UnixToDate(ut int64) string {
	if ut == 0 {
		return "0y 0d 0h"
	}

	now := time.Now()
	givenTime := time.Unix(int64(ut), 0)
	duration := now.Sub(givenTime)

	years := int(duration.Hours() / (24 * 365))
	days := int(duration.Hours()/24) % 365
	hours := int(duration.Hours()) % 24

	return fmt.Sprintf("%dy %dd %dh", int(years), int(days), int(hours))
}
