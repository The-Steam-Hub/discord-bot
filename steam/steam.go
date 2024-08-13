package steam

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type Steam struct {
	Key string
}

type VanityResponse struct {
	Vanity Vanity `json:"response"`
}

type Vanity struct {
	SteamID string `json:"steamid"`
}

func (s *Steam) ParseSteamID(input string) (string, error) {
	// Check if the input is already a SteamID64
	if _, err := strconv.ParseUint(input, 10, 64); err == nil {
		return input, nil
	}

	// Check if the input is a Steam URL
	if strings.HasPrefix(input, "https://steamcommunity.com") {
		return s.GetIDFromURL(input), nil
	}

	// Check if the input is SteamID3
	if strings.HasPrefix(input, "[U:1:") {
		return SteamID3ToSteamID64(input)
	}

	// Check if the input is SteamID
	if strings.HasPrefix(input, "STEAM_") {
		return SteamIDToSteamID64(input)
	}

	return "", errors.New("unknown Steam ID format")
}

func (s *Steam) GetIDFromURL(url string) string {
	vanityRegex := regexp.MustCompile(`https:\/\/steamcommunity\.com\/id\/([^\/]+)`)
	IDRegex := regexp.MustCompile(`https:\/\/steamcommunity\.com\/profiles\/(\d+)`)

	vanityMatch := vanityRegex.FindStringSubmatch(url)
	IDMatch := IDRegex.FindStringSubmatch(url)

	steamID := ""
	if len(vanityMatch) > 1 {
		vanity, err := s.ResolveVanityURL(vanityMatch[1])
		if err != nil {
			return ""
		}
		steamID = vanity.SteamID
	}
	if len(IDMatch) > 1 {
		steamID = IDMatch[1]
	}

	return steamID
}

func (s *Steam) ResolveVanityURL(vanityURL string) (Vanity, error) {
	url := fmt.Sprintf("https://api.steampowered.com/ISteamUser/ResolveVanityURL/v1/?key=%s&vanityurl=%s", s.Key, vanityURL)
	resp, err := http.Get(url)
	if err != nil {
		return Vanity{}, err
	}

	if resp.StatusCode != 200 {
		return Vanity{}, fmt.Errorf("HTTP request failed with status code %d", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return Vanity{}, err
	}

	VanityResponse := VanityResponse{}
	json.Unmarshal(b, &VanityResponse)
	return VanityResponse.Vanity, nil
}

func SteamID64ToSteamID(steamID64 uint64) string {
	universe := (steamID64 >> 56) & 0xFF
	accountID := steamID64 & 0xFFFFFFFF
	authServer := accountID % 2
	accountNumber := (accountID - authServer) / 2
	return fmt.Sprintf("STEAM_%d:%d:%d", universe, authServer, accountNumber)
}

func SteamID64ToSteamID3(steamID64 uint64) string {
	accountID := steamID64 & 0xFFFFFFFF
	return fmt.Sprintf("[U:1:%d]", accountID)
}

func SteamID3ToSteamID64(steamID3 string) (string, error) {
	var accountID uint32
	_, err := fmt.Sscanf(steamID3, "[U:1:%d]", &accountID)
	if err != nil {
		return "", err
	}
	steamID64 := uint64(0x110000100000000) | uint64(accountID)
	return strconv.FormatUint(steamID64, 10), nil
}

func SteamIDToSteamID64(steamID string) (string, error) {
	var universe, authServer, accountNumber uint32
	_, err := fmt.Sscanf(steamID, "STEAM_%d:%d:%d", &universe, &authServer, &accountNumber)
	if err != nil {
		return "", err
	}
	accountID := accountNumber*2 + authServer
	steamID64 := (uint64(universe) << 56) | (1 << 52) | (1 << 32) | uint64(accountID)

	return strconv.FormatUint(steamID64, 10), nil
}
