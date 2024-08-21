package steam

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

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
}

func (s Steam) GetPlayerSummaries(ID ...string) ([]Player, error) {
	baseURL, _ := url.Parse(SteamAPIISteamUser)
	baseURL.Path += "GetPlayerSummaries/v0002"

	params := url.Values{}
	params.Add("key", s.Key)
	params.Add("steamids", strings.Join(ID, ","))
	params.Add("format", "json")
	baseURL.RawQuery = params.Encode()

	resp, err := http.Get(baseURL.String())
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

	var response struct {
		Players struct {
			Players []Player `json:"players"`
		} `json:"response"`
	}

	json.Unmarshal(b, &response)

	// Steam will still return a 200 if the user is not found
	// so we need to check if the response is empty
	if len(response.Players.Players) == 0 {
		return []Player{}, ErrUserNotFound
	}

	return response.Players.Players, nil
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
	baseURL, _ := url.Parse(SteamAPIISteamUser)
	baseURL.Path += "GetPlayerBans/v1"

	params := url.Values{}
	params.Add("key", s.Key)
	params.Add("steamids", p.SteamID)
	params.Add("format", "json")
	baseURL.RawQuery = params.Encode()

	resp, err := http.Get(baseURL.String())
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

	var response struct {
		Players []Player `json:"players"`
	}

	json.Unmarshal(b, &response)
	p.CommunityBanned = response.Players[0].CommunityBanned
	p.VACBanned = response.Players[0].VACBanned
	p.NumOfVacBans = response.Players[0].NumOfVacBans
	p.DaysSinceLastBan = response.Players[0].DaysSinceLastBan
	p.NumOfGameBans = response.Players[0].NumOfGameBans
	p.EconomyBan = response.Players[0].EconomyBan
	return nil
}

func GetBadges(s Steam, p *Player) error {
	baseURL, _ := url.Parse(SteamAPIIPlayerService)
	baseURL.Path += "GetBadges/v1"

	params := url.Values{}
	params.Add("key", s.Key)
	params.Add("steamid", p.SteamID)
	params.Add("format", "json")
	baseURL.RawQuery = params.Encode()

	resp, err := http.Get(baseURL.String())
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

	var response struct {
		Player struct {
			PlayerXP                   int `json:"player_xp"`
			PlayerLevel                int `json:"player_level"`
			PlayerXPNeededToLevelUp    int `json:"player_xp_needed_to_level_up"`
			PlayerXPNeededCurrentLevel int `json:"player_xp_needed_current_level"`
		} `json:"response"`
	}

	json.Unmarshal(b, &response)
	p.PlayerXP = response.Player.PlayerXP
	p.PlayerLevel = response.Player.PlayerLevel
	p.PlayerXPNeededToLevelUp = response.Player.PlayerXPNeededToLevelUp
	p.PlayerXPNeededCurrentLevel = response.Player.PlayerXPNeededCurrentLevel
	return nil
}

func GetSteamLevelDistribution(s Steam, p *Player) error {
	baseURL, _ := url.Parse(SteamAPIIPlayerService)
	baseURL.Path += "GetSteamLevelDistribution/v1"

	params := url.Values{}
	params.Add("key", s.Key)
	params.Add("player_level", strconv.Itoa(p.PlayerLevel))
	params.Add("format", "json")
	baseURL.RawQuery = params.Encode()

	resp, err := http.Get(baseURL.String())
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

	var response struct {
		Player struct {
			PlayerLevelPercentile float64 `json:"player_level_percentile"`
		} `json:"response"`
	}

	json.Unmarshal(b, &response)
	p.PlayerLevelPercentile = response.Player.PlayerLevelPercentile
	return nil
}

// Status returns the player's status as an emoji
func (p Player) Status() string {
	var statusEmoji string
	switch p.PersonaState {
	case 0:
		statusEmoji = "⚫" // Black circle for Offline
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
